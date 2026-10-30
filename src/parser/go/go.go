package parser_go

import (
	"fmt"
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
			outputItems = append(outputItems, parseObject(sch, args))
		} else if sch.Enums != nil {
			outputItems = append(outputItems, parseEnum(sch))
		} else {
			fmt.Println(3)
		}
	}
	return outputItems
}

func parseEnum(sch *schema.Schema) *p.ParserOutputItem {
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

func parseObject(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	structName := sch.PascalName
	copyFnFields := utils.Arr[string]()
	getPropNameAtPosCases := utils.Arr[string]()
	gettersSetters := utils.Arr[string]()
	structFields := utils.Arr[string]()
	structFieldsInit := utils.Arr[string]()
	structFieldsInitParams := utils.Arr[string]()
	setPropAtPosCases := utils.Arr[string]()
	packFnBody := utils.Arr[string]()
	toJsonBody := utils.Arr[string]()
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
		setPropAtPosCases.PushMany(setPropFnWrapper(prop, posStr, rootTyp).Items...)
		toJsonBody.Push(toJsonBodyWrapper(prop))
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
		strings.Join(toJsonBody.Items, "\n"),
	)
	output := p.ParserOutputItem{}
	output.Path = strings.ReplaceAll(sch.RPath, ".", "_")
	output.Path = strings.ReplaceAll(output.Path, "/", "_")
	output.Path = strings.ReplaceAll(output.Path, "-", "_")
	output.Path = "obj_" + output.Path + ".go"
	output.Content = structData
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

func setPropFnHMapPrimWrapper(prop *schema.SchemaProp, posStr string, unpackFn string) string {
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

func setPropFnWrapper(prop *schema.SchemaProp, posStr string, rootTyp string) *utils.Array[string] {
	setPropCases := utils.Arr[string]()
	primWrapperFn := func() {
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
		primWrapperFn()
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
			setPropCases.Push(setPropFnHMapPrimWrapper(
				prop,
				posStr,
				"UnpackString",
			))
		case "i32":
			setPropCases.Push(setPropFnHMapPrimWrapper(
				prop,
				posStr,
				"UnpackInt32",
			))
		case "i64":
			setPropCases.Push(setPropFnHMapPrimWrapper(
				prop,
				posStr,
				"UnpackInt64",
			))
		case "u32":
			setPropCases.Push(setPropFnHMapPrimWrapper(
				prop,
				posStr,
				"UnpackUint32",
			))
		case "u64":
			setPropCases.Push(setPropFnHMapPrimWrapper(
				prop,
				posStr,
				"UnpackUint64",
			))
		case "bool":
			setPropCases.Push(setPropFnHMapPrimWrapper(
				prop,
				posStr,
				"UnpackBool",
			))
		case "f32":
			divider := strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseHmapFloatReq(
						divider,
						posStr,
						prop.GoName,
						"UnpackInt32",
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseHmapFloatOpt(
						divider,
						posStr,
						prop.GoName,
						hmapRootTyp,
						"UnpackInt32",
					),
				)
			}
		case "f64":
			divider := strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
			if prop.Required {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseHmapFloatReq(
						divider,
						posStr,
						prop.GoName,
						"UnpackInt64",
					),
				)
			} else {
				setPropCases.Push(
					string_templates.GoStructSetPropAtPosCaseHmapFloatOpt(
						divider,
						posStr,
						prop.GoName,
						hmapRootTyp,
						"UnpackInt64",
					),
				)
			}
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

func toJsonBodyPrimWrapper(jsonKey string, prop *schema.SchemaProp,
	sliceFnReq string, sliceFnOpt string,
	formatValueFnReq string, formatValueFnOpt string,
	methodCall string,
) string {
	if prop.Array {
		if prop.Required {
			return string_templates.GoToJsonPrimitiveArr(
				prop.GoName,
				jsonKey,
				sliceFnReq,
			)
		} else {
			return string_templates.GoToJsonPrimitiveArr(
				prop.GoName,
				jsonKey,
				sliceFnOpt,
			)
		}
	} else {
		if prop.Required {
			return string_templates.GoToJsonPrimitiveReq(
				formatValueFnReq,
				prop.GoName,
				jsonKey,
				methodCall,
			)
		} else {
			return string_templates.GoToJsonPrimitiveOpt(
				formatValueFnOpt,
				prop.GoName,
				jsonKey,
				methodCall,
			)
		}
	}
}

func toJsonBodyWrapper(prop *schema.SchemaProp) string {
	jsonKey := prop.Name
	if prop.JsonName != nil {
		jsonKey = *prop.JsonName
	}
	switch prop.Typ {
	case "string":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"quoteStrings",
			"quoteStringPointers",
			"jsonQuote",
			"jsonQuote",
			"",
		)
	case "i32":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"int32SliceToStringSlice",
			"int32SliceToStringSliceRef",
			"numToStr",
			"numToStr",
			"",
		)
	case "i64":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"int64SliceToStringSlice",
			"int64SliceToStringSliceRef",
			"numToStr",
			"numToStr",
			"",
		)
	case "u32":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"uint32SliceToStringSlice",
			"uint32SliceToStringSliceRef",
			"numToStr",
			"numToStr",
			"",
		)
	case "u64":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"uint64SliceToStringSlice",
			"uint64SliceToStringSliceRef",
			"numToStr",
			"numToStr",
			"",
		)
	case "f32":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"float32SliceToStringSlice",
			"float32SliceToStringSliceRef",
			"numToStr",
			"numToStr",
			"",
		)
	case "f64":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"float64SliceToStringSlice",
			"float64SliceToStringSliceRef",
			"numToStr",
			"numToStr",
			"",
		)
	case "bool":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"boolSliceToStringSlice",
			"boolSliceToStringSliceRef",
			"numToStr",
			"numToStr",
			"",
		)
	case "enum":
		return toJsonBodyPrimWrapper(
			jsonKey,
			prop,
			"enumSliceToStringSlice",
			"enumSliceToStringSliceRef",
			"jsonQuote",
			"jsonQuote",
			".ToStr()",
		)
	case "object":
		if prop.Array {
			if prop.Required {
				return string_templates.GoToJsonPrimitiveArr(
					prop.GoName,
					jsonKey,
					"objSliceToStringSlice",
				)
			} else {
				return string_templates.GoToJsonPrimitiveArr(
					prop.GoName,
					jsonKey,
					"objSliceToStringSlice",
				)
			}
		} else {
			if prop.Required {
				return string_templates.GoToJsonObjReq(
					prop.GoName,
					jsonKey,
				)
			} else {
				return string_templates.GoToJsonObjOpt(
					prop.GoName,
					jsonKey,
				)
			}
		}
	case "bytes":
		if prop.Array {
			if prop.Required {
				return string_templates.GoToJsonPrimitiveArr(
					prop.GoName,
					jsonKey,
					"bytesSliceToStringSlice",
				)
			} else {
				return string_templates.GoToJsonPrimitiveArr(
					prop.GoName,
					jsonKey,
					"bytesSliceToStringSliceRef",
				)
			}
		} else {
			if prop.Required {
				return string_templates.GoToJsonBytesReq(
					prop.GoName,
					jsonKey,
				)
			} else {
				return string_templates.GoToJsonBytesOpt(
					prop.GoName,
					jsonKey,
				)
			}
		}
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
				return string_templates.GoToJsonPrimitiveArr(
					prop.GoName,
					jsonKey,
					"hmapSliceToStringSlice",
				)
			} else {
				return string_templates.GoToJsonPrimitiveArr(
					prop.GoName,
					jsonKey,
					"hmapSliceToStringSliceRef",
				)
			}
		} else {
			if prop.Required {
				return string_templates.GoToJsonHmapReq(
					formatValueFn,
					prop.GoName,
					methodCall,
				)
			} else {
				return string_templates.GoToJsonHmapOpt(
					formatValueFn,
					prop.GoName,
					methodCall,
				)
			}
		}
	}
	return "// Unknown type: " + prop.Typ
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
