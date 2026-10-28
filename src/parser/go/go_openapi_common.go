package parser_go

import (
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

func GoOpenApiSchemas() []*schema.Schema {
	return []*schema.Schema{
		{
			Name:       "OpenApiObjectSchema",
			PascalName: "OpenApiObjectSchema",
			RPath:      "openapi_object_schema",
			Props: []*schema.SchemaProp{
				{
					Name:     "format",
					GoName:   "Format",
					Typ:      "string",
					Required: false,
				},
				{
					Name:     "typ",
					GoName:   "Typ",
					JsonName: utils.ToRef("type"),
					Typ:      "string",
					Required: false,
				},
				{
					Name:     "enum",
					GoName:   "Enum",
					Typ:      "string",
					Required: false,
					Array:    true,
				},
				{
					Name:     "ref",
					GoName:   "Ref",
					Typ:      "string",
					Required: false,
					JsonName: utils.ToRef("$ref"),
				},
				{
					Name:     "schema",
					GoName:   "Schema",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiObjectSchema"),
					Required: false,
				},
				{
					Name:     "required",
					GoName:   "Required",
					Typ:      "string",
					Required: false,
					Array:    true,
				},
				{
					Name:     "properties",
					GoName:   "Properties",
					Typ:      "hmap",
					HMapTyp:  utils.ToRef("object"),
					Ref:      utils.ToRef(".OpenApiObjectSchema"),
					Required: false,
				},
				{
					Name:     "items",
					GoName:   "Items",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiObjectSchema"),
					Required: false,
				},
			},
		},

		{
			Name:       "OpenApiRouteParamType",
			PascalName: "OpenApiRouteParamType",
			RPath:      "openapi_route_param_type",
			Enums: []*schema.SchemaEnum{
				{
					Name:   "PATH",
					Value:  utils.ToRef("path"),
					GoName: "Path",
				},
				{
					Name:   "QUERY",
					Value:  utils.ToRef("query"),
					GoName: "Query",
				},
				{
					Name:   "HEADER",
					Value:  utils.ToRef("header"),
					GoName: "Header",
				},
				{
					Name:   "COOKIE",
					Value:  utils.ToRef("cookie"),
					GoName: "Cookie",
				},
			},
		},

		{
			Name:       "OpenApiRouteParam",
			PascalName: "OpenApiRouteParam",
			RPath:      "openapi_route_param",
			Props: []*schema.SchemaProp{
				{
					Name:     "in",
					GoName:   "In",
					Typ:      "enum",
					Ref:      utils.ToRef(".OpenApiRouteParamType"),
					Required: true,
				},
				{
					Name:     "name",
					GoName:   "Name",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "description",
					GoName:   "Description",
					Typ:      "string",
					Required: false,
				},
				{
					Name:     "required",
					GoName:   "Required",
					Typ:      "bool",
					Required: true,
				},
				{
					Name:     "schema",
					GoName:   "Schema",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiObjectSchema"),
					Required: false,
				},
			},
		},

		{
			Name:       "OpenApiRouteRequestBody",
			PascalName: "OpenApiRouteRequestBody",
			RPath:      "openapi_route_request_body",
			Props: []*schema.SchemaProp{
				{
					Name:     "description",
					GoName:   "Description",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "content",
					GoName:   "Content",
					Typ:      "hmap",
					HMapTyp:  utils.ToRef("object"),
					Ref:      utils.ToRef(".OpenApiObjectSchema"),
					Required: true,
				},
			},
		},

		{
			Name:       "OpenApiRouteResponse",
			PascalName: "OpenApiRouteResponse",
			RPath:      "openapi_route_response",
			Props: []*schema.SchemaProp{
				{
					Name:     "description",
					GoName:   "Description",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "content",
					GoName:   "Content",
					Typ:      "hmap",
					HMapTyp:  utils.ToRef("object"),
					Ref:      utils.ToRef(".OpenApiObjectSchema"),
					Required: true,
				},
			},
		},

		{
			Name:       "OpenApiRoute",
			PascalName: "OpenApiRoute",
			RPath:      "openapi_route",
			Props: []*schema.SchemaProp{
				{
					Name:     "tags",
					GoName:   "Tags",
					Typ:      "string",
					Required: false,
					Array:    true,
				},
				{
					Name:     "summary",
					GoName:   "Summary",
					Typ:      "string",
					Required: false,
				},
				{
					Name:     "description",
					GoName:   "Description",
					Typ:      "string",
					Required: false,
				},
				{
					Name:     "parameters",
					GoName:   "Parameters",
					Typ:      "object",
					Required: false,
					Array:    true,
					Ref:      utils.ToRef(".OpenApiRouteParam"),
				},
				{
					Name:     "requestBody",
					GoName:   "RequestBody",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiRouteRequestBody"),
					Required: false,
				},
				{
					Name:     "responses",
					GoName:   "Responses",
					Typ:      "hmap",
					HMapTyp:  utils.ToRef("object"),
					Ref:      utils.ToRef(".OpenApiRouteResponse"),
					Required: false,
				},
			},
		},

		{
			Name:       "OpenApiPath",
			PascalName: "OpenApiPath",
			RPath:      "openapi_path",
			Props: []*schema.SchemaProp{
				{
					Name:     "get",
					GoName:   "Get",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiRoute"),
					Required: false,
				},
				{
					Name:     "post",
					GoName:   "Post",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiRoute"),
					Required: false,
				},
				{
					Name:     "put",
					GoName:   "Put",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiRoute"),
					Required: false,
				},
				{
					Name:     "delete",
					GoName:   "Delete",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiRoute"),
					Required: false,
				},
				{
					Name:     "head",
					GoName:   "Head",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiRoute"),
					Required: false,
				},
			},
		},

		{
			Name:       "OpenApiInfoContent",
			PascalName: "OpenApiInfoContent",
			RPath:      "openapi_info_content",
			Props: []*schema.SchemaProp{
				{
					Name:     "name",
					GoName:   "Name",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "url",
					GoName:   "URL",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "email",
					GoName:   "Email",
					Typ:      "string",
					Required: true,
				},
			},
		},

		{
			Name:       "OpenApiInfo",
			PascalName: "OpenApiInfo",
			RPath:      "openapi_info",
			Props: []*schema.SchemaProp{
				{
					Name:     "title",
					GoName:   "Title",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "description",
					GoName:   "Description",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "version",
					GoName:   "Version",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "contact",
					GoName:   "Contact",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiInfoContent"),
					Required: true,
				},
			},
		},

		{
			Name:       "OpenApiServer",
			PascalName: "OpenApiServer",
			RPath:      "openapi_server",
			Props: []*schema.SchemaProp{
				{
					Name:     "url",
					GoName:   "URL",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "description",
					GoName:   "Description",
					Typ:      "string",
					Required: true,
				},
			},
		},

		{
			Name:       "OpenApiComponents",
			PascalName: "OpenApiComponents",
			RPath:      "openapi_components",
			Props: []*schema.SchemaProp{
				{
					Name:     "schemas",
					GoName:   "Schemas",
					Typ:      "hmap",
					HMapTyp:  utils.ToRef("object"),
					Ref:      utils.ToRef(".OpenApiObjectSchema"),
					Required: true,
				},
				{
					Name:     "securitySchemes",
					GoName:   "SecuritySchemes",
					Typ:      "hmap",
					HMapTyp:  utils.ToRef("object"),
					Ref:      utils.ToRef(".OpenApiObjectSchema"),
					Required: true,
				},
			},
		},

		{
			Name:       "OpenApiConfig",
			PascalName: "OpenApiConfig",
			RPath:      "openapi_config",
			Props: []*schema.SchemaProp{
				{
					Name:     "openapi",
					GoName:   "OpenApi",
					Typ:      "string",
					Required: true,
				},
				{
					Name:     "info",
					GoName:   "Info",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiInfo"),
					Required: true,
				},
				{
					Name:     "servers",
					GoName:   "Servers",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiServer"),
					Required: true,
					Array:    true,
				},
				{
					Name:     "paths",
					GoName:   "Paths",
					Typ:      "hmap",
					HMapTyp:  utils.ToRef("object"),
					Ref:      utils.ToRef(".OpenApiPath"),
					Required: true,
				},
				{
					Name:     "components",
					GoName:   "Components",
					Typ:      "object",
					Ref:      utils.ToRef(".OpenApiComponents"),
					Required: true,
				},
			},
		},
	}
}

const GoOpenApiUtils = `
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
	if s == "" {
		return s
	}
	s = strings.ToLower(s[0:1]) + s[1:]
	var i = 0
	for i < len(s) {
		shouldSplit := false
		if s[i] == '_' || s[i] == ' ' || s[i] == '-' {
			shouldSplit = true
		}
		if shouldSplit {
			if i+1 >= len(s) {
				return s[:i]
			}
			s = s[:i] + strings.ToUpper(s[i+1:i+2]) + s[i+2:]
			i += 2
		} else {
			i++
		}
	}
	return s
}

func EnumToOpenApiSchema(name string, enum []*string) (string, OpenApiObjectSchema) {
	return name, OpenApiObjectSchema{
		Typ: stringRef("string"),
		Enum: enum,
	}
}

func StructToOpenApiRef(obj any) string {
	typ := reflect.TypeOf(obj)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return fmt.Sprintf("#/components/schemas/%s", typ.Name())
}

func StructToOpenApiSchema(obj any) (string, OpenApiObjectSchema) {
	outputType := "object"
	outputProps := make(map[string]OpenApiObjectSchema)
	typ := reflect.TypeOf(obj)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	requiredPropNames := []*string{}
	for i := 0; i < typ.NumField(); i++ {
		key := typ.Field(i).Name
		fieldType := typ.Field(i).Type
		byteType := fieldType
		for byteType.Kind() == reflect.Pointer || byteType.Kind() == reflect.Slice {
			byteType = byteType.Elem()
		}
		if byteType.Kind() == reflect.Uint8 {
			propType := "string"
			format := "byte"
			isArr := fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() != reflect.Uint8
			if isArr {
				arrType := "array"
				outputProps[toCamelCase(key)] = OpenApiObjectSchema{Typ: &arrType, Items: &OpenApiObjectSchema{Typ: &propType, Format: &format}}
			} else {
				outputProps[toCamelCase(key)] = OpenApiObjectSchema{Typ: &propType, Format: &format}
			}
			if typ.Field(i).Tag.Get("minibin") == "required" {
				requiredPropNames = append(requiredPropNames, stringRef(toCamelCase(key)))
			}
			continue
		}
		valueType := fieldType
		for valueType.Kind() == reflect.Pointer || valueType.Kind() == reflect.Slice {
			valueType = valueType.Elem()
		}
		valueTypeName := valueType.Name()
		camelKey := toCamelCase(key)
		primitiveTypes := [...]string{
			"string",
			"int32",
			"int64",
			"uint32",
			"uint64",
			"float32",
			"float64",
			"bool",
		}
		if slices.Contains(primitiveTypes[:], valueTypeName) {
			propType := fmt.Sprintf("%s", valueType)
			isArr := fieldType.Kind() == reflect.Slice
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
					Typ: &arrType,
					Items: &OpenApiObjectSchema{
						Typ: &propType,
					},
				}
			} else {
				outputProps[camelKey] = OpenApiObjectSchema{
					Typ: &propType,
				}
			}
		} else {
			propType := fmt.Sprintf("%s", valueType)
			isArr := fieldType.Kind() == reflect.Slice
			parts := strings.Split(propType, ".")
			propType = parts[len(parts)-1]
			propType = strings.ReplaceAll(propType, "[]", "")
			propType = strings.ReplaceAll(propType, "*", "")
			typeRef := fmt.Sprintf("#/components/schemas/%s", propType)
			if isArr {
				propType := "array"
				outputProps[camelKey] = OpenApiObjectSchema{
					Typ: &propType,
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
		if typ.Field(i).Tag.Get("minibin") == "required" {
			requiredPropNames = append(requiredPropNames, &camelKey)
		}
	}
	return typ.Name(), OpenApiObjectSchema{
		Typ:        &outputType,
		Required:   requiredPropNames,
		Properties: &outputProps,
	}
}
`
