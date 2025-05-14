use crossterm::style::{Attribute, Color};

use crate::{renderer::Renderer, tmux::Session, VERSION};

const QUICK_HELP: &'static str = r"
Quick help   ↑ k: up   ↓ j: down   ENTER: select   d: delete   r: rename session
";

struct SessionPaneItem {
    name: String,
    session: Option<Session>,
}

pub struct SessionPane {
    items: Vec<SessionPaneItem>,
    selection: usize,
    selection_row: u16,
}

impl SessionPane {
    pub fn new(sessions: Vec<Session>, selection: usize) -> Self {
        let mut items: Vec<SessionPaneItem> = sessions
            .into_iter()
            .map(|s| SessionPaneItem {
                name: s.to_string(),
                session: Some(s),
            })
            .collect();

        items.push(SessionPaneItem {
            name: String::from(" + new session"),
            session: None,
        });

        Self {
            items,
            selection,
            selection_row: 0,
        }
    }

    pub fn render(&mut self, renderer: &mut Renderer) {
        let col = renderer.term_size().0 as usize;
        let seperator = "=".repeat(col);

        renderer.println(&seperator);

        renderer.println(&format!(
            "{}tmux-sessionpane{} {}",
            Attribute::Bold,
            Attribute::Reset,
            VERSION
        ));

        String::from(QUICK_HELP).trim().lines().for_each(|line| {
            renderer.println(line);
        });

        renderer.println(&seperator).move_to_next_line(1);

        for i in 0..self.items.len() {
            let item = &self.items[i];
            let content = &item.name;

            if let Some(session) = &item.session {
                if session.is_attached {
                    renderer.set_attribute(Attribute::Bold);
                }
            } else {
                renderer
                    .set_foreground_color(Color::Grey)
                    .set_attribute(Attribute::Italic);
            }

            if self.selection == i {
                renderer
                    .set_attribute(Attribute::Bold)
                    .set_foreground_color(Color::Magenta)
                    .set_background_color(Color::DarkGrey);

                self.selection_row = renderer.current_position().1;
            }

            let space = " ".repeat(col - content.len());
            renderer
                .println(&format!("{}{}", content, space))
                .set_attribute(Attribute::Reset);
        }
    }

    pub fn select_next(&mut self) {
        if self.selection == self.items.len() - 1 {
            self.selection = 0;
        } else {
            self.selection += 1;
        }
    }

    pub fn select_prev(&mut self) {
        if self.selection == 0 {
            self.selection = self.items.len() - 1;
        } else {
            self.selection -= 1;
        }
    }

    pub fn get_selected_row(&self) -> u16 {
        self.selection_row
    }

    pub fn get_selected_session(&self) -> Option<&Session> {
        self.items[self.selection].session.as_ref()
    }

    pub fn pop_selected_session(&mut self) -> Option<Session> {
        if self.items[self.selection].session.is_some() {
            let session = self.items.remove(self.selection).session;

            self.selection = self.selection.min(0).max(self.items.len() - 1);

            return session;
        }

        None
    }

    pub fn rename_selected_session(&mut self, name: &String) {
        let item = &mut self.items[self.selection];

        if let Some(session) = item.session.as_mut() {
            session.name = name.clone();
            item.name = session.to_string();
        }
    }
}
