#![allow(dead_code)]

use std::collections::HashMap;

use crate::types::{ERDColumn, ERDEntity, ERDRelationship};

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

    table_changes.sort_by(|a, b| a.name.cmp(&b.name));

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
    relationships: Vec<RelationshipChange>,
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

#[derive(Debug)]
pub struct RelationshipChange {
    id: String,
    name: String,
    change_type: ChangeType,
    end1_cardinality: Option<Change>,
    end2_cardinality: Option<Change>,
    end1_reference: Option<Change>,
    end2_reference: Option<Change>,
    changes: Vec<Change>,
    // tags            []TagChange
}

fn diff_entity(a: &ERDEntity, b: &ERDEntity) -> Option<TableChange> {
    let mut tc = TableChange {
        id: a.element._id.clone(),
        name: a.element.name.clone(),
        change_type: ChangeType::Modify,
        changes: Vec::new(),
        columns: Vec::new(),
        relationships: Vec::new(),
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
    tc.relationships = diff_relationships(a, b);
    // 	tc.Tags = diffTags(a.GetTags(), b.GetTags())

    if tc.changes.len() + tc.columns.len() + tc.relationships.len() == 0 {
        return None;
    }
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
        relationships: Vec::new(),
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

fn diff_relationships(a: &ERDEntity, b: &ERDEntity) -> Vec<RelationshipChange> {
    let mut changes: Vec<RelationshipChange> = Vec::new();

    // get relationship map
    let a_rels = a.get_relationship_map();
    let b_rels = b.get_relationship_map();

    let mut existing_map: HashMap<String, bool> = HashMap::new();

    for (id, a_rel) in a_rels {
        let b_rel = b_rels.get(&id);
        if let Some(b_rel) = b_rel {
            existing_map.insert(id, true);
            if let Some(r) = diff_relationship(a_rel, b_rel) {
                changes.push(r);
            }
        } else {
            // new A Relationhsip
            let r = whole_relationship_change(a_rel, ChangeType::Add);
            changes.push(r);
        }
    }

    for (id, b_rel) in b_rels.iter() {
        if existing_map.contains_key(id) {
            continue; // already been diffed
        }

        let r = whole_relationship_change(b_rel, ChangeType::Remove);
        changes.push(r);
    }

    changes
}

fn diff_relationship(a: &ERDRelationship, b: &ERDRelationship) -> Option<RelationshipChange> {
    let mut r = RelationshipChange {
        id: a._id.clone(),
        name: a.end2.reference._ref.clone(),
        change_type: ChangeType::Modify,
        end1_cardinality: None,
        end2_cardinality: None,
        end1_reference: None,
        end2_reference: None,
        changes: Vec::new(),
    };

    let mut change = false;

    if a.end1.get_cardinality() != b.end1.get_cardinality() {
        change = true;
        r.end1_cardinality = Some(Change {
            name: String::from("end1.cardinality"),
            change_type: ChangeType::Modify,
            value: a.end1.get_cardinality(),
            old: b.end1.get_cardinality(),
        });
    }

    if a.end2.get_cardinality() != b.end2.get_cardinality() {
        change = true;
        r.end2_cardinality = Some(Change {
            name: String::from("end2.cardinality"),
            change_type: ChangeType::Modify,
            value: a.end2.get_cardinality(),
            old: b.end2.get_cardinality(),
        });
    }

    if a.end1.reference._ref != b.end1.reference._ref {
        change = true;
        r.end1_reference = Some(Change {
            name: String::from("end1.reference"),
            change_type: ChangeType::Modify,
            value: a.end1.reference._ref.clone(),
            old: b.end1.reference._ref.clone(),
        });
    }

    if a.end2.reference._ref != b.end2.reference._ref {
        change = true;
        r.end2_reference = Some(Change {
            name: String::from("end2.reference"),
            change_type: ChangeType::Modify,
            value: a.end2.reference._ref.clone(),
            old: b.end2.reference._ref.clone(),
        });
    }

    if a.documentation != b.documentation {
        change = true;
        r.changes.push(Change {
            name: String::from("documentation"),
            change_type: ChangeType::Modify,
            value: a.documentation.clone().unwrap_or(String::from("")),
            old: b.documentation.clone().unwrap_or(String::from("")),
        });
    }

    // 	r.Tags = diffTags(a.GetTags(), b.GetTags())

    // optional relationship fields
    // if !change && r.rags) == 0 {
    if !change {
        return None;
    }

    Some(r)
}

fn whole_relationship_change(c: &ERDRelationship, change_type: ChangeType) -> RelationshipChange {
    let mut r = RelationshipChange {
        id: c._id.clone(),
        name: c.end2.reference._ref.clone(),
        change_type,
        end1_cardinality: None,
        end2_cardinality: None,
        end1_reference: None,
        end2_reference: None,
        changes: Vec::new(),
    };

    r.end1_cardinality = Some(Change {
        name: String::from("end1.cardinality"),
        change_type,
        value: c.end1.get_cardinality(),
        old: String::from(""),
    });
    r.end2_cardinality = Some(Change {
        name: String::from("end2.cardinality"),
        change_type,
        value: c.end2.get_cardinality(),
        old: String::from(""),
    });
    r.end1_reference = Some(Change {
        name: String::from("end1.reference"),
        change_type,
        value: c.end1.reference._ref.clone(),
        old: String::from(""),
    });
    r.end2_reference = Some(Change {
        name: String::from("end2.reference"),
        change_type,
        value: c.end2.reference._ref.clone(),
        old: String::from(""),
    });

    // optional relationship fields
    if let Some(d) = &c.documentation {
        r.changes.push(Change {
            name: String::from("documentation"),
            change_type,
            value: d.clone(),
            old: String::from(""),
        });
    }

    r
}
