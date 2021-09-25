package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	ITEMS_FILE      = "database/items.json"
	SAVED_LOGS_FILE = "database/savedlogs.json"
)

type Database struct{}

func (db *Database) GetAllItems() []Item {
	file, err := os.Open(ITEMS_FILE)

	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	items := []Item{}
	var item Item

	decoder := json.NewDecoder(file)
	decoder.Token()

	for decoder.More() {
		decoder.Decode(&item)
		items = append(items, item)
	}

	return items
}

func inArray(items []Item, item *Item) (found bool) {
	found = false
	for _, i := range items {
		if i == *item {
			found = true
			return
		}
	}
	return
}

// Get all items that belong to a category
func (db *Database) GetItems(it ItemCategory) (items []Item) {
	if it == ICRecent {
		entries := db.getAllEntries()

		for i := len(entries) - 1; i >= 0; i-- {
			item := db.GetItem(entries[i].Item)

			if !inArray(items, item) {
				items = append(items, *item) // this is bad design :/
			}

			if len(items) == RECENT_ITEMS_LIMIT {
				return items
			}
		}

		return items
	}

	allItems := db.GetAllItems()

	for _, item := range allItems {
		if it == ICAll || item.Type == fmt.Sprint(it) {
			items = append(items, item)
		}
	}

	return items
}

// Get item by name
func (db *Database) GetItem(name string) *Item {
	items := db.GetItems(ICAll)
	for _, item := range items {
		if item.Name == name {
			return &item
		}
	}
	return nil
}

// Get all entries that exist
func (db *Database) getAllEntries() []Entry {
	var entries []Entry
	data, err := ioutil.ReadFile(SAVED_LOGS_FILE)

	if err != nil {
		log.Panicln(err)
	}

	err = json.Unmarshal(data, &entries)

	if err != nil {
		log.Panicln(err)
	}

	return entries
}

func (db *Database) EntryCount() int {
	return len(db.getAllEntries())
}

func (db *Database) GetLatestEntry() Entry {
	var allEntries []Entry = db.getAllEntries()
	return allEntries[len(allEntries)-1]
}

// Get total hours logged for item
func (db *Database) GetHoursLogged(name string) int {
	var result = 0
	for _, entry := range db.getAllEntries() {
		if entry.Item == name {
			result += entry.Hours
		}
	}
	return result
}

// Get entries by date
func (db *Database) GetEntries(d time.Time) []Entry {
	var entries []Entry

	for _, e := range db.getAllEntries() {
		if e.Date.Year() == d.Year() && e.Date.YearDay() == d.YearDay() {
			entries = append(entries, e)
		}
	}

	return entries
}

// Get total hours logged for a given day
func (db *Database) GetTotalHours(date time.Time) int {
	entries := db.GetEntries(date)
	hours := 0
	for _, entry := range entries {
		hours += entry.Hours
	}
	return hours
}

func (db *Database) GetTotalHoursForRange(start time.Time, end time.Time) (hours int) {
	t0 := start.Truncate(24 * time.Hour)
	t1 := end.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	for d := t0; d.Before(t1); d = d.AddDate(0, 0, 1) {
		hours += db.GetTotalHours(d)
	}
	return
}

func (db *Database) SaveEntry(entry Entry) {
	var entries []Entry

	if _, err := os.Stat(SAVED_LOGS_FILE); os.IsNotExist(err) {
		os.OpenFile(SAVED_LOGS_FILE, os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		entries = db.getAllEntries()
	}

	// remove entry if exists
	for i, e := range entries {
		if e.Date.Year() == entry.Date.Year() && e.Date.YearDay() == entry.Date.YearDay() &&
			e.Item == entry.Item {
			entries = append(entries[:i], entries[i+1:]...)
			break
		}
	}

	// don't save entries with 0 hours
	// this allows deleting an entry by udating hours to 0
	if entry.Hours != 0 {
		entries = append(entries, entry)
	}

	file, _ := os.OpenFile(SAVED_LOGS_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()

	writedata, _ := json.MarshalIndent(entries, "", "\t")
	file.WriteString(string(writedata))
}

// Check if entry with same date and for same item exists
func (db *Database) EntryExists(d time.Time, item string) bool {
	for _, e := range db.getAllEntries() {
		if e.Date.Year() == d.Year() && e.Date.YearDay() == d.YearDay() && e.Item == item {
			return true
		}
	}
	return false
}
