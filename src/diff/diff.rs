#![allow(dead_code)]

use std::collections::HashMap;

use crate::types::{ColumnLength::Num, ColumnLength::Str, ERDColumn, ERDEntity};

pub fn diff_tables(
    project_a: HashMap<String, ERDEntity>,
    project_b: HashMap<String, ERDEntity>,
) -> Vec<TableChange> {
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
                let tc = whole_table_change(ele, ChangeType::Add);
                table_changes.push(tc);
            }
        }
    }

    for (key, ele) in project_b {
        if existing_map.contains_key(&key) {
            continue; // already been diffed
        }

        // new table in master, so mark as "removed" in the diff
        let tc = whole_table_change(ele, ChangeType::Remove);
        table_changes.push(tc);
    }

    // TODO: sorted?
    table_changes.sort_by(|a, b| a.name.cmp(&b.name));
    // sort.Slice(tableChanges, func(i, j int) bool {
    // 	return tableChanges[i].Name < tableChanges[j].Name
    // })

    table_changes
}

#[derive(Debug, Copy, Clone)]
pub enum ChangeType {
    Add,
    Remove,
    Modify,
}

#[derive(Debug)]
pub struct TableChange {
    pub id: String,
    pub change_type: ChangeType,
    pub name: String,
    pub columns: Vec<ColumnChange>,
    // relationships:[]RelationshipChange,
    pub changes: Vec<Change>,
    // tags      :   []TagChange,
    // data_model: Node,
    // diagram: Node,
}

#[derive(Debug)]
pub struct Change {
    name: String,
    change_type: ChangeType,
    value: String,
    old: String,
}

#[derive(Debug)]
pub struct ColumnChange {
    pub id: String,
    pub name: String,
    pub change_type: ChangeType,
    pub changes: Vec<Change>,
    // pub tags: Vec<TagChange>,
}

fn diff_entity(a: &ERDEntity, b: &ERDEntity) -> Option<TableChange> {
    let mut tc = TableChange {
        id: a.element._id.clone(),
        name: a.element.name.clone(),
        change_type: ChangeType::Modify,
        changes: Vec::new(),
        columns: Vec::new(),
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

    if tc.changes.len() + tc.columns.len() == 0 {
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

fn whole_table_change(e: ERDEntity, change_type: ChangeType) -> TableChange {
    let mut tc = TableChange {
        id: e.element._id.clone(),
        name: e.element.name.clone(),
        change_type,
        changes: Vec::new(),
        columns: Vec::new(),
    };

    // optional table fields
    if let Some(d) = e.element.documentation {
        tc.changes.push(Change {
            name: String::from("documentation"),
            change_type,
            value: d.clone(),
            old: String::from(""),
        })
    }

    for col in e.columns {
        let cc = whole_column_change(col, change_type);
        tc.columns.push(cc);
    }

    // for _, tag := range e.Tags {
    // 	cc := wholeTagChange(tag, changeType)
    // 	tc.Tags = append(tc.Tags, cc)
    // }

    tc
}

fn whole_column_change(c: ERDColumn, change_type: ChangeType) -> ColumnChange {
    let mut cc = ColumnChange {
        id: c.element._id.clone(),
        name: c.element.name.clone(),
        change_type,
        changes: Vec::new(),
    };

    cc.changes.push(Change {
        name: String::from("type"),
        value: c.column_type.clone(),
        old: String::from(""),
        change_type,
    });

    // optional column fields
    if let Some(d) = c.element.documentation {
        cc.changes.push(Change {
            name: String::from("documentation"),
            change_type,
            value: d,
            old: String::from(""),
        })
    }

    if let Some(d) = c.length {
        let val = match d {
            Str(b) => b,
            Num(e) => e.to_string(),
        };

        cc.changes.push(Change {
            name: String::from("length"),
            change_type,
            value: val,
            old: String::from(""),
        })
    }

    if let Some(p) = c.primary_key {
        if p {
            cc.changes.push(Change {
                name: String::from("primaryKey"),
                change_type,
                value: String::from("true"),
                old: String::from(""),
            })
        }
    }
    if let Some(f) = c.foreign_key {
        if f {
            cc.changes.push(Change {
                name: String::from("foreignKey"),
                change_type,
                value: String::from("true"),
                old: String::from(""),
            })
        }
    }

    cc
}
