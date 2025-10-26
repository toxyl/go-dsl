package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/toxyl/math"
)

func parseParam(value string, types map[string]string) metaParam {
	line := strings.TrimSpace(value)
	if line == "" {
		return metaParam{}
	}
	line = regexp.MustCompile(`\t`).ReplaceAllString(line, " ")
	line = regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")

	parts := strings.Split(line, " ")

	if len(parts) < 3 { // name,  default, description
		return metaParam{}
	}

	// we determine the type from the first part
	typ := types[parts[0]]
	if typ == "" {
		typ = "any"
	}

	if typ == "error" {
		return metaParam{
			name: strings.TrimSpace(parts[0]),
			typ:  strings.TrimSpace(typ),
			desc: strings.TrimSpace(strings.Join(parts[1:], " ")),
		}
	}

	if typ == "bool" {
		return metaParam{
			name: strings.TrimSpace(parts[0]),
			typ:  strings.TrimSpace(typ),
			def:  parseValue(parts[1]),
			desc: strings.TrimSpace(strings.Join(parts[2:], " ")),
		}
	}

	if typ == "string" {
		line = strings.TrimSpace(line)
		startQuote := strings.Index(line[len(parts[0]):], `"`)
		if startQuote >= 0 {
			startQuote += len(parts[0])
			endQuote := strings.Index(line[startQuote+1:], `"`)
			if endQuote >= 0 {
				endQuote += startQuote + 1
				parts[1] = line[startQuote+1 : endQuote]
				// Adjust parts slice to have description after the quoted string
				parts = append([]string{parts[0], parts[1]}, strings.Fields(line[endQuote+1:])...)
			}
		}

		return metaParam{
			name: strings.TrimSpace(parts[0]),
			typ:  strings.TrimSpace(typ),
			def:  parseValue(parts[1]),
			desc: strings.TrimSpace(strings.Join(parts[2:], " ")),
		}
	}

	if len(parts) < 5 {
		return metaParam{}
	}
	unit := strings.TrimSpace(parts[1])
	if unit == "-" || unit == "none" {
		unit = ""
	}
	param := metaParam{
		name: strings.TrimSpace(parts[0]),
		unit: unit,
		typ:  strings.TrimSpace(typ),
		def:  parseValue(parts[3]),
		desc: strings.TrimSpace(strings.Join(parts[4:], " ")),
	}

	// parse min/max/default
	if parts[2] != "" {
		if parts[2] == "-" {
			param.min = nil
			param.max = nil
		} else {
			e := strings.Split(parts[2], "..")
			if len(e) == 2 {
				param.min = parseValue(e[0])
				param.max = parseValue(e[1])
			}
		}
	}

	return param
}

func parseValue(s string) any {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// test for infinity
	if s == "Inf" || s == "+Inf" || s == "-Inf" {
		if s == "+Inf" {
			return math.MaxFloat64
		}
		return math.SmallestNonzeroFloat64
	}

	// try parsing as int
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}

	// try parsing as float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	// try parsing as bool
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}

	// return as string
	return s
}
