

build:
	go build -buildmode=c-archive -o runnc.a main.go

	gcc -pthread main.c runnc.a -o main
