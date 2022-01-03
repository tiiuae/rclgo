/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/kivilahtio/go-re/v0"
)

func ucFirst(s string) string { return strings.Title(s) }

func snakeToCamel(in string) string {
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

func camelToSnake(in string) string {
	tmp := []rune(in)
	sb := strings.Builder{}
	sb.Grow(len(tmp))

	ucSequenceLength := 0 //Special semantics for consecutive UC characters

	for i := 0; i < len(tmp); i++ {
		if unicode.IsUpper(tmp[i]) || (ucSequenceLength > 0 && unicode.IsNumber(tmp[i])) {
			ucSequenceLength++

			if i == 0 {
				sb.WriteRune(unicode.ToLower(tmp[i]))
			} else if ucSequenceLength == 1 {
				sb.WriteRune('_')
				sb.WriteRune(unicode.ToLower(tmp[i]))
			} else if i+1 >= len(tmp) {
				sb.WriteRune(unicode.ToLower(tmp[i]))
			} else if unicode.IsUpper(tmp[i+1]) || unicode.IsNumber(tmp[i+1]) {
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

func commentSerializer(lineComment string, preComments *strings.Builder) string {
	if lineComment != "" && preComments.Len() > 0 {
		return lineComment + `. ` + preComments.String()
	}
	if lineComment != "" {
		return lineComment
	}
	if preComments.Len() > 0 {
		return preComments.String()
	}
	return ""
}

/*
The simple linux mkdir -p without all the Go-fuzz
*/
func mkdir_p(destFilePath string) (*os.File, error) {
	_, err := os.Stat(destFilePath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(filepath.Dir(destFilePath), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return nil, err
	}
	return destFile, nil
}

func normalizeMsgDefaultArrayValue(defaultsField string) string {
	return re.Ss(defaultsField, `s!(?:^\[)|(?:\]$)!!gsm`) // So much fun with regexp love! Accurately trim leading and following [] without possible side-effects
}

func defaultValueSanitizer(ros2type string, defaultValue string) string {
	switch ros2type {
	// CSV parser removes the double quotes only, here we invoke the defaults parsing directly, and need to deal with double quotations manually
	case "string", "wstring", "U16String":
		if defaultValue != "" {
			re.S(&defaultValue, `s!(?:^")|(?:"$)!!gsm`)
		}
	}
	return defaultValueSanitizer_(ros2type, defaultValue)
}

func splitMsgDefaultArrayValues(ros2type string, defaultsField string) []string {
	defaultsField = normalizeMsgDefaultArrayValue(defaultsField)
	csv := csv.NewReader(strings.NewReader(defaultsField))
	csv.LazyQuotes = true
	csv.TrimLeadingSpace = true
	defaultValues, err := csv.Read()
	if err != nil && err != io.EOF {
		fmt.Printf("%+v", err)
	}
	switch ros2type {
	// ROS2 string defaults CAN be quoted differently than how Golang MUST be quoted.
	case "string", "wstring", "U16String":
		for i := range defaultValues {
			defaultValues[i] = defaultValueSanitizer_(ros2type, defaultValues[i])
		}
	}

	return defaultValues
}

func defaultValueSanitizer_(ros2type, defaultValue string) string {
	switch ros2type {
	// ROS2 string defaults CAN be quoted differently than how Golang MUST be quoted.
	case "string", "wstring", "U16String":
		if defaultValue != "" {
			re.S(&defaultValue, `s!(?:^')|(?:'$)!!gsm`)
			re.S(&defaultValue, `s!(?:\\)?"!\"!gsm`)
			re.S(&defaultValue, `s!(?:\\)?'!'!gsm`)
			re.S(&defaultValue, `s!(?:^)|(?:$)!"!gsm`)
		} else {
			defaultValue = "\"\""
		}
	}
	return defaultValue
}

func srvNameFromSrvMsgName(s string) string {
	re.S(&s, `s/_(Request|Response)$//`)
	return s
}

func cReturnCodeNameToGo(n string) string {
	re.S(&n, `s/^RCL_RET_//`)
	re.S(&n, `s/^RMW_RET_/RMW_/`)
	return snakeToCamel(strings.ToLower(n))
}

type stringSet map[string]struct{}

func (s stringSet) Add(str string) {
	s[str] = struct{}{}
}

func (s stringSet) AddFrom(s2 stringSet) {
	for key := range s2 {
		s[key] = struct{}{}
	}
}

/* So many ways to skin a ROS2 defaults field
var splitMsgDefaultArrayValues_re = regexp.MustCompile(`((:?^|,)(:?\s*".*?"\s*|.*?)(:?,|$))`)
var splitMsgDefaultArrayValues_re = regexp.MustCompile(`((?<=^|,)?.*?(?=,|$))`) // Where are lookahead/lookbehind?

func splitMsgDefaultArrayValues_re(defaults string) []string {
	return splitMsgDefaultArrayValues_re.FindAllString(normalizeMsgDefaultArrayValue(defaults), -1)
}
func splitMsgDefaultArrayValues_simple(defaults string) []string {
	vals := strings.Split(normalizeMsgDefaultArrayValue(defaults), ",") // This is very unoptimal since it doesn't support quoted fields containing ','
	if len(vals) == 1 && vals[0] == "" {
		return []string{}
	}
	return vals
}
*/
/*
Blatantly copied from https://golang.org/src/encoding/csv/reader.go
Then further modified to fit the specific ROS2 parsing needs.
The way the ROS2 message defaults are handled is just very unique and parsing it simply as .csv will lose important information about the data type of the
incoming field.
*/
/*
func splitMsgDefaultArrayValues(line string) ([]string, error) {

	var state_inDBQuote bool
	var state_inSQuote bool

	fields := list.New()
	sb := strings.Builder{}
	sb.Grow(10)
	for i, c := range line {
		switch c {
		case '"':
			if state_inDBQuote {

			}

		case "'":

		case '\\':

		case ',':
			fields.PushBack(sb.String())
			sb.Reset()
			sb.Grow(10)
		default:
			sb.WriteRune(c)
		}
		fields.PushBack(sb.String())
		line = bytes.TrimLeftFunc(line, unicode.IsSpace)


		if len(line) == 0 || line[0] != '"' {

			// Non-quoted string field

			i := bytes.IndexRune(line, r.Comma)

			field := line

			if i >= 0 {

				field = field[:i]

			} else {

				field = field[:len(field)-lengthNL(field)]

			}

			// Check to make sure a quote does not appear in field.

			if !r.LazyQuotes {

				if j := bytes.IndexByte(field, '"'); j >= 0 {

					col := utf8.RuneCount(fullLine[:len(fullLine)-len(line[j:])])

					err = &ParseError{StartLine: recLine, Line: r.numLine, Column: col, Err: ErrBareQuote}

					break parseField

				}

			}

			r.recordBuffer = append(r.recordBuffer, field...)

			r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))

			if i >= 0 {

				line = line[i+commaLen:]

				continue parseField

			}

			break parseField

		} else {

			// Quoted string field

			line = line[quoteLen:]

			for {

				i := bytes.IndexByte(line, '"')

				if i >= 0 {

					// Hit next quote.

					r.recordBuffer = append(r.recordBuffer, line[:i]...)

					line = line[i+quoteLen:]

					switch rn := nextRune(line); {

					case rn == '"':

						// `""` sequence (append quote).

						r.recordBuffer = append(r.recordBuffer, '"')

						line = line[quoteLen:]

					case rn == r.Comma:

						// `",` sequence (end of field).

						line = line[commaLen:]

						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))

						continue parseField

					case lengthNL(line) == len(line):

						// `"\n` sequence (end of line).

						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))

						break parseField

					case r.LazyQuotes:

						// `"` sequence (bare quote).

						r.recordBuffer = append(r.recordBuffer, '"')

					default:

						// `"*` sequence (invalid non-escaped quote).

						col := utf8.RuneCount(fullLine[:len(fullLine)-len(line)-quoteLen])

						err = &ParseError{StartLine: recLine, Line: r.numLine, Column: col, Err: ErrQuote}

						break parseField

					}

				} else if len(line) > 0 {

					// Hit end of line (copy all data so far).

					r.recordBuffer = append(r.recordBuffer, line...)

					if errRead != nil {

						break parseField

					}

					line, errRead = r.readLine()

					if errRead == io.EOF {

						errRead = nil

					}

					fullLine = line

				} else {

					// Abrupt end of file (EOF or error).

					if !r.LazyQuotes && errRead == nil {

						col := utf8.RuneCount(fullLine)

						err = &ParseError{StartLine: recLine, Line: r.numLine, Column: col, Err: ErrQuote}

						break parseField

					}

					r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))

					break parseField

				}

			}

		}

	}

	if err == nil {

		err = errRead

	}

	// Create a single string and create slices out of it.

	// This pins the memory of the fields together, but allocates once.

	str := string(r.recordBuffer) // Convert to string once to batch allocations

	dst = dst[:0]

	if cap(dst) < len(r.fieldIndexes) {

		dst = make([]string, len(r.fieldIndexes))

	}

	dst = dst[:len(r.fieldIndexes)]

	var preIdx int

	for i, idx := range r.fieldIndexes {

		dst[i] = str[preIdx:idx]

		preIdx = idx

	}

	// Check or update the expected fields per record.

	if r.FieldsPerRecord > 0 {

		if len(dst) != r.FieldsPerRecord && err == nil {

			err = &ParseError{StartLine: recLine, Line: recLine, Err: ErrFieldCount}

		}

	} else if r.FieldsPerRecord == 0 {

		r.FieldsPerRecord = len(dst)

	}

	return dst, err

}
*/
