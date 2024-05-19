default: out/build

clean:
	rm -rf out

test: painter
	go test ./painter/...

out/build: cmd/painter/main.go
	mkdir -p out
	go build -o out/build ./cmd/painter/main.go