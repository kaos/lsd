package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	CF_IGNORE_CASE = (1 << iota)
)

func compare(lhs, rhs interface{}) int {
	return comparef(lhs, rhs, 0)
}

func comparef(lhs, rhs interface{}, flags int) int {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case int:
			return comparef(lhs, float64(rhs), flags)
		case int64:
			return comparef(lhs, float64(rhs), flags)
		case float64:
			switch true {
			case lhs < rhs:
				return -1
			case lhs > rhs:
				return +1
			default:
				return 0
			}
		}
	case int:
		return comparef(float64(lhs), rhs, flags)
	case int64:
		return comparef(float64(lhs), rhs, flags)
	case string:
		switch rhs := rhs.(type) {
		case string:
			if (CF_IGNORE_CASE & flags) > 0 {
				lhs = strings.ToLower(lhs)
				rhs = strings.ToLower(rhs)
			}
			switch true {
			case lhs < rhs:
				return -1
			case lhs > rhs:
				return +1
			default:
				return 0
			}
		}
	case []interface{}:
		// contains test, meaning that we should return -1 if not found, +1 if found, unless it was the only item, in which case we return 0
		if len(lhs) == 0 {
			// special case, empty list equals empty value
			v, ok := rhs.(string)
			if ok && len(v) == 0 {
				return 0
			}
		}
		for _, s := range lhs {
			if s == rhs {
				if len(lhs) == 1 {
					return 0
				}
				return +1
			}
		}
		return -1
	}

	panic(fmt.Sprintf("attempt to compare %T with %T", lhs, rhs))
}

func regexp_match(lhs, rhs interface{}) bool {
	return regexp_matchf(lhs, rhs, "")
}

func regexp_matchf(lhs, rhs interface{}, flags string) bool {
	switch lhs := lhs.(type) {
	case string:
		switch rhs := rhs.(type) {
		case string:
			if len(flags) > 0 {
				rhs = fmt.Sprintf("(?%s:%s)", flags, rhs)
			}
			res, err := regexp.MatchString(rhs, lhs)
			if err != nil {
				panic(err)
			}
			return res
		}
	}

	panic(fmt.Sprintf("regexp on non-string attributes, %T ~ %T", lhs, rhs))
}
