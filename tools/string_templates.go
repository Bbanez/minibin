package main

import (
	"sort"
	"strings"

	"github.com/bbanez/minibin/src/utils"
)

func GenerateGoStringTemplates() {
	fs := utils.NewFS(nil)
	tree, err := fs.FileTree(
		[]string{"string_templates"},
		[]string{},
	)
	if err != nil {
		panic(err)
	}
	for _, file := range tree.Items {
		data := fs.ReadString("string_templates", file.RelPath)
		if data.Error != nil {
			panic(data.Error)
		}
		vars := utils.KeepUniqueStrings(
			utils.AllTextBetween(data.Value, "{{.", "}}"),
		)
		sort.Strings(vars)
		fileName := strings.ReplaceAll(
			strings.ReplaceAll(file.RelPath, "/", "_"),
			".",
			"_",
		) + ".go"
		fnName := strings.Replace(
			utils.ToPascalCase(strings.Replace(fileName, ".go", "", 1)),
			"StringTemplates",
			"",
			1,
		)
		fn := "package string_templates\n\n" +
			"func " + fnName + "(\n"
		mutatedData := strings.TrimSuffix(data.Value, "\n")
		for _, k := range vars {
			mutatedData = strings.ReplaceAll(mutatedData, "{{."+k+"}}", "` + "+k+" + `")
			fn += "    " + k + " string,\n"
		}
		fn += ") string {\n" +
			"    return `" + mutatedData +
			"`\n}\n"
		fs.Write([]byte(fn), "src", "string_templates", fileName).Unwrap()
	}
}
