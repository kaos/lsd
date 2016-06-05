package main

import (
	"fmt"
)

type Filterer interface {
	Filter(Object) bool
}

type FilterFunc func(interface{}) bool
type FilterExpression struct {
	field string
	op    FilterFunc
}

func (e *FilterExpression) Filter(obj Object) bool {
	value, ok := obj[e.field]
	if !ok {
		panic(fmt.Sprintf("Unknown field '%s' in filter", e.field))
	}
	return e.op == nil || e.op(value)
}

type FilterGrouper interface {
	Group([]Filterer)
}

type FilterGroup struct{ filters []Filterer }
type OrFilter struct{ FilterGroup }
type AndFilter struct{ FilterGroup }
type NegateFilter struct{ FilterGroup }

func (j *FilterGroup) Group(filters []Filterer) {
	j.filters = make([]Filterer, len(filters))
	copy(j.filters, filters)
}

func (e *OrFilter) Filter(obj Object) bool {
	for _, f := range e.filters {
		if f.Filter(obj) {
			return true
		}
	}
	return false
}

func (e *AndFilter) Filter(obj Object) bool {
	for _, f := range e.filters {
		if !f.Filter(obj) {
			return false
		}
	}
	return true
}

func (e *NegateFilter) Filter(obj Object) bool {
	if len(e.filters) != 1 {
		panic(e)
	}

	return !e.filters[0].Filter(obj)
}
