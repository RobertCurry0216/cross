run:
	go run ./cmd/cross -f ./crossword.puz

clean:
	test -d ./bin && gtrash put ./bin
	mkdir ./bin

build: clean
	go build -o ./bin ./cmd/cross

.PHONY: run
