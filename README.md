# tmux-sessionizer

### Installation

`tmux.conf`

```
bind 'S' new-window "which tmux-sessionizer 2>&1 > /dev/null && tmux-sessionizer || { echo 'tmux-sessionizer does not exist\n' ; read -s -k '?Press any key to continue...' }"
```
