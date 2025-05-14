use std::{io::stdout, process::exit};

use crossterm::{
    cursor::{Hide, Show},
    event::{self, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use tmux_sessionpane::{
    renderer::Renderer,
    session_pane::SessionPane,
    tmux::{kill_session, list_sessions, new_session, open_session, rename_session},
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

    let mut old_input = String::new();
    let mut input = String::new();
    let mut input_mode = false;

    let mut renderer = Renderer::new();

    'event_loop: loop {
        renderer.clear_term();

        session_pane.render(&mut renderer);

        if input_mode {
            renderer
                .show_cursor()
                .move_to_row(session_pane.get_selected_row())
                .move_to_col(input.len() as u16);
        } else {
            renderer.hide_cursor();
        }

        if let Event::Key(key) = event::read().unwrap() {
            match key.code {
                KeyCode::Char(c) if input_mode => {
                    input.push(c);
                    session_pane.rename_selected_session(&input);
                }
                KeyCode::Backspace if input_mode => {
                    input.pop();
                    session_pane.rename_selected_session(&input);
                }
                KeyCode::Esc if input_mode => {
                    input_mode = false;
                    session_pane.rename_selected_session(&old_input);
                    rename_session(&old_input, &old_input).unwrap();
                }
                KeyCode::Enter if input_mode => {
                    input_mode = false;
                    match rename_session(&old_input, &input) {
                        Ok(true) => session_pane.rename_selected_session(&input),
                        _ => session_pane.rename_selected_session(&old_input),
                    };
                }
                KeyCode::Char('r') => {
                    if let Some(session) = session_pane.get_selected_session() {
                        input_mode = true;
                        input = session.name.clone();
                        old_input = session.name.clone();
                    }
                }
                KeyCode::Esc | KeyCode::Char('q') => break 'event_loop,
                KeyCode::Up | KeyCode::Char('k') => session_pane.select_prev(),
                KeyCode::Down | KeyCode::Char('j') => session_pane.select_next(),
                KeyCode::Enter => match session_pane.get_selected_session() {
                    Some(session) => {
                        stop_raw_mode();
                        open_session(session).unwrap();
                        break 'event_loop;
                    }
                    None => {
                        new_session();
                        match list_sessions().unwrap().last() {
                            Some(session) => {
                                stop_raw_mode();
                                open_session(session).unwrap();
                                break 'event_loop;
                            }
                            None => {
                                eprintln!("error creating new session");
                            }
                        }
                    }
                },
                KeyCode::Char('d') => {
                    if let Some(session) = session_pane.pop_selected_session() {
                        kill_session(&session);
                    }
                }
                _ => {}
            }
        }
    }

    stop_raw_mode();
}
