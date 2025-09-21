# tmux-sessionpane

### Installation

1. Install binary
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/LiddleChild/tmux-sessionpane/refs/heads/main/scripts/install.sh)"
```

2. Put key binding in `tmux.conf`
```
bind 'S' new-window "which tmux-sessionpane 2>&1 > /dev/null && tmux-sessionpane || { echo 'tmux-sessionpane does not exist\n' ; read -s -k '?Press any key to continue...' }"
```
