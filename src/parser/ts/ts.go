package parser_ts

import (
	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/string_templates"
	"github.com/bbanez/minibin/src/utils"
)

func Parse(schemas []*schema.Schema, args *utils.Args) []*p.ParserOutputItem {
	mainData := p.ParserOutputItem{
		Path:    "minibin_generated.ts",
		Content: string_templates.TsCommon(),
	}
	outputItems := []*p.ParserOutputItem{&mainData}
	for i := range schemas {
		sch := schemas[i]
		if sch.Props != nil {
			mainData.Content += parseObject(sch).Content + "\n\n"
		} else if sch.Enums != nil {
			mainData.Content += parseEnum(sch).Content + "\n\n"
		}
	}
	return outputItems
}
