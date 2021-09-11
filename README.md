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
## Priority 1
- See info for currently selected item:
    - Total hours logged to date (MUST HAVE FOR CR'S)
- Have the current category of item displayed at top (make it obvious)
- Handle file not exists
- Handle long list of text (items and entries). Maybe tab between each window and then let the user scroll?
- Functionality to easily amend any day's logged hours
    - Simon: At the moment you can do this by reentering the entry. So to amend CR-1234 you would add another entry for CR-1234 with the correct hours and this will update it. Thoughts on this? I feel like there is probably a better way of doing it
    - Simon: Also deleteing is the same as amending except enter the new hours as 0
## Priority 2
- Add indication that user is not on current week
- Total hours logged this week
    - Breakdown of the hours into days 
- Total hours logged today
    - Breakdown of the hours into items
- Add a favorites category ( & functionality to add / remove favorite items)
    - OR instead of Favorites, maybe a 'Most Recent' category. Idk about you guys, but I tend to rotate between about ~5 items at any one time and this saves the user the effort of setting favorites.
        - Simon: I prefer 'Most Recent' because as you said there is effort in setting favourites
## Priority 3
- Commands from CLI
    - View todays logs
    - Quick add log for item
    - Quick add hours for all favorites
- Automatic entry for a date range (FOR ADMIN ONLY ? TBC /W Alison) OR Copy previous days entries
    - Make the user confirm their entry (date + all items /w hours) & option to 'Save' & option to 'Cancel'
- Change workflow to be 'multi-entry':
    1. Select day
    2. Select item(s)
    3. Then take focus to a different grid with all items from Step 2 (like the lower grid in existing Timesheet)
    4. One-by-one for each item, user can enter hours
    5. Make the user confirm their entry (date + all items /w hours) & option to 'Save' & option to 'Back' / 'Amend' (Take back to Step 4, retain the data initially entered etc.)
- Be able to sign out of whiteboard
- Notes for a CR