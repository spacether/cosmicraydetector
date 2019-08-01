# runs the go program
run:
	go run cmd/cosmicraydetector/main.go

runbin:
	build/cosmicraydetector

build:
	go build -o build/cosmicraydetector -i cmd/cosmicraydetector/main.go