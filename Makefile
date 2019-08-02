PACKAGE  = cosmicraydetector

# runs the go program
run:
	go run cmd/main.go

runbin:
	build/cosmicraydetector

build:
	go build -o build/$(PACKAGE) -i cmd/main.go

test:
	go test ./cmd -v

dep-init:
	dep init