package parser_go

const GoOpenApiModel = "package minibin\n\n" +
	"\n" +
	"type OpenApiObjectSchema struct {\n" +
	"	Type       *string                         `json:\"type,omitempty\"`\n" +
	"	Enum       *[]string                       `json:\"enum,omitempty\"`\n" +
	"	Ref        *string                         `json:\"$ref,omitempty\"`\n" +
	"	Schema     *OpenApiObjectSchema            `json:\"schema,omitempty\"`\n" +
	"	Required   *[]string                       `json:\"required,omitempty\"`\n" +
	"	Properties *map[string]OpenApiObjectSchema `json:\"properties,omitempty\"`\n" +
	"	Items      *OpenApiObjectSchema            `json:\"items,omitempty\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiRouteParamType string\n" +
	"\n" +
	"const (\n" +
	"	OpenApiRouteParamTypePath   OpenApiRouteParamType = \"path\"\n" +
	"	OpenApiRouteParamTypeQuery  OpenApiRouteParamType = \"query\"\n" +
	"	OpenApiRouteParamTypeHeader OpenApiRouteParamType = \"header\"\n" +
	"	OpenApiRouteParamTypeCookie OpenApiRouteParamType = \"cookie\"\n" +
	")\n" +
	"\n" +
	"type OpenApiRouteParam struct {\n" +
	"	In          OpenApiRouteParamType `json:\"in\"`\n" +
	"	Name        string                `json:\"name\"`\n" +
	"	Description string                `json:\"description\"`\n" +
	"	Required    bool                  `json:\"required\"`\n" +
	"	Schema      *OpenApiObjectSchema  `json:\"schema,omitempty\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiRouteRequestBody struct {\n" +
	"	Description string                         `json:\"description\"`\n" +
	"	Content     map[string]OpenApiObjectSchema `json:\"content\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiRouteResponse struct {\n" +
	"	Description string                         `json:\"description\"`\n" +
	"	Content     map[string]OpenApiObjectSchema `json:\"content\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiRoute struct {\n" +
	"	Tags        *[]string                        `json:\"tags,omitempty\"`\n" +
	"	Summary     string                           `json:\"summary,omitempty\"`\n" +
	"	Description string                           `json:\"description,omitempty\"`\n" +
	"	Parameters  *[]OpenApiRouteParam             `json:\"parameters,omitempty\"`\n" +
	"	RequestBody *OpenApiRouteRequestBody         `json:\"requestBody,omitempty\"`\n" +
	"	Responses   *map[uint32]OpenApiRouteResponse `json:\"responses,omitempty\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiPath struct {\n" +
	"	Get    *OpenApiRoute `json:\"get,omitempty\"`\n" +
	"	Post   *OpenApiRoute `json:\"post,omitempty\"`\n" +
	"	Put    *OpenApiRoute `json:\"put,omitempty\"`\n" +
	"	Delete *OpenApiRoute `json:\"delete,omitempty\"`\n" +
	"	Head   *OpenApiRoute `json:\"head,omitempty\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiInfoContent struct {\n" +
	"	Name  string `json:\"name\"`\n" +
	"	URL   string `json:\"url\"`\n" +
	"	Email string `json:\"email\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiInfo struct {\n" +
	"	Title       string             `json:\"title\"`\n" +
	"	Description string             `json:\"description\"`\n" +
	"	Version     string             `json:\"version\"`\n" +
	"	Contact     OpenApiInfoContent `json:\"contact\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiServer struct {\n" +
	"	URL         string `json:\"url\"`\n" +
	"	Description string `json:\"description\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiComponents struct {\n" +
	"	Schemas         map[string]OpenApiObjectSchema `json:\"schemas\"`\n" +
	"	SecuritySchemes map[string]OpenApiObjectSchema `json:\"securitySchemes\"`\n" +
	"}\n" +
	"\n" +
	"type OpenApiConfig struct {\n" +
	"	OpenApi    string                 `json:\"openapi\"`\n" +
	"	Info       OpenApiInfo            `json:\"info\"`\n" +
	"	Servers    []OpenApiServer        `json:\"servers\"`\n" +
	"	Paths      map[string]OpenApiPath `json:\"paths\"`\n" +
	"	Components OpenApiComponents      `json:\"components\"`\n" +
	"}\n"

const GoOpenApiSchema = `
package minibin

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func stringRef(s string) *string {
	return &s
}

func toCamelCase(s string) string {
	s = strings.ToLower(s[0:1]) + s[1:]
	var i = 0
	for i < len(s) {
		shouldSplit := false
		if s[i] == '_' || s[i] == ' ' || s[i] == '-' {
			shouldSplit = true
		}
		if shouldSplit {
			s = s[:i] + strings.ToUpper(s[i+1:i+2]) + s[i+2:]
			i += 2
		} else {
			i++
		}
	}
	return s
}

func EnumToOpenApiSchema(name string, enum []string) (string, OpenApiObjectSchema) {
	return name, OpenApiObjectSchema{
		Type: stringRef("string"),
		Enum: &enum,
	}
}

func StructToOpenApiRef(obj any) string {
	return fmt.Sprintf("#/components/schemas/%s", reflect.TypeOf(obj).Name())
}

func StructToOpenApiSchema(obj any) (string, OpenApiObjectSchema) {
	outputType := "object"
	outputProps := make(map[string]OpenApiObjectSchema)
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	requiredPropNames := []string{}
	for i := 0; i < val.NumField(); i++ {
		key := typ.Field(i).Name
		field := val.Field(i)
		value := field.Interface()
		valueType := field.Type()
		valueTypeName := field.Type().Name()
		camelKey := toCamelCase(key)
		primitiveTypes := [...]string{
			"string",
			"uint32",
			"uint64",
			"float32",
			"float64",
			"bool",
		}
		if slices.Contains(primitiveTypes[:], valueTypeName) {
			propType := fmt.Sprintf("%s", valueType)
			isArr := strings.Contains(propType, "[]")
			propType = strings.ReplaceAll(propType, "[]", "")
			propType = strings.ReplaceAll(propType, "*", "")
			if strings.Contains(propType, "int") || strings.Contains(propType, "float") {
				propType = "number"
			} else if propType == "bool" {
				propType = "boolean"
			}
			if isArr {
				arrType := "array"
				outputProps[camelKey] = OpenApiObjectSchema{
					Type: &arrType,
					Items: &OpenApiObjectSchema{
						Type: &propType,
					},
				}
			} else {
				outputProps[camelKey] = OpenApiObjectSchema{
					Type: &propType,
				}
			}
		} else {
			propType := fmt.Sprintf("%s", valueType)
			isArr := strings.Contains(valueTypeName, "[]")
			parts := strings.Split(propType, ".")
			propType = parts[len(parts)-1]
			propType = strings.ReplaceAll(propType, "[]", "")
			propType = strings.ReplaceAll(propType, "*", "")
			typeRef := fmt.Sprintf("#/components/schemas/%s", propType)
			if isArr {
				propType := "array"
				outputProps[camelKey] = OpenApiObjectSchema{
					Type: &propType,
					Items: &OpenApiObjectSchema{
						Ref: &typeRef,
					},
				}
			} else {
				outputProps[camelKey] = OpenApiObjectSchema{
					Ref: &typeRef,
				}
			}
		}
		valueStr := fmt.Sprintf("%v", value)
		if valueStr != "<nil>" {
			requiredPropNames = append(requiredPropNames, camelKey)
		}
	}
	return reflect.TypeOf(obj).Name(), OpenApiObjectSchema{
		Type:       &outputType,
		Required:   &requiredPropNames,
		Properties: &outputProps,
	}
}
`
