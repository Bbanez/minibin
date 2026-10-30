package parser_go

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/string_templates"
	"github.com/bbanez/minibin/src/utils"
)

func validGoIdentifier(name string) bool {
	if name == "" || (!unicode.IsLetter(rune(name[0])) && name[0] != '_') {
		return false
	}
	for _, r := range name[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	_, keyword := map[string]struct{}{
		"break":       {},
		"default":     {},
		"func":        {},
		"interface":   {},
		"select":      {},
		"case":        {},
		"defer":       {},
		"go":          {},
		"map":         {},
		"struct":      {},
		"chan":        {},
		"else":        {},
		"goto":        {},
		"package":     {},
		"switch":      {},
		"const":       {},
		"fallthrough": {},
		"if":          {},
		"range":       {},
		"type":        {},
		"continue":    {},
		"for":         {},
		"import":      {},
		"return":      {},
		"var":         {},
	}[name]
	return !keyword
}

func validSchemaText(value string) bool {
	return !strings.ContainsAny(value, "\r\n`")
}

func goRefType(ref *string, path string) string {
	if ref == nil {
		panic(fmt.Errorf("missing ref in %s", path))
	}
	parts := strings.Split(*ref, ".")
	if len(parts) != 2 || !validGoIdentifier(parts[1]) {
		panic(fmt.Errorf("invalid ref %q in %s", *ref, path))
	}
	return parts[1]
}

func goString(value string) string {
	return fmt.Sprintf("%q", value)
}

func goTagValue(value string) string {
	quoted := goString(value)
	return quoted[1 : len(quoted)-1]
}

func Parse(schemas []*schema.Schema, args *utils.Args) []*p.ParserOutputItem {
	schemas = append(schemas, GoOpenApiSchemas()...)
	outputItems := []*p.ParserOutputItem{
		{
			Path:    "minibin__common.go",
			Content: Common,
		},
		// {
		// 	Path:    "minibin__openapi_model.go",
		// 	Content: GoOpenApiModel,
		// },
		{
			Path:    "minibin__openapi_utils.go",
			Content: GoOpenApiUtils,
		},
	}
	for i := range schemas {
		sch := schemas[i]
		if sch.Props != nil {
			outputItems = append(outputItems, parseObjectV2(sch, args))
		} else if sch.Enums != nil {
			outputItems = append(outputItems, parseEnum(sch, args))
		} else {
			fmt.Println(3)
		}
	}
	return outputItems
}

func parseEnum(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	if !validGoIdentifier(sch.PascalName) {
		panic(fmt.Errorf("invalid Go schema name %q in %s", sch.Name, sch.RPath))
	}
	output := p.ParserOutputItem{}
	cont := "package minibin\n\nimport \"fmt\"\n\n"
	cont += fmt.Sprintf("type %s string\n\nconst (\n", sch.PascalName)
	toStrFn := fmt.Sprintf(
		"\nfunc (o %s) ToStr() string {\n    switch o {\n",
		sch.PascalName,
	)
	fromStrFn := fmt.Sprintf(
		"\n"+
			"func %sFromStr(v string) (%s, error) {\n"+
			"    switch v {\n",
		sch.PascalName, sch.PascalName,
	)
	toOpenApi := fmt.Sprintf(
		"\n"+
			"func %sOpenApi() (string, OpenApiObjectSchema) {\n"+
			"    return EnumToOpenApiSchema(\n"+
			"        \"%s\",\n"+
			"        []string{\n",
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
		toOpenApi += fmt.Sprintf("            %s,\n", goString(value))
		if len(enum.GoName) > longestNameLen {
			longestNameLen = len(enum.GoName)
		}
		cont += fmt.Sprintf("    %s$name%s = %s\n", enum.GoName, sch.PascalName, goString(value))
		toStrFn += fmt.Sprintf("    case %s:\n        return %s\n", enum.GoName, goString(value))
		fromStrFn += fmt.Sprintf("    case %s:\n        return %s, nil\n", goString(value), enum.GoName)
	}
	toOpenApi += "" +
		"        },\n" +
		"    )\n" +
		"}\n"
	cont += ")\n\n"
	toStrFn += "    default:\n        panic(fmt.Errorf(\"Invalid " + sch.PascalName + ": %s\", o))\n"
	toStrFn += "    }\n}\n"
	fromStrFn += "    default:\n        return \"\", fmt.Errorf(\"invalid " + sch.PascalName + ": %s\", v)\n"
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
	output.Content = cont + toStrFn + fromStrFn + toOpenApi
	return &output
}

// Function will be removed in the future. Use parseObjectV2 instead.
//
// Deprecated: This function is deprecated and will be removed in the future. Use parseObjectV2 instead.
func parseObject(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	if !validGoIdentifier(sch.PascalName) {
		panic(fmt.Errorf("invalid Go schema name %q in %s", sch.Name, sch.RPath))
	}
	if len(sch.Props) > 256 {
		panic(fmt.Errorf("Go minibin supports at most 256 properties per object: %s has %d", sch.RPath, len(sch.Props)))
	}
	output := p.ParserOutputItem{}
	goNames := make(map[string]struct{}, len(sch.Props))
	oStruct := "type " + sch.PascalName + " struct {\n"
	fns := ""
	toJsonFn := ""
	packFn := fmt.Sprintf(
		"\n"+
			"func (o *%s) Pack() []byte {\n    result := []byte{}\n",
		sch.PascalName,
	)
	setPropFn := fmt.Sprintf(
		"\n"+
			"func (o *%s) SetPropAtPos(pos int, v any, level string) error {\n"+
			"    switch pos {\n"+
			"",
		sch.PascalName,
	)
	getPropNameFn := fmt.Sprintf(
		"\n"+
			"func (o *%s) GetPropNameAtPos(pos int) string {\n"+
			"    switch pos {\n"+
			"",
		sch.PascalName,
	)
	copyFn := fmt.Sprintf(
		"\n"+
			"func (o *%s) Copy() *%s {\n"+
			"    output := %s{}\n"+
			"",
		sch.PascalName, sch.PascalName, sch.PascalName,
	)
	primitiveCopyFnWrapper := func(propName string, arr bool, required bool, typ string) string {
		if arr {
			if required {
				return fmt.Sprintf(""+
					"    output.%s = make([]%s, len(o.%s))\n"+
					"    copy(output.%s, o.%s)\n"+
					"",
					propName, typ, propName, propName, propName,
				)
			} else {
				return fmt.Sprintf(""+
					"    output.%s = make([]*%s, len(o.%s))\n"+
					"    for i, item := range o.%s {\n"+
					"		if item == nil {\n"+
					"           output.%s[i] = nil\n"+
					"	   } else {\n"+
					"		   itemCopy := *item\n"+
					"           output.%s[i] = &itemCopy\n"+
					"       }\n"+
					"    }\n"+
					"",
					propName, typ, propName, propName, propName, propName,
				)
			}
		} else {
			if required {
				return fmt.Sprintf("    output.%s = o.%s\n", propName, propName)
			} else {
				return fmt.Sprintf(""+
					"    if o.%s != nil {\n"+
					"        itemCopy := *o.%s\n"+
					"        output.%s = &itemCopy\n"+
					"    } else {\n"+
					"        output.%s = nil\n"+
					"    }\n"+
					"",
					propName, propName, propName, propName,
				)
			}
		}
	}

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
				"        if item == nil { panic(\"nil array element cannot be encoded\") }\n"+
				"        result = append(result, %s(%sitem%s, %d)...)\n"+
				"    }\n",
				propName, propName, name, ptr, callFn, pos,
			)
		}
		return fmt.Sprintf(""+
			"    if o.%s != nil {\n"+
			"        result = append(result, %s(%so.%s%s, %d)...)\n"+
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
					"        d, ok := v.(%s)\n"+
					"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
					"        o.%s = append(o.%s, %sd)\n"+
					"",
				pos, propType, propName, propName, pointer,
			)
		}
		return fmt.Sprintf(
			""+
				"    case %d:\n"+
				"        d, ok := v.(%s)\n"+
				"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
				"        o.%s = %sd\n"+
				"",
			pos, propType, propName, pointer,
		)
	}

	importStrconv := true
	toJsonWrapper := func(prop *schema.SchemaProp) {
		jsonKey := prop.Name
		if prop.JsonName != nil {
			jsonKey = *prop.JsonName
		}
		switch prop.Typ {
		case "string":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"quoteStrings",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"quoteStringPointers",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"jsonQuote",
						prop.GoName,
						jsonKey,
						"",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"jsonQuote",
						prop.GoName,
						jsonKey,
						"",
					)
				}
			}
			return
		case "i32":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"int32SliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"int32SliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				}
			}
			return
		case "i64":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"int64SliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"int64SliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				}
			}
			return
		case "u32":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"uint32SliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"uint32SliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				}
			}
			return
		case "u64":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"uint64SliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"uint64SliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				}
			}
			return
		case "f32":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"float32SliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"float32SliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				}
			}
			return
		case "f64":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"float64SliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"float64SliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				}
			}
			return
		case "bool":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"boolSliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"boolSliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"numToStr",
						prop.GoName,
						jsonKey,
						"",
					)
				}
			}
			return
		case "enum":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"enumSliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"enumSliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveReq(
						"jsonQuote",
						prop.GoName,
						jsonKey,
						".ToStr()",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveOpt(
						"jsonQuote",
						prop.GoName,
						jsonKey,
						".ToStr()",
					)
				}
			}
			return
		case "object":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"objSliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"objSliceToStringSlice",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonObjReq(
						prop.GoName,
						jsonKey,
					)
				} else {
					toJsonFn += string_templates.GoToJsonObjOpt(
						prop.GoName,
						jsonKey,
					)
				}
			}
			return
		case "bytes":
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"bytesSliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"bytesSliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonBytesReq(
						prop.GoName,
						jsonKey,
					)
				} else {
					toJsonFn += string_templates.GoToJsonBytesOpt(
						prop.GoName,
						jsonKey,
					)
				}
			}
			return
		case "hmap":
			formatValueFn := ""
			methodCall := ""
			switch *prop.HMapTyp {
			case "string":
				formatValueFn = "jsonQuote"
			case "i32", "i64", "u32", "u64", "f32", "f64", "bool":
				formatValueFn = "numToStr"
			case "enum":
				formatValueFn = "jsonQuote"
				methodCall = ".ToStr()"
			case "object":
				methodCall = ".ToJson()"
			}
			if prop.Array {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"hmapSliceToStringSlice",
					)
				} else {
					toJsonFn += string_templates.GoToJsonPrimitiveArr(
						prop.GoName,
						jsonKey,
						"hmapSliceToStringSliceRef",
					)
				}
			} else {
				if prop.Required {
					toJsonFn += string_templates.GoToJsonHmapReq(
						formatValueFn,
						prop.GoName,
						methodCall,
					)
				} else {
					toJsonFn += string_templates.GoToJsonHmapOpt(
						formatValueFn,
						prop.GoName,
						methodCall,
					)
				}
			}
			return
		}
	}
	_ = toJsonWrapper

	longestNameLen := 1
	longestTypeLen := 1
	for i := range sch.Props {
		prop := sch.Props[i]
		if !validGoIdentifier(prop.GoName) {
			panic(fmt.Errorf("invalid Go property name %q in %s.props[%d]", prop.Name, sch.RPath, i))
		}
		if !validGoIdentifier(prop.Name) {
			panic(fmt.Errorf("invalid Go constructor parameter %q in %s.props[%d]", prop.Name, sch.RPath, i))
		}
		if _, exists := goNames[prop.GoName]; exists {
			panic(fmt.Errorf("duplicate Go property name %q in %s", prop.GoName, sch.RPath))
		}
		goNames[prop.GoName] = struct{}{}
		if !validSchemaText(prop.Desc) || (prop.BsonName != nil && (!validSchemaText(*prop.BsonName) || strings.ContainsAny(*prop.BsonName, "\\\""))) {
			panic(fmt.Errorf("invalid description or BSON name in %s.props[%d]", sch.RPath, i))
		}
		typ := ""
		switch prop.Typ {
		case "string":
			typ = "string"
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			if prop.Required {
				packFn += packWrapperRequired("PackString", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackString", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "i32":
			typ = "int32"
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			if prop.Required {
				packFn += packWrapperRequired("PackInt32", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackInt32", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "i64":
			typ = "int64"
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			if prop.Required {
				packFn += packWrapperRequired("PackInt64", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackInt64", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "u32":
			typ = "uint32"
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			if prop.Required {
				packFn += packWrapperRequired("PackUint32", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackUint32", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "u64":
			typ = "uint64"
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			if prop.Required {
				packFn += packWrapperRequired("PackUint64", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackUint64", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "f32":
			typ = "float32"
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			if prop.Required {
				if prop.Array {
					packFn += fmt.Sprintf(
						""+
							"    for i := range o.%s {\n"+
							"        item := o.%s[i]\n"+
							"        result = append(result, PackFloat32(item, %d, %f)...)\n"+
							"    }\n",
						prop.GoName, prop.GoName, i, prop.Decimals,
					)
				} else {
					packFn += fmt.Sprintf(
						""+
							"    result = append(result, PackFloat32(o.%s, %d, %f)...)\n",
						prop.GoName, i, prop.Decimals,
					)
				}
			} else {
				if prop.Array {
					packFn += fmt.Sprintf(""+
						"    for i := range o.%s {\n"+
						"        item := o.%s[i]\n"+
						"        if item != nil {\n"+
						"            result = append(result, PackFloat32(*item, %d, %f)...)\n"+
						"        }\n"+
						"    }\n",
						prop.GoName, prop.GoName, i, prop.Decimals,
					)
				} else {
					packFn += fmt.Sprintf(""+
						"    if o.%s != nil {\n"+
						"        result = append(result, PackFloat32(*o.%s, %d, %f)...)\n"+
						"    }\n",
						prop.GoName, prop.GoName, i, prop.Decimals,
					)
				}
			}
			pointer := ""
			if !prop.Required {
				pointer = "&"
			}
			if prop.Array {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        ud, ok := v.(int32)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        d := float32(ud) / %f\n"+
						"        o.%s = append(o.%s, %sd)\n"+
						"",
					i, prop.Decimals, prop.GoName, prop.GoName, pointer,
				)
			} else {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        ud, ok := v.(int32)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        d := float32(ud) / %f\n"+
						"        o.%s = %sd\n"+
						"",
					i, prop.Decimals, prop.GoName, pointer,
				)
			}
		case "f64":
			typ = "float64"
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			if prop.Required {
				if prop.Array {
					packFn += fmt.Sprintf(
						""+
							"    for i := range o.%s {\n"+
							"        item := o.%s[i]\n"+
							"        result = append(result, PackFloat64(item, %d, %f)...)\n"+
							"    }\n",
						prop.GoName, prop.GoName, i, prop.Decimals,
					)
				} else {
					packFn += fmt.Sprintf(
						""+
							"    result = append(result, PackFloat64(o.%s, %d, %f)...)\n",
						prop.GoName, i, prop.Decimals,
					)
				}
			} else {
				if prop.Array {
					packFn += fmt.Sprintf(""+
						"    for i := range o.%s {\n"+
						"        item := o.%s[i]\n"+
						"        if item != nil {\n"+
						"            result = append(result, PackFloat64(*item, %d, %f)...)\n"+
						"        }\n"+
						"    }\n",
						prop.GoName, prop.GoName, i, prop.Decimals,
					)
				} else {
					packFn += fmt.Sprintf(""+
						"    if o.%s != nil {\n"+
						"        result = append(result, PackFloat64(*o.%s, %d, %f)...)\n"+
						"    }\n",
						prop.GoName, prop.GoName, i, prop.Decimals,
					)
				}
			}
			pointer := ""
			if !prop.Required {
				pointer = "&"
			}
			if prop.Array {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        ud, ok := v.(int64)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        d := float64(ud) / %f\n"+
						"        o.%s = append(o.%s, %sd)\n"+
						"",
					i, prop.Decimals, prop.GoName, prop.GoName, pointer,
				)
			} else {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        ud, ok := v.(int64)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        d := float64(ud) / %f\n"+
						"        o.%s = %sd\n"+
						"",
					i, prop.Decimals, prop.GoName, pointer,
				)
			}
		case "bool":
			typ = "bool"
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			if prop.Required {
				packFn += packWrapperRequired("PackBool", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackBool", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "enum":
			typ = goRefType(prop.Ref, fmt.Sprintf("%s.props[%d]", sch.RPath, i))
			copyFn += primitiveCopyFnWrapper(prop.GoName, prop.Array, prop.Required, typ)
			packFn += packWrapperRequired("PackEnum", prop.GoName, i, prop.Array, ".ToStr()")
			if prop.Array {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        d, ok := v.(string)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        enum, err := %sFromStr(d)\n"+
						"        if err != nil { return err }\n"+
						"        o.%s = append(o.%s, enum)\n"+
						"",
					i, typ, prop.GoName, prop.GoName,
				)
			} else {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        d, ok := v.(string)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        enum, err := %sFromStr(d)\n"+
						"        if err != nil { return err }\n"+
						"        o.%s = enum\n"+
						"",
					i, typ, prop.GoName,
				)
			}
		case "object":
			typ = goRefType(prop.Ref, fmt.Sprintf("%s.props[%d]", sch.RPath, i))
			if prop.Array {
				copyFn += fmt.Sprintf(""+
					"    output.%s = make([]*%s, len(o.%s))\n"+
					"    for i, item := range o.%s {\n"+
					"        if item != nil { output.%s[i] = item.Copy() }\n"+
					"    }\n"+
					"",
					prop.GoName, typ, prop.GoName, prop.GoName, prop.GoName,
				)
			} else {
				if prop.Required {
					copyFn += fmt.Sprintf("    output.%s = *o.%s.Copy()\n", prop.GoName, prop.GoName)
				} else {
					copyFn += fmt.Sprintf(""+
						"    if o.%s != nil {\n"+
						"        output.%s = o.%s.Copy()\n"+
						"    } else {\n"+
						"        output.%s = nil\n"+
						"    }\n"+
						"",
						prop.GoName, prop.GoName, prop.GoName, prop.GoName,
					)
				}
			}
			if prop.Required {
				packFn += packWrapperRequired("PackObject", prop.GoName, i, prop.Array, ".Pack()")
			} else {
				packFn += packWrapperOptional("PackObject", prop.GoName, i, prop.Array, ".Pack()")
			}
			if prop.Array {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        d, ok := v.([]byte)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        lvl := level+\".\"+o.GetPropNameAtPos(pos)\n"+
						"        obj, err := Unpack%s(d, &lvl)\n"+
						"        if err != nil { return err }\n"+
						"        o.%s = append(o.%s, obj)\n"+
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
						"        d, ok := v.([]byte)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        lvl := level+\".\"+o.GetPropNameAtPos(pos)\n"+
						"        obj, err := Unpack%s(d, &lvl)\n"+
						"        if err != nil { return err }\n"+
						"        o.%s = %sobj\n"+
						"",
					i, typ, prop.GoName, deref,
				)
			}
		case "bytes":
			typ = "[]byte"
			if prop.Array {
				if prop.Required {
					copyFn += fmt.Sprintf("    output.%s = make([][]byte, len(o.%s))\n    for i, item := range o.%s { output.%s[i] = cloneBytes(item) }\n", prop.GoName, prop.GoName, prop.GoName, prop.GoName)
				} else {
					copyFn += fmt.Sprintf("    output.%s = make([]*[]byte, len(o.%s))\n    for i, item := range o.%s { if item != nil { itemCopy := cloneBytes(*item); output.%s[i] = &itemCopy } }\n", prop.GoName, prop.GoName, prop.GoName, prop.GoName)
				}
			} else if prop.Required {
				copyFn += fmt.Sprintf("    output.%s = cloneBytes(o.%s)\n", prop.GoName, prop.GoName)
			} else {
				copyFn += fmt.Sprintf("    if o.%s != nil { itemCopy := cloneBytes(*o.%s); output.%s = &itemCopy }\n", prop.GoName, prop.GoName, prop.GoName)
			}
			if prop.Required {
				packFn += packWrapperRequired("PackBytes", prop.GoName, i, prop.Array, "")
			} else {
				packFn += packWrapperOptional("PackBytes", prop.GoName, i, prop.Array, "")
			}
			setPropFn += setPropWrapperNormal(prop.GoName, typ, i, prop.Array, prop.Required)
		case "hmap":
			typ = "map[string]"
			if *prop.HMapTyp == "object" || *prop.HMapTyp == "enum" {
				if prop.Ref == nil {
					panic(fmt.Errorf(
						"Invalid HMap value type '%s' found in: %s.props[%d]",
						*prop.HMapTyp, sch.RPath, i,
					))
				}
				typ += goRefType(prop.Ref, fmt.Sprintf("%s.props[%d]", sch.RPath, i))
			} else {
				switch *prop.HMapTyp {
				case "string":
					typ += "string"
				case "i32":
					typ += "int32"
				case "i64":
					typ += "int64"
				case "u32":
					typ += "uint32"
				case "u64":
					typ += "uint64"
				case "f32":
					typ += "float32"
				case "f64":
					typ += "float64"
				case "bool":
					typ += "bool"
				case "bytes":
					typ += "[]byte"
				default:
					panic(fmt.Errorf(
						"Invalid HMap value type '%s' found in: %s.props[%d]",
						*prop.HMapTyp, sch.RPath, i,
					))
				}
			}
			// TODO: Implement HMap copy, pack, and set property logic
			// panic("HMap copy, pack, and set property logic not implemented yet")
			valuePacker := ""
			valueUnpacker := ""
			decimalDevider := ""
			derefVal := ""
			switch *prop.HMapTyp {
			case "string":
				valuePacker = "PackString(v, 0)"
				valueUnpacker = "UnpackString(v, 2)"
			case "i32":
				valuePacker = "PackInt32(v, 0)"
				valueUnpacker = "UnpackInt32(v, 2)"
			case "i64":
				valuePacker = "PackInt64(v, 0)"
				valueUnpacker = "UnpackInt64(v, 2)"
			case "u32":
				valuePacker = "PackUint32(v, 0)"
				valueUnpacker = "UnpackUint32(v, 2)"
			case "u64":
				valuePacker = "PackUint64(v, 0)"
				valueUnpacker = "UnpackUint64(v, 2)"
			case "f32":
				valuePacker = fmt.Sprintf("PackFloat32(v, 0, %f)",
					prop.Decimals,
				)
				valueUnpacker = "UnpackInt32(v, 2)"
				decimalDevider = fmt.Sprintf(" / %f", math.Pow(10, float64(prop.Decimals)))
			case "f64":
				valuePacker = fmt.Sprintf("PackFloat64(v, 0, %f)",
					prop.Decimals,
				)
				valueUnpacker = "UnpackInt64(v, 2)"
				decimalDevider = fmt.Sprintf(" / %f", math.Pow(10, float64(prop.Decimals)))
			case "bool":
				valuePacker = "PackBool(v, 0)"
				valueUnpacker = "UnpackBool(v, 2)"
			case "bytes":
				valuePacker = "PackBytes(v, 0)"
				valueUnpacker = "UnpackBytes(v, 2)"
			case "object":
				valuePacker = "v.Pack()"
				valueUnpacker = "Unpack" + goRefType(prop.Ref, fmt.Sprintf("%s.props[%d]", sch.RPath, i)) + "(v, nil)"
				derefVal = "*"
			case "enum":
				valuePacker = "v.Pack()"
				valueUnpacker = "Unpack" + goRefType(prop.Ref, fmt.Sprintf("%s.props[%d]", sch.RPath, i)) + "(v, nil)"
				derefVal = "*"
			default:
				panic(fmt.Errorf(
					"Invalid HMap value type '%s' found in: %s.props[%d]",
					*prop.HMapTyp, sch.RPath, i,
				))
			}
			derefPrt := ""
			prt := ""
			if !prop.Required {
				derefPrt = "*"
				prt = "&"
				packFn += fmt.Sprintf(""+
					"   if o.%s != nil",
					prop.GoName)
			}
			createIfNil := fmt.Sprintf(""+
				"        if o.%s == nil {\n"+
				"            o.%s = %s%s{}\n"+
				"        }\n",
				prop.GoName, prop.GoName, prt, typ,
			)
			packFn += fmt.Sprintf(""+
				"    {\n"+
				"        for k, v := range %so.%s {\n"+
				"            result = append(result, PackHMap(k, %s, %d)...)\n"+
				"        }\n"+
				"    }\n",
				derefPrt, prop.GoName, valuePacker, i)
			setPropFn += fmt.Sprintf(""+
				"    case %d:\n"+
				"        d, ok := v.(map[string][]byte)\n"+
				"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
				createIfNil+
				"        for k, v := range d {\n"+
				"			val, _ := %s\n"+
				"           (%so.%s)[k] = %sval%s\n"+
				"		 }\n"+
				"",
				i, valueUnpacker, derefPrt, prop.GoName, derefVal, decimalDevider,
			)
		default:
			panic(fmt.Errorf(
				"Invalid type '%s' found in: %s.props[%d]",
				prop.Typ, sch.RPath, i,
			))
		}
		toJsonWrapper(prop)
		getPropNameFn += fmt.Sprintf(
			""+
				"    case %d:\n"+
				"        return %s\n"+
				"",
			i, goString(prop.Name),
		)
		if !prop.Required || (prop.Array && prop.Typ == "object") {
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
		requiredTag := ""
		if prop.Required {
			requiredTag = " minibin:\"required\""
		}
		desc := "    "
		if prop.Desc != "" {
			desc = fmt.Sprintf("    // %s\n    ", prop.Desc)
		}
		jsonName := goTagValue(prop.Name)
		if prop.JsonName != nil {
			jsonName = *prop.JsonName
		}
		oStruct += fmt.Sprintf(
			"%s%s$name%s$type`json:\"%s,omitempty\"%s%s`\n",
			desc, prop.GoName, typ, jsonName, bson, requiredTag,
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

	fns += fmt.Sprintf("\n"+
		"func %sFromJson(jsonStr string) (*%s, error) {\n"+
		"    result := %s{}\n"+
		"    err := json.Unmarshal([]byte(jsonStr), &result)\n"+
		"    if err != nil {\n"+
		"        return nil, err\n"+
		"    }\n"+
		"    return &result, nil\n"+
		"}\n"+
		"",
		sch.PascalName, sch.PascalName, sch.PascalName,
	)
	toJsonFn = string_templates.GoToJsonFnWrapper(
		sch.PascalName,
		toJsonFn,
	)
	setPropFn += fmt.Sprintf(
		"" +
			"    default:\n" +
			"        return fmt.Errorf(\"unknown property position %%d for %%s\", pos, level)\n" +
			"    }\n" +
			"    return nil\n" +
			"}\n",
	)
	getPropNameFn += fmt.Sprintf(
		"" +
			"    default:\n" +
			"        return \"__unknown__\"+\"[\"+strconv.FormatInt(int64(pos), 10)+\"]\"\n" +
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
			"func Unpack%s(b []byte, level *string) (*%s, error) {\n"+
			"    if level == nil {\n"+
			"        l := \"%s\"\n"+
			"        level = &l"+
			"    }\n"+
			"    result := %s{}\n"+
			"    err := unpack(&result, b, *level, strings.Count(*level, \".\"))\n"+
			"    if err != nil {\n"+
			"       return nil, err\n"+
			"    }\n"+
			"    return &result, nil\n"+
			"}\n",
		sch.PascalName, sch.PascalName, sch.PascalName, sch.PascalName,
	)
	copyFn += "    return &output\n}\n"
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
	importPacks := "import (\n"
	importPacks += "\t\"encoding/json\"\n"
	importPacks += "\t\"fmt\"\n"
	importPacks += "\t\"strings\"\n"
	if importStrconv {
		importPacks += "\t\"strconv\"\n"
	}
	importPacks += ")\n\n"
	output.Content = "package minibin\n\n" + importPacks
	output.Content += oStruct + newFn + fns + toJsonFn + setPropFn + getPropNameFn + packFn + unpackFn + copyFn
	return &output
}

func packPrimFn(
	prop *schema.SchemaProp,
	packFn string,
	posStr string,
	methodCall string,
) string {
	if prop.Array {
		if prop.Required {
			return string_templates.GoStructPackPrimArrReq(
				methodCall,
				packFn,
				posStr,
				prop.GoName,
			)
		} else {
			return string_templates.GoStructPackPrimArrOpt(
				methodCall,
				packFn,
				posStr,
				prop.GoName,
			)
		}
	} else {
		if prop.Required {
			return string_templates.GoStructPackPrimReq(
				methodCall,
				packFn,
				posStr,
				prop.GoName,
			)
		} else {
			return string_templates.GoStructPackPrimOpt(
				methodCall,
				packFn,
				posStr,
				prop.GoName,
			)
		}
	}
}

func packHmapFn(
	prop *schema.SchemaProp,
	packFn string,
	methodCall string,
	divider string,
) string {
	if prop.Required {
		return string_templates.GoStructPackHmapReq(
			divider,
			methodCall,
			packFn,
			prop.GoName,
		)
	} else {
		return string_templates.GoStructPackHmapOpt(
			divider,
			methodCall,
			packFn,
			prop.GoName,
		)
	}
}

func packProp(sch *schema.Schema, prop *schema.SchemaProp, posStr string) *utils.Array[string] {
	packFnBody := utils.Arr[string]()
	switch prop.Typ {
	case "string":
		packFnBody.Push(
			packPrimFn(prop, "PackString", posStr, ""),
		)
	case "i32":
		packFnBody.Push(
			packPrimFn(prop, "PackInt32", posStr, ""),
		)
	case "i64":
		packFnBody.Push(
			packPrimFn(prop, "PackInt64", posStr, ""),
		)
	case "u32":
		packFnBody.Push(
			packPrimFn(prop, "PackUint32", posStr, ""),
		)
	case "u64":
		packFnBody.Push(
			packPrimFn(prop, "PackUint64", posStr, ""),
		)
	case "bool":
		packFnBody.Push(
			packPrimFn(prop, "PackBool", posStr, ""),
		)
	case "object":
		packFnBody.Push(
			packPrimFn(prop, "PackObject", posStr, ".Pack()"),
		)
	case "f32":
		divider := strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
		if prop.Array {
			if prop.Required {
				packFnBody.Push(
					string_templates.GoStructPackFloatArrReq(
						divider,
						"PackFloat32",
						posStr,
						prop.GoName,
					),
				)
			} else {
				packFnBody.Push(
					string_templates.GoStructPackFloatArrOpt(
						divider,
						"PackFloat32",
						posStr,
						prop.GoName,
					),
				)
			}
		} else {
			if prop.Required {
				packFnBody.Push(
					string_templates.GoStructPackFloatReq(
						divider,
						"PackFloat32",
						posStr,
						prop.GoName,
					),
				)
			} else {
				packFnBody.Push(
					string_templates.GoStructPackFloatOpt(
						divider,
						"PackFloat32",
						posStr,
						prop.GoName,
					),
				)
			}
		}
	case "f64":
		divider := strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
		if prop.Array {
			if prop.Required {
				packFnBody.Push(
					string_templates.GoStructPackFloatArrReq(
						divider,
						"PackFloat64",
						posStr,
						prop.GoName,
					),
				)
			} else {
				packFnBody.Push(
					string_templates.GoStructPackFloatArrOpt(
						divider,
						"PackFloat64",
						posStr,
						prop.GoName,
					),
				)
			}
		} else {
			if prop.Required {
				packFnBody.Push(
					string_templates.GoStructPackFloatReq(
						divider,
						"PackFloat64",
						posStr,
						prop.GoName,
					),
				)
			} else {
				packFnBody.Push(
					string_templates.GoStructPackFloatOpt(
						divider,
						"PackFloat64",
						posStr,
						prop.GoName,
					),
				)
			}
		}
	case "hmap":
		divider := ""
		if prop.Decimals > 1 {
			divider = "," + strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
		}
		if prop.HMapTyp == nil {
			panic(fmt.Errorf(
				"Missing HMap value type found in: %s.props[%s]",
				sch.RPath, posStr,
			))
		}
		rootHmapTyp := "map[string]" + miniTypToGoTyp(
			*prop.HMapTyp, nil, prop.Ref,
		)
		hmapTyp := rootHmapTyp
		if !prop.Required || (prop.Array && prop.Typ == "object") {
			hmapTyp = "*" + hmapTyp
		}
		if prop.Array {
			hmapTyp = "[]" + hmapTyp
		}
		switch *prop.HMapTyp {
		case "string":
			packFnBody.Push(packHmapFn(prop, "PackString", "", divider))
		case "i32":
			packFnBody.Push(packHmapFn(prop, "PackInt32", "", divider))
		case "i64":
			packFnBody.Push(packHmapFn(prop, "PackInt64", "", divider))
		case "u32":
			packFnBody.Push(packHmapFn(prop, "PackUint32", "", divider))
		case "u64":
			packFnBody.Push(packHmapFn(prop, "PackUint64", "", divider))
		case "f32":
			packFnBody.Push(packHmapFn(prop, "PackFloat32", "", divider))
		case "f64":
			packFnBody.Push(packHmapFn(prop, "PackFloat64", "", divider))
		case "bool":
			packFnBody.Push(packHmapFn(prop, "PackBool", "", divider))
		case "bytes":
			packFnBody.Push(packHmapFn(prop, "PackBytes", "", divider))
		case "object":
			packFnBody.Push(packHmapFn(prop, "PackObject", ".Pack()", divider))
		case "enum":
			packFnBody.Push(packHmapFn(prop, "PackEnum", "", divider))
		default:
			panic(fmt.Errorf(
				"Invalid HMap value type '%s' found in: %s.props[%s]",
				*prop.HMapTyp, sch.RPath, posStr,
			))
		}
	}
	return &packFnBody
}

func setPropFnHMapWrapper(prop *schema.SchemaProp, posStr string, unpackFn string) string {
	valueTyp := miniTypToGoTyp(*prop.HMapTyp, nil, prop.Ref)
	if prop.Required {
		return string_templates.GoStructSetPropAtPosCaseHmapPrimReq(
			posStr,
			prop.GoName,
			unpackFn,
		)
	} else {
		return string_templates.GoStructSetPropAtPosCaseHmapPrimOpt(
			posStr,
			prop.GoName,
			valueTyp,
			unpackFn,
		)
	}
}

func setPropFnWrapper(prop *schema.SchemaProp, posStr string, rootTyp string, typ string) *utils.Array[string] {
	setPropCases := utils.Arr[string]()
	primWrapperFn := func(typ string) {
		if prop.Array {
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCasePrimArrReq(
						posStr,
						prop.GoName,
						rootTyp,
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCasePrimArrOpt(
						posStr,
						prop.GoName,
						rootTyp,
					),
				)
			}
		} else {
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCasePrimReq(
						posStr,
						prop.GoName,
						rootTyp,
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCasePrimOpt(
						posStr,
						prop.GoName,
						rootTyp,
					),
				)
			}
		}
	}
	switch prop.Typ {
	case "string", "i32", "i64", "u32", "u64", "bool":
		primWrapperFn(rootTyp)
	case "f32":
		divider := strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
		if prop.Array {
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseFloatArrReq(
						divider,
						"float32",
						posStr,
						prop.GoName,
						"int32",
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseFloatArrOpt(
						divider,
						"float32",
						posStr,
						prop.GoName,
						"int32",
					),
				)
			}
		} else {
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseFloatReq(
						divider,
						"float32",
						posStr,
						prop.GoName,
						"int32",
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseFloatOpt(
						divider,
						"float32",
						posStr,
						prop.GoName,
						"int32",
					),
				)
			}
		}
	case "f64":
		divider := strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
		if prop.Array {
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseFloatArrReq(
						divider,
						"float64",
						posStr,
						prop.GoName,
						"int64",
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseFloatArrOpt(
						divider,
						"float64",
						posStr,
						prop.GoName,
						"int64",
					),
				)
			}
		} else {
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseFloatReq(
						divider,
						"float64",
						posStr,
						prop.GoName,
						"int64",
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseFloatOpt(
						divider,
						"float64",
						posStr,
						prop.GoName,
						"int64",
					),
				)
			}
		}
	case "enum":
		if prop.Array {
			setPropCases.Push(
				string_templates.GoStructSetPropAtPosCaseEnumArr(
					rootTyp,
					posStr,
					prop.GoName,
				),
			)
		} else {
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseEnumReq(
						rootTyp,
						posStr,
						prop.GoName,
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseEnumOpt(
						rootTyp,
						posStr,
						prop.GoName,
					),
				)
			}
		}
	case "object":
		if prop.Array {
			setPropCases.Push(
				string_templates.GoStructSetPropAtPosCaseObjArr(
					rootTyp,
					posStr,
					prop.GoName,
				),
			)
		} else {
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseObjReq(
						rootTyp,
						posStr,
						prop.GoName,
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseObjOpt(
						rootTyp,
						posStr,
						prop.GoName,
					),
				)
			}
		}
	case "hmap":
		hmapRootTyp := "map[string]" + miniTypToGoTyp(*prop.HMapTyp, nil, prop.Ref)
		hmapTyp := hmapRootTyp
		if !prop.Required || (prop.Array && prop.Typ == "object") {
			hmapTyp = "*" + hmapTyp
		}
		if prop.Array {
			hmapTyp = "[]" + hmapTyp
		}
		switch *prop.HMapTyp {
		case "string":
			setPropCases.Push(setPropFnHMapWrapper(
				prop,
				posStr,
				"UnpackString",
			))
		case "i32":
			setPropCases.Push(setPropFnHMapWrapper(
				prop,
				posStr,
				"UnpackInt32",
			))
		case "i64":
			setPropCases.Push(setPropFnHMapWrapper(
				prop,
				posStr,
				"UnpackInt64",
			))
		case "u32":
			setPropCases.Push(setPropFnHMapWrapper(
				prop,
				posStr,
				"UnpackUint32",
			))
		case "u64":
			setPropCases.Push(setPropFnHMapWrapper(
				prop,
				posStr,
				"UnpackUint64",
			))
		case "bool":
			setPropCases.Push(setPropFnHMapWrapper(
				prop,
				posStr,
				"UnpackBool",
			))
		}
	}
	return &setPropCases
}

