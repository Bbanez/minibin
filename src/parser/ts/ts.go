package parser_ts

import (
	"fmt"

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
			// outputItems = append(outputItems, parseEnum(sch, args))
		} else {
			fmt.Println(3)
		}
	}
	return outputItems
}

func parseObject(sch *schema.Schema, args *utils.Args) *p.ParserOutputItem {
	output := p.ParserOutputItem{}
	// classData := fmt.Sprintf(
	// 	"export class %s {\n",
	// 	sch.PascalName,
	// )
	// packFn := ""+
	// "    pack() ArrayBuffer {\n",
	// "         "
	return &output
}
