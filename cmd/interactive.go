package cmd

import (
	"backlog/models"
	"backlog/storage"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type reloadMsg struct {
	backlog *models.Backlog
	err     error
}

type addItemMsg struct {
	item models.BacklogItem
	err  error
}

type updateItemMsg struct {
	item models.BacklogItem
	err  error
}

type moveItemMsg struct {
	item models.BacklogItem
	err  error
}

type deleteItemMsg struct {
	itemTitle string
	err       error
}

type model struct {
	backlog        *models.Backlog
	storage        *storage.Storage
	cursor         int
	selectedCol    int // 0=todo, 1=in-progress, 2=done
	items          [3][]models.BacklogItem
	err            error
	message        string
	showHelp       bool
	addMode        bool
	inputs         []textinput.Model
	focusIndex     int
	viewMode       bool
	viewingItem    *models.BacklogItem
	editInputs     []textinput.Model
	editFocus      int
	searchMode     bool
	searchInput    textinput.Model
	searchQuery    string
	terminalWidth  int
	terminalHeight int
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			Width(35)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	statusRibbonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#1E1E1E")).
				Padding(0, 1)

	statusTabActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Bold(true).
				Padding(0, 2)

	statusTabInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 2)
)

const (
	tagIcon = "\U0001F3F7 " // label/tag
	dueIcon = "\U000023F0 " // alarm clock
)

func initialModel(backlog *models.Backlog, store *storage.Storage) model {
	m := model{
		backlog: backlog,
		storage: store,
		cursor:  0,
		// Start in the IN PROGRESS column by default for more relevant view
		selectedCol:    1,
		showHelp:       true,
		addMode:        false,
		inputs:         make([]textinput.Model, 4),
		editInputs:     make([]textinput.Model, 4),
		terminalWidth:  120, // Default, will be updated by Init
		terminalHeight: 30,  // Default, will be updated by Init
	}
	m.organizeItems()
	m.initInputs()
	m.initEditInputs()
	m.initSearchInput()
	return m
}

func (m *model) initInputs() {
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
		t.CharLimit = 200

		switch i {
		case 0:
			t.Placeholder = "Title (required)"
			t.Focus()
			t.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
			t.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
		case 1:
			t.Placeholder = "Description"
			t.CharLimit = 500
		case 2:
			t.Placeholder = "Due date (DD-MM-YYYY)"
			t.CharLimit = 10
		case 3:
			t.Placeholder = "Tags (comma-separated)"
		}

		m.inputs[i] = t
	}
}

func (m *model) initEditInputs() {
	var t textinput.Model
	for i := range m.editInputs {
		t = textinput.New()
		t.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
		t.CharLimit = 200

		switch i {
		case 0:
			t.Placeholder = "Title"
			t.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
			t.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
		case 1:
			t.Placeholder = "Description"
			t.CharLimit = 500
		case 2:
			t.Placeholder = "Due date (DD-MM-YYYY)"
			t.CharLimit = 10
		case 3:
			t.Placeholder = "Tags (comma-separated)"
		}

		m.editInputs[i] = t
	}
}

func (m *model) initSearchInput() {
	m.searchInput = textinput.New()
	m.searchInput.Placeholder = "Search items..."
	m.searchInput.CharLimit = 100
	m.searchInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	m.searchInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	m.searchInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
}

func (m *model) organizeItems() {
	m.items[0] = []models.BacklogItem{} // todo
	m.items[1] = []models.BacklogItem{} // in-progress
	m.items[2] = []models.BacklogItem{} // done

	for _, item := range m.backlog.Items {
		// Filter by search query if active
		if m.searchQuery != "" {
			if !m.matchesSearch(item) {
				continue
			}
		}

		switch item.Status {
		case models.StatusTodo:
			m.items[0] = append(m.items[0], item)
		case models.StatusInProgress:
			m.items[1] = append(m.items[1], item)
		case models.StatusDone:
			m.items[2] = append(m.items[2], item)
		}
	}
}

