package tools

import (
	"fmt"

	"github.com/bbanez/minibin/src/parser"
	parser_cpp "github.com/bbanez/minibin/src/parser/cpp"
	parser_go "github.com/bbanez/minibin/src/parser/go"
	parser_ts "github.com/bbanez/minibin/src/parser/ts"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

func Generate(args *utils.Args) {
	schemas := schema.Read(args.Input)
	var output []*parser.ParserOutputItem
	switch args.Lang {
	case "go":
		output = parser_go.Parse(schemas, args)
	case "ts":
		output = parser_ts.Parse(schemas, args)
	case "cpp":
		output = parser_cpp.Parse(schemas, args)
	default:
		args.PrintHelp()
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
