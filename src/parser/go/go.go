package parser_go

import (
	"fmt"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/string_templates"
	"github.com/bbanez/minibin/src/utils"
)

func Parse(schemas []*schema.Schema, args *utils.Args) []*p.ParserOutputItem {
	schemas = append(schemas, GoOpenApiSchemas()...)
	outputItems := []*p.ParserOutputItem{
		{
			Path:    "minibin__common.go",
			Content: string_templates.GoCommon(),
		},
		{
			Path:    "minibin__openapi_utils.go",
			Content: string_templates.GoOpenapiCommon(),
		},
	}
	for i := range schemas {
		sch := schemas[i]
		if sch.Props != nil {
			outputItems = append(outputItems, parseObject(sch))
		} else if sch.Enums != nil {
			outputItems = append(outputItems, parseEnum(sch))
		} else {
			fmt.Println(3)
		}
	}
	return outputItems
}
