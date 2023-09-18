.PHONY: build clean

build:
	go build -o test .\cmd\main\main.go

clean:
	rm -f test