func copyProp(prop *schema.SchemaProp, rootTyp string, nameIndentStr string) *utils.Array[string] {
	copyFnFields := utils.Arr[string]()
	switch prop.Typ {
	case "string", "i32", "i64", "u32", "u64", "f32", "f64", "bool", "enum":
		if prop.Array {
			if prop.Required {
				copyFnFields.Push(
					string_templates.GoStructCopyFieldPrimArrReq(
						prop.GoName,
						rootTyp,
					),
				)
			} else {
				copyFnFields.Push(
					string_templates.GoStructCopyFieldPrimArrOpt(
						prop.GoName,
						rootTyp,
					),
				)
			}
		} else {
			if prop.Required {
				copyFnFields.Push(
					string_templates.GoStructCopyFieldReq(
						nameIndentStr,
						"",
						prop.GoName,
					),
				)
			} else {
				copyFnFields.Push(
					string_templates.GoStructCopyFieldOpt(
						nameIndentStr,
						"",
						prop.GoName,
					),
				)
			}
		}
	case "object":
		if prop.Array {
			copyFnFields.Push(
				string_templates.GoStructCopyFieldObjArrOpt(
					prop.GoName,
					rootTyp,
				),
			)
		} else {
			if prop.Required {
				copyFnFields.Push(
					string_templates.GoStructCopyFieldReq(
						nameIndentStr+"*",
						".Copy()",
						prop.GoName,
					),
				)
			} else {
				copyFnFields.Push(
					string_templates.GoStructCopyField(
						nameIndentStr,
						".Copy()",
						prop.GoName,
					),
				)
			}
		}
	}
	return &copyFnFields
}

