package main

import (
	"fmt"

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
	var output []*parser.ParserOutputItem
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
	fs := utils.NewFS(&args.Output)
	for i := range output {
		item := output[i]
		fs.Write([]byte(item.Content), item.Path)
	}
}
