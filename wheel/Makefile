.PHONY: all
all: build flash

.PHONY: build
build:
	tinygo build -target=bluepill -o main.bin main.go

.PHONY: flash
flash:
	tinygo flash -target=bluepill

