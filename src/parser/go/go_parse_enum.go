package parser_go

import (
	"strings"

	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/string_templates"
	"github.com/bbanez/minibin/src/utils"
)

func parseEnum(sch *schema.Schema) *p.ParserOutputItem {
	name := sch.PascalName
	upperName := utils.ToUpperSnakeCase(name)
	fromStrCases := utils.Arr[string]()
	toStrCases := utils.Arr[string]()
	valueList := utils.Arr[string]()
	values := utils.Arr[string]()

	longestNameLen := 1
	for _, enum := range sch.Enums {
		enum.GoName = upperName + "_" + utils.ToUpperSnakeCase(enum.Name)
		if len(enum.GoName) > longestNameLen {
			longestNameLen = len(enum.GoName)
		}
	}

	for _, enum := range sch.Enums {
		var value string
		if enum.Value != nil {
			value = *enum.Value
		} else {
			value = enum.Name
		}
		enumItemName := upperName + "_" + utils.ToUpperSnakeCase(enum.Name)
		indent := longestNameLen - len(enumItemName) + 1
		indentStr := utils.RepeatString(" ", indent)
		fromStrCases.Push(
			string_templates.GoEnumFromStrCase(
				enumItemName,
				value,
			),
		)
		toStrCases.Push(
			string_templates.GoEnumToStrCase(
				enumItemName,
				value,
			),
		)
		valueList.Push(
			string_templates.GoEnumValueListItem(
				value,
			),
		)
		values.Push(
			string_templates.GoEnumValue(
				indentStr,
				enumItemName,
				name,
				value,
			),
		)
	}

	enumData := string_templates.GoEnum(
		strings.Join(fromStrCases.Items, "\n"),
		name,
		strings.Join(toStrCases.Items, "\n"),
		strings.Join(valueList.Items, "\n"),
		strings.Join(values.Items, "\n"),
	)
	output := p.ParserOutputItem{}
	output.Path = strings.ReplaceAll(sch.RPath, ".", "_")
	output.Path = strings.ReplaceAll(output.Path, "/", "_")
	output.Path = "enum_" + output.Path + ".go"
	output.Content = enumData
	return &output
}
