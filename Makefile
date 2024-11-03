clean:
	rm -rf generated
all:
	mkdir generated
	go run main.go > generated/out.ll
	llc generated/out.ll -o generated/out.s
	clang generated/out.s -o generated/out
	./generated/out

