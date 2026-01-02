package main

import (
	"fmt"
	"strings"

	"github.com/bbanez/minibin/src/parser"
	parser_go "github.com/bbanez/minibin/src/parser/go"
	parser_ts "github.com/bbanez/minibin/src/parser/ts"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

func main() {
	args := utils.GetArgs()
	fmt.Println(args)
	schemas := schema.Read(args.Input)
	fs := utils.NewFS(&args.Output)
	var output []*parser.ParserOutputItem
	if args.Clear {
		files := fs.ListFiles("")
		if files.Error != nil {
			panic(files.Error)
		}
		for i := range files.Value {
			filePath := files.Value[i]
			if strings.HasPrefix(filePath, "obj_") ||
				strings.HasPrefix(filePath, "enum_") ||
				strings.HasPrefix(filePath, "minibin__") {
				fs.Delete(strings.Split(filePath, fs.Slash)...)
			}
		}
		return
	}
	switch args.Lang {
	case "go":
		output = parser_go.Parse(schemas, &args)
	case "ts":
		output = parser_ts.Parse(schemas, &args)
	default:
		panic(
			fmt.Errorf("Invalid language provided: %s", args.Lang),
		)
	}
	for i := range output {
		item := output[i]
		fs.Write([]byte(item.Content), item.Path)
	}
}
