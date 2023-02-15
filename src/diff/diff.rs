#![allow(dead_code)]

use std::collections::HashMap;

use crate::types::ERDEntity;

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
    // columns    :  []ColumnChange,
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

fn whole_table_change(e: ERDEntity, change_type: ChangeType) -> TableChange {
    let mut tc = TableChange {
        id: e.element._id.clone(),
        name: e.element.name.clone(),
        change_type,
        changes: Vec::new(),
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

    // for _, col := range e.Columns {
    // 	cc := wholeColumnChange(col, changeType)
    // 	tc.Columns = append(tc.Columns, cc)
    // }

    // for _, tag := range e.Tags {
    // 	cc := wholeTagChange(tag, changeType)
    // 	tc.Tags = append(tc.Tags, cc)
    // }

    tc
}
