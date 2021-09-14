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
- For edit mode: `go run .`
- For view mode: `go run . v`

## Todo
## Priority 1
- See info for currently selected item:
    - Total hours logged to date (MUST HAVE FOR CR'S)
- Have the current category of item displayed at top (make it obvious)
    - `Simon:` I tried to get this at the top left of the view so it would look something like the below, but it just does not want to
```
┌───────────────────Admin─┐
```
- Handle file not exists
- Handle long list of text (items and entries). Maybe tab between each window and then let the user scroll? Using key events would be way easier though like Page Up/Down
- Functionality to easily amend any day's logged hours
    - `Simon:` At the moment you can do this by reentering the entry. So to amend CR-1234 you would add another entry for CR-1234 with the correct hours and this will update it. Thoughts on this? I feel like there is probably a better way of doing it
    - `Andre:` Seems to be fine the way it is based on meeting with Alison. But one small thing that could improve it, prepopulate the input field with the existing entry so the user knows what they are amending/ deleting?
    - `Simon:` Also deleteing is the same as amending except enter the new hours as 0
- Option to do automatic full-day entry for a date range (FOR Annual Leave, Personal Leave, Public Holidays)
## Priority 2
- View Mode
    - Total hours logged for a each day
    - Total hours logged for the week
    - Help text for controls
    - Month view?
    - Type in date to go to
    - `Simon:` I believe the below are also part of the view mode?
    - Total hours logged this week
        - Breakdown of the hours into days 
    - Total hours logged today
        - Breakdown of the hours into items
- Add a favorites category ( & functionality to add / remove favorite items)
    - OR instead of Favorites, maybe a 'Most Recent' category. Idk about you guys, but I tend to rotate between about ~5 items at any one time and this saves the user the effort of setting favorites.
        - `Simon:` I prefer 'Most Recent' because as you said there is effort in setting favourites
- Error/ warning messages
    - No error messages for entering greater than 24 hours on current Timesheet. It will simply input 24 if you put >24 
    - CR size 'Please talk to team leader' messages: 80 for small, 200(?) for medium, ??? for large
- Be able to type in date to go to
- Undo shortcut (preferably using ctrl-z)
## Priority 3
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
- Be able to sign out of whiteboard
- Notes for a CR
- Ability to hide/unhide UI tooltips
- Allow program to take parameters (which in turn allows for desktop shortcuts)
- Add 'open in JIRA' shortcut
