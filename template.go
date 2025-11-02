package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62"))

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("235"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// Tab constants
const (
	TabMain = iota
	TabTwo
	TabThree
)

// Focus constants for main tab
const (
	FocusLeft = iota
	FocusRight
)

// Left panel focus areas
const (
	FocusMethodList = iota
	FocusURL
	FocusParams
	FocusAuth
	FocusHeaders
	FocusBody
)

// Right panel focus areas
const (
	FocusResponse = iota
	FocusHeaders2
	FocusRaw
)

type Model struct {
	// Navigation
	currentTab int
	leftFocus  int
	rightFocus int
	panelFocus int // 0 = left, 1 = right

	// Tab 1 - Method list (4 options)
	methodList     []string
	methodSelected int

	// Tab 1 - Text inputs
	urlInput    textinput.Model
	paramsInput textinput.Model

	// Tab 1 - Auth section (with 3 option list)
	authInput    textinput.Model
	authList     []string
	authSelected int

	// Tab 1 - Headers and Body
	headersInput textinput.Model
	bodyInput    textinput.Model

	// Tab 1 - Body type section (with 8 option list)
	bodyTypeList     []string
	bodyTypeSelected int

	// Tab 1 - Right side viewports
	statusDisplay1 string
	statusDisplay2 string
	viewport1      viewport.Model
	viewport2      viewport.Model
	viewport3      viewport.Model

	// Saved values
	savedURL         *string
	savedParams      *string
	savedAuth        *string
	savedHeaders     *string
	savedBody        *string
	savedMethodTypes [4]bool
	savedAuthTypes   [3]bool
	savedBodyTypes   [8]bool

	// Window dimensions
	width  int
	height int
	ready  bool
}

func initialTemplateModel() Model {
	// Method list (4 options)
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	// Auth list (3 options)
	authTypes := []string{"None", "Bearer", "Basic"}

	// Body type list (8 options)
	bodyTypes := []string{"JSON", "XML", "Form", "Text", "HTML", "Raw"}

	// Initialize text inputs with default width (will be adjusted by WindowSizeMsg)
	urlInput := textinput.New()
	urlInput.Placeholder = "Enter URL"
	urlInput.Width = 40

	paramsInput := textinput.New()
	paramsInput.Placeholder = "Query parameters"
	paramsInput.Width = 40

	authInput := textinput.New()
	authInput.Placeholder = "Auth token/credentials"
	authInput.Width = 40

	headersInput := textinput.New()
	headersInput.Placeholder = "Custom headers"
	headersInput.Width = 40

	bodyInput := textinput.New()
	bodyInput.Placeholder = "Request body"
	bodyInput.Width = 40

	return Model{
		currentTab:       TabMain,
		methodList:       methods,
		methodSelected:   0,
		authList:         authTypes,
		authSelected:     0,
		bodyTypeList:     bodyTypes,
		bodyTypeSelected: 0,
		urlInput:         urlInput,
		paramsInput:      paramsInput,
		authInput:        authInput,
		headersInput:     headersInput,
		bodyInput:        bodyInput,
		panelFocus:       FocusLeft,
		leftFocus:        -1,
		rightFocus:       -1,
		statusDisplay1:   "200",
		statusDisplay2:   "OK",
		width:            80,
		height:           24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate viewport dimensions based on window size
		vpWidth := (msg.Width / 2) - 10
		if vpWidth < 20 {
			vpWidth = 20
		}

		// Calculate viewport height based on available space
		// Account for: top border (1), tabs (2), help (2), padding (4)
		availableHeight := msg.Height - 9
		vpHeight := (availableHeight - 6) / 3 // Divide among 3 viewports with titles
		if vpHeight < 5 {
			vpHeight = 5
		}

		// Update text input widths
		inputWidth := (msg.Width / 2) - 15
		if inputWidth < 20 {
			inputWidth = 20
		}
		m.urlInput.Width = inputWidth
		m.paramsInput.Width = inputWidth
		m.authInput.Width = inputWidth
		m.headersInput.Width = inputWidth
		m.bodyInput.Width = inputWidth

		if !m.ready {
			// Initialize viewports
			m.viewport1 = viewport.New(vpWidth, vpHeight)
			m.viewport1.SetContent("Response body will appear here")

			m.viewport2 = viewport.New(vpWidth, vpHeight)
			m.viewport2.SetContent("Headers will appear here")

			m.viewport3 = viewport.New(vpWidth, vpHeight)
			m.viewport3.SetContent("Raw response will appear here")

			m.ready = true
		} else {
			// Update viewport dimensions
			m.viewport1.Width = vpWidth
			m.viewport1.Height = vpHeight
			m.viewport2.Width = vpWidth
			m.viewport2.Height = vpHeight
			m.viewport3.Width = vpWidth
			m.viewport3.Height = vpHeight
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// Tab navigation (carousel)
		case "ctrl+n":
			m.currentTab = (m.currentTab + 1) % 3
			return m, nil

		case "ctrl+p":
			m.currentTab = (m.currentTab - 1 + 3) % 3
			return m, nil

		// Panel switching (left/right)
		case "ctrl+r":
			if m.currentTab == TabMain {
				m.panelFocus = (m.panelFocus + 1) % 2
			}
			return m, nil

		// Tab key for right panel viewport navigation
		case "tab":
			if m.currentTab == TabMain && m.panelFocus == FocusRight {
				m.rightFocus = (m.rightFocus + 1) % 3
			}
			return m, nil

		// Left panel navigation
		case "ctrl+g":
			if m.currentTab == TabMain && m.panelFocus == FocusLeft {
				// Always focus on method list first
				if m.leftFocus != FocusMethodList {
					m.leftFocus = FocusMethodList
					m.blurAllInputs()
				} else {
					// If already focused, cycle through methods
					m.methodSelected = (m.methodSelected + 1) % len(m.methodList)
				}
			}
			return m, nil

		case "ctrl+u":
			if m.currentTab == TabMain && m.panelFocus == FocusLeft {
				m.leftFocus = FocusURL
				m.blurAllInputs()
				m.urlInput.Focus()
			}
			return m, nil

		case "ctrl+q":
			if m.currentTab == TabMain && m.panelFocus == FocusLeft {
				m.leftFocus = FocusParams
				m.blurAllInputs()
				m.paramsInput.Focus()
			}
			return m, nil

		case "ctrl+a":
			if m.currentTab == TabMain && m.panelFocus == FocusLeft {
				if m.leftFocus == FocusAuth {
					// Cycle through auth types
					m.authSelected = (m.authSelected + 1) % len(m.authList)
				} else {
					m.leftFocus = FocusAuth
					m.blurAllInputs()
					m.authInput.Focus()
				}
			}
			return m, nil

		case "ctrl+h":
			if m.currentTab == TabMain && m.panelFocus == FocusLeft {
				m.leftFocus = FocusHeaders
				m.blurAllInputs()
				m.headersInput.Focus()
			}
			return m, nil

		case "ctrl+b":
			if m.currentTab == TabMain && m.panelFocus == FocusLeft {
				if m.leftFocus == FocusBody {
					// Cycle through body types
					m.bodyTypeSelected = (m.bodyTypeSelected + 1) % len(m.bodyTypeList)
				} else {
					m.leftFocus = FocusBody
					m.blurAllInputs()
					m.bodyInput.Focus()
				}
			}
			return m, nil

		// Backspace to unfocus within current panel
		case "backspace":
			if m.currentTab == TabMain {
				if m.panelFocus == FocusLeft {
					m.blurAllInputs()
					// Stay on left panel but don't focus any specific item
					m.leftFocus = -1
				} else if m.panelFocus == FocusRight {
					// Unfocus viewport on right panel
					m.rightFocus = -1
				}
			}
			return m, nil

		// Method list navigation (h/l for carousel)
		case "h":
			if m.currentTab == TabMain && m.panelFocus == FocusLeft && m.leftFocus == FocusMethodList {
				m.methodSelected = (m.methodSelected - 1 + len(m.methodList)) % len(m.methodList)
			}
			return m, nil

		case "l":
			if m.currentTab == TabMain && m.panelFocus == FocusLeft && m.leftFocus == FocusMethodList {
				m.methodSelected = (m.methodSelected + 1) % len(m.methodList)
			}
			return m, nil

		// Save values
		case "ctrl+s":
			m.saveValues()
			return m, nil

		// Run command (placeholder)
		case "ctrl+ ":
			// TODO: Run command with saved values
			return m, m.runCommand()
		}

		// Update focused text input
		if m.currentTab == TabMain && m.panelFocus == FocusLeft {
			switch m.leftFocus {
			case FocusURL:
				m.urlInput, cmd = m.urlInput.Update(msg)
				cmds = append(cmds, cmd)
			case FocusParams:
				m.paramsInput, cmd = m.paramsInput.Update(msg)
				cmds = append(cmds, cmd)
			case FocusAuth:
				m.authInput, cmd = m.authInput.Update(msg)
				cmds = append(cmds, cmd)
			case FocusHeaders:
				m.headersInput, cmd = m.headersInput.Update(msg)
				cmds = append(cmds, cmd)
			case FocusBody:
				m.bodyInput, cmd = m.bodyInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

		// Update viewports when right panel is focused
		if m.currentTab == TabMain && m.panelFocus == FocusRight && m.rightFocus >= 0 {
			switch m.rightFocus {
			case FocusResponse:
				m.viewport1, cmd = m.viewport1.Update(msg)
				cmds = append(cmds, cmd)
			case FocusHeaders2:
				m.viewport2, cmd = m.viewport2.Update(msg)
				cmds = append(cmds, cmd)
			case FocusRaw:
				m.viewport3, cmd = m.viewport3.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) blurAllInputs() {
	m.urlInput.Blur()
	m.paramsInput.Blur()
	m.authInput.Blur()
	m.headersInput.Blur()
	m.bodyInput.Blur()
}

func (m *Model) saveValues() {
	m.savedURL = m.urlInput.Value()
	m.savedParams = m.paramsInput.Value()
	m.savedAuth = m.authInput.Value()
	m.savedHeaders = m.headersInput.Value()
	m.savedBody = m.bodyInput.Value()

	// Save auth type selections
	for i := range m.savedMethodTypes {
		m.savedMethodTypes[i] = (i == m.methodSelected)
	}
	for i := range m.savedAuthTypes {
		m.savedAuthTypes[i] = (i == m.authSelected)
	}

	// Save body type selections
	for i := range m.savedBodyTypes {
		m.savedBodyTypes[i] = (i == m.bodyTypeSelected)
	}
}

func (m Model) runCommand() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement command execution with saved values
		return nil
	}
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// Top border with "Burrow" on the right
	topBorder := m.renderTopBorder()

	// Tabs
	tabs := m.renderTabs()

	// Content based on current tab
	var content string
	switch m.currentTab {
	case TabMain:
		content = m.renderMainTab()
	case TabTwo:
		content = m.renderTabTwo()
	case TabThree:
		content = m.renderTabThree()
	}

	// Help text
	help := m.renderHelp()

	// Ensure content fits within window
	view := fmt.Sprintf("%s\n%s\n%s\n%s", topBorder, tabs, content, help)

	// Constrain to window dimensions
	style := lipgloss.NewStyle().
		MaxWidth(m.width).
		MaxHeight(m.height)

	return style.Render(view)
}

func (m Model) renderTopBorder() string {
	borderWidth := m.width - 4
	if borderWidth < 20 {
		borderWidth = 20
	}

	borderLine := strings.Repeat("─", borderWidth-8)
	burrowText := titleStyle.Render("Burrow")

	return fmt.Sprintf("┌%s%s┐", borderLine, burrowText)
}

func (m Model) renderTabs() string {
	tabs := []string{"Main", "Tab 2", "Tab 3"}
	var renderedTabs []string

	for i, tab := range tabs {
		if i == m.currentTab {
			renderedTabs = append(renderedTabs, selectedStyle.Render(fmt.Sprintf("[ %s ]", tab)))
		} else {
			renderedTabs = append(renderedTabs, fmt.Sprintf("  %s  ", tab))
		}
	}

	return strings.Join(renderedTabs, " ")
}

func (m Model) renderMainTab() string {
	leftPanel := m.renderLeftPanel()
	rightPanel := m.renderRightPanel()

	// Calculate panel widths dynamically
	panelWidth := (m.width / 2) - 4
	if panelWidth < 30 {
		panelWidth = 30
	}

	leftStyle := lipgloss.NewStyle().
		Width(panelWidth).
		Height(m.height-8).
		Padding(0, 1)

	rightStyle := lipgloss.NewStyle().
		Width(panelWidth).
		Height(m.height-8).
		Padding(0, 1)

	if m.panelFocus == FocusLeft {
		leftStyle = leftStyle.Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("205"))
	} else {
		rightStyle = rightStyle.Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("205"))
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftPanel),
		rightStyle.Render(rightPanel),
	)
}

func (m Model) renderLeftPanel() string {
	var sections []string

	// Method list (4 options)
	methodListStr := m.renderMethodList()
	sections = append(sections, methodListStr)

	// URL input
	urlSection := titleStyle.Render("URL:") + "\n" + m.renderInput(m.urlInput, m.leftFocus == FocusURL)
	sections = append(sections, urlSection)

	// Params input
	paramsSection := titleStyle.Render("Params:") + "\n" + m.renderInput(m.paramsInput, m.leftFocus == FocusParams)
	sections = append(sections, paramsSection)

	// Auth section with 3 option list
	authListStr := m.renderAuthList()
	authSection := titleStyle.Render("Auth:") + "\n" + authListStr + "\n" + m.renderInput(m.authInput, m.leftFocus == FocusAuth)
	sections = append(sections, authSection)

	// Headers input
	headersSection := titleStyle.Render("Headers:") + "\n" + m.renderInput(m.headersInput, m.leftFocus == FocusHeaders)
	sections = append(sections, headersSection)

	// Body section with 8 option list
	bodyTypeListStr := m.renderBodyTypeList()
	bodySection := titleStyle.Render("Body:") + "\n" + bodyTypeListStr + "\n" + m.renderInput(m.bodyInput, m.leftFocus == FocusBody)
	sections = append(sections, bodySection)

	return strings.Join(sections, "\n\n")
}

func (m Model) renderMethodList() string {
	var methods []string
	for i, method := range m.methodList {
		if i == m.methodSelected && m.leftFocus == FocusMethodList && m.panelFocus == FocusLeft {
			methods = append(methods, selectedStyle.Render(fmt.Sprintf("[%s]", method)))
		} else if i == m.methodSelected {
			methods = append(methods, focusedStyle.Render(method))
		} else {
			methods = append(methods, method)
		}
	}
	return strings.Join(methods, "  ")
}

func (m Model) renderAuthList() string {
	var authTypes []string
	for i, authType := range m.authList {
		if i == m.authSelected {
			authTypes = append(authTypes, selectedStyle.Render(fmt.Sprintf("[%s]", authType)))
		} else {
			authTypes = append(authTypes, authType)
		}
	}
	return strings.Join(authTypes, "  ")
}

func (m Model) renderBodyTypeList() string {
	var bodyTypes []string
	for i, bodyType := range m.bodyTypeList {
		if i == m.bodyTypeSelected {
			bodyTypes = append(bodyTypes, selectedStyle.Render(fmt.Sprintf("[%s]", bodyType)))
		} else {
			bodyTypes = append(bodyTypes, bodyType)
		}
	}
	// Split into two rows for better display
	half := len(bodyTypes) / 2
	row1 := strings.Join(bodyTypes[:half], "  ")
	row2 := strings.Join(bodyTypes[half:], "  ")
	return row1 + "\n" + row2
}

func (m Model) renderInput(input textinput.Model, focused bool) string {
	if focused {
		return focusedStyle.Render(input.View())
	}
	return input.View()
}

func (m Model) renderRightPanel() string {
	var sections []string

	// Title
	sections = append(sections, titleStyle.Render("Response"))

	// Two small status displays
	statusLine := fmt.Sprintf("Status: %s  %s", m.statusDisplay1, m.statusDisplay2)
	sections = append(sections, statusLine)

	// Viewport 1 - Response body
	vp1Title := titleStyle.Render("Body:")
	vp1Render := m.viewport1.View()
	if m.rightFocus == FocusResponse && m.panelFocus == FocusRight {
		vp1Render = focusedStyle.Render(vp1Render)
	}
	sections = append(sections, vp1Title+"\n"+vp1Render)

	// Viewport 2 - Headers
	vp2Title := titleStyle.Render("Headers:")
	vp2Render := m.viewport2.View()
	if m.rightFocus == FocusHeaders2 && m.panelFocus == FocusRight {
		vp2Render = focusedStyle.Render(vp2Render)
	}
	sections = append(sections, vp2Title+"\n"+vp2Render)

	// Viewport 3 - Raw
	vp3Title := titleStyle.Render("Raw:")
	vp3Render := m.viewport3.View()
	if m.rightFocus == FocusRaw && m.panelFocus == FocusRight {
		vp3Render = focusedStyle.Render(vp3Render)
	}
	sections = append(sections, vp3Title+"\n"+vp3Render)

	return strings.Join(sections, "\n\n")
}

func (m Model) renderTabTwo() string {
	contentHeight := m.height - 10
	if contentHeight < 5 {
		contentHeight = 5
	}

	content := titleStyle.Render("Tab 2 - Placeholder") + "\n\n" +
		helpStyle.Render("Content coming soon...")

	return lipgloss.NewStyle().
		Width(m.width - 4).
		Height(contentHeight).
		Render(content)
}

func (m Model) renderTabThree() string {
	contentHeight := m.height - 10
	if contentHeight < 5 {
		contentHeight = 5
	}

	content := titleStyle.Render("Tab 3 - Placeholder") + "\n\n" +
		helpStyle.Render("Content coming soon...")

	return lipgloss.NewStyle().
		Width(m.width - 4).
		Height(contentHeight).
		Render(content)
}

func (m Model) renderHelp() string {
	helpText := []string{
		"ctrl+n/p: next/prev tab",
		"ctrl+r: switch panel",
		"tab: switch viewport",
		"ctrl+g: method",
		"ctrl+u: url",
		"ctrl+q: params",
		"ctrl+a: auth",
		"ctrl+h: headers",
		"ctrl+b: body",
		"h/l: navigate lists",
		"backspace: unfocus",
		"ctrl+s: save",
		"ctrl+space: run",
		"q: quit",
	}
	return helpStyle.Render(strings.Join(helpText, " • "))
}

func main() {
	p := tea.NewProgram(initialTemplateModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
