build:
	go build -ldflags="-s -w" -o self-got

install:
	make build
	sudo mv self-got /usr/local/bin/
