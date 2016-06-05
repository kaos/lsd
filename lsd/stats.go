package main

import (
	"fmt"
	"math"
)

type Statser interface {
	Filterer
	Update(src, dst Object, key string)
}

type Stats struct {
	*FilterExpression
	update StatsFunc
}

// src is nil to indicate that any last calculations should be done, after all objects have been presented
type StatsFunc func(s *Stats, src, dst Object, key string)

func (s *Stats) Update(src, dst Object, key string) {
	s.update(s, src, dst, key)
}

func (j *FilterGroup) Update(src, dst Object, key string) {
	// a stats group (stats join by StatsAnd / StatsOr always only counts
	StatsCount(nil, src, dst, key)
}

var stats_funcs = map[string]StatsFunc{
	"sum": StatsSum,
	"min": StatsCompare(-1),
	"max": StatsCompare(1),
	"avg": StatsAvg,
	"std": StatsStd,
}

func NewStatser(fun, field string) Statser {
	update, ok := stats_funcs[fun]
	if !ok {
		return nil
	}

	return &Stats{
		FilterExpression: &FilterExpression{field: field},
		update:           update,
	}
}

func StatsCount(s *Stats, src, dst Object, key string) {
	if src == nil {
		return
	}

	value := 0
	stat, present := dst[key]
	if present {
		value = stat.(int)
	}
	value++
	dst[key] = value
}

func StatsSum(s *Stats, src, dst Object, key string) {
	if src == nil {
		return
	}

	value := FieldValue(src[s.field])
	stat, present := dst[key]
	if present {
		value += FieldValue(stat)
	}
	dst[key] = value
}

func StatsCompare(test int) StatsFunc {
	return func(s *Stats, src, dst Object, key string) {
		if src == nil {
			return
		}

		value := src[s.field]
		stat, present := dst[key]
		if present {
			if compare(stat, value) == test {
				value = stat
			}
		}
		dst[key] = value
	}
}

func StatsAvg(s *Stats, src, dst Object, key string) {
	type Ctx struct {
		count float64
		total float64
	}

	var ctx Ctx
	stat, present := dst["_"+key]
	if present {
		ctx = stat.(Ctx)
	}

	if src == nil {
		if ctx.count > 0 {
			dst[key] = ctx.total / ctx.count
		}
		return
	}

	ctx.count++
	ctx.total += FieldValue(src[s.field])
	dst["_"+key] = ctx
}

func StatsStd(s *Stats, src, dst Object, key string) {
	type Ctx struct {
		values []float64
		total  float64
	}

	var ctx Ctx
	stat, present := dst["_"+key]
	if present {
		ctx = stat.(Ctx)
	}

	if src != nil {
		value := FieldValue(src[s.field])
		ctx.values = append(ctx.values, value)
		ctx.total += value
		dst["_"+key] = ctx
		return
	}

	if len(ctx.values) == 0 {
		return
	}

	avg := ctx.total / float64(len(ctx.values))
	dev := 0.0
	for _, v := range ctx.values {
		d := v - avg
		dev += d * d
	}
	dst[key] = math.Sqrt(dev / float64(len(ctx.values)))
}

func FieldValue(value interface{}) float64 {
	switch value := value.(type) {
	case int:
		return float64(value)
	case int64:
		return float64(value)
	case float64:
		return value
	case nil:
		return 0
	default:
		panic(fmt.Sprintf("stat on non numeric field: %T", value))
	}
}