func (m model) matchesSearch(item models.BacklogItem) bool {
	query := m.searchQuery
	// Search in title, description, and tags
	if strings.Contains(strings.ToLower(item.Title), query) {
		return true
	}
	if strings.Contains(strings.ToLower(item.Description), query) {
		return true
	}
	for _, tag := range item.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func (m model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle add mode separately for key messages only. Non-key
	// messages (like addItemMsg) should still be processed by the
	// general handler below so that the list is refreshed after a
	// successful add.
	if m.addMode {
		if _, ok := msg.(tea.KeyMsg); ok {
			return m.updateAddMode(msg)
		}
	}

	// Handle view mode separately for key messages only. Non-key
	// messages (like updateItemMsg) should still be processed by the
	// general handler below so that we can exit view mode after
	// saving.
	if m.viewMode {
		if _, ok := msg.(tea.KeyMsg); ok {
			return m.updateViewMode(msg)
		}
	}

	// Handle search mode separately
	if m.searchMode {
		return m.updateSearchMode(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		m.message = ""
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "?":
			m.showHelp = !m.showHelp

		case "s":
			m.searchMode = true
			m.searchInput.Focus()
			return m, nil

		case "a":
			m.addMode = true
			m.focusIndex = 0
			m.inputs[0].Focus()
			return m, nil

		case "enter":
			// Show detail view for selected item
			if len(m.items[m.selectedCol]) > 0 && m.cursor < len(m.items[m.selectedCol]) {
				m.viewMode = true
				m.viewingItem = &m.items[m.selectedCol][m.cursor]
				// Populate edit inputs with current values
				m.editInputs[0].SetValue(m.viewingItem.Title)
				m.editInputs[1].SetValue(m.viewingItem.Description)
				m.editInputs[2].SetValue(m.viewingItem.DueDate)
				m.editInputs[3].SetValue(strings.Join(m.viewingItem.Tags, ", "))
				m.editFocus = 0
				m.editInputs[0].Focus()
			}
			return m, nil

		case "t":
			// View TODO items in the main panel
			m.selectedCol = 0
			m.cursor = 0
			return m, nil

		case "i":
			// View IN PROGRESS items in the main panel
			m.selectedCol = 1
			m.cursor = 0
			return m, nil

		case "c":
			// View DONE (completed) items in the main panel
			m.selectedCol = 2
			m.cursor = 0
			return m, nil

		case "tab":
			// Cycle status view forward via the bottom ribbon tabs
			m.selectedCol = (m.selectedCol + 1) % 3
			m.cursor = 0
			return m, nil

		case "shift+tab":
			// Cycle status view backward
			m.selectedCol = (m.selectedCol + 2) % 3
			m.cursor = 0
			return m, nil

		case "left":
			if m.selectedCol > 0 {
				m.selectedCol--
				m.cursor = 0
			}

		case "right":
			if m.selectedCol < 2 {
				m.selectedCol++
				m.cursor = 0
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			maxItems := len(m.items[m.selectedCol])
			if m.cursor < maxItems-1 {
				m.cursor++
			}

		case "1":
			return m, m.moveItemToStatus(models.StatusTodo)

		case "2":
			return m, m.moveItemToStatus(models.StatusInProgress)

		case "3":
			return m, m.moveItemToStatus(models.StatusDone)

		case "d":
			return m, m.deleteCurrentItem()

		case "r":
			return m, m.reloadData()
		}

	case reloadMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.backlog = msg.backlog
			m.organizeItems()
			m.cursor = 0
			m.message = "Reloaded data"
		}

	case addItemMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			// Backlog was already updated and saved in submitNewItem; just
			// re-organize the in-memory view and exit add mode.
			m.organizeItems()
			m.message = fmt.Sprintf("Added '%s'", msg.item.Title)
			// Exit add mode and clear the form now that the item is saved
			m.addMode = false
			for i := range m.inputs {
				m.inputs[i].SetValue("")
			}
			m.focusIndex = 0
		}

	case updateItemMsg:
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("ERROR: %v", msg.err)
		} else {
			// Find and update the item in backlog
			for i := range m.backlog.Items {
				if m.backlog.Items[i].ID == msg.item.ID {
					m.backlog.Items[i] = msg.item
					break
				}
			}
			m.organizeItems()
			m.message = fmt.Sprintf("Updated '%s'", msg.item.Title)
			// Exit view mode after successful update
			m.viewMode = false
			m.viewingItem = nil
		}

	case moveItemMsg:
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("ERROR: %v", msg.err)
		} else {
			m.organizeItems()
			m.cursor = 0
			if msg.item.Title != "" {
				m.message = fmt.Sprintf("Moved '%s' to %s", msg.item.Title, string(msg.item.Status))
			} else {
				m.message = "Moved item"
			}
		}

	case deleteItemMsg:
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("ERROR: %v", msg.err)
		} else {
			m.organizeItems()
			if m.cursor >= len(m.items[m.selectedCol]) && m.cursor > 0 {
				m.cursor--
			}
			if msg.itemTitle != "" {
				m.message = fmt.Sprintf("Deleted '%s'", msg.itemTitle)
			} else {
				m.message = "Deleted item"
			}
		}
	}

	return m, nil
}

