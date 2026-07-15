package parser_go

import (
	"fmt"
	"strings"
	"unicode"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
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
	return true
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
	outputItems := []*p.ParserOutputItem{
		{
			Path:    "minibin__common.go",
			Content: Common,
		},
		{
			Path:    "minibin__openapi_model.go",
			Content: GoOpenApiModel,
		},
		{
			Path:    "minibin__openapi_schema.go",
			Content: GoOpenApiSchema,
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
			"func %sFromStr(v string) %s {\n"+
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
		fromStrFn += fmt.Sprintf("    case %s:\n        return %s\n", goString(value), enum.GoName)
	}
	toOpenApi += "" +
		"        },\n" +
		"    )\n" +
		"}\n"
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
	output.Content = cont + toStrFn + fromStrFn + toOpenApi
	return &output
}

func parseObject(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	if !validGoIdentifier(sch.PascalName) {
		panic(fmt.Errorf("invalid Go schema name %q in %s", sch.Name, sch.RPath))
	}
	if len(sch.Props) > 256 {
		panic(fmt.Errorf("Go minibin supports at most 256 properties per object: %s has %d", sch.RPath, len(sch.Props)))
	}
	output := p.ParserOutputItem{}
	oStruct := "type " + sch.PascalName + " struct {\n"
	fns := ""
	toJsonFn := fmt.Sprintf("\nfunc (o *%s) ToJson() string {\n    result := []string{}\n", sch.PascalName)
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
				"        if item != nil {\n"+
				"            result = append(result, %s(%sitem%s, %d)...)\n"+
				"        }\n"+
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
		jsonKey := goString("\"" + prop.Name + "\":")
		var reqArrStr string
		var reqStr string
		var optArrStr string
		var optStr string
		switch prop.Typ {
		case "string":
			reqArrStr = fmt.Sprintf("strings.Join(quoteStrings(o.%s), \",\")", prop.GoName)
			reqStr = "strconv.Quote(o." + prop.GoName + ")"
			optArrStr = fmt.Sprintf("strings.Join(quoteStringPointers(o.%s), \",\")", prop.GoName)
			optStr = "strconv.Quote(*o." + prop.GoName + ")"
		case "i32":
			reqArrStr = fmt.Sprintf("\" + strings.Join(int32SliceToStringSlice(o.%s), \",\") + \"", prop.GoName)
			reqStr = fmt.Sprintf("strconv.FormatInt(int64(o.%s), 10)", prop.GoName)
			optArrStr = fmt.Sprintf("\" + strings.Join(int32SliceToStringSliceRef(o.%s), \",\") + \"", prop.GoName)
			optStr = fmt.Sprintf("strconv.FormatInt(int64(*o.%s), 10)", prop.GoName)
			if !prop.Array {
				importStrconv = true
			}
		case "i64":
			reqArrStr = fmt.Sprintf("\" + strings.Join(int64SliceToStringSlice(o.%s), \",\") + \"", prop.GoName)
			reqStr = fmt.Sprintf("strconv.FormatInt(o.%s, 10)", prop.GoName)
			optArrStr = fmt.Sprintf("\" + strings.Join(int64SliceToStringSliceRef(o.%s), \",\") + \"", prop.GoName)
			optStr = fmt.Sprintf("strconv.FormatInt(*o.%s, 10)", prop.GoName)
			if !prop.Array {
				importStrconv = true
			}
		case "u32":
			reqArrStr = fmt.Sprintf("\" + strings.Join(uint32SliceToStringSlice(o.%s), \",\") + \"", prop.GoName)
			reqStr = fmt.Sprintf("strconv.FormatUint(uint64(o.%s), 10)", prop.GoName)
			optArrStr = fmt.Sprintf("\" + strings.Join(uint32SliceToStringSliceRef(o.%s), \",\") + \"", prop.GoName)
			optStr = fmt.Sprintf("strconv.FormatUint(uint64(*o.%s), 10)", prop.GoName)
			if !prop.Array {
				importStrconv = true
			}
		case "u64":
			reqArrStr = fmt.Sprintf("\" + strings.Join(uint64SliceToStringSlice(o.%s), \",\") + \"", prop.GoName)
			reqStr = fmt.Sprintf("strconv.FormatUint(o.%s, 10)", prop.GoName)
			optArrStr = fmt.Sprintf("\" + strings.Join(uint64SliceToStringSliceRef(o.%s), \",\") + \"", prop.GoName)
			optStr = fmt.Sprintf("strconv.FormatUint(*o.%s, 10)", prop.GoName)
			if !prop.Array {
				importStrconv = true
			}
		case "f32":
			reqArrStr = fmt.Sprintf(""+
				"\" + strings.Join(float32SliceToStringSlice(o.%s), \",\") + \"",
				prop.GoName,
			)
			reqStr = fmt.Sprintf("strconv.FormatFloat(float64(o.%s), 'f', -1, 32)", prop.GoName)
			optArrStr = fmt.Sprintf(""+
				"\" + strings.Join(float32SliceToStringSliceRef(o.%s), \",\") + \"",
				prop.GoName,
			)
			optStr = fmt.Sprintf("strconv.FormatFloat(float64(*o.%s), 'f', -1, 32)", prop.GoName)
			if !prop.Array {
				importStrconv = true
			}
		case "f64":
			reqArrStr = fmt.Sprintf(""+
				"\" + strings.Join(float64SliceToStringSlice(o.%s), \",\") + \"",
				prop.GoName,
			)
			reqStr = fmt.Sprintf("strconv.FormatFloat(o.%s, 'f', -1, 64)", prop.GoName)
			optArrStr = fmt.Sprintf(""+
				"\" + strings.Join(float64SliceToStringSliceRef(o.%s), \",\") + \"",
				prop.GoName,
			)
			optStr = fmt.Sprintf("strconv.FormatFloat(*o.%s, 'f', -1, 64)", prop.GoName)
			if !prop.Array {
				importStrconv = true
			}
		case "bool":
			reqArrStr = fmt.Sprintf(""+
				"\" + strings.Join(boolSliceToStringSlice(o.%s), \",\") + \"",
				prop.GoName,
			)
			reqStr = fmt.Sprintf("strconv.FormatBool(o.%s)", prop.GoName)
			optArrStr = fmt.Sprintf(""+
				"\" + strings.Join(boolSliceToStringSliceRef(o.%s), \",\") + \"",
				prop.GoName,
			)
			optStr = fmt.Sprintf("strconv.FormatBool(*o.%s)", prop.GoName)
			if !prop.Array {
				importStrconv = true
			}
		case "enum":
			reqArrStr = fmt.Sprintf(""+
				"strings.Join(enumSliceToJSONStringSlice(o.%s), \",\")",
				prop.GoName,
			)
			reqStr = fmt.Sprintf("strconv.Quote(o.%s.ToStr())", prop.GoName)
			optArrStr = fmt.Sprintf(""+
				"strings.Join(enumSliceToJSONStringSlice(o.%s), \",\")",
				prop.GoName,
			)
			optStr = fmt.Sprintf("strconv.Quote((*o.%s).ToStr())", prop.GoName)
		case "object":
			reqArrStr = fmt.Sprintf(""+
				"\" + strings.Join(objSliceToStringSlice(o.%s), \",\") + \"",
				prop.GoName,
			)
			reqStr = fmt.Sprintf("o.%s.ToJson()", prop.GoName)
			optArrStr = fmt.Sprintf(""+
				"\" + strings.Join(objSliceToStringSlice(o.%s), \",\") + \"",
				prop.GoName,
			)
			optStr = fmt.Sprintf("(*o.%s).ToJson()", prop.GoName)
		case "bytes":
			reqArrStr = fmt.Sprintf(""+
				"\" + strings.Join(bytesSliceToStringSlice(o.%s), \",\") + \"",
				prop.GoName,
			)
			reqStr = fmt.Sprintf("\"[\" + strings.Join(bytesToStringSlice(o.%s), \",\") + \"]\"", prop.GoName)
			optArrStr = fmt.Sprintf(""+
				"\" + strings.Join(bytesSliceToStringSliceRef(o.%s), \",\") + \"",
				prop.GoName,
			)
			optStr = fmt.Sprintf("\"[\" + strings.Join(bytesToStringSlice(*o.%s), \",\") + \"]\"", prop.GoName)
		}
		if prop.Required {
			if prop.Array {
				toJsonFn += fmt.Sprintf(""+
					"    if len(o.%s) > 0 {\n"+
					"        result = append(result, %s + \"[\" + %s + \"]\")\n"+
					"    } else {\n"+
					"        result = append(result, %s + \"[]\")\n"+
					"    }\n"+
					"",
					prop.GoName, jsonKey, reqArrStr, jsonKey,
				)
			} else {
				toJsonFn += fmt.Sprintf(""+
					"    result = append(result, %s + %s)\n",
					jsonKey, reqStr,
				)
			}
		} else {
			toJsonFn += fmt.Sprintf(""+
				"    if o.%s != nil {\n"+
				"    ",
				prop.GoName,
			)
			if prop.Array {
				toJsonFn += fmt.Sprintf(""+
					"        if len(o.%s) > 0 {\n"+
					"            result = append(result, %s + \"[\" + %s + \"]\")\n"+
					"        } else {\n"+
					"            result = append(result, %s + \"[]\")\n"+
					"        }\n"+
					"",
					prop.GoName, jsonKey, optArrStr, jsonKey,
				)
			} else {
				toJsonFn += fmt.Sprintf(""+
					"        result = append(result, %s + %s)\n",
					jsonKey, optStr,
				)
			}
			toJsonFn += "    }\n"
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
						"        o.%s = append(o.%s, %sFromStr(d))\n"+
						"",
					i, prop.GoName, prop.GoName, typ,
				)
			} else {
				setPropFn += fmt.Sprintf(
					""+
						"    case %d:\n"+
						"        d, ok := v.(string)\n"+
						"        if !ok { return fmt.Errorf(\"invalid datatype for %%s\", level) }\n"+
						"        o.%s = %sFromStr(d)\n"+
						"",
					i, prop.GoName, typ,
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
		oStruct += fmt.Sprintf(
			"%s%s$name%s$type`json:\"%s,omitempty\"%s%s`\n",
			desc, prop.GoName, typ, goTagValue(prop.Name), bson, requiredTag,
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
	toJsonFn += "    return \"{\" + strings.Join(result, \",\") + \"}\"\n}\n"
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
