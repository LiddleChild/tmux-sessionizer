use std::io::stdout;

use crossterm::{
    cursor::{MoveToNextLine, MoveToRow},
    execute,
    style::{Attribute, Color, ResetColor, SetBackgroundColor, SetForegroundColor},
    terminal::{window_size, Clear, ClearType},
};

use crate::{tmux::Session, VERSION};

const QUICK_HELP: &'static str = r"
Quick help   ↑ k: up   ↓ j: down   ENTER: select   d: delete
";

struct SessionPaneItem {
    name: String,
    session: Option<Session>,
}

pub struct SessionPane {
    items: Vec<SessionPaneItem>,
    selection: usize,
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

        Self { items, selection }
    }

    pub fn render(&self) {
        let col = window_size().unwrap().columns as usize;

        execute!(stdout(), MoveToRow(0), Clear(ClearType::All)).unwrap();

        let seperator = "=".repeat(col);

        print!("{}", seperator);
        execute!(stdout(), MoveToNextLine(1)).unwrap();

        print!(
            "{}tmux-sessionizer{} {}",
            Attribute::Bold,
            Attribute::Reset,
            VERSION
        );
        execute!(stdout(), MoveToNextLine(1)).unwrap();

        String::from(QUICK_HELP).trim().lines().for_each(|line| {
            print!("{line}");
            execute!(stdout(), MoveToNextLine(1)).unwrap();
        });

        print!("{}", seperator);
        execute!(stdout(), MoveToNextLine(2)).unwrap();

        for (i, item) in self.items.iter().enumerate() {
            let mut content = item.name.clone();

            let space = " ".repeat(col - content.len());

            if let Some(session) = &item.session {
                if session.is_attached {
                    content = format!("{}{}", Attribute::Bold, item.name);
                }
            } else {
                execute!(stdout(), SetForegroundColor(Color::Grey)).unwrap();
                content = format!("{}{}", Attribute::Italic, item.name);
            }

            if self.selection == i {
                print!("{}", Attribute::Bold);
                execute!(
                    stdout(),
                    SetForegroundColor(Color::Magenta),
                    SetBackgroundColor(Color::DarkGrey),
                )
                .unwrap();
            }

            print!("{}{}", content, space);
            execute!(stdout(), MoveToNextLine(1), ResetColor).unwrap();
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
}
