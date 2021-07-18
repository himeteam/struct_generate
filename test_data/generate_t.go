package test_data

type test_embed_struct struct {
	eA int
}

// test_struct struct comment
type test_struct struct {
	test_embed_struct

	a int `tag2:"test1"` // comment1
	B int64
	c string
	D *string
	E bool `json:"e" tag2:"test2"` // comment2
	f []byte
	G map[string]bool
	H []map[interface{}]interface{}
	I map[string][]*string
	J map[string]func(J bool) (r []int, err error)
	K struct{ A bool }
}

type interfaceType interface {
	A()
	B(bool)
	C(p1 bool) bool
	D(p1 bool, p2 bool) (ret bool, err error)
}

type funcType func(p1 bool, p2 bool) (ret bool, err error)

// A1 comment2
func (test_struct) A1() bool {
	return true
}

// A2 comment2
func (test_struct) A2() {
	return
}