func (m model) updateAddMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.addMode = false
			// Clear inputs
			for i := range m.inputs {
				m.inputs[i].SetValue("")
			}
			m.inputs[0].Focus()
			return m, nil

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Submit on enter from last field
			if s == "enter" && m.focusIndex == len(m.inputs)-1 {
				return m, m.submitNewItem()
			}

			// Cycle through inputs
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
					m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = lipgloss.NewStyle()
					m.inputs[i].TextStyle = lipgloss.NewStyle()
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
	}

	// Show add form if in add mode
	if m.addMode {
		return m.renderAddForm()
	}

	// Show detail view if in view mode
	if m.viewMode && m.viewingItem != nil {
		return m.renderDetailView()
	}

	// Show search mode
	if m.searchMode {
		return m.renderSearchMode()
	}

	var s strings.Builder
	var headerBuilder strings.Builder

	// Title
	title := "BACKLOG KANBAN BOARD"
	if m.searchQuery != "" {
		title += fmt.Sprintf(" (filtered: '%s')", m.searchQuery)
	}
	headerBuilder.WriteString(titleStyle.Render(title) + "\n\n")

	header := headerBuilder.String()
	s.WriteString(header)

	// Single main panel taking up the full terminal width.
	// Subtract a couple of columns so the right border stays inside the
	// visible area even on tight terminals.
	panelWidth := m.terminalWidth - 2
	if panelWidth < 10 {
		panelWidth = m.terminalWidth
	}

	// Determine the title based on the currently selected column
	currentTitle := "TODO"
	switch m.selectedCol {
	case 0:
		currentTitle = "TODO"
	case 1:
		currentTitle = "IN PROGRESS"
	case 2:
		currentTitle = "DONE"
	}

	// Choose a height for the main panel so that message + stats + help + ribbon sit near the bottom
	panelHeight := 0
	if m.terminalHeight > 0 {
		headerHeight := lipgloss.Height(header)
		statsHeight := 2 // "Total ..." + blank line
		// If we have a status message, reserve a line for it just above Total
		if m.message != "" {
			statsHeight++
		}
		helpLines := 1
		if m.showHelp {
			helpLines = 2 // views + actions
		}
		ribbonHeight := 1 // status ribbon at the very bottom

		// Leave at least 1 line of margin at the bottom
		panelHeight = m.terminalHeight - headerHeight - statsHeight - helpLines - ribbonHeight - 1
		if panelHeight < 5 {
			panelHeight = 5
		}
	}

	// Render only the currently selected status in a single large box
	mainPanel := m.renderColumnWithSize(currentTitle, m.selectedCol, panelWidth, panelHeight)
	s.WriteString(mainPanel)
	s.WriteString("\n\n")

	// Status / info message just above the Total line
	if m.message != "" {
		s.WriteString(messageStyle.Render(m.message) + "\n")
	}

	// Stats
	total := len(m.backlog.Items)
	s.WriteString(fmt.Sprintf("Total: %d items (%d todo, %d in-progress, %d done)\n\n",
		total, len(m.items[0]), len(m.items[1]), len(m.items[2])))

	// Help legend split across two lines to avoid overflowing
	if m.showHelp {
		viewsNav := "Views: t=todo i=in-progress c=done (or tab/shift+tab/left/right) | Navigation: up/down items"
		actions := "Actions: Enter=edit s=search a=add 1=todo 2-in-progress 3=done d=delete r=reload | ?=help q=quit"
		s.WriteString(helpStyle.Render(viewsNav) + "\n")
		s.WriteString(helpStyle.Render(actions))
	} else {
		s.WriteString(helpStyle.Render("Press '?' for help"))
	}

	// Status ribbon at the very bottom indicating the current view
	s.WriteString("\n")
	s.WriteString(m.renderStatusRibbon())

	return s.String()
}

