obu:
	@go build -o ./bin/obu ./obu/main.go
	@./bin/obu

receiver:
	@go build -o ./bin/receiver ./obu-data-receiver/main.go
	@./bin/receiver

.PHONY: obu