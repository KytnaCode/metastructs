package testdata

type SubType struct {
	Abc string
	Num float64
}

//go:generate go run github.com/kytnacode/metastructs/cmd to-map -t ToMapBench
type ToMapBench struct {
	Str1, Str2, Str3 string
	Bool             bool
	Float            float64
	Int              int
	Uint             uint
	Slice            []string
	Map              map[float64]string
	AnonymousStruct  struct {
		Field1 string
		Field2 float64
	}
}
