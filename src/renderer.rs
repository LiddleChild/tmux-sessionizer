use std::io::stdout;

use crossterm::{
    cursor::{Hide, MoveToColumn, MoveToNextLine, MoveToRow, Show},
    execute,
    style::{Attribute, Color, SetBackgroundColor, SetForegroundColor},
    terminal::{self, Clear, ClearType},
};

pub struct Renderer {
    row: u16,
    col: u16,
}

impl Renderer {
    pub fn new() -> Self {
        execute!(stdout(), MoveToRow(0), MoveToColumn(0)).unwrap();
        Self { row: 0, col: 0 }
    }

    pub fn print(&mut self, s: &str) -> &mut Self {
        print!("{}", s);
        self
    }

    pub fn println(&mut self, s: &str) -> &mut Self {
        self.print(s).move_to_next_line(1)
    }

    pub fn clear_term(&mut self) -> &mut Self {
        self.move_to_row(0).move_to_col(0);
        execute!(stdout(), Clear(ClearType::All)).unwrap();
        self
    }

    pub fn term_size(&self) -> (u16, u16) {
        terminal::size().unwrap()
    }

    pub fn move_to_row(&mut self, row: u16) -> &mut Self {
        self.row = row;
        execute!(stdout(), MoveToRow(row)).unwrap();
        self
    }

    pub fn move_to_col(&mut self, col: u16) -> &mut Self {
        self.col = col;
        execute!(stdout(), MoveToColumn(col)).unwrap();
        self
    }

    pub fn move_to_next_line(&mut self, amount: u16) -> &mut Self {
        self.row += amount;
        execute!(stdout(), MoveToNextLine(amount)).unwrap();
        self
    }

    pub fn set_attribute(&mut self, attr: Attribute) -> &mut Self {
        print!("{}", attr);
        self
    }

    pub fn set_background_color(&mut self, color: Color) -> &mut Self {
        execute!(stdout(), SetBackgroundColor(color)).unwrap();
        self
    }

    pub fn set_foreground_color(&mut self, color: Color) -> &mut Self {
        execute!(stdout(), SetForegroundColor(color)).unwrap();
        self
    }

    pub fn current_position(&self) -> (u16, u16) {
        (self.col, self.row)
    }

    pub fn show_cursor(&mut self) -> &mut Self {
        execute!(stdout(), Show).unwrap();
        self
    }

    pub fn hide_cursor(&mut self) -> &mut Self {
        execute!(stdout(), Hide).unwrap();
        self
    }
}
