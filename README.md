# TImeSheet CLI
A CLI helper for TimeSheet

## Instructions
- In the `data` folder: 
    - Make a copy of the `items.json.base` file and name it `items.json`
    - Make a copy of the `savedlogs.json.base` file and name it `savedlogs.json`
- For edit mode: `go run .`
- For view mode: `go run . v`

## Todo

## Opportunities
- Null fields for items will read the values from previous item in json. (i.e. item without description in items.json file will show description of previous item from .json file) [FIX IN FUTURE]
- Selected category is not saved after submitting entry

## Priority 1
- Rework data structure to match JIRA timesheet data [IN PROGRESS - AMS]
- See info for currently selected item: [TO DO]
    - Total hours logged to date (MUST HAVE FOR CR'S)
- Option to do automatic full-day entry for a date range (FOR Annual Leave, Personal Leave, Public Holidays) [IN PROGRESS - AMS]
- Be able to hide weekend in weekly view

## Priority 2
- Functionality to easily amend any day's logged hours [DONE - POLISH]
    - Prepopulate input field for existing entries
- View Mode
    - Highlight currently selected day and go to this day when switching back to edit mode
    - Month view? [GOOD TO HAVE]
    - Type in date to go to [TO DO]
- Be able to type in date to go to [TO DO]
- Notes for a CR [GOOD TO HAVE, MUST HAVE LATER]
- Make 'Sign out of Whiteboard' work with actual whiteboard

## Priority 3
- Error/ warning messages [NOT TOO IMPORTANT]
    - No error messages for entering greater than 24 hours on current Timesheet. It will simply input 24 if you put >24 
    - CR size 'Please talk to team leader' messages: 80 for small, 200(?) for medium, ??? for large
- Add indication that user is not on current week .Done but need to highlight the text, this is actually a bigger change than expected :/
    - `Andre:` Seems like this one doesn't need to be changed too much from what it already is based on Alison's response
- Commands from CLI
    - View todays logs
    - Quick add log for item
    - Quick add hours for all favorites
- Change workflow to be 'multi-entry':
    1. Select day
    2. Select item(s)
    3. Then take focus to a different grid with all items from Step 2 (like the lower grid in existing Timesheet)
    4. One-by-one for each item, user can enter hours
    5. Make the user confirm their entry (date + all items /w hours) & option to 'Save' & option to 'Back' / 'Amend' (Take back to Step 4, retain the data initially entered etc.)
- Handle file not exists
- Handle long list of text (items and entries). Maybe tab between each window and then let the user scroll? Using key events would be way easier though like Page Up/Down

## Presentation Notes
- Any above items that we want to have but didn't have time, add them to a 'road map' to show our plans for features
