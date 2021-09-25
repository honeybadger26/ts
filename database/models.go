package database

import "time"

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
	Name        string
	Description string
	Size        string
	TotalHours  float32
	URL         string
	Type        string
}

type Entry struct {
	Date  time.Time
	Item  string
	Hours int // should be a float
}
