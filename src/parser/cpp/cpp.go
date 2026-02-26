package parser_cpp

import (
	"strings"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

func Parse(schemas []*schema.Schema, args *utils.Args) []*p.ParserOutputItem {
	hFile := CommonFunctionsH + "\n\n"
	cFile := ""
	for _, sch := range schemas {
		if sch.Props != nil {
			h, c := parseObject(sch)
			hFile += h + "\n\n"
			cFile += c + "\n\n"
		}
	}
	hFile += "#endif // MINIBIN_H"
	outputItems := []*p.ParserOutputItem{
		{
			Path:    "minibin.hpp",
			Content: hFile,
		},
		{
			Path:    "minibin.cpp",
			Content: cFile,
		},
	}
	return outputItems
}

func parseObject(sch *schema.Schema) (string, string) {
	cFile := ""
	aName := sch.PascalName
	aConstructorArgs := []string{}
	aProps := []string{}
	for _, prop := range sch.Props {
		typ := ""
		switch prop.Typ {
		case "string":
			typ = "std::string"
		case "i32":
			typ = "int32_t"
		case "i64":
			typ = "int64_t"
		case "u32":
			typ = "uint32_t"
		case "u64":
			typ = "uint64_t"
		case "f32":
			typ = "float"
		case "f64":
			typ = "double"
		case "bool":
			typ = "bool"
		case "enum":
			typ = strings.Split(*prop.Ref, ".")[1]
		case "object":
			typ = strings.Split(*prop.Ref, ".")[1]
		case "bytes":
			typ = "std::vector<uint8_t>"
		}
		if prop.Array {
			typ = "std::vector<" + typ + ">"
		}
		if !prop.Required {
			typ += "*"
		}
		aConstructorArgs = append(aConstructorArgs, typ+" "+prop.Name)
		aProps = append(aProps, "    "+typ+" "+prop.Name+";")
	}
	hFile := strings.ReplaceAll(HClass, "@name", aName)
	hFile = strings.ReplaceAll(hFile, "@constructorArgs", strings.Join(aConstructorArgs, ", "))
	hFile = strings.ReplaceAll(hFile, "@props", strings.Join(aProps, "\n"))
	return hFile, cFile
}
