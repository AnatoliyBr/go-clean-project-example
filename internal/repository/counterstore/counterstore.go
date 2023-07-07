package counterstore

import (
	"math"

	"github.com/AnatoliyBr/go-clean-project-example/internal/entity"
)

type CounterStore struct {
	counters entity.Counters
}

func NewCounterStore(initvals map[string]int) *CounterStore {
	return &CounterStore{
		counters: initvals,
	}
}

func (cs *CounterStore) Set(name string, val int) int {
	if val < 0 {
		return -1
	} else {
		cs.counters[name] = val
		return val
	}
}

func (cs *CounterStore) Get(name string) int {
	if val, ok := cs.counters[name]; ok {
		return val
	} else {
		return -1
	}
}

func (cs *CounterStore) Inc(name string) int {
	if _, ok := cs.counters[name]; ok {
		if cs.counters[name] == math.MaxInt {
			return -2
		} else {
			cs.counters[name]++
			return cs.counters[name]
		}
	} else {
		return -1
	}
}

func (cs *CounterStore) Dec(name string) int {
	if _, ok := cs.counters[name]; ok {
		if cs.counters[name] == 0 {
			return -2
		} else {
			cs.counters[name]--
			return cs.counters[name]
		}
	} else {
		return -1
	}
}
