.PHONY: bin

bin: cmd/report/report cmd/unproductive/unproductive
	cp cmd/report/report bin/
	cp cmd/unproductive/unproductive bin/

cmd/report/report: cmd/report/main.go
	cd cmd/report ; go build
	strip cmd/report/report

cmd/unproductive/unproductive: cmd/unproductive/main.go
	cd cmd/unproductive ; go build
	strip cmd/unproductive/unproductive
