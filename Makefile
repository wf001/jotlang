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
	go run main.go > $(dir)/out.ll
	llc $(dir)/out.ll -o $(dir)/out.s
	clang $(dir)/out.s -o $(dir)/out

run:
	@echo "----------------------\n"
	@./$(dir)/out
	@echo $?
	@echo "----------------------\n"

