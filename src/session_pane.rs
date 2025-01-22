use std::io::stdout;

use crossterm::{
    cursor::{MoveToNextLine, MoveToRow},
    execute,
    style::{Attribute, Color, ResetColor, SetBackgroundColor, SetForegroundColor},
    terminal::{window_size, Clear, ClearType},
};

use crate::tmux::Session;

pub struct SessionPane {
    sessions: Vec<Session>,
    selection: usize,
}

impl SessionPane {
    pub fn new(sessions: Vec<Session>, selection: usize) -> Self {
        Self {
            sessions,
            selection,
        }
    }

    pub fn render(&self) {
        execute!(stdout(), MoveToRow(0), Clear(ClearType::All)).unwrap();

        for (i, session) in self.sessions.iter().enumerate() {
            let mut content = session.to_string();

            let col = window_size().unwrap().columns;
            let space = " ".repeat(col as usize - content.len());

            if session.is_attached {
                content = format!("{}{content}", Attribute::Bold);
            }

            if self.selection == i {
                execute!(
                    stdout(),
                    SetForegroundColor(Color::Magenta),
                    SetBackgroundColor(Color::DarkGrey)
                )
                .unwrap();
            }

            print!("{content}{}", space);
            execute!(stdout(), MoveToNextLine(1), ResetColor).unwrap();
        }
    }

    pub fn select_next(&mut self) {
        if self.selection == self.sessions.len() - 1 {
            self.selection = 0;
        } else {
            self.selection += 1;
        }
    }

    pub fn select_prev(&mut self) {
        if self.selection == 0 {
            self.selection = self.sessions.len() - 1;
        } else {
            self.selection -= 1;
        }
    }

    pub fn get_current_session(&self) -> &Session {
        &self.sessions[self.selection]
    }
}
