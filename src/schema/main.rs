use crate::{
    math::pow::pow,
    schema::typ::{self, Schema},
    utils::{fs::FS, string::to_pascal_case},
};
use serde_json;

pub fn read(path_parts: Vec<String>) -> Vec<Schema> {
    if path_parts.len() == 0 {
        panic!("Invalid path provided.");
    }
    let fs = FS::new(Some(path_parts));
    let tree = fs.file_tree(&[], &[]);
    match tree {
        Ok(tree) => {
            let mut schemas: Vec<Schema> = vec![];
            for i in 0..tree.len() {
                let file_path = tree[i].rel_path.clone();
                let slash = fs.slash.clone();
                let file_path_parts: Vec<&str> = file_path.split(&slash).collect();
                let bytes_result = fs.read(&file_path_parts);
                match bytes_result {
                    Ok(bytes) => {
                        let schema_str = String::from_utf8(bytes).unwrap();
                        let tmp_schemas: Vec<Schema> = serde_json::from_str(&schema_str).unwrap();
                        for j in 0..tmp_schemas.len() {
                            let mut sch = (tmp_schemas[j]).clone();
                            sch.r_path =
                                file_path.clone().replace(".json", "") + "." + sch.name.as_str();
                            if sch.name == "" {
                                panic!("\nMissing property 'name' in {}\n\n", sch.r_path);
                            }
                            if sch.props.len() == 0 && sch.enums.len() == 0 {
                                panic!(
                                    "\nMissing property \"props\" and \"enums\" in {} but schema must contain one of them.",
                                    sch.r_path,
                                )
                            }
                            sch.pascal_name = to_pascal_case(&sch.name);
                            if sch.props.len() > 0 {
                                for k in 0..sch.props.len() {
                                    let mut prop = (sch.props[k]).clone();
                                    if prop.name == "" {
                                        panic!(
                                            "\nMissing property 'name' in {}.props[{}]\n\n",
                                            sch.r_path, k
                                        );
                                    }
                                    prop.go_name = prop.name[0..1].to_uppercase() + &prop.name[1..];
                                    if prop.typ == "" {
                                        panic!(
                                            "\nMissing property 'typ' in {}.props[{}]\n\n",
                                            sch.r_path, k
                                        );
                                    }
                                    if !typ::SCHEMA_PROP_ALLOWED_TYPES.contains(&prop.typ.as_str())
                                    {
                                        panic!(
                                            "\nInvalid property 'typ' in {}.props[{}].typ: {}. Allowed types are: {:?}\n\n",
                                            sch.r_path,
                                            k,
                                            prop.typ,
                                            typ::SCHEMA_PROP_ALLOWED_TYPES
                                        );
                                    }
                                    if let Some(dec) = prop.decimals {
                                        prop.decimals = Some(pow(10.0, dec));
                                    }
                                    if prop.typ == "enum" {
                                        prop.required = true;
                                    }
                                    if prop.typ == "hmap" {
                                        if prop.hmap_typ == "" {}
                                    }
                                }
                            }
                        }
                    }
                    Err(err) => {
                        panic!("Failed to read schema file {}: {}", file_path, err);
                    }
                }
            }
            schemas
        }
        Err(err) => {
            panic!("Failed to read schema files: {}", err);
        }
    }
}
