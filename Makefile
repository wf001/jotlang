dir = generated/$(shell date +%s)

all:
	mkdir -p $(dir)
	make build dir=$(dir)
	make run dir=$(dir) -i
	rm -rf modo

clean:
	rm -rf generated
	rm -rf modo

build:
	mkdir -p $(dir)
	go build ./cmd/modo

build-debug:
	mkdir -p $(dir)
	go build -a -p 1 -x -work ./cmd/modo

run:
	@echo "----------------------\n"
	./modo run --debug --exec "1+2"
	@echo $?
	@echo "----------------------\n"


test-all:
	make test-compiler
	make test-go

test-compiler:
	./script/test-lite.sh
	./script/test-full.sh

test-full-compiler:
	./script/test-full.sh

test-lite-compiler:
	./script/test-lite.sh

test-go:
	go test -v ./... | tc
