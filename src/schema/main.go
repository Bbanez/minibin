package schema

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/bbanez/minibin/src/utils"
)

func Read(pathParts []string) []*Schema {
	fs := utils.NewFS(&pathParts)
	tree, err := fs.FileTree([]string{}, []string{})
	if err != nil {
		fmt.Println("Failed to list schema files in input directory:", strings.Join(pathParts, "/"))
		panic(err)
	}
	schemas := []*Schema{}
	for i := range tree.Items {
		filePath := tree.Items[i].RelPath
		filePathParts := strings.Split(filePath, "/")
		bytesResult := fs.Read(filePathParts...)
		if bytesResult.Error != nil {
			panic("Failed to read " + fs.Path(filePathParts...) + ": " + bytesResult.Error.Error())
		}
		tmpSchemas := []*Schema{}
		err := json.Unmarshal(bytesResult.Value, &tmpSchemas)
		if err != nil {
			panic("Failed to deserialize " + fs.Path(filePathParts...) + ": " + err.Error())
		}
		for j := range tmpSchemas {
			sch := tmpSchemas[j]
			sch.RPath = strings.Replace(filePath, ".json", "", 1) + "." + sch.Name
			if sch.Name == "" {
				panic(
					fmt.Errorf(
						"\nMissing property \"name\" in: %s\n\n",
						sch.RPath,
					),
				)
			}
			if sch.Props == nil && sch.Enums == nil {
				panic(
					fmt.Errorf(
						"\nMissing property \"props\" and \"enums\" in %s but schema must contain one of them.",
						sch.RPath,
					),
				)
			}
			sch.PascalName = utils.ToPascalCase(sch.Name)
			if sch.Props != nil {
				for k := range sch.Props {
					prop := sch.Props[k]
					if prop.Name == "" {
						panic(
							fmt.Errorf(
								"\nMissing property \"name\" in %s.props[%d]",
								sch.RPath, k,
							),
						)
					}
					prop.GoName = strings.ToUpper(prop.Name[0:1]) + prop.Name[1:]
					if prop.Typ == "" {
						panic(
							fmt.Errorf(
								"\nMissing property \"typ\" in %s.props[%d]",
								sch.RPath, k,
							),
						)
					}
					if !slices.Contains(SchemaPropAllowedTypes, prop.Typ) {
						panic(
							fmt.Errorf(
								"\nType \"%s\" is not allowed in %s.props[%d].typ -> Allowed: %s",
								prop.Typ, sch.RPath, k, strings.Join(SchemaPropAllowedTypes, ", "),
							),
						)
					}
					prop.Decimals = float32(math.Pow10(int(prop.Decimals)))
					if prop.Typ == "enum" {
						prop.Required = true
					}
					if prop.Typ == "hmap" {
						if prop.HMapTyp == nil {
							panic(
								fmt.Errorf(
									"\nMissing property \"hmap_typ\" in %s.props[%d]",
									sch.RPath, k,
								),
							)
						}
						if !slices.Contains(SchemaPropAllowedTypes, prop.Typ) {
							panic(
								fmt.Errorf(
									"\nType \"%s\" is not allowed in %s.props[%d].hmapTyp -> Allowed: %s",
									prop.Typ, sch.RPath, k, strings.Join(SchemaPropAllowedTypes, ", "),
								),
							)
						}
					}
				}
			}
		}
		schemas = append(schemas, tmpSchemas...)
	}
	return schemas
}
