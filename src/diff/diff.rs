#![allow(dead_code)]

use std::collections::HashMap;

use crate::types::{ERDColumn, ERDEntity};

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
    pub name: String,
    pub change_type: ChangeType,
    pub value: String,
    pub old: String,
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

    tc.columns = diff_columns(a, b);
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

    for col in e.columns.iter() {
        let cc = whole_column_change(col, change_type);
        tc.columns.push(cc);
    }

    // for _, tag := range e.Tags {
    // 	cc := wholeTagChange(tag, changeType)
    // 	tc.Tags = append(tc.Tags, cc)
    // }

    tc
}

fn whole_column_change(c: &ERDColumn, change_type: ChangeType) -> ColumnChange {
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
    if let Some(d) = &c.element.documentation {
        cc.changes.push(Change {
            name: String::from("documentation"),
            change_type,
            value: d.clone(),
            old: String::from(""),
        })
    }

    if let Some(d) = &c.length {
        let val = d.to_string();

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

fn diff_columns(a: &ERDEntity, b: &ERDEntity) -> Vec<ColumnChange> {
    let mut changes: Vec<ColumnChange> = Vec::new();

    // get column map
    let a_columns = a.get_column_map();
    let b_columns = b.get_column_map();

    let mut existing_map: HashMap<String, bool> = HashMap::new();

    for (id, a_col) in a_columns {
        let b_col = b_columns.get(&id);
        if let Some(b_col) = b_col {
            existing_map.insert(id, true);

            if let Some(cc) = diff_column(&a_col, b_col) {
                changes.push(cc);
            }
        } else {
            // new A column
            let cc = whole_column_change(a_col, ChangeType::Add);
            changes.push(cc);
        }
    }

    for (id, b_col) in b_columns {
        if existing_map.contains_key(&id) {
            continue; // already been diffed
        }

        let cc = whole_column_change(b_col, ChangeType::Remove);
        changes.push(cc);
    }

    changes
}

fn diff_column(a: &ERDColumn, b: &ERDColumn) -> Option<ColumnChange> {
    let mut cc = ColumnChange {
        id: a.element._id.clone(),
        name: a.element.name.clone(),
        change_type: ChangeType::Modify,
        changes: Vec::new(),
    };

    compare(&mut cc, "name", &a.element.name, &b.element.name);
    compare(
        &mut cc,
        "documentation",
        &a.element.documentation,
        &b.element.documentation,
    );
    compare(&mut cc, "type", &a.column_type, &b.column_type);
    compare(&mut cc, "primaryKey", &a.primary_key, &b.primary_key);
    compare(&mut cc, "foreignKey", &a.foreign_key, &b.foreign_key);
    compare(&mut cc, "nullable", &a.nullable, &b.nullable);
    compare(&mut cc, "unique", &a.unique, &b.unique);
    compare(&mut cc, "length", &a.length, &b.length);

    // for tags
    // cc.Tags = diffTags(a.GetTags(), b.GetTags())

    if cc.changes.len() == 0 {
        // if len(cc.Changes)+len(cc.Tags) == 0 {
        return None;
    }

    Some(cc)
}

fn compare<T>(cc: &mut ColumnChange, name: &str, a: &T, b: &T)
where
    T: std::fmt::Debug + std::cmp::PartialEq,
{
    if a == b {
        return;
    }

    let c = Change {
        name: name.to_string(),
        change_type: ChangeType::Modify,
        value: format!("{:#?}", a),
        old: format!("{:#?}", b),
    };

    cc.changes.push(c);
}
