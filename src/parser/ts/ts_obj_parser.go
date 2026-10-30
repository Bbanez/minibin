package parser_ts

import (
	"fmt"
	"strconv"
	"strings"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/string_templates"
	"github.com/bbanez/minibin/src/utils"
)

func parseObjectV2(sch *schema.Schema) *p.ParserOutputItem {
	output := p.ParserOutputItem{
		Path:    fmt.Sprintf("obj_%s.ts", sch.PascalName),
		Content: "",
	}

	copyValues := utils.Arr[string]()
	emptyValues := utils.Arr[string]()
	gettersSetters := utils.Arr[string]()
	mapValues := utils.Arr[string]()
	packCases := utils.Arr[string]()
	setPropAtPosCases := utils.Arr[string]()
	values := utils.Arr[string]()

	name := sch.PascalName

	longestNameLen := 1
	longestTypeLen := 1
	for _, prop := range sch.Props {
		typ := miniTypToTsTyp(prop.Typ, prop.HMapTyp, prop.Ref)
		if len(prop.GoName) > longestNameLen {
			longestNameLen = len(prop.GoName)
		}
		if len(typ) > longestTypeLen {
			longestTypeLen = len(typ)
		}
	}

	for pos, prop := range sch.Props {
		posStr := strconv.Itoa(pos)
		rootTyp := miniTypToTsTyp(prop.Typ, prop.HMapTyp, prop.Ref)
		typ := rootTyp
		indentNameStr := utils.RepeatString(" ", longestNameLen-len(prop.Name)+1)
		if prop.Array && prop.Typ != "hmap" {
			typ = "Array<" + rootTyp + ">"
		}
		if !prop.Required {
			typ += " | undefined"
		}
		desc := ""
		if prop.Desc != "" {
			desc += " /** " + prop.Desc + " */"
		}
		values.Push(string_templates.TsStructValue(
			desc,
			indentNameStr,
			prop.Name,
			typ,
		))
		mapValues.Push(string_templates.TsStructMapValue(
			indentNameStr,
			prop.Name,
		))
		gettersSetters.Push(string_templates.TsStructGetSet(
			utils.ToPascalCase(prop.Name),
			prop.Name,
			typ,
		))

		switch prop.Typ {
		case "string":
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"\"\"",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				"",
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				"",
				"packString",
			))
		case "i32":
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"0",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				"",
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				"",
				"packInt32",
			))
		case "i64":
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"BigInt(0)",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				"",
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				"",
				"packInt64",
			))
		case "u32":
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"0",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				"",
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				"",
				"packUint32",
			))
		case "u64":
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"BigInt(0)",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				"",
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				"",
				"packUint64",
			))
		case "f32":
			divider := strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"0",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				" / "+divider,
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				", "+divider,
				"",
				"packFloat32",
			))
		case "f64":
			divider := strconv.FormatFloat(float64(prop.Decimals), 'f', -1, 32)
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"0",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				" / "+divider,
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				", "+divider,
				"",
				"packFloat64",
			))
		case "bool":
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"false",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				"",
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				"",
				"packBool",
			))
		case "bytes":
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				"[]",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				"",
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				"",
				"packBytes",
			))
		case "enum":
			copyValues.Push(copyValue(
				prop,
				"",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				rootTyp+"Default()",
			))
			setPropAtPosCases.Push(setPropAtPosCasePrim(
				prop,
				posStr,
				rootTyp,
				"",
			))
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				"",
				"packEnum",
			))
		case "object":
			copyValues.Push(copyValue(
				prop,
				".copy()",
			))
			emptyValues.Push(emptyValuePrim(
				prop,
				indentNameStr,
				rootTyp+".newEmpty()",
			))
			if prop.Array {
				if prop.Required {
					setPropAtPosCases.Push(string_templates.TsStructSetPropAtPosObjArrReq(
						posStr,
						prop.Name,
						rootTyp,
					))
				} else {
					setPropAtPosCases.Push(string_templates.TsStructSetPropAtPosObjArrOpt(
						posStr,
						prop.Name,
						rootTyp,
					))
				}
			} else {
				setPropAtPosCases.Push(string_templates.TsStructSetPropAtPosObj(
					posStr,
					prop.Name,
					rootTyp,
				))
			}
			packCases.Push(packPropCasePrim(
				prop,
				posStr,
				"",
				".pack()",
				"packObject",
			))
		case "hmap":
			if prop.Required {
				emptyValues.Push(string_templates.TsStructEmptyValue(
					indentNameStr,
					prop.Name,
					"{}",
				))
			} else {
				emptyValues.Push(string_templates.TsStructEmptyValue(
					indentNameStr,
					prop.Name,
					"undefined",
				))
			}
			switch *prop.HMapTyp {
			case "string":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "i32":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "i64":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "u32":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "u64":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "f32":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "f64":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "bool":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "enum":
				copyValues.Push(copyHmapValue(
					prop,
					"",
				))
			case "object":
				copyValues.Push(copyHmapValue(
					prop,
					".copy()",
				))
			}
		}
	}

	output.Content = string_templates.TsStruct(
		strings.Join(copyValues.Items, "\n"),
		strings.Join(emptyValues.Items, "\n"),
		strings.Join(gettersSetters.Items, "\n\n"),
		strings.Join(mapValues.Items, "\n"),
		name,
		strings.Join(packCases.Items, "\n"),
		strings.Join(setPropAtPosCases.Items, "\n"),
		strings.Join(values.Items, "\n"),
	)
	return &output
}

