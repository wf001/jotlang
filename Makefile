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

run:
	@echo "----------------------\n"
	./modo run --debug --exec "1"
	@echo $?
	@echo "----------------------\n"


test-all:
	make test-compiler
	make test-go

test-compiler:
	./script/test.sh

test-go:
	go test -v ./...
