package store

import (
	"fmt"
)

type valueSlice []interface{}

type Filter interface {
	Condition() string
	Values() valueSlice
}

type Direction bool

const (
	Ascending  Direction = true
	Descending           = false
)

///////////////////////////////////////////////////////////

type NullFilter struct{}

func (f NullFilter) Condition() string {
	return "1=1"
}
func (f NullFilter) Values() valueSlice {
	return valueSlice{}
}

///////////////////////////////////////////////////////////

type IdFilter struct {
	Id string
}

func (f IdFilter) Condition() string {
	return "id LIKE ? || '%'"
}
func (f IdFilter) Values() valueSlice {
	return valueSlice{f.Id}
}

///////////////////////////////////////////////////////////

type ExactIdFilter struct {
	Id string
}

func (f ExactIdFilter) Condition() string {
	return "id = ?"
}

func (f ExactIdFilter) Values() valueSlice {
	return valueSlice{f.Id}
}

///////////////////////////////////////////////////////////

type UrlFilter struct {
	Url string
}

func (f UrlFilter) Condition() string {
	return "url = ?"
}

func (f UrlFilter) Values() valueSlice {
	return valueSlice{f.Url}
}

///////////////////////////////////////////////////////////

type TagFilter struct {
	Tag string
}

func (f TagFilter) Condition() string {
	return "tag = ?"
}
func (f TagFilter) Values() valueSlice {
	return valueSlice{f.Tag}
}

///////////////////////////////////////////////////////////

type UntaggedFilter struct{}

func (f UntaggedFilter) Condition() string {
	return "tag IS NULL"
}
func (f UntaggedFilter) Values() valueSlice {
	return valueSlice{}
}

///////////////////////////////////////////////////////////

type RemoteFilter struct {
	Filter Filter
}

func (f RemoteFilter) Condition() string {
	return f.Filter.Condition() + " AND url IS NOT NULL"
}
func (f RemoteFilter) Values() valueSlice {
	return f.Filter.Values()
}

///////////////////////////////////////////////////////////

type RandomOrderer struct {
	Filter Filter
}

func (f RandomOrderer) Condition() string {
	return f.Filter.Condition() + " ORDER BY random()"
}
func (f RandomOrderer) Values() valueSlice {
	return f.Filter.Values()
}

///////////////////////////////////////////////////////////

type DateOrderer struct {
	Filter    Filter
	Direction Direction
}

func (f DateOrderer) Condition() string {
	if f.Direction == Ascending {
		return f.Filter.Condition() + " ORDER BY added_at ASC"
	} else {
		return f.Filter.Condition() + " ORDER BY added_at DESC"
	}
}
func (f DateOrderer) Values() valueSlice {
	return f.Filter.Values()
}

///////////////////////////////////////////////////////////

type Limiter struct {
	Filter Filter
	Limit  int
}

func (f Limiter) Condition() string {
	return f.Filter.Condition() + fmt.Sprintf(" LIMIT %d", f.Limit)
}
func (f Limiter) Values() valueSlice {
	return f.Filter.Values()
}
