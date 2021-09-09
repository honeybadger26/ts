# TImeSheet CLI
A CLI helper for TimeSheet

## Instructions
- Create a `data` folder
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
- Total hours logged this week
    - Breakdown of the hours into days 
- Total hours logged today
    - Breakdown of the hours into items
- Switch between item types (My work, My teams work, other)
    - Have the current category of item displayed at top (make it obvious)
    - Add a favorites category ( & functionality to add / remove favorite items)
    - OR instead of Favorites, maybe a 'Most Recent' category. I tend to rotate between about ~5 items at any one time and this wouldn't require user to set favorites.
- Restrict list of entries for particular day
- See info for currently selected item:
    - Total hours logged to date (MUST HAVE FOR CR'S)
- Handle file not exists
- Commands from CLI
    - View todays logs
    - Quick add log for item
    - Quick add hours for all favorites
- Automatic entry for a date range (FOR ADMIN ONLY ? TBC /W Alison) OR Copy previous days entries
    - Make the user confirm their entry (date + all items /w hours) & option to 'Save' & option to 'Cancel'
- Handle long list of items
- Switch between days of the week
- Switch between different weeks
- Functionality to easily amend any day's logged hours
- Change workflow to be 'multi-entry':
    1. Select day
    2. Select item(s)
    3. Then take focus to a different grid with all items from Step 2 (like the lower grid in existing Timesheet)
    4. One-by-one for each item, user can enter hours
    5. Make the user confirm their entry (date + all items /w hours) & option to 'Save' & option to 'Back' / 'Amend' (Take back to Step 4, retain the data initially entered etc.)
