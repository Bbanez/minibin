package tools

import (
	"fmt"
	"strings"

	"github.com/bbanez/minibin/src/utils"
)

func ClearOutput(args *utils.Args) {
	fs := utils.NewFS(nil)
	tree, err := fs.FileTree(
		args.Output,
		[]string{},
	)
	if err != nil {
		panic(err)
	}
	for i := range tree.Items {
		node := tree.Items[i]
		if strings.HasPrefix(node.RelPath, "obj_") ||
			strings.HasPrefix(node.RelPath, "enum_") ||
			strings.HasPrefix(node.RelPath, "minibin__") {
			fmt.Println("Deleting: " + node.RelPath)
			fs.Delete(node.AbsPath).Unwrap()
		}
	}
}