func packPropCasePrim(
	prop *schema.SchemaProp, posStr string,
	decimals string,
	methodCall string, packMethod string,
) string {
	if prop.Array {
		if prop.Required {
			return string_templates.TsStructPackPrimArrReq(
				decimals,
				methodCall,
				packMethod,
				posStr,
				prop.Name,
			)
		} else {
			return string_templates.TsStructPackPrimArrOpt(
				decimals,
				methodCall,
				packMethod,
				posStr,
				prop.Name,
			)
		}
	} else {
		if prop.Required {
			return string_templates.TsStructPackPrimReq(
				decimals,
				methodCall,
				packMethod,
				posStr,
				prop.Name,
			)
		} else {
			return string_templates.TsStructPackPrimOpt(
				decimals,
				methodCall,
				packMethod,
				posStr,
				prop.Name,
			)
		}
	}
}

func setPropAtPosCasePrim(
	prop *schema.SchemaProp, posStr string,
	rootType string, divider string,
) string {
	if prop.Array {
		if prop.Required {
			return string_templates.TsStructSetPropAtPosPrimArrReq(
				divider,
				posStr,
				prop.Name,
				rootType,
			)
		} else {
			return string_templates.TsStructSetPropAtPosPrimArrOpt(
				divider,
				posStr,
				prop.Name,
				rootType,
			)
		}
	} else {
		return string_templates.TsStructSetPropAtPosPrim(
			divider,
			posStr,
			prop.Name,
			rootType,
		)
	}
}

func emptyValuePrim(prop *schema.SchemaProp, indent string, value string) string {
	if prop.Array {
		if prop.Required {
			return string_templates.TsStructEmptyValue(
				indent,
				prop.Name,
				"[]",
			)
		} else {
			return string_templates.TsStructEmptyValue(
				indent,
				prop.Name,
				"undefined",
			)
		}
	} else {
		if prop.Required {
			return string_templates.TsStructEmptyValue(
				indent,
				prop.Name,
				value,
			)
		} else {
			return string_templates.TsStructEmptyValue(
				indent,
				prop.Name,
				"undefined",
			)
		}
	}
}

func copyHmapValue(
	prop *schema.SchemaProp,
	methodCall string,
) string {
	if prop.Required {
		return string_templates.TsStructCopyValueHmapReq(
			methodCall,
			prop.Name,
		)
	} else {
		return string_templates.TsStructCopyValueHmapOpt(
			methodCall,
			prop.Name,
		)
	}
}

func copyValue(
	prop *schema.SchemaProp,
	methodCall string,
) string {
	if prop.Array {
		if prop.Required {
			return string_templates.TsStructCopyValueReqArr(
				methodCall,
				prop.Name,
			)
		} else {
			return string_templates.TsStructCopyValueOptArr(
				methodCall,
				prop.Name,
			)
		}
	} else {
		if prop.Required {
			return string_templates.TsStructCopyValueReq(
				methodCall,
				prop.Name,
			)
		} else {
			return string_templates.TsStructCopyValueOpt(
				methodCall,
				prop.Name,
			)
		}
	}
}

func tsRefType(ref *string, path string) string {
	if ref == nil {
		panic(fmt.Errorf("missing ref in %s", path))
	}
	parts := strings.Split(*ref, ".")
	if len(parts) != 2 {
		panic(fmt.Errorf("invalid ref %q in %s", *ref, path))
	}
	return parts[1]
}

func miniTypToTsTyp(miniTyp string, hmapTyp *string, ref *string) string {
	switch miniTyp {
	case "string":
		return "string"
	case "i32":
		return "number"
	case "i64":
		return "bigint"
	case "u32":
		return "number"
	case "u64":
		return "bigint"
	case "f32":
		return "number"
	case "f64":
		return "number"
	case "bool":
		return "boolean"
	case "enum", "object":
		if ref == nil {
			panic(fmt.Errorf("ref is nil for type %s", miniTyp))
		}
		return tsRefType(ref, "")
	case "bytes":
		return "number[]"
	case "hmap":
		if hmapTyp == nil {
			panic(fmt.Errorf("hmapTyp is nil for type %s", miniTyp))
		}
		typ := "Record<string, "
		switch *hmapTyp {
		case "string":
			typ += "string"
		case "i32":
			typ += "number"
		case "i64":
			typ += "bigint"
		case "u32":
			typ += "number"
		case "u64":
			typ += "bigint"
		case "f32":
			typ += "number"
		case "f64":
			typ += "number"
		case "bool":
			typ += "bool"
		case "bytes":
			typ += "number[]"
		case "enum", "object":
			if ref == nil {
				panic(fmt.Errorf("ref is nil for type %s", miniTyp))
			}
			typ += tsRefType(ref, "")
		default:
			panic(fmt.Errorf("invalid hmapTyp: %s", *hmapTyp))
		}
		typ += ">"
		return typ
	default:
		panic(fmt.Errorf("invalid miniTyp: %s", miniTyp))
	}
}
