use std::{io::stdout, process::exit};

use crossterm::{
    cursor::{Hide, Show},
    event::{self, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use tmux_sessionizer::{
    session_pane::SessionPane,
    tmux::{list_sessions, open_session},
};

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

    let mut session_pane = SessionPane::new(sessions);

    'event_loop: loop {
        session_pane.render();

        if let Event::Key(key) = event::read().unwrap() {
            match key.code {
                KeyCode::Char('q') | KeyCode::Esc => break 'event_loop,
                KeyCode::Up | KeyCode::Char('j') => session_pane.select_next(),
                KeyCode::Down | KeyCode::Char('k') => session_pane.select_prev(),
                KeyCode::Enter => {
                    disable_raw_mode().unwrap();
                    open_session(session_pane.get_current_session()).unwrap();
                    break 'event_loop;
                }
                _ => {}
            }
        }
    }

    stop_raw_mode();
}
