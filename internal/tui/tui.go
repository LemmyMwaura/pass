package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lemmyMwaura/pass/internal/account"
	pwgen "github.com/lemmyMwaura/pass/internal/password"
	"github.com/lemmyMwaura/pass/internal/vault"
)

type screen int

const (
	screenWelcome screen = iota
	screenAuth
	screenVault
	screenAdd
	screenDetail
	screenGenerate
	screenConfirmDelete
)

type authMode int

const (
	authLogin authMode = iota
	authCreate
)

type menuItem struct {
	title, desc string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type entryItem struct {
	entry vault.Entry
}

func (i entryItem) Title() string { return i.entry.Service }
func (i entryItem) Description() string {
	return i.entry.Username
}
func (i entryItem) FilterValue() string { return i.entry.Service }

type model struct {
	screen   screen
	authMode authMode

	welcome list.Model
	entries list.Model

	authInputs []textinput.Model
	authFocus  int

	addInputs []textinput.Model
	addFocus  int

	vault     *vault.Vault
	selected  *vault.Entry
	status    string
	errMsg    string
	generated string
	quitting  bool
	width     int
	height    int
}

// Run starts the Bubble Tea password manager UI.
func Run() error {
	m := newModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newModel() model {
	welcomeItems := []list.Item{
		menuItem{"Login", "Unlock an existing vault"},
		menuItem{"Create account", "Create a new encrypted vault"},
		menuItem{"Quit", "Exit the password manager"},
	}
	welcome := list.New(welcomeItems, list.NewDefaultDelegate(), 0, 0)
	welcome.Title = "pass"
	welcome.SetShowStatusBar(false)
	welcome.SetFilteringEnabled(false)
	welcome.Styles.Title = titleStyle

	entries := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	entries.Title = "Vault"
	entries.SetShowStatusBar(true)
	entries.SetFilteringEnabled(true)
	entries.Styles.Title = titleStyle

	return model{
		screen:     screenWelcome,
		welcome:    welcome,
		entries:    entries,
		authInputs: newAuthInputs(authLogin),
		addInputs:  newAddInputs(),
	}
}

func newAuthInputs(mode authMode) []textinput.Model {
	username := textinput.New()
	username.Placeholder = "username"
	username.CharLimit = 64
	username.Width = 40
	username.Focus()

	password := textinput.New()
	password.Placeholder = "master password"
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.CharLimit = 128
	password.Width = 40

	inputs := []textinput.Model{username, password}
	if mode == authCreate {
		confirm := textinput.New()
		confirm.Placeholder = "confirm master password"
		confirm.EchoMode = textinput.EchoPassword
		confirm.EchoCharacter = '•'
		confirm.CharLimit = 128
		confirm.Width = 40
		inputs = append(inputs, confirm)
	}
	return inputs
}

func newAddInputs() []textinput.Model {
	mk := func(placeholder string, password bool) textinput.Model {
		ti := textinput.New()
		ti.Placeholder = placeholder
		ti.CharLimit = 256
		ti.Width = 40
		if password {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		}
		return ti
	}

	service := mk("service (e.g. github)", false)
	service.Focus()
	username := mk("username / email", false)
	password := mk("password (empty = generate)", true)
	notes := mk("notes (optional)", false)
	return []textinput.Model{service, username, password, notes}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := 2, 4
		m.welcome.SetSize(msg.Width-h, msg.Height-v)
		m.entries.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			m.lockVault()
			return m, tea.Quit
		}
	}

	switch m.screen {
	case screenWelcome:
		return m.updateWelcome(msg)
	case screenAuth:
		return m.updateAuth(msg)
	case screenVault:
		return m.updateVault(msg)
	case screenAdd:
		return m.updateAdd(msg)
	case screenDetail:
		return m.updateDetail(msg)
	case screenGenerate:
		return m.updateGenerate(msg)
	case screenConfirmDelete:
		return m.updateConfirmDelete(msg)
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var body string
	switch m.screen {
	case screenWelcome:
		body = m.welcome.View()
	case screenAuth:
		body = m.viewAuth()
	case screenVault:
		body = m.viewVault()
	case screenAdd:
		body = m.viewAdd()
	case screenDetail:
		body = m.viewDetail()
	case screenGenerate:
		body = m.viewGenerate()
	case screenConfirmDelete:
		body = m.viewConfirmDelete()
	}

	footer := helpStyle.Render(m.helpText())
	status := ""
	if m.errMsg != "" {
		status = errorStyle.Render(m.errMsg) + "\n"
	} else if m.status != "" {
		status = successStyle.Render(m.status) + "\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left, status+body, "", footer)
}

func (m model) helpText() string {
	switch m.screen {
	case screenWelcome:
		return "↑/↓ navigate • enter select • q quit"
	case screenAuth:
		return "tab next field • enter submit • esc back"
	case screenVault:
		return "↑/↓ navigate • / filter • enter view • a add • d delete • g generate • q lock & quit"
	case screenAdd:
		return "tab next field • enter save • esc cancel"
	case screenDetail:
		return "esc back • d delete"
	case screenGenerate:
		return "n new password • esc back"
	case screenConfirmDelete:
		return "y confirm • n/esc cancel"
	default:
		return ""
	}
}

func (m *model) lockVault() {
	if m.vault != nil {
		m.vault.Lock()
		m.vault = nil
	}
}

func (m *model) refreshEntries() {
	items := make([]list.Item, 0)
	if m.vault != nil {
		for _, e := range m.vault.List() {
			items = append(items, entryItem{entry: e})
		}
	}
	m.entries.SetItems(items)
}

func (m model) updateWelcome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			item, ok := m.welcome.SelectedItem().(menuItem)
			if !ok {
				return m, nil
			}
			m.errMsg = ""
			m.status = ""
			switch item.title {
			case "Login":
				m.authMode = authLogin
				m.authInputs = newAuthInputs(authLogin)
				m.authFocus = 0
				m.screen = screenAuth
				return m, textinput.Blink
			case "Create account":
				m.authMode = authCreate
				m.authInputs = newAuthInputs(authCreate)
				m.authFocus = 0
				m.screen = screenAuth
				return m, textinput.Blink
			case "Quit":
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.welcome, cmd = m.welcome.Update(msg)
	return m, cmd
}

func (m model) updateAuth(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.errMsg = ""
			m.screen = screenWelcome
			return m, nil
		case "tab", "shift+tab", "up", "down":
			m.authFocus = cycleFocus(m.authFocus, len(m.authInputs), msg.String() == "shift+tab" || msg.String() == "up")
			cmds := make([]tea.Cmd, len(m.authInputs))
			for i := range m.authInputs {
				if i == m.authFocus {
					cmds[i] = m.authInputs[i].Focus()
				} else {
					m.authInputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		case "enter":
			return m.submitAuth()
		}
	}

	cmds := make([]tea.Cmd, len(m.authInputs))
	for i := range m.authInputs {
		m.authInputs[i], cmds[i] = m.authInputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) submitAuth() (tea.Model, tea.Cmd) {
	username := strings.TrimSpace(m.authInputs[0].Value())
	password := m.authInputs[1].Value()

	var (
		v   *vault.Vault
		err error
	)

	if m.authMode == authCreate {
		confirm := ""
		if len(m.authInputs) > 2 {
			confirm = m.authInputs[2].Value()
		}
		v, err = account.Create(username, password, confirm)
	} else {
		v, err = account.Login(username, password)
	}

	if err != nil {
		m.errMsg = err.Error()
		return m, nil
	}

	m.vault = v
	m.errMsg = ""
	if m.authMode == authCreate {
		m.status = fmt.Sprintf("Account %q created", username)
	} else {
		m.status = fmt.Sprintf("Unlocked vault for %q", username)
	}
	m.refreshEntries()
	m.screen = screenVault
	return m, nil
}

func (m model) viewAuth() string {
	title := "Login"
	if m.authMode == authCreate {
		title = "Create account"
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	labels := []string{"Username", "Master password"}
	if m.authMode == authCreate {
		labels = append(labels, "Confirm password")
	}

	for i, input := range m.authInputs {
		b.WriteString(labelStyle.Render(labels[i]))
		b.WriteString("\n")
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}
	return boxStyle.Render(b.String())
}

func (m model) updateVault(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.entries.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.entries, cmd = m.entries.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.lockVault()
			m.quitting = true
			return m, tea.Quit
		case "a":
			m.status = ""
			m.errMsg = ""
			m.addInputs = newAddInputs()
			m.addFocus = 0
			m.screen = screenAdd
			return m, textinput.Blink
		case "g":
			m.generated = pwgen.Generate(16, true, true)
			m.status = ""
			m.errMsg = ""
			m.screen = screenGenerate
			return m, nil
		case "d":
			item, ok := m.entries.SelectedItem().(entryItem)
			if !ok {
				m.errMsg = "no entry selected"
				return m, nil
			}
			e := item.entry
			m.selected = &e
			m.errMsg = ""
			m.status = ""
			m.screen = screenConfirmDelete
			return m, nil
		case "enter":
			item, ok := m.entries.SelectedItem().(entryItem)
			if !ok {
				return m, nil
			}
			e := item.entry
			m.selected = &e
			m.errMsg = ""
			m.status = ""
			m.screen = screenDetail
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.entries, cmd = m.entries.Update(msg)
	return m, cmd
}

func (m model) viewVault() string {
	header := subtitleStyle.Render(fmt.Sprintf("user: %s  •  %d entries", m.vault.Username, len(m.vault.List())))
	return header + "\n" + m.entries.View()
}

func (m model) updateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.errMsg = ""
			m.screen = screenVault
			return m, nil
		case "tab", "shift+tab", "up", "down":
			m.addFocus = cycleFocus(m.addFocus, len(m.addInputs), msg.String() == "shift+tab" || msg.String() == "up")
			cmds := make([]tea.Cmd, len(m.addInputs))
			for i := range m.addInputs {
				if i == m.addFocus {
					cmds[i] = m.addInputs[i].Focus()
				} else {
					m.addInputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		case "enter":
			service := strings.TrimSpace(m.addInputs[0].Value())
			username := strings.TrimSpace(m.addInputs[1].Value())
			password := m.addInputs[2].Value()
			notes := strings.TrimSpace(m.addInputs[3].Value())
			if password == "" {
				password = pwgen.Generate(16, true, true)
			}
			if _, err := m.vault.Add(service, username, password, notes); err != nil {
				m.errMsg = err.Error()
				return m, nil
			}
			m.status = fmt.Sprintf("Saved %q", service)
			m.errMsg = ""
			m.refreshEntries()
			m.screen = screenVault
			return m, nil
		}
	}

	cmds := make([]tea.Cmd, len(m.addInputs))
	for i := range m.addInputs {
		m.addInputs[i], cmds[i] = m.addInputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) viewAdd() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add entry"))
	b.WriteString("\n\n")
	labels := []string{"Service", "Username", "Password", "Notes"}
	for i, input := range m.addInputs {
		b.WriteString(labelStyle.Render(labels[i]))
		b.WriteString("\n")
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}
	return boxStyle.Render(b.String())
}

func (m model) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.selected = nil
			m.screen = screenVault
			return m, nil
		case "d":
			m.screen = screenConfirmDelete
			return m, nil
		}
	}
	return m, nil
}

