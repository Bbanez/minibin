package schema

var SchemaPropAllowedTypes = []string{
	"string",
	"i32",
	"i64",
	"u32",
	"u64",
	"f32",
	"f64",
	"bool",
	"object",
	"enum",
	"bytes",
}

type SchemaProp struct {
	Desc     string
	Name     string
	GoName   string
	Typ      string
	GoTyp    string
	Ref      *string
	Required bool
	Array    bool
	BsonName *string
}

type SchemaEnum struct {
	Name   string
	GoName string
	Value  *string
}

type Schema struct {
	RPath      string
	Name       string
	PascalName string
	Typ        string
	Props      []*SchemaProp
	Enums      []*SchemaEnum
}