func parseObjectV2(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	structName := sch.PascalName
	copyFnFields := utils.Arr[string]()
	getPropNameAtPosCases := utils.Arr[string]()
	gettersSetters := utils.Arr[string]()
	structFields := utils.Arr[string]()
	structFieldsInit := utils.Arr[string]()
	structFieldsInitParams := utils.Arr[string]()
	setPropAtPosCases := utils.Arr[string]()
	packFnBody := utils.Arr[string]()
	toJsonBody := ""
	goNames := make(map[string]struct{}, len(sch.Props))

	longestNameLen := 1
	longestTypeLen := 1
	for i, prop := range sch.Props {
		if !validGoIdentifier(prop.GoName) {
			panic(
				fmt.Errorf(
					"invalid Go property name %q in %s.props[%d]",
					prop.Name, sch.RPath, i,
				),
			)
		}
		if !validGoIdentifier(prop.Name) {
			panic(
				fmt.Errorf(
					"invalid Go constructor parameter %q in %s.props[%d]",
					prop.Name, sch.RPath, i,
				),
			)
		}
		if _, exists := goNames[prop.GoName]; exists {
			panic(
				fmt.Errorf(
					"duplicate Go property name %q in %s",
					prop.GoName, sch.RPath,
				),
			)
		}
		goNames[prop.GoName] = struct{}{}
		if !validSchemaText(prop.Desc) ||
			(prop.BsonName != nil && (!validSchemaText(*prop.BsonName) ||
				strings.ContainsAny(*prop.BsonName, "\\\""))) {
			panic(
				fmt.Errorf(
					"invalid description or BSON name in %s.props[%d]",
					sch.RPath, i,
				),
			)
		}
		typ := miniTypToGoTyp(prop.Typ, prop.HMapTyp, prop.Ref)
		if len(prop.GoName) > longestNameLen {
			longestNameLen = len(prop.GoName)
		}
		if len(typ) > longestTypeLen {
			longestTypeLen = len(typ)
		}
	}

	for pos, prop := range sch.Props {
		posStr := strconv.Itoa(pos)
		rootTyp := miniTypToGoTyp(prop.Typ, prop.HMapTyp, prop.Ref)
		typ := rootTyp
		if !prop.Required || (prop.Array && prop.Typ == "object") {
			typ = "*" + typ
		}
		if prop.Array && prop.HMapTyp == nil {
			typ = "[]" + typ
		}
		prop.GoTyp = typ
		nameIndent := longestNameLen - len(prop.GoName) + 1
		nameIndentStr := utils.RepeatString(" ", nameIndent)
		typeIndent := longestTypeLen - len(typ) + 1
		typeIndentStr := utils.RepeatString(" ", typeIndent)
		jsonPropName := prop.Name
		if prop.JsonName != nil {
			jsonPropName = *prop.JsonName
		}
		getPropNameAtPosCases.Push(
			string_templates.GoStructGetPropNameAtPosCase(
				posStr,
				prop.Name,
			),
		)
		gettersSetters.Push(
			string_templates.GoStructPropGetSet(
				prop.GoName,
				sch.PascalName,
				typ,
			),
		)
		propDesc := prop.Desc
		if propDesc != "" {
			propDesc = "// " + prop.Desc
		}
		structFields.Push(
			string_templates.GoStructField(
				"`",
				typeIndentStr,
				propDesc,
				jsonPropName,
				nameIndentStr,
				"",
				prop.GoName,
				typ,
			),
		)
		structFieldsInit.Push(
			string_templates.GoStructFieldInit(
				nameIndentStr,
				prop.GoName,
			),
		)
		structFieldsInitParams.Push(
			string_templates.GoStructFieldInitParam(
				propDesc,
				typeIndentStr,
				nameIndentStr,
				prop.GoName,
				typ,
			),
		)

		copyFnFields.PushMany(copyProp(prop, rootTyp, nameIndentStr).Items...)
		packFnBody.PushMany(packProp(sch, prop, posStr).Items...)
		setPropAtPosCases.PushMany(setPropFnWrapper(prop, posStr, rootTyp, typ).Items...)
	}

	structData := string_templates.GoStruct(
		strings.Join(copyFnFields.Items, "\n"),
		strings.Join(getPropNameAtPosCases.Items, "\n"),
		strings.Join(gettersSetters.Items, "\n\n"),
		strings.Join(packFnBody.Items, "\n"),
		strings.Join(setPropAtPosCases.Items, "\n"),
		strings.Join(structFields.Items, "\n"),
		strings.Join(structFieldsInit.Items, "\n"),
		strings.Join(structFieldsInitParams.Items, "\n"),
		structName,
		toJsonBody,
	)
	output := p.ParserOutputItem{}
	output.Path = strings.ReplaceAll(sch.RPath, ".", "_")
	output.Path = strings.ReplaceAll(output.Path, "/", "_")
	output.Path = strings.ReplaceAll(output.Path, "-", "_")
	output.Path = "obj_" + output.Path + ".go"
	output.Content = structData
	return &output
}

