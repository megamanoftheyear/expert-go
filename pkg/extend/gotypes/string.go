package gotypes

type String struct{ str string }

func NewString(str string) *String { return &String{str: str} }
func (str *String) Len() int       { return len(str.str) }
func (str *String) String() string { return str.str }
func (str *String) WithPrefix(prefix string) *String {
	str.str = prefix + str.str

	return str
}

func (str *String) WithSuffix(suffix string) *String {
	str.str += suffix

	return str
}
