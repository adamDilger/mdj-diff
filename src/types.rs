#![allow(dead_code)]

use std::{collections::HashMap, fs, path::Path, process::Command};

use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct Project {
    _type: String,
    _id: String,
    name: String,

    #[serde(rename = "ownedElements")]
    owned_elements: Vec<Node>,
}

impl Project {
    pub fn from_file(path: &str, cwd: &str) -> Project {
        let contents = fs::read_to_string(Path::new(cwd).join(path))
            .expect("Should have been able to read the file");

        let project: Project = serde_json::from_str(&contents).unwrap();

        project
    }

    pub fn from_git(commit: Option<&str>, path: &str, cwd: &str) -> Project {
        if let None = commit {
            return Project::from_file(path, cwd);
        }

        let cmd = match commit {
            Some(c) => (format!("show"), format!("{}:{}", c, path)),
            None => (format!("diff"), format!("{}", path)),
        };

        let output = Command::new("git")
            .args(["-C", cwd, &cmd.0, &cmd.1])
            .output()
            .expect("failed to execute process");

        let out = String::from_utf8_lossy(&output.stdout);
        let project: Project = serde_json::from_str(&out).unwrap();

        project
    }

    pub fn get_entity_map(self) -> HashMap<String, ERDEntity> {
        let mut entities: HashMap<String, ERDEntity> = HashMap::new();

        for ele in self.owned_elements {
            if let Node::ERDDataModel(data_model) = ele {
                for ele in data_model.element.owned_elements.unwrap() {
                    if let Node::ERDEntity(entity) = ele {
                        entities.insert(entity.element._id.clone(), entity);
                    }
                }
            }
        }

        entities
    }
}

#[derive(Deserialize, Debug)]
#[serde(tag = "_type")]
pub enum Node {
    ERDDataModel(ERDDataModel),
    ERDDiagram(ERDDiagram),
    ERDEntity(ERDEntity),
    ERDRelationship(ERDRelationship),
}

#[derive(Deserialize, Debug)]
pub struct ERDDataModel {
    #[serde(flatten)]
    element: Element,
}

#[derive(Deserialize, Debug)]
pub struct ERDDiagram {
    #[serde(rename = "defaultDiagram", default)]
    default_diagram: bool,

    #[serde(flatten)]
    element: Element,
}

#[derive(Deserialize, Debug)]
pub struct ERDEntity {
    #[serde(flatten)]
    pub element: Element,

    pub columns: Vec<ERDColumn>,
}

#[derive(Deserialize, Debug)]
pub struct ERDColumn {
    #[serde(flatten)]
    pub element: Element,

    #[serde(rename = "type")]
    pub column_type: String,

    #[serde(rename = "referenceTo")]
    pub reference_to: Option<Ref>,
    #[serde(rename = "primaryKey")]
    pub primary_key: Option<bool>,
    #[serde(rename = "foreignKey")]
    pub foreign_key: Option<bool>,

    pub nullable: Option<bool>,
    pub unique: Option<bool>,

    pub length: Option<ColumnLength>,
}

#[derive(Deserialize, Debug)]
#[serde(untagged)]
pub enum ColumnLength {
    Str(String),
    Num(u32),
}

#[derive(Deserialize, Debug)]
pub struct Element {
    pub _id: String,
    pub _parent: Ref,
    pub name: String,
    pub documentation: Option<String>,

    #[serde(rename = "ownedElements")]
    pub owned_elements: Option<Vec<Node>>,
}

#[derive(Deserialize, Debug)]
pub struct ERDRelationship {
    _id: String,
    _parent: Ref,
    name: Option<String>,

    documentation: Option<String>,

    end1: ERDRelationshipEnd,
    end2: ERDRelationshipEnd,
}

#[derive(Deserialize, Debug)]
pub struct ERDRelationshipEnd {
    _id: String,
    _parent: Ref,
    reference: Ref,
    cardinality: String,
}

#[derive(Deserialize, Debug)]
pub struct Ref {
    #[serde(rename = "$ref")]
    _ref: String,
}

#[derive(Deserialize, Debug)]
pub struct Tag {
    kind: String,
    value: String,

    #[serde(flatten)]
    element: Element,
}
