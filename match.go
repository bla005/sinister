package sinister

import (
	"strings"
)

const (
	slash              = 0x2f
	leftSquareBracket  = 0x5b
	rightSquareBracket = 0x5d
	whitespace         = 0x20
)

func encodePath(text string) int {
	r := 0
	for i, c := range text {
		r += (i ^ i*int(c) + i)
	}
	return r
}

func isMatch(a, b string) bool {
	return a == b
}

func isNumeric(n rune) bool {
	if n >= 0x30 && n <= 0x39 {
		return true
	}
	return false
}

func isAZ(c rune) bool {
	if c >= 0x61 && c <= 0x7a {
		return true
	}
	return false
}
func isRuneValid(c rune) bool {
	if c == leftSquareBracket || c == rightSquareBracket || c == slash || isAZ(c) {
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

func validatePath(path string) ([]string, string, int) {
	if !(len(path) > 0 && path[0] == slash) {
		panic("invalid path")
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
					l = 0
					r = 0
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
	return params, nPath.String(), encodePath(nPath.String())
}

// func stripLastSlash(path string) string {
// 	if !isMainPath(path) {
// 		if path[len(path)-1] == 47 {
// 			return path[:len(path)-1]
// 		}
// 	}
// 	return path
// }
//
// func formatRequestPath(path string) (string, []string) {
// 	if isMainPath(path) {
// 		return path, []string{}
// 	}
// 	var p []string
// 	var f strings.Builder
// 	if len(path) > 0 {
// 		if path[len(path)-1] == slash {
// 			path = path[:len(path)-1]
// 		}
// 		if path[0] == slash {
// 			path = path[1:]
// 		}
// 		s := strings.Split(path, "/")
// 		for i, w := range s {
// 			if isNumeric(rune(w[0])) {
// 				p = append(p, w)
// 				s[i] = "#"
// 			}
// 		}
// 		f.WriteString("/")
// 		f.WriteString(strings.Join(s, "/"))
// 	}
// 	return f.String(), p
// }

func formatReqPath(path string) (string, []string, int, bool) {
	if len(path) == 0 || len(path) > 0 && path[0] != slash {
		return "", []string{}, 0, false
	}
	if isMainPath(path) {
		return path, []string{}, 0, true
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
					return "", []string{}, 0, false
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
			return "", []string{}, 0, false
		}
	}
	if r.Len() == 0 {
		r.WriteString(path)
	}
	rs := r.String()
	if rs[len(rs)-1] == slash {
		rs = r.String()[:r.Len()-1]
	}
	return rs, p, encodePath(rs), true
}

// func isReqPathValid(path string) bool {
// 	if len(path) == 0 {
// 		return false
// 	}
// 	for _, c := range path {
// 		// if !((c >= 97 && c <= 122) || c == 47 || c >= 48 && c <= 57) {
// 		if isAZ(c) || isNumeric(c) || c == slash {
// 			return true
// 		}
// 	}
// 	return false
// }
