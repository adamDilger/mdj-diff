use crate::{diff::diff::TableChange, types::Project};

// const WIDTH: usize = 80;
const EQUALS_LINE: &str =
    "================================================================================";

pub trait Printer {
    fn render_title(&self, changes: &Vec<TableChange>);

    fn render_changeset_header(name: String, change_count: usize);

    fn render_entity_header(&self, tc: &TableChange);

    fn render_table(&self, a: &Project, b: &Project, changes: &Vec<TableChange>);

    fn print(&self, a: &Project, b: &Project, changes: &Vec<TableChange>);
}

pub struct TextPrinter {}

impl Printer for TextPrinter {
    fn render_title(&self, changes: &Vec<TableChange>) {
        print!("{}\nMDJ Diff\n{}\n\n", EQUALS_LINE, EQUALS_LINE);
        print!("Entities changed: {}\n\n", changes.len());
    }

    fn render_changeset_header(name: String, change_count: usize) {
        print!("{} {}:\n", name, change_count)
    }

    fn render_entity_header(&self, change: &TableChange) {
        let mut o = String::new();
        o += EQUALS_LINE;
        o += &format!("{} - {}\n", change.name, change.change_type.to_string());
        o += EQUALS_LINE;

        change.change_type.print_color(&o);
    }

    fn render_table(&self, _a: &Project, _b: &Project, _changes: &Vec<TableChange>) {}

    fn print(&self, a: &Project, b: &Project, changes: &Vec<TableChange>) {
        self.render_title(changes);

        for change in changes.iter() {
            self.render_entity_header(change);
            // 	printTableAttributes(p, change)

            println!();

            if change.columns.len() > 0 {
                TextPrinter::render_changeset_header(String::from("Columns"), change.columns.len());
            }

            // 	for _, column := range change.Columns {
            // 		printColumnChanges(p, column)
            // 		fmt.Println()
            // 	}

            // 	if len(change.Relationships) > 0 {
            // 		p.RenderChangesetHeader("Relationships", len(change.Relationships))
            // 	}

            // 	for _, rel := range change.Relationships {
            // 		printRelationshipChanges(p, A, B, rel)
            // 		fmt.Println()
            // 	}

            // 	if len(change.Tags) > 0 {
            // 		p.RenderChangesetHeader("Tags", len(change.Tags))
            // 	}

            // 	for _, tag := range change.Tags {
            // 		renderChanges(p, tag.Name, tag.Changes)
            // 		fmt.Println()
            // 	}
        }
    }
}