func (m model) renderColumn(title string, colIndex int) string {
	return m.renderColumnWithSize(title, colIndex, 40, 0)
}

func (m model) renderColumnWithWidth(title string, colIndex int, width int) string {
	return m.renderColumnWithSize(title, colIndex, width, 0)
}

func (m model) renderColumnWithSize(title string, colIndex int, width int, height int) string {
	var content strings.Builder

	// Items
	items := m.items[colIndex]
	if len(items) == 0 {
		content.WriteString(helpStyle.Render("(empty)"))
	} else {
		for i, item := range items {
			itemStr := m.formatItemWithWidth(item, width-4)
			if colIndex == m.selectedCol && i == m.cursor {
				itemStr = selectedStyle.Render("> " + itemStr)
			} else {
				itemStr = "  " + itemStr
			}
			content.WriteString(itemStr + "\n")
		}
	}

	style := columnStyle.Width(width)
	if height > 0 {
		style = style.Height(height)
	}
	if colIndex == m.selectedCol {
		style = style.BorderForeground(lipgloss.Color("#FF00FF"))
	}

	return style.Render(content.String())
}

func (m model) renderStatusRibbon() string {
	labels := []string{"Todo", "In-Progress", "Done"}
	var tabs []string

	for i, label := range labels {
		style := statusTabInactiveStyle
		if i == m.selectedCol {
			style = statusTabActiveStyle
		}
		tabs = append(tabs, style.Render(label))
	}

	ribbon := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	width := m.terminalWidth
	if width <= 0 {
		return ribbon
	}

	return statusRibbonStyle.Width(width).Align(lipgloss.Center).Render(ribbon)
}

func (m model) formatItem(item models.BacklogItem) string {
	return m.formatItemWithWidth(item, 36)
}

func (m model) formatItemWithWidth(item models.BacklogItem, width int) string {
	var parts []string

	// Calculate available width for title
	titleWidth := width - 20 // Reserve space for tags and date

	// Title
	title := item.Title
	if len(title) > titleWidth {
		title = title[:titleWidth-3] + "..."
	}
	parts = append(parts, title)

	// Tags
	if len(item.Tags) > 0 {
		tags := strings.Join(item.Tags, ",")
		maxTagWidth := 15
		if len(tags) > maxTagWidth {
			tags = tags[:maxTagWidth-3] + "..."
		}
		parts = append(parts, tagIcon+tags)
	}

	// Due date
	if item.DueDate != "" {
		parts = append(parts, dueIcon+item.DueDate)
	}

	return strings.Join(parts, " | ")
}

func (m model) renderSearchMode() string {
	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(60)

	var s strings.Builder

	s.WriteString(titleStyle.Render("SEARCH") + "\n\n")
	s.WriteString("Enter search query:\n\n")
	s.WriteString(m.searchInput.View() + "\n\n")
	s.WriteString(helpStyle.Render("Enter: search | Esc: cancel") + "\n")

	return searchStyle.Render(s.String())
}

func (m *model) moveItemToStatus(status models.Status) tea.Cmd {
	return func() tea.Msg {
		if len(m.items[m.selectedCol]) == 0 {
			return nil
		}

		item := m.items[m.selectedCol][m.cursor]

		// Update in backlog
		var updated models.BacklogItem
		for i := range m.backlog.Items {
			if m.backlog.Items[i].ID == item.ID {
				m.backlog.Items[i].Status = status
				m.backlog.Items[i].UpdatedAt = time.Now()
				updated = m.backlog.Items[i]
				break
			}
		}

		// Save
		if err := m.storage.Save(m.backlog); err != nil {
			return moveItemMsg{err: err}
		}

		// If we didn't find the item for some reason, still send a message to
		// trigger a re-organize and status update based on our best guess.
		if updated.ID == "" {
			updated = item
			updated.Status = status
		}

		return moveItemMsg{item: updated, err: nil}
	}
}

