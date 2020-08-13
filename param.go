package sinister

import "strconv"

// Param ...
type Param struct {
	Name  string
	Value string
}

// URLParam ...
type URLParam string

// Int ...
func (p URLParam) Int() (int, error) {
	n, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, ErrInvalidParam
	}
	return n, nil
}

// Int64 ...
func (p URLParam) Int64() (int64, error) {
	n, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, ErrInvalidParam
	}
	m := int64(n)
	return m, nil
}

// String ...
func (p URLParam) String() string {
	return string(p)
}

// Bytes ...
func (p URLParam) Bytes() []byte {
	return []byte(p)
}
