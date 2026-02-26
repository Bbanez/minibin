package parser_cpp

import (
	p "github.com/bbanez/minibin/src/parser"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

func Parse(schemas []*schema.Schema, args *utils.Args) []*p.ParserOutputItem {
	hFile := ""
	cFile := ""
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

func parseObject(sch *schema.Schema, args *utils.Args) (string, string) {
	hFile := ""
	cFile := ""
	for i, prop := range sch.Props {

	}
	return hFile, cFile
}
