use serde::Deserialize;
use std::process::Command;

fn main() {
    let cwd = "/Users/adamdilger/dev/geo/gep-scripts";
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
    ERDDataModel,
    Element,
}

#[derive(Deserialize, Debug)]
struct ERDDataModel {
    _type: String,
    _id: String,
    _parent: Ref,
    name: String,
}

#[derive(Deserialize, Debug)]
struct Element {
    _type: String,
    _id: String,
    name: String,
}

#[derive(Deserialize, Debug)]
struct Ref {
    #[serde(rename = "$ref")]
    _ref: String,
}
