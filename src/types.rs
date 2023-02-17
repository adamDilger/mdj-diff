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

impl Node {
    pub fn get_tag_map(&self) -> HashMap<String, &Tag> {
        let mut out: HashMap<String, &Tag> = HashMap::new();

        let tags = match self {
            Node::ERDDataModel(e) => &e.element.tags,
            Node::ERDDiagram(e) => &e.element.tags,
            Node::ERDEntity(e) => &e.element.tags,
            Node::ERDRelationship(e) => &e.tags,
        };

        for t in tags.iter() {
            out.insert(t.element.name.clone(), t);
        }

        out
    }
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

impl ERDEntity {
    pub fn get_column_map(&self) -> HashMap<String, &ERDColumn> {
        let mut out: HashMap<String, &ERDColumn> = HashMap::new();

        for c in self.columns.iter() {
            out.insert(c.element._id.clone(), c);
        }

        out
    }

    pub fn get_relationship_map(&self) -> HashMap<String, &ERDRelationship> {
        let mut out: HashMap<String, &ERDRelationship> = HashMap::new();

        if let Some(oe) = &self.element.owned_elements {
            for c in oe.iter() {
                if let Node::ERDRelationship(c) = c {
                    out.insert(c._id.clone(), c);
                }
            }
        }

        out
    }
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

#[derive(Deserialize, Debug, PartialEq)]
#[serde(untagged)]
pub enum ColumnLength {
    Str(String),
    Num(u32),
}

impl ColumnLength {
    pub fn to_string(&self) -> String {
        match self {
            ColumnLength::Str(s) => s.clone(),
            ColumnLength::Num(n) => n.to_string(),
        }
    }
}

#[derive(Deserialize, Debug)]
pub struct Element {
    pub _id: String,
    pub _parent: Ref,
    pub name: String,
    pub documentation: Option<String>,

    #[serde(default = "default_tags")]
    pub tags: Vec<Tag>,

    #[serde(rename = "ownedElements")]
    pub owned_elements: Option<Vec<Node>>,
}

fn default_tags() -> Vec<Tag> {
    Vec::new()
}

#[derive(Deserialize, Debug)]
pub struct ERDRelationship {
    pub _id: String,
    pub _parent: Ref,
    pub name: Option<String>,

    #[serde(default = "default_tags")]
    pub tags: Vec<Tag>,

    pub documentation: Option<String>,

    pub end1: ERDRelationshipEnd,
    pub end2: ERDRelationshipEnd,
}

#[derive(Deserialize, Debug)]
pub struct ERDRelationshipEnd {
    pub _id: String,
    pub _parent: Ref,
    pub reference: Ref,
    pub cardinality: String,
}

impl ERDRelationshipEnd {
    pub fn get_cardinality(&self) -> String {
        if self.cardinality == "" {
            return String::from("1");
        }

        self.cardinality.clone()
    }
}

#[derive(Deserialize, Debug)]
pub struct Ref {
    #[serde(rename = "$ref")]
    pub _ref: String,
}

#[derive(Deserialize, Debug)]
pub struct Tag {
    pub kind: String,
    pub value: Option<String>,

    #[serde(flatten)]
    pub element: Element,
}
