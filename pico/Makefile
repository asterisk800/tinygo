current_dir := $(notdir $(shell pwd))
init:
	go mod init github.com/asterisk800/tinygo/$(current_dir)
	git init
	go mod tidy
flash:
	tinygo flash -target=pico