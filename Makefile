all: capture

capture: capture.cpp
	g++ $< -o $@ -lX11

sender/sender: sender/sender.go
	(cd sender; go build -o sender)

run_fast: capture sender/sender
	./capture | GOMAXPROCS=2 ./sender/sender -host=$(DIODERHOST)
