package main

import (
	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/sirupsen/logrus"
)

type MemoryStore struct {
	data map[int]float64
}

func (m *MemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	logrus.Infof("Inserted distance: %v and total distance is: %v", d.Value, m.data[d.OBUID])
	return nil
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}
