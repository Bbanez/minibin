package parser_go

import (
	"fmt"
	"strings"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

func Parse(schemas []*schema.Schema, args *utils.Args) []*p.ParserOutputItem {
	outputItems := []*p.ParserOutputItem{
		{
			Path:    "minibin__common.go",
			Content: Common,
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

func parseEnum(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	output := p.ParserOutputItem{}
	cont := "package minibin\n\nimport \"fmt\"\n\n"
	cont += fmt.Sprintf("type %s string\n\nconst (\n", sch.PascalName)
	toStrFn := fmt.Sprintf(
		"\nfunc (o %s) ToStr() string {\n    switch o {\n",
		sch.PascalName,
	)
	fromStrFn := fmt.Sprintf(
		"\n"+
			"func %sFromStr(v string) %s {\n"+
			"    switch v {\n",
		sch.PascalName, sch.PascalName,
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
		toStrFn += fmt.Sprintf(
			"    case %s:\n        return \"%s\"\n",
			enum.GoName, value,
		)
		fromStrFn += fmt.Sprintf(
			"    case \"%s\":\n        return %s\n",
			value, enum.GoName,
		)
	}
	cont += ")\n\n"
	toStrFn += "    default:\n        panic(fmt.Errorf(\"Invalid " + sch.PascalName + ": %s\", o))\n"
	toStrFn += "    }\n}\n"
	fromStrFn += "    default:\n        panic(fmt.Errorf(\"Invalid " + sch.PascalName + ": %s\", v))\n"
	fromStrFn += "    }\n}\n"
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
	output.Path = "enum_" + output.Path + ".go"
	output.Content = cont + toStrFn + fromStrFn
	return &output
}

func parseObject(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	output := p.ParserOutputItem{}
	oStruct := "type " + sch.PascalName + " struct {\n"
	fns := ""
	packFn := fmt.Sprintf(
		"\n"+
			"func (o *%s) Pack() []byte {\n    result := []byte{}\n",
		sch.PascalName,
	)
	setPropFn := fmt.Sprintf(
		"\n"+
			"func (o *%s) SetPropAtPos(pos int, v any) {\n"+
			"    switch pos {\n"+
			"",
		sch.PascalName,
	)

	newFn := fmt.Sprintf(
		"func New%s(\n",
		sch.PascalName,
	)
	newFnBody := ""

	packWrapperRequired := func(name string, propName string, pos int, arr bool, callFn string) string {
		if arr {
			return fmt.Sprintf(
				""+
					"    for i := range o.%s {\n"+
					"        item := o.%s[i]\n"+
					"        result = append(result, %s(item%s, %d)...)\n"+
					"    }\n",
				propName, propName, name, callFn, pos,
			)
		}
		return fmt.Sprintf(
			""+
				"    result = append(result, %s(o.%s%s, %d)...)\n",
			name, propName, callFn, pos,
		)
	}

	packWrapperOptional := func(
		name string,
		propName string,
		pos int,
		arr bool,
		callFn string,
	) string {
		ptr := ""
		if callFn == "" {
			ptr = "*"
		}
		if arr {
			return fmt.Sprintf(""+
				"    for i := range o.%s {\n"+
				"        item := o.%s[i]\n"+
				"        if item != nil {\n"+
				"            result = append(result, %s(%sitem%s, %d)...)\n"+
				"        } else {\n"+
				"            result = append(result, 0)\n"+
				"        }\n"+
				"    }\n",
				propName, propName, name, ptr, callFn, pos,
			)
		}
		return fmt.Sprintf(""+
			"    if o.%s != nil {\n"+
			"        result = append(result, %s(%so.%s%s, %d)...)\n"+
			"    } else {\n"+
			"        result = append(result, 0)\n"+
			"    }\n",
			propName, name, ptr, propName, callFn, pos,
		)
	}

	setPropWrapperNormal := func(propName string, propType string, pos int, arr bool, required bool) string {
		pointer := ""
		if !required {
			pointer = "&"
		}
		if arr {
			return fmt.Sprintf(
				""+
					"    case %d:\n"+
					"        d := v.(%s)\n"+
					"        o.%s = append(o.%s, %sd)\n"+
					"",
				pos, propType, propName, propName, pointer,
			)
		}
		return fmt.Sprintf(
			""+
				"    case %d:\n"+
				"        d := v.(%s)\n"+
				"        o.%s = %sd\n"+
				"",
			pos, propType, propName, pointer,
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
				packFn += packWrapperRequired("PackString", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackString", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "i32":
			typ = "int32"
			if prop.Required {
				packFn += packWrapperRequired("PackInt32", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackInt32", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "i64":
			typ = "int64"
			if prop.Required {
				packFn += packWrapperRequired("PackInt64", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackInt64", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "u32":
			typ = "uint32"
			if prop.Required {
				packFn += packWrapperRequired("PackUint32", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackUint32", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "u64":
			typ = "uint64"
			if prop.Required {
				packFn += packWrapperRequired("PackUint64", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackUint64", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "f32":
			typ = "float32"
			if prop.Required {
				packFn += packWrapperRequired("PackFloat32", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackFloat32", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "f64":
			typ = "float64"
			if prop.Required {
				packFn += packWrapperRequired("PackFloat64", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackFloat64", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "bool":
			typ = "bool"
			if prop.Required {
				packFn += packWrapperRequired("PackBool", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackBool", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "enum":
			typ = strings.Split(*prop.Ref, ".")[1]
			packFn += packWrapperRequired("PackString", prop.GoName, i, prop.Array, ".ToStr()")
			if prop.Array {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        d := v.(string)\n"+
						"        o.%s = append(o.%s, %sFromStr(d))\n"+
						"",
					i, prop.GoName, prop.GoName, typ,
				)
			} else {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        d := v.(string)\n"+
						"        o.%s = %sFromStr(d)\n"+
						"",
					i, prop.GoName, typ,
				)
			}
		case "object":
			typ = strings.Split(*prop.Ref, ".")[1]
			if prop.Required {
				packFn += packWrapperRequired("PackObject", prop.GoName, i, prop.Array, ".Pack()")
			} else {
				packFn += packWrapperOptional("PackObject", prop.GoName, i, prop.Array, ".Pack()")
			}
			if prop.Array {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        d := v.([]byte)\n"+
						"        obj, err := Unpack%s(d)\n"+
						"        if err == nil {\n"+
						"            o.%s = append(o.%s, obj)\n"+
						"        }\n"+
						"",
					i, typ, prop.GoName, prop.GoName,
				)
			} else {
				deref := ""
				if prop.Required {
					deref = "*"
				}
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        d := v.([]byte)\n"+
						"        obj, err := Unpack%s(d)\n"+
						"        if err == nil {\n"+
						"            o.%s = %sobj\n"+
						"        }\n"+
						"",
					i, typ, prop.GoName, deref,
				)
			}
		case "bytes":
			typ = "[]byte"
			if prop.Required {
				packFn += packWrapperRequired("PackBytes", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackBytes", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		default:
			panic(fmt.Errorf(
				"Invalid type '%s' found in: %s.props[%d]",
				prop.Typ, sch.RPath, i,
			))
		}
		if (!prop.Required && prop.Typ != "enum") || (prop.Array && prop.Typ == "object") {
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
	setPropFn += fmt.Sprintf(
		"" +
			"    }\n" +
			"}\n",
	)
	packFn += fmt.Sprintf(
		"" +
			"    res, err := Compress(result)\n" +
			"    if err != nil {\n" +
			"        panic(err)\n" +
			"    }\n" +
			"    return res\n" +
			"}\n",
	)
	unpackFn := fmt.Sprintf(
		"\n"+
			"func Unpack%s(b []byte) (*%s, error) {\n"+
			"    result := %s{}\n"+
			"    err := Unpack(&result, b)\n"+
			"    if err != nil {\n"+
			"       return nil, err\n"+
			"    }\n"+
			"    return &result, nil\n"+
			"}\n",
		sch.PascalName, sch.PascalName, sch.PascalName,
	)
	oStruct += "}\n\n"
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
		newFn += fmt.Sprintf(
			"    %s%s%s,\n",
			prop.Name, strings.Repeat(" ", nameDelta), prop.GoTyp,
		)
		newFnBody += fmt.Sprintf(
			"    o.%s = %s\n",
			prop.GoName, prop.Name,
		)
	}
	newFn += fmt.Sprintf(
		") *%s {\n    o := %s{}\n%s    return &o\n}\n\n",
		sch.PascalName, sch.PascalName, newFnBody,
	)
	output.Path = strings.ReplaceAll(sch.RPath, ".", "_")
	output.Path = strings.ReplaceAll(output.Path, "/", "_")
	output.Path = strings.ReplaceAll(output.Path, "-", "_")
	output.Path = "obj_" + output.Path + ".go"
	output.Content = "package minibin\n\n"
	output.Content += oStruct + newFn + fns + setPropFn + packFn + unpackFn
	return &output
}
