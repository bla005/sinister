package sinister

import (
	"fmt"
	"strings"
)

const (
	slash              = 0x2f
	leftSquareBracket  = 0x5b
	rightSquareBracket = 0x5d
	whitespace         = 0x20
)

func isMatch(a, b string) bool {
	return a == b
}

func isAZ(c rune) bool {
	if c >= 0x61 && c <= 0x7a {
		return true
	}
	return false
}

func isNumeric(n rune) bool {
	if n >= 0x30 && n <= 0x39 {
		return true
	}
	return false
}

func isRuneValid(c rune) bool {
	if c == leftSquareBracket || c == rightSquareBracket || c == slash || isAZ(c) || isNumeric(c) {
		return true
	}
	return false
}
func isRuneValidV2(c rune) bool {
	if c == slash || isAZ(c) || isNumeric(c) {
		return true
	}
	return false
}
func isMainPath(path string) bool {
	if len(path) == 1 && path[0] == slash {
		return true
	}
	return false
}

func validatePath(path, method string) ([]string, string) {
	if len(path) == 0 {
		panic("invalid path: too short")
	}
	if path[0] != slash {
		panic("invalid path: has to start with a slash")
	}

	l := 0
	r := 0
	o := false
	s := 0
	params := make([]string, 0)
	var nPath strings.Builder
	n := 0
	for i, c := range path {
		if isRuneValid(c) {
			if c == slash {
				if c == slash && i-s == 1 {
					panic("invalid path: too many slashes")
				}
				s = i
			}
			if c == leftSquareBracket && i == len(path)-1 {
				panic("invalid path")
			}
			if o {
				switch {
				case i-l == 1 && !isAZ(c):
					panic("invalid path: use a-z letters for param name")
				case i-l > 1 && c == leftSquareBracket:
					panic("invalid path: bracket left unclosed")
				case i-l > 1 && c == rightSquareBracket:
					nPath.WriteString("#")
					n = i + 1
					r = i
					params = append(params, path[l+1:r])
					// l = 0
					// r = 0
					o = false
				case i == len(path)-1:
					panic("invalid path: closing bracket not found")
				}
			} else {
				switch {
				case c == leftSquareBracket:
					if i-s != 1 {
						panic("invalid path: use a slash to start a sub-path")
					}
					l = i
					o = true
					nPath.WriteString(path[n:i])
				case c == rightSquareBracket:
					panic("invalid path: opening bracket is missing")
				case r > 0 && i-r == 1 && c != slash:
					panic("smth wrong")
				}
			}
		} else {
			panic("invalid path: illegal characters, (a-z, /, [, ]) are allowed")
		}
	}
	nPath.WriteString(path[n:])
	if nPath.Len() == 0 {
		nPath.WriteString(path)
	}
	return params, fmt.Sprintf("%s_%s", nPath.String(), method)
}

func validateRequestPath(path, method string) (string, []string, bool) {
	if len(path) == 0 {
		return "", []string{}, false
	}
	if path[0] != slash {
		return "", []string{}, false
	}

	if isMainPath(path) {
		return path, []string{}, true
	}
	digit := 0
	found := false
	slashIndex := 0
	n := 0
	r := strings.Builder{}
	p := []string{}
	for i, c := range path {
		if isRuneValidV2(c) {
			if c == slash {
				if c == slash && i-slashIndex == 1 {
					return "", []string{}, false
				}
				slashIndex = i
			}
			if found {
				if c == slash {
					r.WriteString("#/")
					p = append(p, path[digit:i])
					n = i + 1
					digit = 0
					found = false
				}
			} else {
				if isNumeric(c) {
					found = true
					digit = i
					r.WriteString(path[n:i])
				}
			}
			if isNumeric(c) && i == len(path)-1 {
				r.WriteString("#/")
				p = append(p, path[digit:])
				digit = 0
				found = false
			}
		} else {
			return "", []string{}, false
		}
	}
	if r.Len() == 0 {
		r.WriteString(path)
	}
	rs := r.String()
	if rs[len(rs)-1] == slash {
		rs = r.String()[:r.Len()-1]
	}
	return fmt.Sprintf("%s_%s", rs, method), p, true
}
