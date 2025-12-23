package parser

import (
	"fmt"
	"strings"

	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

func GoParser(schemas []*schema.Schema, args *utils.Args) []*ParserOutputItem {
	outputItems := []*ParserOutputItem{
		{
			Path:    "minibin__common.go",
			Content: GoCommon,
		},
	}
	for i := range schemas {
		sch := schemas[i]
		if sch.Props != nil {
			outputItems = append(outputItems, parseObject(sch, args))
		} else if sch.Enums != nil {
			outputItems = append(outputItems, parseEnum(sch, args))
		} else {
			fmt.Println(3)
		}
	}
	return outputItems
}

func parseEnum(sch *schema.Schema, args *utils.Args) *ParserOutputItem {
	output := ParserOutputItem{}
	cont := "package minibin\n\nimport \"fmt\"\n\n"
	cont += fmt.Sprintf("type %s string\n\nconst (\n", sch.PascalName)
	fns := fmt.Sprintf(
		"\nfunc (o %s) ToStr() string {\n    switch o {\n",
		sch.PascalName,
	)
	upperName := utils.ToUpperSnakeCase(sch.PascalName)
	longestNameLen := 1
	for i := range sch.Enums {
		enum := sch.Enums[i]
		enum.GoName = upperName + "_" + utils.ToUpperSnakeCase(enum.Name)
		var value string
		if enum.Value != nil {
			value = *enum.Value
		} else {
			value = enum.Name
		}
		if len(enum.GoName) > longestNameLen {
			longestNameLen = len(enum.GoName)
		}
		cont += fmt.Sprintf(
			"    %s$name%s = \"%s\"\n",
			enum.GoName, sch.PascalName, value,
		)
		fns += fmt.Sprintf(
			"    case %s:\n        return \"%s\"\n",
			enum.GoName, value,
		)
	}
	cont += ")\n\n"
	fns += "    default:\n        panic(fmt.Errorf(\"Invalid " + sch.PascalName + ": %s\", o))\n"
	fns += "    }\n}\n"
	for i := range sch.Enums {
		enum := sch.Enums[i]
		nameDelta := longestNameLen - len(enum.GoName) + 1
		nameSpackes := strings.Repeat(" ", nameDelta)
		cont = strings.Replace(
			cont,
			enum.GoName+"$name",
			enum.GoName+nameSpackes,
			1,
		)
	}
	output.Path = strings.ReplaceAll(sch.RPath, ".", "_")
	output.Path = strings.ReplaceAll(output.Path, "/", "_")
	output.Path += ".go"
	output.Content = cont + fns
	return &output
}

func parseObject(sch *schema.Schema, args *utils.Args) *ParserOutputItem {
	output := ParserOutputItem{}
	oStruct := "type " + sch.PascalName + " struct {\n"
	fns := ""
	packFn := fmt.Sprintf(
		"func (o *%s) Pack() []byte {\n    result := []byte{}\n",
		sch.PascalName,
	)
	packWrapperRequired := func(name string, propName string) string {
		return fmt.Sprintf("    result = append(result, %s(o.%s)...)\n", name, propName)
	}
	packWrapperOptional := func(name string, propName string) string {
		return fmt.Sprintf(""+
			"    if o.%s != nil {\n"+
			"        result = append(result, %s(*o.%s)...)\n"+
			"    } else {\n"+
			"        result = append(result, 0)\n"+
			"    }\n",
			propName, name, propName,
		)
	}
	longestNameLen := 1
	longestTypeLen := 1
	for i := range sch.Props {
		prop := sch.Props[i]
		typ := ""
		switch prop.Typ {
		case "string":
			typ = "string"
			if prop.Required {
				packFn += packWrapperRequired("PackString", prop.GoName)
			} else {
				packFn += packWrapperOptional("PackString", prop.GoName)
			}
		case "i32":
			typ = "int32"
			packFn += fmt.Sprintf(
				"    result = append(result, PackInt32(o.%s)...)\n",
				prop.GoName,
			)
		case "i64":
			typ = "int64"
			packFn += fmt.Sprintf(
				"    result = append(result, PackInt64(o.%s)...)\n",
				prop.GoName,
			)
		case "u32":
			typ = "uint32"
			packFn += fmt.Sprintf(
				"    result = append(result, PackUint32(o.%s)...)\n",
				prop.GoName,
			)
		case "u64":
			typ = "uint64"
			packFn += fmt.Sprintf(
				"    result = append(result, PackUint64(o.%s)...)\n",
				prop.GoName,
			)
		case "f32":
			typ = "float32"
			packFn += fmt.Sprintf(
				"    result = append(result, PackFloat32(o.%s)...)\n",
				prop.GoName,
			)
		case "f64":
			typ = "float64"
			packFn += fmt.Sprintf(
				"    result = append(result, PackFloat64(o.%s)...)\n",
				prop.GoName,
			)
		case "bool":
			typ = "bool"
			packFn += fmt.Sprintf(
				"    result = append(result, PackBool(o.%s)...)\n",
				prop.GoName,
			)
		case "object":
			typ = strings.Split(*prop.Ref, ".")[1]
		default:
			panic(fmt.Errorf(
				"Invalid type '%s' found in: %s.props[%d]",
				prop.Typ, sch.RPath, i,
			))
		}
		if !prop.Required {
			typ = "*" + typ
		}
		if prop.Array {
			typ = "[]" + typ
		}
		prop.GoTyp = typ
		if len(prop.GoName) > longestNameLen {
			longestNameLen = len(prop.GoName)
		}
		if len(typ) > longestTypeLen {
			longestTypeLen = len(typ)
		}
		bson := ""
		if args.InjectBson {
			if prop.BsonName != nil {
				bson = fmt.Sprintf(" bson:\"%s\"", *prop.BsonName)
			} else {
				bson = fmt.Sprintf(" bson:\"%s\"", prop.Name)
			}
		}
		oStruct += fmt.Sprintf(
			"    %s$name%s$type`json:\"%s,omitempty\"%s`\n",
			prop.GoName, typ, prop.Name, bson,
		)
		fns += fmt.Sprintf(
			"func (o *%s) Get%s() %s {\n    return o.%s\n}\n",
			sch.PascalName, prop.GoName, typ, prop.GoName,
		)
		fns += fmt.Sprintf(
			"func (o *%s) Set%s(v %s) {\n    o.%s = v\n}\n",
			sch.PascalName, prop.GoName, typ, prop.GoName,
		)
	}
	packFn += "    return result\n}\n"
	oStruct += "}\n"
	for i := range sch.Props {
		prop := sch.Props[i]
		nameDelta := longestNameLen - len(prop.GoName) + 1
		typDelta := longestTypeLen - len(prop.GoTyp) + 1
		base1 := prop.GoName + "$name" + prop.GoTyp
		oStruct = strings.Replace(
			oStruct,
			base1+"$type",
			base1+strings.Repeat(" ", typDelta),
			1,
		)
		oStruct = strings.Replace(
			oStruct,
			prop.GoName+"$name",
			prop.GoName+strings.Repeat(" ", nameDelta),
			1,
		)
	}
	output.Path = strings.ReplaceAll(sch.RPath, ".", "_")
	output.Path = strings.ReplaceAll(output.Path, "/", "_")
	output.Path += ".go"
	output.Content = "package minibin\n\n"
	output.Content += oStruct + fns + packFn
	return &output
}
