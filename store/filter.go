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

type RandomOrderer struct {
	Filter Filter
	Single bool
}

func (f RandomOrderer) Condition() string {
	if f.Single {
		return f.Filter.Condition() + " AND images._ROWID_ >= (abs(random()) % (SELECT max(_ROWID_) FROM images)) LIMIT 1"
	} else {
		return f.Filter.Condition() + " ORDER BY random()"
	}
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
