package gogen

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func parseNamedCaptureGroupsRegex(textRow string, regexpStr *regexp.Regexp) (map[string]string, error) {
	subexpNames := regexpStr.SubexpNames()
	namedCaptureGroups := make(map[string]string, len(subexpNames))
	matches := regexpStr.FindAllStringSubmatch(textRow, -1)
	if matches == nil {
		return namedCaptureGroups, fmt.Errorf("Unable to parse text '%s' using regexp '%s'\n", textRow, regexpStr)
	}
	for _, match := range matches {
		for groupIdx, group := range match {
			namedCaptureGroups[subexpNames[groupIdx]] = group
		}
	}
	return namedCaptureGroups, nil
}

func ucFirst(s string) string { return strings.Title(s) }

func SnakeToCamel(in string) string {
	tmp := []rune(in)
	tmp[0] = unicode.ToUpper(tmp[0])
	for i := 0; i < len(tmp); i++ {
		if tmp[i] == '_' {
			tmp[i+1] = unicode.ToUpper(tmp[i+1])
			tmp = append(tmp[:i], tmp[i+1:]...)
			i--
		}
	}
	return string(tmp)
}

func CamelToSnake(in string) string {
	tmp := []rune(in)
	sb := strings.Builder{}
	sb.Grow(len(tmp))

	ucSequenceLength := 0 //Special semantics for consecutive UC characters

	for i := 0; i < len(tmp); i++ {
		if unicode.IsUpper(tmp[i]) {
			ucSequenceLength++

			if i == 0 {
				sb.WriteRune(unicode.ToLower(tmp[i]))
			} else if ucSequenceLength == 1 {
				sb.WriteRune('_')
				sb.WriteRune(unicode.ToLower(tmp[i]))
			} else if i+1 >= len(tmp) {
				sb.WriteRune(unicode.ToLower(tmp[i]))
			} else if unicode.IsUpper(tmp[i+1]) {
				sb.WriteRune(unicode.ToLower(tmp[i]))
			} else if ucSequenceLength >= 2 {
				sb.WriteRune('_')
				sb.WriteRune(unicode.ToLower(tmp[i]))
			} else {
				sb.WriteRune('_')
				sb.WriteRune(unicode.ToLower(tmp[i]))
			}
		} else {
			ucSequenceLength = 0
			sb.WriteRune(tmp[i])
		}
	}
	return sb.String()
}

func normalizeMsgDefaultArrayValue(defaults string) string {
	return strings.Trim(defaults, "[]")
}

/*
var splitMsgDefaultArrayValues_re = regexp.MustCompile(`((:?^|,)(:?\s*".*?"\s*|.*?)(:?,|$))`)
var splitMsgDefaultArrayValues_re = regexp.MustCompile(`((?<=^|,)?.*?(?=,|$))`) // Where are lookahead/lookbehind?

func SplitMsgDefaultArrayValues_csv(defaults string) []string {
	csv := csv.NewReader(strings.NewReader(normalizeMsgDefaultArrayValue(defaults)))
	csv.LazyQuotes = true
	csv.TrimLeadingSpace = true
	defaultValues, err := csv.Read()
	if err != nil {
		fmt.Printf("%+v", err)
	}
	return defaultValues
}
func SplitMsgDefaultArrayValues_re(defaults string) []string {
	return splitMsgDefaultArrayValues_re.FindAllString(normalizeMsgDefaultArrayValue(defaults), -1)
}
*/
func SplitMsgDefaultArrayValues(defaults string) []string {
	vals := strings.Split(normalizeMsgDefaultArrayValue(defaults), ",") // This is very unoptimal since it doesn't support quoted fields containing ','
	if len(vals) == 1 && vals[0] == "" {
		return []string{}
	}
	return vals
}

func ValOrNil(val string) string {
	if val == "" {
		return "nil"
	} else {
		return val
	}
}
