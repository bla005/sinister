package sinister

import (
	"strings"
)

func isMatch(path string, requestPath string) bool {
	return path == requestPath
}

func isNumeric(c byte) bool {
	if c >= 48 && c <= 57 {
		return true
	}
	return false
}

const (
	Slash              = 47
	LeftSquareBracket  = 91
	RightSquareBracket = 93
)

func validCh(c rune) bool {
	if c == LeftSquareBracket || c == RightSquareBracket || c == Slash || (c >= 97 && c <= 122) {
		return true
	}
	return false
}

func isAZ(c rune) bool {
	if c >= 97 && c <= 122 {
		return true
	}
	return false
}

func isMainPath(path string) bool {
	if len(path) == 1 && path[0] == Slash {
		return true
	}
	return false
}

func stripLastSlash(path string) string {
	if !isMainPath(path) {
		if path[len(path)-1] == 47 {
			return path[:len(path)-1]
		}
	}
	return path
}

func formatRequestPath(path string) (string, []string) {
	if isMainPath(path) {
		return path, []string{}
	}
	var p []string
	var f strings.Builder
	if len(path) > 0 {
		if path[len(path)-1] == Slash {
			path = path[:len(path)-1]
		}
		if path[0] == Slash {
			path = path[1:]
		}
		s := strings.Split(path, "/")
		for i, w := range s {
			if isNumeric(w[0]) {
				p = append(p, w)
				s[i] = "#"
			}
		}
		f.WriteString("/")
		f.WriteString(strings.Join(s, "/"))
	}
	return f.String(), p
}

func validatePath(path string) ([]string, string) {
	path = strings.ReplaceAll(path, " ", "")
	// path = stripLastSlash(path)
	if !(len(path) > 0 && path[0] == Slash) {
		panic("invalid path")
	}
	if path[len(path)-1] == LeftSquareBracket {
		panic("invalid path")
	}

	leftBracketIdx := 0
	rightBracketIdx := 0
	isBracketOpen := false
	slashIdx := 0
	// params := make([]*param, 0)
	params := make([]string, 0)
	var nPath strings.Builder
	for i, c := range path {
		if validCh(c) {
			if c == Slash {
				slashIdx = i
			}
			// search for first left bracket
			if c == LeftSquareBracket && !isBracketOpen {
				if i-slashIdx != 1 {
					// if the char before the left bracket is not a slash
					panic("invalid path: use a slash to end path")
				}
				leftBracketIdx = i
				isBracketOpen = true
			} else if isBracketOpen && i-leftBracketIdx == 1 && !isAZ(c) {
				// if a bracket was opened but next char is not a letter
				panic("invalid path: use a-z letters for param name")
			} else if isBracketOpen && i-leftBracketIdx > 1 && c == LeftSquareBracket {
				// if a left bracket exists and we passed at least 1 char and we found another left bracket
				panic("invalid path: bracket left unclosed")
			} else if isBracketOpen && i-leftBracketIdx > 1 && c == RightSquareBracket {
				// if left bracket exists and we passed at least one char and the current char is a right bracket param is valid
				rightBracketIdx = i
				r = append(r, path[leftBracketIdx+1:rightBracketIdx])
				// params = append(params, &param{name: path[leftBracketIdx+1 : rightBracketIdx], pos: leftBracketIdx + 1})
				if nPath.Len() > 0 {
					nPath.WriteString("#")
				} else {
					nPath.WriteString(path[:leftBracketIdx])
					nPath.WriteString("#")
				}
				leftBracketIdx = 0
				isBracketOpen = false
				rightBracketIdx = 0
			} else if isBracketOpen && i == len(path)-1 {
				// bracket is open and we passed the whole string
				panic("invalid path: closing bracket not found")
			} else if !isBracketOpen && c == RightSquareBracket {
				// there is a square bracket but it was never opened
				panic("invalid path: opening bracket is missing")
			} else if c == Slash && i-slashIdx == 1 {
				// current char is slash and the char before it was a slash too
				panic("invalid path: too many slashes")
			}
		} else {
			panic("invalid path: illegal characters, (a-z, /, [, ]) are allowed")
		}
	}
	if nPath.Len() == 0 {
		nPath.WriteString(path)
	}
	return params, nPath.String()
}

func isReqPathValid(path string) bool {
	if len(path) == 0 {
		return false
	}
	for _, c := range path {
		if !((c >= 97 && c <= 122) || c == 47 || c >= 48 && c <= 57) {
			return false
		}
	}
	return true
}
