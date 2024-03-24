run:
	go build && ./sink -config-path="./example-config.json"

run-washtub:
	@go build
	@./sink -washtub=127.0.0.1:9000 -config-path="./example-config.json"
	