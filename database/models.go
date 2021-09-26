package database

import (
	"database/sql"
	"time"
)

const (
	RECENT_ITEMS_LIMIT = 5
)

type ItemCategory int

const (
	ICRecent ItemCategory = iota
	ICAll
	ICCR
	ICAdmin
)

func (it ItemCategory) String() string {
	switch it {
	case ICRecent:
		return "Recent"
	case ICAll:
		return "All"
	case ICCR:
		return "CR"
	case ICAdmin:
		return "Admin"
	}
	return ""
}

func (it ItemCategory) GetNext() ItemCategory {
	categories := []ItemCategory{
		ICRecent,
		ICAll,
		ICCR,
		ICAdmin,
	}
	return categories[(int(it+1) % len(categories))]
}

type Item struct {
	ItemId      int
	Name        string
	Description sql.NullString
	Size        sql.NullString
	URL         sql.NullString
	Type        string
}

type Entry struct {
	Date  time.Time
	Item  string
	Hours int // should be a float
}