func (m *model) deleteCurrentItem() tea.Cmd {
	return func() tea.Msg {
		if len(m.items[m.selectedCol]) == 0 {
			return nil
		}

		item := m.items[m.selectedCol][m.cursor]

		// Remove from backlog
		newItems := make([]models.BacklogItem, 0, len(m.backlog.Items))
		for _, i := range m.backlog.Items {
			if i.ID != item.ID {
				newItems = append(newItems, i)
			}
		}
		m.backlog.Items = newItems

		// Save
		if err := m.storage.Save(m.backlog); err != nil {
			return deleteItemMsg{err: err}
		}

		return deleteItemMsg{itemTitle: item.Title, err: nil}
	}
}

func (m model) reloadData() tea.Cmd {
	return func() tea.Msg {
		backlog, err := m.storage.Load()
		return reloadMsg{
			backlog: backlog,
			err:     err,
		}
	}
}

func (m model) renderAddForm() string {
	var s strings.Builder

	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(80)

	s.WriteString(titleStyle.Render("ADD NEW ITEM") + "\n\n")

	labels := []string{"Title:", "Description:", "Due Date:", "Tags:"}
	for i := range m.inputs {
		s.WriteString(labels[i] + "\n")
		s.WriteString(m.inputs[i].View() + "\n\n")
	}

	s.WriteString(helpStyle.Render("Tab/Shift+Tab: navigate | Enter: submit | Esc: cancel") + "\n")

	return formStyle.Render(s.String())
}

func (m model) submitNewItem() tea.Cmd {
	return func() tea.Msg {
		// Validate title
		title := strings.TrimSpace(m.inputs[0].Value())
		if title == "" {
			return addItemMsg{err: fmt.Errorf("title is required")}
		}

		// Get other fields
		description := strings.TrimSpace(m.inputs[1].Value())
		dueDate := strings.TrimSpace(m.inputs[2].Value())
		tagsStr := strings.TrimSpace(m.inputs[3].Value())

		// Validate due date if provided
		if dueDate != "" && !isValidDateFormat(dueDate) {
			return addItemMsg{err: fmt.Errorf("invalid date format. Use DD-MM-YYYY")}
		}

		// Parse tags
		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		// Create new item
		item := models.BacklogItem{
			ID:          generateID(),
			Title:       title,
			Description: description,
			DueDate:     dueDate,
			Tags:        tags,
			Status:      models.StatusTodo,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Add to backlog
		m.backlog.Items = append(m.backlog.Items, item)

		// Save
		if err := m.storage.Save(m.backlog); err != nil {
			return addItemMsg{err: err}
		}

		// Clear inputs
		for i := range m.inputs {
			m.inputs[i].SetValue("")
		}
		m.inputs[0].Focus()
		m.addMode = false

		return addItemMsg{item: item, err: nil}
	}
}

func (m model) updateViewMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			// Save the edited item automatically on Esc and exit
			m.message = ""
			return m, m.saveEditedItem()

		case "q":
			// Treat 'q' like Esc in the detail view: save and exit
			m.message = ""
			return m, m.saveEditedItem()

		case "tab", "shift+tab", "up", "down":
			// Cycle through inputs
			if key == "up" || key == "shift+tab" {
				m.editFocus--
			} else {
				m.editFocus++
			}

			if m.editFocus > len(m.editInputs)-1 {
				m.editFocus = 0
			} else if m.editFocus < 0 {
				m.editFocus = len(m.editInputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.editInputs))
			for i := 0; i < len(m.editInputs); i++ {
				if i == m.editFocus {
					cmds[i] = m.editInputs[i].Focus()
					m.editInputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
					m.editInputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
				} else {
					m.editInputs[i].Blur()
					m.editInputs[i].PromptStyle = lipgloss.NewStyle()
					m.editInputs[i].TextStyle = lipgloss.NewStyle()
				}
			}

			return m, tea.Batch(cmds...)

		default:
			// Handle character input for the focused text input
			cmd := m.updateEditInputs(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *model) updateEditInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.editInputs))

	for i := range m.editInputs {
		m.editInputs[i], cmds[i] = m.editInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) saveEditedItem() tea.Cmd {
	return func() tea.Msg {
		if m.viewingItem == nil {
			return updateItemMsg{err: fmt.Errorf("no item to update")}
		}

		// Validate title
		title := strings.TrimSpace(m.editInputs[0].Value())
		if title == "" {
			return updateItemMsg{err: fmt.Errorf("title is required")}
		}

		// Get other fields
		description := strings.TrimSpace(m.editInputs[1].Value())
		dueDate := strings.TrimSpace(m.editInputs[2].Value())
		tagsStr := strings.TrimSpace(m.editInputs[3].Value())

		// Validate due date if provided
		if dueDate != "" && !isValidDateFormat(dueDate) {
			return updateItemMsg{err: fmt.Errorf("invalid date format. Use DD-MM-YYYY")}
		}

		// Parse tags
		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		// Update item
		updatedItem := *m.viewingItem
		updatedItem.Title = title
		updatedItem.Description = description
		updatedItem.DueDate = dueDate
		updatedItem.Tags = tags
		updatedItem.UpdatedAt = time.Now()

		// Update in backlog
		for i := range m.backlog.Items {
			if m.backlog.Items[i].ID == updatedItem.ID {
				m.backlog.Items[i] = updatedItem
				break
			}
		}

		// Save
		if err := m.storage.Save(m.backlog); err != nil {
			return updateItemMsg{err: err}
		}

		return updateItemMsg{item: updatedItem, err: nil}
	}
}

