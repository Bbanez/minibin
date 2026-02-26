package parser_cpp

import (
	"fmt"
	"strings"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

type hClassFile struct {
	Name        string
	Content     string
	ParentCount int
}

type cppClassFile struct {
	Name        string
	Content     string
	ParentCount int
}

func Parse(schemas []*schema.Schema, args *utils.Args) []*p.ParserOutputItem {
	hFile := CommonFunctionsH + "\n\n"
	cFile := CommonFunctionsCPP + "\n\n"
	hClassFiles := []*hClassFile{}
	cppClassFiles := []*cppClassFile{}
	for _, sch := range schemas {
		if sch.Props != nil {
			h, c := parseObject(sch)
			hClassFiles = append(hClassFiles, &h)
			cppClassFiles = append(cppClassFiles, &c)
		} else if sch.Enums != nil {
			h, c := parseEnim(sch)
			hClassFiles = append(hClassFiles, &h)
			cFile += c + "\n\n"
		}
	}
	for _, hClassFile := range hClassFiles {
		hFile += hClassFile.Content + "\n\n"
	}
	hFile += "#endif // MINIBIN_H"
	for _, cppClassFile := range cppClassFiles {
		cFile += cppClassFile.Content + "\n\n"
	}
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

func parseEnim(sch *schema.Schema) (hClassFile, string) {
	cFile := ""
	aName := sch.PascalName
	aEnumValues := []string{}
	hFile := hClassFile{
		Name:        aName,
		Content:     "",
		ParentCount: 0,
	}
	for _, enum := range sch.Enums {
		var value string
		if enum.Value != nil {
			value = *enum.Value
		} else {
			value = enum.Name
		}
		aEnumValues = append(aEnumValues, "    "+value)
	}
	hFile.Content = strings.ReplaceAll(HEnim, "@name", aName)
	hFile.Content = strings.ReplaceAll(hFile.Content, "@enumValues", strings.Join(aEnumValues, ",\n"))
	return hFile, cFile
}

func parseObject(sch *schema.Schema) (hClassFile, cppClassFile) {
	aName := sch.PascalName
	cFile := cppClassFile{
		Name:        aName,
		Content:     "",
		ParentCount: 0,
	}
	hFile := hClassFile{
		Name:        aName,
		Content:     "",
		ParentCount: 0,
	}
	aConstructorArgs := []string{}
	aProps := []string{}
	aConstructorBody := []string{}
	aPosToPropName := []string{}
	aEmptyConstructorArgs := []string{}
	aPackProps := []string{}
	aUnpackProps := []string{}
	for i, prop := range sch.Props {
		typ := ""
		val := ""
		aPackProp := ""
		aUnpackProp := ""
		switch prop.Typ {
		case "string":
			typ = "std::string"
			val = "\"\""
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packString(this->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, prop.Name, prop.Name,
			)
			ptr := ""
			if !prop.Required {
				ptr = "&"
			}
			assignVal := fmt.Sprintf(
				"result.%s = %sv",
				prop.Name, ptr,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackString(b, atByte, lenD);\n"+
					"		 %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "i32":
			typ = "int32_t"
			val = "0"
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packInt32(this->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, prop.Name, prop.Name,
			)
			ptr := ""
			if !prop.Required {
				ptr = "&"
			}
			assignVal := fmt.Sprintf(
				"result.%s = %sv",
				prop.Name, ptr,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackInt32(b, atByte, lenD);\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "i64":
			typ = "int64_t"
			val = "0"
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packInt64(this->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = v",
				prop.Name,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackInt64(b, atByte, lenD);\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "u32":
			typ = "uint32_t"
			val = "0"
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packUint32(this->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = v",
				prop.Name,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackUint32(b, atByte, lenD);\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "u64":
			typ = "uint64_t"
			val = "0"
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packUint64(this->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = v",
				prop.Name,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackUint64(b, atByte, lenD);\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "f32":
			typ = "float"
			val = "0.0f"
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packFloat32(this->%s, %d, %d.0f);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, int(prop.Decimals), prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = float(v) / %.1f",
				prop.Name, prop.Decimals,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(float(v) / %d.0f)",
					prop.Name, int(prop.Decimals),
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackUint32(b, atByte, lenD);\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "f64":
			typ = "double"
			val = "0.0"
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packFloat64(this->%s, %d, %d.0f);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, int(prop.Decimals), prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = double(v) / %.1f",
				prop.Name, prop.Decimals,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(double(v) / %.1f)",
					prop.Name, prop.Decimals,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackUint64(b, atByte, lenD);\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "bool":
			typ = "bool"
			val = "false"
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packBool(this->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = v",
				prop.Name,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackBool(b, atByte, lenD);\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "enum":
			typ = strings.Split(*prop.Ref, ".")[1]
			hFile.ParentCount++
			val = fmt.Sprintf("(%s)0", typ)
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packEnum(%sToString(this->%s), %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, typ, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = %sFromString(v)",
				prop.Name, typ,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(%sFromString(v))",
					prop.Name, typ,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackEnum(b, atByte, lenD);\n"+
					"		 %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		case "object":
			typ = strings.Split(*prop.Ref, ".")[1]
			hFile.ParentCount++
			val = fmt.Sprintf("%s()", typ)
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packObject(this->%s.pack(), %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = unpack%s(v, &l)",
				prop.Name, typ,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(unpack%s(v, &l))",
					prop.Name, typ,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackObject(b, atByte, lenD);\n"+
					"        std::string l = lvl+\".%s\";\n"+
					"        %s;\n"+
					"        atByte = ab;",
				typ, assignVal,
			)
		case "bytes":
			typ = "std::vector<uint8_t>"
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packBytes(this->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := fmt.Sprintf(
				"result.%s = v",
				prop.Name,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        auto [v, ab] = _unpackBytes(b, atByte, lenD);\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
		}
		if prop.Array {
			typ = "std::vector<" + typ + ">"
			val = typ + "()"
			aPackProp = fmt.Sprintf(
				""+
					"    for(int i = 0; i < this->%s.size(); i++) {\n"+
					"        %s\n"+
					"    }",
				prop.Name,
				strings.ReplaceAll(aPackProp,
					"this->"+prop.Name, "this->"+prop.Name+"[i]"),
			)
			aUnpackProp = strings.ReplaceAll(
				aUnpackProp,
				fmt.Sprintf("result.%s = %sValue",
					prop.Name, prop.Name),
				fmt.Sprintf("result.%s.push_back(%sValue)",
					prop.Name, prop.Name),
			)
		}
		if !prop.Required {
			typ += "*"
			val = "nullptr"
			if prop.Array {
				tmp := strings.ReplaceAll(aPackProp,
					"this->"+prop.Name+"[i]", "this->"+prop.Name+"->at(i)")
				tmp = strings.ReplaceAll(tmp,
					"this->"+prop.Name+".size", "this->"+prop.Name+"->size")
				aPackProp = fmt.Sprintf(
					""+
						"    if (this->%s != nullptr) {\n"+
						"        %s\n"+
						"    }",
					prop.Name,
					tmp,
				)
			} else {
				aPackProp = fmt.Sprintf(
					""+
						"    if (this->%s != nullptr) {\n"+
						"        %s\n"+
						"    }",
					prop.Name,
					strings.ReplaceAll(aPackProp,
						"this->"+prop.Name, "*this->"+prop.Name),
				)
			}
			tmp := strings.ReplaceAll(aUnpackProp, "        ", "             ")
			tmp = strings.ReplaceAll(tmp,
				fmt.Sprintf("result.%s = %sValue", prop.Name, prop.Name),
				fmt.Sprintf("result.%s = &%sValue", prop.Name, prop.Name),
			)
			tmp = strings.ReplaceAll(tmp,
				fmt.Sprintf("result.%s.push", prop.Name),
				fmt.Sprintf("result.%s->push", prop.Name),
			)
			aUnpackProp = fmt.Sprintf(
				""+
					"        if (result.%s != nullptr) {\n"+
					"            result.%s = new %s();\n"+
					"        }\n"+
					"%s",
				prop.Name, prop.Name, strings.ReplaceAll(typ, "*", ""), tmp,
			)
		}
		aUnpackPropIfElse := "if"
		if i > 0 {
			aUnpackPropIfElse = "else if"
		}
		aUnpackProp = fmt.Sprintf(
			""+
				"        %s (pos == %d) {\n"+
				"%s\n"+
				"        }",
			aUnpackPropIfElse, i,
			strings.ReplaceAll(aUnpackProp, "        ", "            "),
		)
		aUnpackProps = append(aUnpackProps, aUnpackProp)
		aPackProps = append(aPackProps, aPackProp)
		val = fmt.Sprintf("    this->%s = %s;", prop.Name, val)
		aEmptyConstructorArgs = append(aEmptyConstructorArgs, val)
		aConstructorArgs = append(aConstructorArgs, typ+" "+prop.Name)
		aProps = append(aProps, "    "+typ+" "+prop.Name+";")
		aConstructorBody = append(aConstructorBody, "    this->"+prop.Name+" = "+prop.Name+";")
		aPosToPropName = append(
			aPosToPropName,
			fmt.Sprintf(
				"    if (pos == %d) return \"%s\";",
				i, prop.Name,
			),
		)
	}
	aUnpackProps = append(
		aUnpackProps,
		""+
			"        else {\n"+
			"            throw std::runtime_error(\"Invalid property position: \" + std::to_string(pos));\n"+
			"        }",
	)
	hFile.Content = strings.ReplaceAll(HClass, "@name", aName)
	hFile.Content = strings.ReplaceAll(hFile.Content, "@constructorArgs", strings.Join(aConstructorArgs, ", "))
	hFile.Content = strings.ReplaceAll(hFile.Content, "@props", strings.Join(aProps, "\n"))

	cFile.Content = strings.ReplaceAll(CClass, "@name", aName)
	cFile.Content = strings.ReplaceAll(cFile.Content, "@constructorArgs", strings.Join(aConstructorArgs, ", "))
	cFile.Content = strings.ReplaceAll(cFile.Content, "@constructorBody", strings.Join(aConstructorBody, "\n"))
	cFile.Content = strings.ReplaceAll(cFile.Content, "@posToPropName", strings.Join(aPosToPropName, "\n"))
	cFile.Content = strings.ReplaceAll(cFile.Content, "@emptyConstructorArgs", strings.Join(aEmptyConstructorArgs, "\n"))
	cFile.Content = strings.ReplaceAll(cFile.Content, "@packProps", strings.Join(aPackProps, "\n"))
	cFile.Content = strings.ReplaceAll(cFile.Content, "@unpackProps", strings.Join(aUnpackProps, "\n"))
	return hFile, cFile
}
