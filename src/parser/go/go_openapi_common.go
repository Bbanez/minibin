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
