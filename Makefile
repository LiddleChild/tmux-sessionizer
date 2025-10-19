.PHONY: install
install:
	go install

debug:
	tail -f ~/.config/tmux-sessionizer/debug.log
