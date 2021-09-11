package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	ITEMS_FILE      = "data/items.json"
	SAVED_LOGS_FILE = "data/savedlogs.json"
)

type ItemCategory int

const (
	All ItemCategory = iota
	CR
	Admin
)
const ITEMTYPE_COUNT = 3

func (it ItemCategory) String() string {
	switch it {
	case All:
		return "All"
	case CR:
		return "CR"
	case Admin:
		return "Admin"
	}
	return ""
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
	Date  string
	Item  string
	Hours int
}

type Database struct{}

func (db *Database) getItems(it ItemCategory) []Item {
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
		if it == All || item.Type == fmt.Sprint(it) {
			items = append(items, item)
		}
	}

	return items
}

func (db *Database) getItem(name string) *Item {
	items := db.getItems(All)
	for _, item := range items {
		if item.Name == name {
			return &item
		}
	}
	return nil
}

// get entry for day
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

func (db *Database) getEntries(date time.Time) []Entry {
	var entries []Entry

	for _, entry := range db.getAllEntries() {
		if entry.Date == date.Format("02/01/2006") {
			entries = append(entries, entry)
		}
	}

	return entries
}

func (db *Database) getTotalHours(date time.Time) int {
	entries := db.getEntries(date)
	hours := 0
	for _, entry := range entries {
		hours += entry.Hours
	}
	return hours
}

func (db *Database) saveEntry(entry Entry) {
	var entries []Entry

	if _, err := os.Stat(SAVED_LOGS_FILE); os.IsNotExist(err) {
		os.OpenFile(SAVED_LOGS_FILE, os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		entries = db.getAllEntries()
	}

	// remove entry if exists
	for i, e := range entries {
		if e.Date == entry.Date && e.Item == entry.Item {
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

func (db *Database) entryExists(date string, item string) bool {
	for _, e := range db.getAllEntries() {
		if e.Date == date && e.Item == item {
			return true
		}
	}
	return false
}
