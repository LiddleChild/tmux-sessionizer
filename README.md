# tmux-sessionpane

### Installation

1. Clone repository

2. Install binary
```bash
cargo install --path .
```

3. Put key binding in `tmux.conf`

```
bind 'S' new-window "which tmux-sessionpane 2>&1 > /dev/null && tmux-sessionpane || { echo 'tmux-sessionpane does not exist\n' ; read -s -k '?Press any key to continue...' }"
```
