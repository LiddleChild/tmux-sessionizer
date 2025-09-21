.PHONY: install
install:
	go install

debug:
	tail -f ~/.config/tmux-sessionpane/debug.log
