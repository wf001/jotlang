dir = generated/$(shell date +%s)

all:
	mkdir -p $(dir)
	make build dir=$(dir)
	make run dir=$(dir) -i

clean:
	rm -rf generated
	rm -rf modo

build:
	mkdir -p $(dir)
	go build ./cmd/modo

run:
	@echo "----------------------\n"
	./modo run --debug -o $(dir)/out
	@echo $?
	@echo "----------------------\n"

