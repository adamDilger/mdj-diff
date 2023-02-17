#![allow(dead_code)]

mod diff;
mod printer;
mod types;

use diff::diff::diff_tables;
use types::Project;

use crate::printer::printer::{Printer, TextPrinter};

fn main() {
    let cwd = "/Users/adamdilger/dev/geo/gep-scripts";
    let path = "gep.mdj";

    let project_one = Project::from_git(Some("master"), path, cwd);
    // println!("{:#?}", project_one.name);

    let project_two: Project = Project::from_git(None, path, cwd);
    // println!("{:#?}", project_two.name);

    let ok = diff_tables(project_one.get_entity_map(), project_two.get_entity_map());
    // println!("{:#?}", ok);

    let tp = TextPrinter {};
    tp.print(&project_one, &project_two, &ok);
}
