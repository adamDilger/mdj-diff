mod diff;
mod types;

use diff::diff::diff_tables;
use types::Project;

fn main() {
    let cwd = "/Users/adamdilger/geo/geometry/gep/gep-scripts";
    let path = "gep.mdj";

    let project_one = Project::from_git(Some("master"), path, cwd);
    // println!("{:#?}", project_one.name);

    let project_two: Project = Project::from_git(None, path, cwd);
    // println!("{:#?}", project_two.name);

    diff_tables(project_one.get_entity_map(), project_two.get_entity_map());
}
