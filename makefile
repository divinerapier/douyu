build: init
	go build -o _output/bin/danmaku github.com/divinerapier/douyu/cmd/danmaku

init:
	mkdir -p _output/bin

clean:
	rm _output/bin/danmaku
	go clean --cache github.com/divinerapier/douyu/cmd/danmaku
