use std::{io::stdout, process::exit};

use crossterm::{
    cursor::{Hide, Show},
    event::{self, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use tmux_sessionizer::{
    session_pane::SessionPane,
    tmux::{list_sessions, new_session, open_session},
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

    let current_session = sessions
        .iter()
        .enumerate()
        .find_map(|(i, s)| if s.is_attached { Some(i) } else { None })
        .unwrap_or(0);

    let mut session_pane = SessionPane::new(sessions, current_session);

    'event_loop: loop {
        session_pane.render();

        if let Event::Key(key) = event::read().unwrap() {
            match key.code {
                KeyCode::Esc | KeyCode::Char('q') => break 'event_loop,
                KeyCode::Up | KeyCode::Char('k') => session_pane.select_prev(),
                KeyCode::Down | KeyCode::Char('j') => session_pane.select_next(),
                KeyCode::Enter => match session_pane.get_current_session() {
                    Some(session) => {
                        disable_raw_mode().unwrap();
                        open_session(session).unwrap();
                        break 'event_loop;
                    }
                    None => {
                        new_session();
                        match list_sessions().unwrap().last() {
                            Some(session) => {
                                disable_raw_mode().unwrap();
                                open_session(session).unwrap();
                                break 'event_loop;
                            }
                            None => {
                                eprintln!("error creating new session");
                            }
                        }
                    }
                },
                _ => {}
            }
        }
    }

    stop_raw_mode();
}
