package parser_ts

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
			Path:    "minibin__common.ts",
			Content: Common,
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
	output := p.ParserOutputItem{
		Path:    fmt.Sprintf("enum_%s.ts", sch.PascalName),
		Content: "",
	}
	enumData := fmt.Sprintf(""+
		"export type %s =\n",
		sch.PascalName,
	)
	defaultValue := fmt.Sprintf(
		"export function %sDefault(): %s {\n",
		sch.PascalName, sch.PascalName,
	)
	for i := range sch.Enums {
		enum := sch.Enums[i]
		val := ""
		if enum.Value != nil {
			val = *enum.Value
		} else {
			val = enum.Name
		}
		if i == 0 {
			defaultValue += fmt.Sprintf(
				"    return '%s';\n",
				val,
			)
		}
		enumData += fmt.Sprintf(
			"    | '%s'\n",
			val,
		)
	}
	defaultValue += "}\n"
	enumData += "\n\n"
	enumData += defaultValue
	output.Content = enumData
	return &output
}

func parseObject(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	output := p.ParserOutputItem{
		Path:    fmt.Sprintf("obj_%s.ts", sch.PascalName),
		Content: "",
	}
	imports := map[string]bool{}
	typData := fmt.Sprintf(
		"export interface %sT {\n",
		sch.PascalName,
	)
	classData := fmt.Sprintf(
		"export class %s {\n",
		sch.PascalName,
	)
	pmData := fmt.Sprintf(
		"const pm: Array<(o: %s, v: unknown) => void> = [\n",
		sch.PascalName,
	)
	constructorData := fmt.Sprintf(
		"    constructor(data: %sT) {\n",
		sch.PascalName,
	)
	copyData := fmt.Sprintf(
		""+
			"    copy(): %s {\n"+
			"        return new %s({\n",
		sch.PascalName, sch.PascalName,
	)
	emptyData := fmt.Sprintf(
		"    static newEmpty(): %s {\n"+
			"        return new %s({\n",
		sch.PascalName, sch.PascalName,
	)
	getSetData := ""
	packData := "" +
		"    pack(): number[] {\n" +
		"        const buf: number[] = [];\n"
	for i := range sch.Props {
		prop := sch.Props[i]
		tsTyp := ""
		packStr := ""
		emptryValue := ""
		switch prop.Typ {
		case "string":
			tsTyp = "string"
			emptryValue = "''"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packString(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "i32":
			tsTyp = "number"
			emptryValue = "0"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packInt32(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "i64":
			tsTyp = "bigint"
			emptryValue = "BigInt(0)"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packInt64(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "u32":
			tsTyp = "number"
			emptryValue = "0"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packUint32(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "u64":
			tsTyp = "bigint"
			emptryValue = "BigInt(0)"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packUint64(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "f32":
			tsTyp = "number"
			emptryValue = "0"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packFloat32(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "f64":
			tsTyp = "number"
			emptryValue = "0"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packFloat64(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "bool":
			tsTyp = "boolean"
			emptryValue = "false"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packBool(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "object":
			if prop.Array {
				copyData += fmt.Sprintf(
					"            %s: this.%s.map(e => e.copy()),\n",
					prop.Name, prop.Name,
				)
			} else {
				if prop.Required {
					copyData += fmt.Sprintf(
						"            %s: this.%s.copy(),\n",
						prop.Name, prop.Name,
					)
				} else {
					copyData += fmt.Sprintf(
						"            %s: this.%s ? this.%s.copy() : undefined,\n",
						prop.Name, prop.Name, prop.Name,
					)
				}
			}
			tsTyp = strings.Split(*prop.Ref, ".")[1]
			if tsTyp != sch.PascalName {
				imp := fmt.Sprintf(
					"import { %s } from './obj_%s'\n",
					tsTyp, tsTyp,
				)
				imports[imp] = true
			}
			emptryValue = fmt.Sprintf(
				"%s.newEmpty()",
				tsTyp,
			)
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packObject(buf, this.%s.pack(), %d);\n",
				prop.Name, i,
			)
		case "enum":
			tsTyp = strings.Split(*prop.Ref, ".")[1]
			if prop.Required {
				imp := fmt.Sprintf(
					"import { type %s, %sDefault } from './enum_%s'\n",
					tsTyp, tsTyp, tsTyp,
				)
				imports[imp] = true
			} else {
				imp := fmt.Sprintf(
					"import type { %s } from './enum_%s'\n",
					tsTyp, tsTyp,
				)
				imports[imp] = true
			}
			emptryValue = fmt.Sprintf(
				"%sDefault()",
				tsTyp,
			)
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packEnum(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		case "bytes":
			tsTyp = "number[]"
			emptryValue = "[]"
			packStr = fmt.Sprintf(
				""+
					"    Minibin.packBytes(buf, this.%s, %d);\n",
				prop.Name, i,
			)
		}
		if prop.Typ != "object" {
			if prop.Array {
				copyData += fmt.Sprintf(
					"            %s: this.%s.map(e => e),\n",
					prop.Name, prop.Name,
				)
			} else {
				copyData += fmt.Sprintf(
					"            %s: this.%s,\n",
					prop.Name, prop.Name,
				)
			}
		}
		if !prop.Required && !prop.Array {
			tsTyp += " | undefined"
			if !prop.Array {
				emptryValue = "undefined"
			}
		}
		if prop.Array {
			tsTyp = fmt.Sprintf("Array<%s>", tsTyp)
			emptryValue = "[]"
			packData += fmt.Sprintf(
				""+
					"        for(let i = 0; i < this.%s.length; i++) {\n"+
					"        %s"+
					"        }\n",
				prop.Name, strings.Replace(
					strings.Replace(
						packStr,
						fmt.Sprintf(", %d)", i),
						fmt.Sprintf("[i], %d)", i),
						1,
					),
					".pack()[i]",
					"[i].pack()",
					1,
				),
			)
		} else {
			packData += fmt.Sprintf(
				""+
					"        if(this.%s !== undefined) {\n"+
					"        %s"+
					"        }\n",
				prop.Name, packStr,
			)
		}
		typData += fmt.Sprintf(
			"    %s: %s;\n",
			prop.Name, tsTyp,
		)
		classData += fmt.Sprintf(
			"    %s: %s;\n",
			prop.Name, tsTyp,
		)
		if prop.Typ == "object" {
			objTyp := strings.Split(*prop.Ref, ".")[1]
			op := ""
			if prop.Array {
				op = fmt.Sprintf("o.%s.push(res)", prop.Name)
			} else {
				op = fmt.Sprintf("o.%s = res", prop.Name)
			}
			pmData += fmt.Sprintf(
				""+
					"    (o, v) => {\n"+
					"        const [res, err] = %s.unpack(v as number[])\n"+
					"        if (err == null && res != null) {\n"+
					"            %s;\n"+
					"        }\n"+
					"    },\n",
				objTyp, op,
			)
		} else {
			op := ""
			if prop.Array {
				op = fmt.Sprintf("o.%s.push(v as %s)", prop.Name, tsTyp)
			} else {
				op = fmt.Sprintf("o.%s = v as %s", prop.Name, tsTyp)
			}
			pmData += fmt.Sprintf(
				""+
					"    (o, v) => {\n"+
					"        %s;\n"+
					"    },\n",
				op,
			)
		}
		constructorData += fmt.Sprintf(
			"        this.%s = data.%s;\n",
			prop.Name, prop.Name,
		)
		emptyData += "            " + prop.Name + ": " + emptryValue + ",\n"
		getSetData += fmt.Sprintf(""+
			"    get%s(): %s {\n"+
			"        return this.%s;\n"+
			"    }\n"+
			"    set%s(v: %s): void {\n"+
			"        this.%s = v;\n"+
			"    }\n\n",
			prop.GoName, tsTyp, prop.Name,
			prop.GoName, tsTyp, prop.Name,
		)
	}
	copyData += "" +
		"        });\n" +
		"    }\n\n"
	typData += "}\n"
	pmData += "];\n\n"
	constructorData += "    }\n"
	emptyData += "" +
		"        });\n" +
		"    }\n\n"
	packData += "" +
		"        return buf;\n" +
		"    }\n"
	classData +=
		constructorData + "\n" +
			copyData + "\n" +
			emptyData + "\n" +
			getSetData +
			fmt.Sprintf("\n"+
				"    setPropAtPos(pos: number, v: unknown): void {\n"+
				"        if(!pm[pos]) {\n"+
				"            return;\n"+
				"        }\n"+
				"        pm[pos](this, v);\n"+
				"    }\n\n") +
			packData + "\n" +
			fmt.Sprintf(""+
				"    static unpack(buffer: number[]): [%s | null, Error | null] {\n"+
				"        const result = %s.newEmpty();\n"+
				"        const err = Minibin.unpack(result, buffer);\n"+
				"        if (err) return [null, err]\n"+
				"        return [result, null];\n"+
				"    }\n",
				sch.PascalName, sch.PascalName,
			) +
			"}\n"
	for imp := range imports {
		output.Content += imp
	}
	output.Content += "import { Minibin } from './minibin__common';\n\n" +
		typData + "\n" + pmData + "\n" + classData
	return &output
}
