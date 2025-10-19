# tmux-sessionizer

### Installation

1. Install binary

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/LiddleChild/tmux-sessionizer/refs/heads/main/scripts/install.sh)"
```

2. Put key binding in `tmux.conf`

```
bind 'S' new-window "which tmux-sessionizer 2>&1 > /dev/null && tmux-sessionizer || { echo 'tmux-sessionizer does not exist\n' ; read -s -k '?Press any key to continue...' }"
```
