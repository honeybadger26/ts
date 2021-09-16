# TImeSheet CLI
A CLI helper for TimeSheet

## Instructions
- Create a `data` folder
- Copy the `` a `items.json` in the `data` folder with the following format:
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
- For edit mode: `go run .`
- For view mode: `go run . v`

## Presentation Notes
- Any above items that we want to have but didn't have time, add them to a 'road map' to show our plans for features

## Todo
## Priority 1
- See info for currently selected item: [TO DO]
    - Total hours logged to date (MUST HAVE FOR CR'S)
- Have the current category of item displayed at top (make it obvious) [IN PROGRESS - SYS]
    - `Simon:` I tried to get this at the top left of the view so it would look something like the below, but it just does not want to
```
┌───────────────────Admin─┐
```
- Option to do automatic full-day entry for a date range (FOR Annual Leave, Personal Leave, Public Holidays) [IN PROGRESS - AMS]
- Be able to sign out of whiteboard [TO DO]
## Priority 2
- Functionality to easily amend any day's logged hours [DONE - POLISH]
    - `Simon:` At the moment you can do this by reentering the entry. So to amend CR-1234 you would add another entry for CR-1234 with the correct hours and this will update it. Thoughts on this? I feel like there is probably a better way of doing it
    - `Andre:` Seems to be fine the way it is based on meeting with Alison. But one small thing that could improve it, prepopulate the input field with the existing entry so the user knows what they are amending/ deleting?
    - `Simon:` Also deleteing is the same as amending except enter the new hours as 0
- View Mode
    - Quick switch between view and edit mode (like with a keyboard shortcut) [IN PROGRESS - SYS]
    - Help text for controls [IN PROGRESS - SYS]
    - Month view? [GOOD TO HAVE]
    - Type in date to go to [TO DO]
    - `Simon:` I believe the below are also part of the view mode?
    - Total hours logged this week [TO DO]
        - Breakdown of the hours into days
        - Warning message, if not reached 40 hrs
- Be able to type in date to go to [TO DO]
- Undo shortcut (preferably using ctrl-z)
- Notes for a CR [GOOD TO HAVE, MUST HAVE LATER]
- Add 'open in JIRA' shortcut [MUST HAVE LATER]
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
- Ability to hide/unhide UI tooltips
- Handle file not exists
- Handle long list of text (items and entries). Maybe tab between each window and then let the user scroll? Using key events would be way easier though like Page Up/Down
