use serde::{Deserialize, Serialize};

pub const SCHEMA_PROP_ALLOWED_TYPES: [&str; 12] = [
    "string", "i32", "i64", "u32", "u64", "f32", "f64", "bool", "object", "enum", "bytes", "hmap",
];

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SchemaProp {
    pub desc: String,
    pub name: String,
    pub go_name: String,
    pub typ: String,
    pub decimals: Option<f32>,
    pub ref_: Option<String>,
    pub required: bool,
    pub array: bool,
    pub bson_name: Option<String>,
    pub json_name: Option<String>,
    pub hmap_typ: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SchemaEnum {
    pub name: String,
    pub go_name: String,
    pub value: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Schema {
    pub r_path: String,
    pub name: String,
    pub pascal_name: String,
    pub typ: String,
    pub props: Vec<SchemaProp>,
    pub enums: Vec<SchemaEnum>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SchemaHMapKeyValue<T> {
    pub key: String,
    pub value: T,
}
