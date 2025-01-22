use std::{io::stdout, process::exit};

use crossterm::{
    cursor::{Hide, MoveToNextLine, MoveToRow, Show},
    event::{self, Event, KeyCode},
    execute,
    style::{Attribute, Stylize},
    terminal::{
        disable_raw_mode, enable_raw_mode, Clear, ClearType, EnterAlternateScreen,
        LeaveAlternateScreen,
    },
};
use tmux_session_pane::tmux::{list_sessions, open_session};

fn start_raw_mode() {
    enable_raw_mode().unwrap();
    execute!(stdout(), EnterAlternateScreen, Hide).unwrap();
}

fn stop_raw_mode() {
    execute!(stdout(), LeaveAlternateScreen, Show).unwrap();
    disable_raw_mode().unwrap();
}

fn main() {
    start_raw_mode();

    let sessions = match list_sessions() {
        Ok(sessions) => sessions,
        Err(msg) => {
            eprintln!("{msg}");
            exit(1);
        }
    };

    execute!(stdout()).unwrap();

    let mut select: usize = 0;

    'event_loop: loop {
        execute!(stdout(), MoveToRow(0), Clear(ClearType::All)).unwrap();

        for (i, session) in sessions.iter().enumerate() {
            let mut content = format!(
                "{} (created at {}) {}",
                session.name,
                session.created_at,
                if session.is_attached {
                    "(attached)"
                } else {
                    ""
                }
            );

            if session.is_attached {
                content = format!("{}{content}{}", Attribute::Bold, Attribute::Reset);
            }

            if select == i {
                content = content.negative().to_string();
            }

            print!("{content}");
            execute!(stdout(), MoveToNextLine(1)).unwrap();
        }

        if let Event::Key(key) = event::read().unwrap() {
            match key.code {
                KeyCode::Char('q') => break 'event_loop,
                KeyCode::Up | KeyCode::Char('k') => {
                    select = select.max(1) - 1;
                }
                KeyCode::Down | KeyCode::Char('j') => {
                    select = select.min(sessions.len() - 2) + 1;
                }
                KeyCode::Enter => {
                    disable_raw_mode().unwrap();
                    open_session(&sessions[select]).unwrap();
                    break 'event_loop;
                }
                _ => {}
            }
        }
    }

    stop_raw_mode();
}
