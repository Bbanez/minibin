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
			h, c := parseEnum(sch)
			hClassFiles = append(hClassFiles, &h)
			cppClassFiles = append(cppClassFiles, &c)
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

func parseEnum(sch *schema.Schema) (hClassFile, cppClassFile) {
	aName := sch.PascalName
	aEnumValues := []string{}
	hFile := hClassFile{
		Name:        aName,
		Content:     "",
		ParentCount: 0,
	}
	cFile := cppClassFile{
		Name:        aName,
		Content:     "",
		ParentCount: 0,
	}
	aEnumToStringCases := []string{}
	aStringToEnumCases := []string{}
	for _, enum := range sch.Enums {
		var value string
		if enum.Value != nil {
			value = *enum.Value
		} else {
			value = enum.Name
		}
		aEnumValues = append(aEnumValues, "    "+value)
		aEnumToStringCases = append(aEnumToStringCases, fmt.Sprintf(
			""+
				"        case %s::%s:\n"+
				"            return \"%s\";",
			aName, value, value,
		))
		aStringToEnumCases = append(aStringToEnumCases, fmt.Sprintf(
			""+
				"        if (s == \"%s\") {\n"+
				"            return %s::%s;\n"+
				"        }",
			value, aName, value,
		))
	}
	hFile.Content = strings.ReplaceAll(HEnum, "@name", aName)
	hFile.Content = strings.ReplaceAll(hFile.Content, "@enumValues", strings.Join(aEnumValues, ",\n"))
	cFile.Content = strings.ReplaceAll(CEnum, "@name", aName)
	cFile.Content = strings.ReplaceAll(cFile.Content, "@enumToStringCases", strings.Join(aEnumToStringCases, "\n"))
	cFile.Content = strings.ReplaceAll(cFile.Content, "@stringToEnumCases", strings.Join(aStringToEnumCases, "\n"))
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
	aPrintStr := []string{}
	aPrintPrep := []string{}
	for i, prop := range sch.Props {
		typ := ""
		val := ""
		aPackProp := ""
		aUnpackProp := ""
		switch prop.Typ {
		case "string":
			typ = "std::string"
			val = "\"\""
			ptrDeref := ""
			if !prop.Required {
				ptrDeref = "*"
			}
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packString(%sthis->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, ptrDeref, prop.Name, i, prop.Name, prop.Name,
			)
			ptr := ""
			if !prop.Required {
				ptr = "&"
			}
			assignVal := fmt.Sprintf(
				"result.%s = %sval",
				prop.Name, ptr,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(val)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<std::string, uint32_t> v = _unpackString(b, atByte, lenD);\n"+
					"		 std::string val = v.a;\n"+
					"		 uint32_t ab = v.b;\n"+
					"		 %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += this->%s[i];\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s;",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? *this->%s : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
		case "i32":
			typ = "int32_t"
			val = "0"
			ptrDeref := ""
			if !prop.Required {
				ptrDeref = "*"
			}
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packInt32(%sthis->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, ptrDeref, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := ""
			if prop.Required {
				assignVal = fmt.Sprintf(
					"result.%s = v",
					prop.Name,
				)
			} else {
				assignVal = fmt.Sprintf(
					"result.%s = &v",
					prop.Name,
				)
			}
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<int32_t, uint32_t> r = _unpackInt32(b, atByte, lenD);\n"+
					"		 int32_t v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += std::to_string(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = std::to_string(this->%s);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? std::to_string(*this->%s) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
		case "i64":
			typ = "int64_t"
			val = "0"
			ptrDeref := ""
			if !prop.Required {
				ptrDeref = "*"
			}
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packInt64(%sthis->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, ptrDeref, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := ""
			if prop.Required {
				assignVal = fmt.Sprintf(
					"result.%s = v",
					prop.Name,
				)
			} else {
				assignVal = fmt.Sprintf(
					"result.%s = &v",
					prop.Name,
				)
			}
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<int64_t, uint32_t> r = _unpackInt64(b, atByte, lenD);\n"+
					"		 int64_t v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += std::to_string(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = std::to_string(this->%s);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? std::to_string(*this->%s) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
		case "u32":
			typ = "uint32_t"
			val = "0"
			ptrDeref := ""
			if !prop.Required {
				ptrDeref = "*"
			}
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packUint32(%sthis->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, ptrDeref, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := ""
			if prop.Required {
				assignVal = fmt.Sprintf(
					"result.%s = v",
					prop.Name,
				)
			} else {
				assignVal = fmt.Sprintf(
					"result.%s = &v",
					prop.Name,
				)
			}
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<uint32_t, uint32_t> r = _unpackUint32(b, atByte, lenD);\n"+
					"		 uint32_t v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += std::to_string(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = std::to_string(this->%s);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? std::to_string(*this->%s) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
		case "u64":
			typ = "uint64_t"
			val = "0"
			ptrDeref := ""
			if !prop.Required {
				ptrDeref = "*"
			}
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packUint64(%sthis->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, ptrDeref, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := ""
			if prop.Required {
				assignVal = fmt.Sprintf(
					"result.%s = v",
					prop.Name,
				)
			} else {
				assignVal = fmt.Sprintf(
					"result.%s = &v",
					prop.Name,
				)
			}
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<uint64_t, uint32_t> r = _unpackUint64(b, atByte, lenD);\n"+
					"		 uint64_t v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += std::to_string(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = std::to_string(this->%s);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? std::to_string(*this->%s) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
		case "f32":
			typ = "float"
			val = "0.0f"
			ptr := ""
			if !prop.Required {
				ptr = "*"
			}
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packFloat32(%sthis->%s, %d, %d.0f);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, ptr, prop.Name, i, int(prop.Decimals), prop.Name, prop.Name,
			)
			assignVal := ""
			if prop.Required {
				assignVal = fmt.Sprintf(
					"result.%s = float(v) / %.1f",
					prop.Name, prop.Decimals,
				)
			} else {
				assignVal = fmt.Sprintf(
					""+
						"%s tmpValue = float(v) / %.1f;\n"+
						"     result.%s = &tmpValue",
					typ, prop.Decimals, prop.Name,
				)
			}
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(float(v) / %d.0f)",
					prop.Name, int(prop.Decimals),
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<int32_t, uint32_t> r = _unpackInt32(b, atByte, lenD);\n"+
					"		 int32_t v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += std::to_string(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = std::to_string(this->%s);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? std::to_string(*this->%s) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
		case "f64":
			typ = "double"
			val = "0.0"
			ptrDeref := ""
			if !prop.Required {
				ptrDeref = "*"
			}
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packFloat64(%sthis->%s, %d, %d.0f);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, ptrDeref, prop.Name, i, int(prop.Decimals), prop.Name, prop.Name,
			)
			assignVal := ""
			if prop.Required {
				assignVal = fmt.Sprintf(
					"result.%s = double(v) / %.1f",
					prop.Name, prop.Decimals,
				)
			} else {
				assignVal = fmt.Sprintf(
					""+
						"%s tmpValue = double(v) / %.1f;\n"+
						"     result.%s = &tmpValue",
					typ, prop.Decimals, prop.Name,
				)
			}
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(double(v) / %.1f)",
					prop.Name, prop.Decimals,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<int64_t, uint32_t> r = _unpackInt64(b, atByte, lenD);\n"+
					"		 int64_t v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += std::to_string(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = std::to_string(this->%s);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? std::to_string(*this->%s) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
		case "bool":
			typ = "bool"
			val = "false"
			ptrDeref := ""
			if !prop.Required {
				ptrDeref = "*"
			}
			aPackProp = fmt.Sprintf(
				""+
					"    std::vector<uint8_t> %sBytes = _packBool(%sthis->%s, %d);\n"+
					"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
				prop.Name, ptrDeref, prop.Name, i, prop.Name, prop.Name,
			)
			assignVal := ""
			if prop.Required {
				assignVal = fmt.Sprintf(
					"result.%s = v",
					prop.Name,
				)
			} else {
				assignVal = fmt.Sprintf(
					"result.%s = &v",
					prop.Name,
				)
			}
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(v)",
					prop.Name,
				)
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<bool, uint32_t> r = _unpackBool(b, atByte, lenD);\n"+
					"		 bool v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += std::to_string(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = std::to_string(this->%s);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? std::to_string(*this->%s) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
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
					"        Tuple<std::string, uint32_t> r = _unpackEnum(b, atByte, lenD);\n"+
					"		 std::string v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"		 %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += %sToString(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, typ, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = %sToString(this->%s);",
						prop.Name, typ, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? %sToString(*this->%s) : \"null\";",
						prop.Name, prop.Name, typ, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
		case "object":
			typ = strings.Split(*prop.Ref, ".")[1]
			hFile.ParentCount++
			val = fmt.Sprintf("%s()", typ)
			if prop.Required || prop.Array {
				aPackProp = fmt.Sprintf(
					""+
						"    std::vector<uint8_t> %sBytes = _packObject(this->%s.pack(), %d);\n"+
						"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
					prop.Name, prop.Name, i, prop.Name, prop.Name,
				)
			} else {
				aPackProp = fmt.Sprintf(
					""+
						"    std::vector<uint8_t> %sBytes = _packObject(this->%s->pack(), %d);\n"+
						"    result.insert(result.end(), %sBytes.begin(), %sBytes.end());",
					prop.Name, prop.Name, i, prop.Name, prop.Name,
				)
			}
			assignVal := fmt.Sprintf(
				"result.%s = unpack%s(v, l)",
				prop.Name, typ,
			)
			if prop.Array {
				assignVal = fmt.Sprintf(
					"result.%s.push_back(unpack%s(v, l))",
					prop.Name, typ,
				)
			} else {
				if !prop.Required {
					assignVal = fmt.Sprintf(
						""+
							"%s tmpValue = unpack%s(v, l);\n"+
							"    result.%s = &tmpValue",
						typ, typ, prop.Name,
					)
				}
			}
			aUnpackProp = fmt.Sprintf(
				""+
					"        Tuple<std::vector<uint8_t>, uint32_t> r = _unpackObject(b, atByte, lenD);\n"+
					"		 std::vector<uint8_t> v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        std::string l = lvl+\".%s\";\n"+
					"        %s;\n"+
					"        atByte = ab;",
				typ, assignVal,
			)
			if prop.Array {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = \"\";\n"+
							"    for(int i = 0; i < this->%s.size(); i++) {\n"+
							"        %sStr += this->%s[i].print(indent * 2);\n"+
							"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
							"    }",
						prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = \"\";\n"+
							"    if (this->%s != nullptr) {\n"+
							"        for(int i = 0; i < this->%s->size(); i++) {\n"+
							"            %sStr += this->%s->at(i).print(indent * 2);\n"+
							"            if (i < this->%s->size() - 1) %sStr += \", \";\n"+
							"        }"+
							"    } else {"+
							"        %sStr = \"null\";\n"+
							"	 }",
						prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
					))
				}
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s.print(indent * 2);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? this->%s->print(indent * 2) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
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
					"        Tuple<std::vector<uint8_t>, uint32_t> r = _unpackBytes(b, atByte, lenD);\n"+
					"		 std::vector<uint8_t> v = r.a;\n"+
					"		 uint32_t ab = r.b;\n"+
					"        %s;\n"+
					"        atByte = ab;",
				assignVal,
			)
			if prop.Array {
				aPrintPrep = append(aPrintPrep, fmt.Sprintf(
					""+
						"    std::string %sStr = \"\";\n"+
						"    for(int i = 0; i < this->%s.size(); i++) {\n"+
						"        %sStr += std::to_string(this->%s[i]);\n"+
						"        if (i < this->%s.size() - 1) %sStr += \", \";\n"+
						"    }",
					prop.Name, prop.Name, prop.Name, prop.Name, prop.Name, prop.Name,
				))
			} else {
				if prop.Required {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = std::to_string(this->%s);",
						prop.Name, prop.Name,
					))
				} else {
					aPrintPrep = append(aPrintPrep, fmt.Sprintf(
						""+
							"    std::string %sStr = this->%s != nullptr ? std::to_string(*this->%s) : \"null\";",
						prop.Name, prop.Name, prop.Name,
					))
				}
			}
			aPrintStr = append(aPrintStr, fmt.Sprintf("indentStr + \"%s: \" + %sStr", prop.Name, prop.Name))
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
					aPackProp,
					// strings.ReplaceAll(aPackProp,
					// 	"this->"+prop.Name, "*this->"+prop.Name),
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
	cFile.Content = strings.ReplaceAll(cFile.Content, "@printPrep", strings.Join(aPrintPrep, "\n"))
	cFile.Content = strings.ReplaceAll(cFile.Content, "@printStr", strings.Join(aPrintStr, " + \"\\n\" + "))
	return hFile, cFile
}
