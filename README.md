# Backlog - Terminal Backlog Manager

A simple and efficient terminal application for managing backlog items with a Kanban-style board view.

## Features

- ✅ Create, update, delete backlog items
- 📋 Kanban-style board view (Todo, In Progress, Done)
- 🎮 **Interactive mode** with keyboard navigation
- 🔍 Search functionality
- 🏷️ Tag support
- 📅 Due date tracking (DD-MM-YYYY format)
- 📦 Archive completed items
- 💾 JSON-based storage in `~/backlog`

## Installation

### Build from source

```bash
go build -o bl
```

### Install globally (optional)

```bash
# Move the binary to a location in your PATH
sudo mv bl /usr/local/bin/
```

## Usage

### Add a new backlog item

```bash
bl add "Task title" --desc "Description" --due "15-12-2025" --tags "tag1,tag2"
```

**Options:**
- `--desc`: Description of the task
- `--due`: Due date in DD-MM-YYYY format
- `--tags`: Comma-separated tags

### List all items (Kanban board view)

**Static view:**
```bash
bl list
```

**Interactive mode:**
```bash
bl list -i
# or
bl list --interactive
```

In interactive mode, you can:
- Navigate between columns with `←` and `→` arrow keys
- Navigate between items with `↑` and `↓` arrow keys
- Press `Enter` to edit the selected item (opens editable detail view)
- Press `s` to search/filter items
- Press `a` to add a new item (opens a form)
- Press `1` to move selected item to TODO
- Press `2` to move selected item to IN PROGRESS
- Press `3` to move selected item to DONE
- Press `d` to delete the selected item
- Press `r` to reload data from disk
- Press `?` to toggle help
- Press `q` to quit

**Adding items in interactive mode:**
When you press `a`, a form appears where you can:
- Use `Tab` or `↓` to move to the next field
- Use `Shift+Tab` or `↑` to move to the previous field
- Press `Enter` on the last field to submit
- Press `Esc` to cancel

**Editing items in interactive mode:**
When you press `Enter` on an item, an editable detail view appears where you can:
- Use `Tab`, `↑`, or `↓` to navigate between fields
- Edit title, description, due date, and tags (you can type any character including 'q')
- Press `Esc` to save changes and return to the board

**Searching in interactive mode:**
When you press `s`, a search prompt appears where you can:
- Type your search query (searches in title, description, and tags)
- Press `Enter` to apply the filter
- Press `Esc` to cancel
- The board will show only matching items with the filter indicator in the title

### Update a backlog item

```bash
bl update <id> --status in-progress
bl update <id> --title "New title" --desc "New description"
bl update <id> --due "20-12-2025" --tags "new,tags"
```

**Options:**
- `--status`: Change status (todo, in-progress, done)
- `--title`: Update title
- `--desc`: Update description
- `--due`: Update due date
- `--tags`: Update tags

### Delete a backlog item

```bash
bl delete <id>
```

### Search for items

```bash
bl search "keyword"
```

Searches in title, description, and tags.

### Archive completed items

```bash
bl archive
```

Moves all items with "done" status to `~/backlog/archive.json`.

## Data Storage

All data is stored in JSON format in the `~/backlog` directory:
- `~/backlog/items.json` - Active backlog items
- `~/backlog/archive.json` - Archived completed items

## Examples

```bash
# Add a new task
bl add "Implement user authentication" --desc "Add login and signup" --due "15-12-2025" --tags "backend,security"

# View the board (static)
bl list

# View the board (interactive mode)
bl list -i

# Move task to in-progress (CLI)
bl update 176457 --status in-progress

# Or use interactive mode to move items with keyboard shortcuts!

# Search for tasks
bl search "authentication"

# Archive completed tasks
bl archive

# Delete a task
bl delete 176457
```

## ID Usage

You can use partial IDs when updating or deleting items. The application will match the first item whose ID starts with the provided prefix.

For example, if an item has ID `1764579489317886000`, you can use:
```bash
bl update 176457 --status done
```

## License

MIT

