package parser_ts

import (
	"fmt"
	"strings"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/string_templates"
	"github.com/bbanez/minibin/src/utils"
)

func parseEnum(sch *schema.Schema) *p.ParserOutputItem {
	name := sch.PascalName
	values := utils.Arr[string]()
	listItems := utils.Arr[string]()
	defaultValue := ""

	for i, enum := range sch.Enums {
		var value string
		if enum.Value != nil {
			value = *enum.Value
		} else {
			value = enum.Name
		}
		if i == 0 {
			defaultValue = value
		}
		values.Push(
			string_templates.TsEnumValue(
				value,
			),
		)
		listItems.Push(
			string_templates.TsEnumValueListItem(
				value,
			),
		)
	}

	enumData := string_templates.TsEnum(
		defaultValue,
		name,
		strings.Join(listItems.Items, "\n"),
		strings.Join(values.Items, "\n"),
	)
	return &p.ParserOutputItem{
		Path:    fmt.Sprintf("enum_%s.ts", sch.PascalName),
		Content: enumData,
	}
}