func miniTypToGoTyp(miniTyp string, hmapTyp *string, ref *string) string {
	switch miniTyp {
	case "string":
		return "string"
	case "i32":
		return "int32"
	case "i64":
		return "int64"
	case "u32":
		return "uint32"
	case "u64":
		return "uint64"
	case "f32":
		return "float32"
	case "f64":
		return "float64"
	case "bool":
		return "bool"
	case "enum", "object":
		if ref == nil {
			panic(fmt.Errorf("ref is nil for type %s", miniTyp))
		}
		return goRefType(ref, "")
	case "bytes":
		return "[]byte"
	case "hmap":
		if hmapTyp == nil {
			panic(fmt.Errorf("hmapTyp is nil for type %s", miniTyp))
		}
		typ := "map[string]"
		switch *hmapTyp {
		case "string":
			typ += "string"
		case "i32":
			typ += "int32"
		case "i64":
			typ += "int64"
		case "u32":
			typ += "uint32"
		case "u64":
			typ += "uint64"
		case "f32":
			typ += "float32"
		case "f64":
			typ += "float64"
		case "bool":
			typ += "bool"
		case "bytes":
			typ += "[]byte"
		case "enum", "object":
			if ref == nil {
				panic(fmt.Errorf("ref is nil for type %s", miniTyp))
			}
			typ += goRefType(ref, "")
		default:
			panic(fmt.Errorf("invalid hmapTyp: %s", *hmapTyp))
		}
		return typ
	default:
		panic(fmt.Errorf("invalid miniTyp: %s", miniTyp))
	}
}
