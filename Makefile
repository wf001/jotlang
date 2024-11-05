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
	llc ./a.ll -o ./a/out.s
	clang ./out.s -o .//out

run:
	@echo "----------------------\n"
	@./$(dir)/out
	@echo $?
	@echo "----------------------\n"

