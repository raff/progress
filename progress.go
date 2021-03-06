package progress

import (
	"fmt"
	"strings"
	"sync"

	ui "github.com/gizak/termui"
)

type Progress struct {
	border bool // objects have borders

	gs []*ui.Gauge     // list of progress bars
	ps []*ui.Paragraph // list of associated info

	header *ui.Paragraph // a "header" text box

	mtext    *ui.Paragraph // a "messages" text box
	messages []string      // list of messages to return
	sync.Mutex
}

type newOption func(p *Progress)

// Messages specify that we want a "messages" section with "m" lines
func Messages(m int) newOption {
	return func(p *Progress) {
		if p.border {
			m += 2
		}
		p.mtext = ui.NewParagraph("")
		p.mtext.Border = p.border
		p.mtext.Height = m
	}
}

// Header specify that we want a "header" section with "h" lines
func Header(h int) newOption {
	return func(p *Progress) {
		padding := 1
		if p.border {
			h += 2
			padding = 0
		}
		p.header = ui.NewParagraph("")
		p.header.Border = p.border
		p.header.PaddingBottom = padding
		p.header.Height = h
	}
}

//
// New creates a new Progress object that manages "n" progress item and a message area with "m" lines
//
func New(n int, border bool, options ...newOption) *Progress {
	p := &Progress{border: border, gs: make([]*ui.Gauge, n)}

	height, padding := 2, 1

	for _, opt := range options {
		opt(p) // process option
	}

	if p.border {
		height, padding = 3, 0
	} else {
		p.ps = make([]*ui.Paragraph, n)
	}

	if p.header != nil {
		ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, p.header)))
	}

	for i := 0; i < n; i++ {
		p.gs[i] = ui.NewGauge()
		p.gs[i].LabelAlign = ui.AlignLeft
		p.gs[i].Height = height
		p.gs[i].Border = p.border
		p.gs[i].Percent = 0
		p.gs[i].PaddingBottom = padding
		p.gs[i].BarColor = ui.ColorGreen

		if p.border {
			ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, p.gs[i])))
		} else {
			p.ps[i] = ui.NewParagraph("")
			p.ps[i].Height = 1
			p.ps[i].Border = false

			// build layout
			ui.Body.AddRows(ui.NewRow(
				ui.NewCol(4, 0, p.ps[i]),
				ui.NewCol(8, 0, p.gs[i])))
		}
	}

	if p.mtext != nil {
		ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, p.mtext)))
	}

	ui.Body.Align()
	ui.Render(ui.Body)

	return p
}

func PercInt(curr, max int) int {
	if curr == 0 || max == 0 {
		return 0
	}

	return curr * 100 / max
}

func PercInt64(curr, max int64) int {
	if curr == 0 || max == 0 {
		return 0
	}

	return int(curr * 100 / max)
}

func PercFloat(curr, max float64) int {
	if curr == 0.0 || max == 0.0 {
		return 0
	}

	return int(curr * 100.0 / max)
}

func (p *Progress) SetHeader(m string) {
	if p.header != nil {
		p.header.Text = m
	}
}

func (p *Progress) SetHeaderf(f string, v ...interface{}) {
	if p.header != nil {
		p.header.Text = fmt.Sprintf(f, v...)
	}
}

func (p *Progress) Set(item int, label string, value int) {
	if p.border {
		p.gs[item].BorderLabel = label
	} else {
		p.ps[item].Text = label
	}
	p.gs[item].Percent = value
	ui.Render(ui.Body)
}

func (p *Progress) SetColor(item int, c ui.Attribute) {
	p.gs[item].BarColor = c
	if p.border {
		p.gs[item].BorderFg = c
		p.gs[item].BorderLabelFg = c
	} else {
		p.ps[item].TextFgColor = c
	}
}

func (p *Progress) AddMessagef(f string, v ...interface{}) {
	p.AddMessage(fmt.Sprintf(f, v...))
}

func (p *Progress) AddMessage(m string) {
	p.Mutex.Lock()
	p.messages = append(p.messages, m)
	p.Mutex.Unlock()

	if p.mtext != nil {
		h := p.mtext.Height
		if p.border {
			h -= 2
		}

		start := len(p.messages) - p.mtext.Height
		if start < 0 {
			start = 0
		}

		p.Mutex.Lock()
		p.mtext.Text = strings.Join(p.messages[start:], "\n")
		p.Mutex.Unlock()
		ui.Render(ui.Body)
	}
}

func (p *Progress) Messages() []string {
	return p.messages
}

func Color(c string) ui.Attribute {
	return ui.StringToAttribute(c)
}
