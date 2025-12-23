package schema

type SchemaProp struct {
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
