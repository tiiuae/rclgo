package gogen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kivilahtio/go-re/v0"
)

var errorTypesCFileMatchingRegexp string

func init() {
	prepareErrorTypesCFileMatchingRegexp()
}

/*
prepareErrorTypesCFileMatchingRegexp is a convenience function to be more easily able to define the C header files to look for error definitions without needing to fiddle with complex regexp
*/
func prepareErrorTypesCFileMatchingRegexp() {
	errorTypesCFileMatchingRegexp = "(" + strings.Join(cErrorTypeFiles, ")|(") + ")"
	re.R(&errorTypesCFileMatchingRegexp, `s!\!!\!!`)
	errorTypesCFileMatchingRegexp = "m!" + errorTypesCFileMatchingRegexp + "!"
}

func (g *Generator) GenerateROS2ErrorTypes() error {
	destFilePath := filepath.Join(g.config.DestPath, "pkg/rclgo/errortypes.gen.go")
	var errorTypes []*ROS2ErrorType

	for _, includeLookupDir := range g.config.RootPaths {
		for tries := 0; tries < 10; tries++ {
			fmt.Printf("Looking for rcl C include files to parse error definitions from '%s'\n", includeLookupDir)

			filepath.Walk(includeLookupDir, func(path string, info os.FileInfo, err error) error { //nolint:errcheck
				if err == nil && re.M(path, errorTypesCFileMatchingRegexp) {
					fmt.Printf("Analyzing: %s\n", path)
					errorTypes, err = generateGolangErrorTypesFromROS2ErrorDefinitionsPath(errorTypes, path)
					if err != nil {
						fmt.Printf("Error converting ROS2 Errors from '%s' to '%s', error: %v\n", path, destFilePath, err)
					}
				}
				return nil
			})

			if len(errorTypes) == 0 {
				includeLookupDir = filepath.Join(includeLookupDir, "..")
				if includeLookupDir == "/" {
					break
				}
			} else {
				break
			}
		}
	}
	if len(errorTypes) == 0 {
		fmt.Printf("Unable to find any rcl C error header files?\n")
		return nil
	}

	fmt.Printf("Generating ROS2 Error definitions: %s\n", destFilePath)
	return g.generateGoFile(
		destFilePath,
		ros2ErrorCodes,
		templateData{
			"errorTypes":  errorTypes,
			"includes":    cErrorTypeFiles,
			"dedupFilter": ros2errorTypesDeduplicationFilter,
		},
	)
}

func generateGolangErrorTypesFromROS2ErrorDefinitionsPath(errorTypes []*ROS2ErrorType, path string) ([]*ROS2ErrorType, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(content), "\n") {
		errType := parseROS2ErrorType(line)
		if errType != nil {
			errorTypes = append(errorTypes, errType)
		}
	}

	return errorTypes, nil
}

var ros2errorTypesCommentsBuffer = strings.Builder{}                  // Collect pre-field comments here to be included in the comments. Flushed on empty lines.
var ros2errorTypesDeduplicationMap = make(map[string]string, 1024)    // Some RMW and RCL error codes overlap, so we need to deduplicate them from the dynamic type casting switch-case
var ros2errorTypesDeduplicationFilter = make(map[string]string, 1024) // Entries ending up here actually filter template entries

func parseROS2ErrorType(row string) *ROS2ErrorType {
	if re.M(row, `m!
		^
		\#define\s+
		(?P<name>(?:RCL|RMW)_RET_\w+)\s+
		(?:(?P<int>\d+)|(?P<reference>\w+))\s*
		(?://\s*(?P<comment>.+))?
		\s*$
	!x`) {
		et := &ROS2ErrorType{
			Name:      re.R0.Z["name"],
			Rcl_ret_t: re.R0.Z["int"],
			Reference: re.R0.Z["reference"],
			Comment:   commentSerializer(re.R0.Z["comment"], &ros2errorTypesCommentsBuffer),
		}
		ros2errorTypesCommentsBuffer.Reset()
		updateROS2errorTypesDeduplicationMap(et.Rcl_ret_t, et.Name)
		return et

	} else if re.M(row, `m!/{2,}\s*(.+)$!`) { // this is a comment line
		if re.R0.S[1] != "" {
			ros2errorTypesCommentsBuffer.WriteString(re.R0.S[1])
		}

		return nil

	} else if row == "" { // do not process empty lines or comment lines
		ros2errorTypesCommentsBuffer.Reset()

		return nil

	}
	return nil
}

func updateROS2errorTypesDeduplicationMap(rcl_ret_t string, name string) {
	_, ok := ros2errorTypesDeduplicationMap[rcl_ret_t]
	if ok {
		ros2errorTypesDeduplicationFilter[name] = rcl_ret_t // If the rcl_ret_t was taken, deduplicate
	} else {
		ros2errorTypesDeduplicationMap[rcl_ret_t] = name // On first match, simply register that the map key is taken
	}
}
