# Interactive Mode Guide

The interactive mode provides a rich, keyboard-driven interface for managing your backlog items.

## Starting Interactive Mode

```bash
bl list -i
```

or

```bash
bl list --interactive
```

## Interface Layout

```
┌─────────────────────────────────────────────────────────────────┐
│                   📋 BACKLOG KANBAN BOARD                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ╭─────────────╮  ╭─────────────╮  ╭─────────────╮            │
│  │ ▶ TODO      │  │ IN PROGRESS │  │    DONE     │            │
│  │             │  │             │  │             │            │
│  │ ▶ Item 1    │  │   Item 3    │  │   Item 5    │            │
│  │   Item 2    │  │   Item 4    │  │             │            │
│  │   Item 4    │  │             │  │             │            │
│  ╰─────────────╯  ╰─────────────╯  ╰─────────────╯            │
│                                                                 │
│  Total: 5 items (3 todo, 1 in-progress, 1 done)               │
│                                                                 │
│  Navigation: ←/→ columns, ↑/↓ items | Actions: 1=todo ...     │
└─────────────────────────────────────────────────────────────────┘
```

## Keyboard Controls

### Navigation
- **`←` (Left Arrow)**: Move to the previous column (left)
- **`→` (Right Arrow)**: Move to the next column (right)
- **`↑` (Up Arrow)**: Move to the previous item in the current column
- **`↓` (Down Arrow)**: Move to the next item in the current column

### Actions
- **`Enter`**: Edit the selected item (opens editable detail view)
- **`s`**: Search/filter items
- **`a`**: Add a new item (opens a form)
- **`1`**: Move the selected item to TODO column
- **`2`**: Move the selected item to IN PROGRESS column
- **`3`**: Move the selected item to DONE column
- **`d`**: Delete the selected item
- **`r`**: Reload data from disk (useful if data was changed externally)

### Other
- **`?`**: Toggle help text on/off
- **`q`**: Quit interactive mode and return to terminal

## Visual Indicators

- **`▶`** next to column name: Currently selected column
- **`▶`** next to item: Currently selected item
- **Purple border**: Selected column is highlighted with a purple border
- **Green message**: Success messages appear at the top after actions
- **Item details**: Each item shows title, tags (🏷), and due date (📅)

## Workflow Example

1. Start interactive mode: `bl list -i`
2. Use `→` to navigate to the TODO column
3. Use `↓` to select an item you want to work on
4. Press `2` to move it to IN PROGRESS
5. When done, press `3` to move it to DONE
6. Press `q` to exit

## Tips

- The interface updates in real-time as you make changes
- All changes are immediately saved to `~/backlog/items.json`
- You can quickly move items between columns without typing commands
- Use `r` to refresh if you're collaborating or editing the JSON file directly
- Press `?` to hide the help text for a cleaner view

## Adding New Items

Press `a` to open the add item form. The form has four fields:

1. **Title** (required) - The name of your backlog item
2. **Description** (optional) - Detailed description
3. **Due Date** (optional) - Format: DD-MM-YYYY
4. **Tags** (optional) - Comma-separated tags

### Form Navigation
- **Tab** or **↓**: Move to next field
- **Shift+Tab** or **↑**: Move to previous field
- **Enter** (on last field): Submit the form
- **Esc**: Cancel and return to board view

The new item will be created with status "todo" and appear in the TODO column.

## Editing Item Details

Press `Enter` on any selected item to open an editable detail view. The detail view shows:

- **ID**: The unique identifier (read-only)
- **Status**: Current status (read-only, use 1/2/3 keys on board to change)
- **Title**: Editable text input
- **Description**: Editable text input
- **Due Date**: Editable text input (DD-MM-YYYY format)
- **Tags**: Editable text input (comma-separated)
- **Created**: Creation timestamp (read-only)
- **Updated**: Last update timestamp (read-only)

### Editing Controls
- **Tab**, **↑**, **↓**: Navigate between fields
- **Esc**: Save changes and return to board

You can type any character (including 'q') in the text fields. All changes are validated and saved to disk immediately when you press Esc.

## Searching and Filtering

Press `s` to open the search prompt. You can:

1. Type your search query (case-insensitive)
2. Press `Enter` to apply the filter
3. Press `Esc` to cancel

The search looks for matches in:
- Item titles
- Item descriptions
- Item tags

When a filter is active:
- The board title shows "(filtered: 'your query')"
- Only matching items are displayed
- Press `s` again and clear the search (or press Esc) to show all items

## Responsive Layout

The kanban board automatically adjusts to your terminal width:
- Columns expand to use available space
- Minimum column width is 30 characters
- Item text truncates intelligently to fit column width

## Advantages Over CLI Commands

| Task | CLI Command | Interactive Mode |
|------|-------------|------------------|
| Add new item | `bl add "title" --desc "..." --due "..." --tags "..."` | Press `a` + fill form |
| Edit item details | `bl update 176457 --title "..." --desc "..." --due "..." --tags "..."` | Navigate + press `Enter` + edit + Esc |
| Move item to in-progress | `bl update 176457 --status in-progress` | Navigate + press `2` |
| Delete an item | `bl delete 176457` | Navigate + press `d` |
| Search items | `bl search "keyword"` | Press `s` + type query + Enter |
| View all items | `bl list` | `bl list -i` (with navigation) |
| Multiple updates | Multiple commands | Quick keyboard shortcuts |

Interactive mode is perfect for:
- Quick status updates
- Adding and editing items on the fly
- Searching and filtering your backlog
- Reviewing your backlog
- Organizing multiple items
- Getting a visual overview of your work
- Working efficiently without remembering item IDs

