obu:
	@go build -o ./bin/obu ./obu
	@./bin/obu

receiver:
	@go build -o ./bin/receiver ./obu-data-receiver
	@./bin/receiver

calculator:
	@go build -o ./bin/distance ./distance-calculator
	@./bin/distance

aggregator:
	@go build -o ./bin/aggregator ./aggregator
	@./bin/aggregator

.PHONY: obu aggregator 