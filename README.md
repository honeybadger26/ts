# TImeSheet CLI
A CLI helper for TimeSheet

## Instructions
- Create a `data` folder and 2 files:
- Create a `items.json` in the `data` folder with the following format:
```
[
    {
        "name": "Item1",
        "description": "Item1 Description",
        "size": "Small",
        "totalhours": 10
    },
    ...
]
```
- Create a `savedlogs.json` in the `data` folder with the following format (empty list):
```
[]
```
- Run `go run .`

## Todo
- Total hours logged today
- Switch between item types (My work, My teams work, other)
- Restrict list of entries for particular day
- See info for currently selected item:
    - Total hours logged to date
- Handle file not exists
- Commands from CLI
    - View todays logs
    - Quick add log for item
- Automatic entry for a date range OR Copy previous days entries
- Handle long list of items