use tmux_session_pane::tmux::list_sessions;

fn main() {
    let sessions = list_sessions();

    println!("{sessions:#?}");
}
