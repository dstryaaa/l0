.PHONY: build clean

build:
	go build -o l0 .\cmd\main\main.go

clean:
	rm -f l0