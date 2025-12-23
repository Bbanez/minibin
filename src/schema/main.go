package schema

import (
	"encoding/json"
	"strings"

	"github.com/bbanez/minibin/src/utils"
)

func Read(path string) []*Schema {
	pathParts := strings.Split(path, "/")
	fs := utils.NewFS(&pathParts)
	filePathsResult := fs.ListFiles("")
	if filePathsResult.Error != nil {
		panic(filePathsResult.Error)
	}
	schemas := []*Schema{}
	for i := range filePathsResult.Value {
		filePath := filePathsResult.Value[i]
		filePathParts := strings.Split(filePath, "/")
		bytesResult := fs.Read(filePathParts...)
		if bytesResult.Error != nil {
			panic(bytesResult.Error)
		}
		tmpSchemas := []*Schema{}
		err := json.Unmarshal(bytesResult.Value, &tmpSchemas)
		if err != nil {
			panic(err)
		}
		for j := range tmpSchemas {
			sch := tmpSchemas[j]
			sch.RPath = strings.Replace(filePath, ".json", "", 1) + "." + sch.Name
			sch.PascalName = utils.ToPascalCase(sch.Name)
			if sch.Props != nil {
				for k := range sch.Props {
					prop := sch.Props[k]
					prop.GoName = strings.ToUpper(prop.Name[0:1]) + prop.Name[1:]
				}
			}
		}
		schemas = append(schemas, tmpSchemas...)
	}
	return schemas
}
