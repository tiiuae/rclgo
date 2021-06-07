package gogen

import (
	"container/list"
	"fmt"
	"io/ioutil"
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
	errorTypesCFileMatchingRegexp = "(" + strings.Join(ROS2_ERROR_TYPES_C_FILES, ")|(") + ")"
	re.R(&errorTypesCFileMatchingRegexp, `s!\!!\!!`)
	errorTypesCFileMatchingRegexp = "m!" + errorTypesCFileMatchingRegexp + "!"
}

func GenerateROS2ErrorTypes(rootPath, destPathPkgRoot string) error {
	destFilePath := filepath.Join(destPathPkgRoot, "..", "errortypes.gen.go")

	ros2ErrorsList := list.New()

	includeLookupDir := rootPath
	for tries := 0; tries < 10; tries++ {
		fmt.Printf("Looking for rcl C include files to parse error definitions from '%s'\n", includeLookupDir)

		filepath.Walk(includeLookupDir, func(path string, info os.FileInfo, err error) error {
			if re.M(path, errorTypesCFileMatchingRegexp) {
				fmt.Printf("Analyzing: %s\n", path)
				md, err := GenerateGolangErrorTypesFromROS2ErrorDefinitionsPath(path)
				if err != nil {
					fmt.Printf("Error converting ROS2 Errors from '%s' to '%s', error: %v\n", path, destFilePath, err)
				}
				ros2ErrorsList.PushBackList(md)
			}
			return nil
		})

		if ros2ErrorsList.Len() == 0 {
			includeLookupDir = filepath.Join(includeLookupDir, "..")
			if includeLookupDir == "/" {
				break
			}
		} else {
			break
		}
	}
	if ros2ErrorsList.Len() == 0 {
		fmt.Printf("Unable to find any rcl C error header files?\n")
		return nil
	}

	ros2ErrorsAry := make([]*ROS2ErrorType, ros2ErrorsList.Len())
	i := 0
	for e := ros2ErrorsList.Front(); e != nil; e = e.Next() {
		ros2ErrorsAry[i] = e.Value.(*ROS2ErrorType)
		i++
	}

	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Generating ROS2 Error definitions: %s\n", destFilePath)
	return ros2ErrorCodes.Execute(destFile, map[string]interface{}{
		"ERRORS":                   ros2ErrorsAry,
		"ROS2_ERROR_TYPES_C_FILES": cErrorTypeFiles,
		"DEDUP_FILTER":             ros2errorTypesDeduplicationFilter,
	})
}

func GenerateGolangErrorTypesFromROS2ErrorDefinitionsPath(path string) (*list.List, error) {
	var errorTypes = list.List{}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(content), "\n") {
		errType, err := ParseROS2ErrorType(line)
		if err != nil {
			return nil, err
		}
		if errType != nil {
			errorTypes.PushBack(errType)
		}
	}

	return &errorTypes, nil
}

var ros2errorTypesCommentsBuffer = strings.Builder{}                  // Collect pre-field comments here to be included in the comments. Flushed on empty lines.
var ros2errorTypesDeduplicationMap = make(map[string]string, 1024)    // Some RMW and RCL error codes overlap, so we need to deduplicate them from the dynamic type casting switch-case
var ros2errorTypesDeduplicationFilter = make(map[string]string, 1024) // Entries ending up here actually filter template entries

func ParseROS2ErrorType(row string) (*ROS2ErrorType, error) {
	if re.M(row, `m!
		^
		\#define\s+
		(?P<name>\w+)\s+
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
		return et, nil

	} else if re.M(row, `m!/{2,}\s*(.+)$!`) { // this is a comment line
		if re.R0.S[1] != "" {
			ros2errorTypesCommentsBuffer.WriteString(re.R0.S[1])
		}

		return nil, nil

	} else if row == "" { // do not process empty lines or comment lines
		ros2errorTypesCommentsBuffer.Reset()

		return nil, nil

	}
	return nil, nil
}

func updateROS2errorTypesDeduplicationMap(rcl_ret_t string, name string) {
	_, ok := ros2errorTypesDeduplicationMap[rcl_ret_t]
	if ok {
		ros2errorTypesDeduplicationFilter[name] = rcl_ret_t // If the rcl_ret_t was taken, deduplicate
	} else {
		ros2errorTypesDeduplicationMap[rcl_ret_t] = name // On first match, simply register that the map key is taken
	}
}
