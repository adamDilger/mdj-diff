#![allow(dead_code)]

use serde::Deserialize;
use std::process::Command;

fn main() {
    let cwd = "/Users/adamdilger/geo/geometry/gep/gep-scripts";
    let commit = "master";
    let path = "gep.mdj";

    let a = format!("{commit}:{path}");

    let output = Command::new("git")
        .args(["-C", cwd, "show", &a])
        .output()
        .expect("failed to execute process");

    let out = String::from_utf8_lossy(&output.stdout);
    let project: Project = serde_json::from_str(&out).unwrap();

    println!("{:#?}", project);
}

#[derive(Deserialize, Debug)]
struct Project {
    _type: String,
    _id: String,
    name: String,

    #[serde(rename = "ownedElements")]
    owned_elements: Vec<Node>,
}

#[derive(Deserialize, Debug)]
#[serde(tag = "_type")]
enum Node {
    ERDDataModel(ERDDataModel),
    ERDDiagram(ERDDiagram),
    ERDEntity(ERDEntity),
    ERDRelationship(ERDRelationship),
}

#[derive(Deserialize, Debug)]
struct ERDDataModel {
    #[serde(flatten)]
    element: Element,
}

#[derive(Deserialize, Debug)]
struct ERDDiagram {
    #[serde(rename = "defaultDiagram", default)]
    default_diagram: bool,

    #[serde(flatten)]
    element: Element,
}

#[derive(Deserialize, Debug)]
struct ERDEntity {
    #[serde(flatten)]
    element: Element,

    columns: Vec<ERDColumn>,
}

#[derive(Deserialize, Debug)]
struct ERDColumn {
    #[serde(flatten)]
    element: Element,

    #[serde(rename = "type")]
    column_type: String,

    #[serde(rename = "referenceTo")]
    reference_to: Option<Ref>,
    #[serde(rename = "primaryKey")]
    primary_key: Option<bool>,
    #[serde(rename = "foreignKey")]
    foreign_key: Option<bool>,

    nullable: Option<bool>,
    unique: Option<bool>,

    length: Option<ColumnLength>,
}

#[derive(Deserialize, Debug)]
#[serde(untagged)]
enum ColumnLength {
    Str(String),
    Num(u32),
}

#[derive(Deserialize, Debug)]
struct Element {
    _id: String,
    _parent: Ref,
    name: String,
    documentation: Option<String>,

    #[serde(rename = "ownedElements")]
    owned_elements: Option<Vec<Node>>,
}

#[derive(Deserialize, Debug)]
struct ERDRelationship {
    _id: String,
    _parent: Ref,
    name: Option<String>,

    documentation: Option<String>,

    end1: ERDRelationshipEnd,
    end2: ERDRelationshipEnd,
}

#[derive(Deserialize, Debug)]
struct ERDRelationshipEnd {
    _id: String,
    _parent: Ref,
    reference: Ref,
    cardinality: String,
}

#[derive(Deserialize, Debug)]
struct Ref {
    #[serde(rename = "$ref")]
    _ref: String,
}

#[derive(Deserialize, Debug)]
struct Tag {
    kind: String,
    value: String,

    #[serde(flatten)]
    element: Element,
}
