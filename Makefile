all:
	- mkdir -p .output/bin
	- go build  -o .output/bin/frack github.com/chmouel/chmoufrack/cli

check:
	gometalinter ./...		
