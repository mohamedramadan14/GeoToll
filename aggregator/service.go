package main

import (
	"fmt"

	"github.com/mohamedramadan14/roads-fees-system/types"
)

type Aggregator interface {
	AggregateDistances(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	Store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		Store: store,
	}
}
func (i *InvoiceAggregator) AggregateDistances(distance types.Distance) error {
	fmt.Printf("processing and inserting distance in the storage: %v\n", distance)
	return i.Store.Insert(distance)
}
