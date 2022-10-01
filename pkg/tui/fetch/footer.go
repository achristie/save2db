package fetch

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type footerModel struct {
	status    Status
	statusMsg string
}

type Status int

const (
	INPROGRESS Status = iota
	COMPLETED
	ERROR
)

type StatusMsg struct {
	status Status
	msg    string
}

func newFooter() footerModel {
	return footerModel{}
}

func (m footerModel) Init() tea.Cmd {
	return nil
}

func (m footerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatusMsg:
		m.statusMsg = msg.msg
		m.status = msg.status
		return m, nil
	}
	return m, nil
}

func (m footerModel) View() string {
	s := fmt.Sprint(m.status) + " " + m.statusMsg
	return s
}

func (s Status) String() string {
	switch s {
	case ERROR:
		return "error"
	case INPROGRESS:
		return "inprogress"
	case COMPLETED:
		return "complete"
	}
	return ""
}
