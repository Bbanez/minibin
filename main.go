package main

import (
	"fmt"

	"github.com/bbanez/minibin/src/parser"
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
		output = parser.GoParser(schemas, &args)
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

//
// func TTest() {
// 	args := utils.GetArgs()
// 	u := minibin.User{
// 		Id:        "1",
// 		CreatedAt: 1,
// 		UpdatedAt: 2,
// 		Name:      utils.StringRef("Bane"),
// 		Email:     "test@test.com",
// 		Role:      minibin.USER_ROLE_ADMIN,
// 	}
// 	fs := utils.NewFS(&args.Output)
// 	fs.Write(u.Pack(), "dump.txt")
// 	packed := fs.Read("dump.txt")
// 	if packed.Error != nil {
// 		panic(packed.Error)
// 	}
// }
//
// func Dooo() {
// }
