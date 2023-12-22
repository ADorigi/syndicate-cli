package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/adorigi/systeminfo"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Tabs       []string
	TabContent []table.Model
	activeTab  int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#EC9F05", Dark: "#EC9F05"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m model) View() string {
	doc := strings.Builder{}

	// doc.WriteString()

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab].View()))

	return docStyle.Render(doc.String())
}

func main() {

	info := systeminfo.CollectStats()
	tabs := []string{"System Information", "Disk Information", "Network Information"}
	// tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}
	columns := []table.Column{
		{Title: "Property", Width: 29},
		{Title: "Value", Width: 29},
	}

	rowsSys := []table.Row{
		{"CPU Name", info.CpuName},
		{"CPU Architecture", info.CpuArch},
		{"Operating System", info.OperatingSystem},
	}

	tableSys := table.New(
		table.WithColumns(columns),
		table.WithRows(rowsSys),
		table.WithFocused(false),
		table.WithHeight(2),
	)

	rowsDisk := []table.Row{
		{"Available Disk Storage: ", fmt.Sprint(info.DiskAvailable) + info.StorageUnit},
		{"Used Disk Storage: ", fmt.Sprint(info.DiskUsed) + info.StorageUnit},
		{"Disk Used %: ", strconv.FormatFloat(info.DiskUsedPercent, 'g', -1, 32) + "%"},
	}

	tableDisk := table.New(
		table.WithColumns(columns),
		table.WithRows(rowsDisk),
		table.WithFocused(false),
		table.WithHeight(3),
	)

	rowsNetwork := []table.Row{
		{"Local IPv4: ", info.LocalIPv4},
		{"Global IP: ", info.GlobalIP},
	}

	tableNetwork := table.New(
		table.WithColumns(columns),
		table.WithRows(rowsNetwork),
		table.WithFocused(false),
		table.WithHeight(2),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(highlightColor).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#EC9F05")).
		Background(lipgloss.Color("#EC9F05")).
		Bold(true)
	tableSys.SetStyles(s)
	tableDisk.SetStyles(s)
	tableNetwork.SetStyles(s)

	m := model{Tabs: tabs, TabContent: []table.Model{tableSys, tableDisk, tableNetwork}}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
