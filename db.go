package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	ITEMS_FILE      = "data/items.json"
	SAVED_LOGS_FILE = "data/savedlogs.json"
)

type Item struct {
	Name        string
	Description string
	Size        string
	TotalHours  float32
}

type Entry struct {
	Date  string
	Item  string
	Hours int
}

type Database struct{}

func (db *Database) getItems() []Item {
	file, err := os.Open(ITEMS_FILE)

	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	items := []Item{}
	var data Item

	decoder := json.NewDecoder(file)
	decoder.Token()

	for decoder.More() {
		decoder.Decode(&data)
		items = append(items, data)
	}

	return items
}

func (db *Database) getItem(name string) *Item {
	items := db.getItems()
	for _, item := range items {
		if item.Name == name {
			return &item
		}
	}
	return nil
}

// get entry for day
func (db *Database) getEntries() []Entry {
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

func (db *Database) getTotalHours() int {
	entries := db.getEntries()
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
		entries = db.getEntries()
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