func (m model) viewDetail() string {
	if m.selected == nil {
		return "No entry selected"
	}
	e := m.selected
	var b strings.Builder
	b.WriteString(titleStyle.Render(e.Service))
	b.WriteString("\n\n")
	b.WriteString(labelStyle.Render("Username") + "\n" + e.Username + "\n\n")
	b.WriteString(labelStyle.Render("Password") + "\n" + e.Password + "\n\n")
	if e.Notes != "" {
		b.WriteString(labelStyle.Render("Notes") + "\n" + e.Notes + "\n")
	}
	return boxStyle.Render(b.String())
}

func (m model) updateGenerate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.screen = screenVault
			return m, nil
		case "n":
			m.generated = pwgen.Generate(16, true, true)
			return m, nil
		}
	}
	return m, nil
}

func (m model) viewGenerate() string {
	return boxStyle.Render(
		titleStyle.Render("Generated password") + "\n\n" + m.generated + "\n",
	)
}

func (m model) updateConfirmDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			if m.selected == nil {
				m.screen = screenVault
				return m, nil
			}
			if err := m.vault.Delete(m.selected.Service); err != nil {
				m.errMsg = err.Error()
				m.screen = screenVault
				return m, nil
			}
			m.status = fmt.Sprintf("Deleted %q", m.selected.Service)
			m.selected = nil
			m.errMsg = ""
			m.refreshEntries()
			m.screen = screenVault
			return m, nil
		case "n", "esc":
			m.screen = screenVault
			if m.selected != nil {
				// came from detail? stay simple: always back to vault list
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) viewConfirmDelete() string {
	name := ""
	if m.selected != nil {
		name = m.selected.Service
	}
	return boxStyle.Render(
		titleStyle.Render("Delete entry") + "\n\n" +
			fmt.Sprintf("Delete %q? This cannot be undone.", name) + "\n",
	)
}

func cycleFocus(current, total int, backward bool) int {
	if total == 0 {
		return 0
	}
	if backward {
		current--
		if current < 0 {
			current = total - 1
		}
		return current
	}
	return (current + 1) % total
}