func (m model) updateSearchMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.searchMode = false
			m.searchQuery = ""
			m.searchInput.SetValue("")
			m.searchInput.Blur()
			m.organizeItems()
			return m, nil

		case "enter":
			m.searchQuery = strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
			m.searchMode = false
			m.searchInput.Blur()
			m.organizeItems()
			m.cursor = 0
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m model) renderDetailView() string {
	if m.viewingItem == nil {
		return ""
	}

	item := m.viewingItem

	detailStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(100)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	statusStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00"))

	var s strings.Builder

	s.WriteString(titleStyle.Render("EDIT ITEM DETAILS") + "\n\n")

	// Show error or message if any
	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("ERROR: %v", m.err)) + "\n\n")
	} else if m.message != "" {
		s.WriteString(messageStyle.Render(m.message) + "\n\n")
	}

	// ID (read-only)
	s.WriteString(labelStyle.Render("ID: "))
	s.WriteString(item.ID + "\n\n")

	// Status (read-only for now)
	s.WriteString(labelStyle.Render("Status: "))
	statusColor := "#00FF00"
	switch item.Status {
	case models.StatusTodo:
		statusColor = "#FFA500"
	case models.StatusInProgress:
		statusColor = "#00BFFF"
	case models.StatusDone:
		statusColor = "#00FF00"
	}
	statusStyle = statusStyle.Foreground(lipgloss.Color(statusColor))
	s.WriteString(statusStyle.Render(string(item.Status)) + "\n\n")

	// Editable fields
	labels := []string{"Title:", "Description:", "Due Date:", "Tags:"}
	for i := range m.editInputs {
		s.WriteString(labelStyle.Render(labels[i]) + "\n")
		s.WriteString(m.editInputs[i].View() + "\n\n")
	}

	// Timestamps
	s.WriteString(labelStyle.Render("Created: "))
	s.WriteString(item.CreatedAt.Format("2006-01-02 15:04:05") + "\n")
	s.WriteString(labelStyle.Render("Updated: "))
	s.WriteString(item.UpdatedAt.Format("2006-01-02 15:04:05") + "\n\n")

	s.WriteString(helpStyle.Render("Tab/up/down: navigate | Esc/q: save and exit") + "\n")

	return detailStyle.Render(s.String())
}
