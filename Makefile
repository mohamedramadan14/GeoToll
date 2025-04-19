obu:
	@go build -o ./bin/obu ./obu
	@./bin/obu

receiver:
	@go build -o ./bin/receiver ./obu-data-receiver
	@./bin/receiver

.PHONY: obu