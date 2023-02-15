#![allow(dead_code)]

use serde::Deserialize;
use std::{collections::HashMap, fs, path::Path, process::Command};

fn main() {
    let cwd = "/Users/adamdilger/geo/geometry/gep/gep-scripts";
    let path = "gep.mdj";

    let project_one = Project::from_git(Some("master"), path, cwd);
    // println!("{:#?}", project_one.name);

    let project_two: Project = Project::from_git(None, path, cwd);
    // println!("{:#?}", project_two.name);

    diff_tables(project_one.get_entity_map(), project_two.get_entity_map());
}

#[derive(Deserialize, Debug)]
struct Project {
    _type: String,
    _id: String,
    name: String,

    #[serde(rename = "ownedElements")]
    owned_elements: Vec<Node>,
}

impl Project {
    fn from_file(path: &str, cwd: &str) -> Project {
        let contents = fs::read_to_string(Path::new(cwd).join(path))
            .expect("Should have been able to read the file");

        let project: Project = serde_json::from_str(&contents).unwrap();

        project
    }

    fn from_git(commit: Option<&str>, path: &str, cwd: &str) -> Project {
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

    fn get_entity_map(self) -> HashMap<String, ERDEntity> {
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

fn diff_tables(project_a: HashMap<String, ERDEntity>, project_b: HashMap<String, ERDEntity>) {
    let mut existing_map: HashMap<String, bool> = HashMap::new();
    let mut table_changes: Vec<TableChange> = Vec::new();

    for (key, ele) in project_a {
        let b_entity = project_b.get(&key);

        match b_entity {
            Some(b_e) => {
                existing_map.insert(key, true);

                let tc = diff_entity(&ele, b_e);

                match tc {
                    Some(tc) => {
                        println!("{:#?}", tc);
                        table_changes.push(tc);
                    }
                    None => (),
                }
            }
            None => {
                // new A table
                // tc := wholeTableChange(e, ChangeTypeAdd)
                // tableChanges = append(tableChanges, tc)
            }
        }
    }

    // for id, e := range B {
    // 	if _, ok := existingMap[id]; ok {
    // 		continue // already been diffed
    // 	}

    // 	// new table in master, so mark as "removed" in the diff
    // 	tc := wholeTableChange(e, ChangeTypeRemove)
    // 	tableChanges = append(tableChanges, tc)
    // }

    // sort.Slice(tableChanges, func(i, j int) bool {
    // 	return tableChanges[i].Name < tableChanges[j].Name
    // })

    // return tableChanges
}

#[derive(Debug)]
enum ChangeType {
    Add,
    Remove,
    Modify,
}

#[derive(Debug)]
struct TableChange {
    id: String,
    change_type: ChangeType,
    name: String,
    // columns    :  []ColumnChange,
    // relationships:[]RelationshipChange,
    changes: Vec<Change>,
    // tags      :   []TagChange,
    // data_model: Node,
    // diagram: Node,
}

#[derive(Debug)]
struct Change {
    name: String,
    change_type: ChangeType,
    value: String,
    old: String,
}

fn diff_entity(a: &ERDEntity, b: &ERDEntity) -> Option<TableChange> {
    let mut tc = TableChange {
        id: a.element._id.clone(),
        name: a.element.name.clone(),
        change_type: ChangeType::Modify,
        changes: Vec::new(),
    };

    if a.element.name != b.element.name {
        let c = Change {
            name: String::from("name"),
            change_type: ChangeType::Modify,
            value: a.element.name.clone(),
            old: b.element.name.clone(),
        };
        tc.changes.push(c);
    }

    if a.element.documentation != b.element.documentation {
        let c = Change {
            name: String::from("documentation"),
            change_type: ChangeType::Modify,
            value: a.element.documentation.clone().unwrap_or(String::from("")),
            old: b.element.documentation.clone().unwrap_or(String::from("")),
        };
        tc.changes.push(c);
    }

    // 	tc.Columns = diffColumns(a, b)
    // 	tc.Relationships = diffRelationships(a, b)
    // 	tc.Tags = diffTags(a.GetTags(), b.GetTags())

    if tc.changes.len() == 0 {
        return None;
    }
    // 		len(tc.Columns)+
    // 		len(tc.Relationships)+
    // 		len(tc.Tags) == 0 {
    // 		return nil
    // 	}
    //
    // 	return tc
    //

    Some(tc)
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
