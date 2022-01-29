package pomodoro

import (
	termui "github.com/gizak/termui/v3"
)

func (ui *Ui) setGridRect() {
	termWidth, termHeight := termui.TerminalDimensions()
	ui.grid.SetRect(calculateDimensions(termWidth, termHeight))
}

func (ui *Ui) setGridElements() {
	ui.grid.Set(
		termui.NewRow(2.0/3,
			termui.NewCol(1, ui.heading),
		),
		termui.NewRow(1.0/3,
			termui.NewCol(1, ui.loader),
		),
	)
}

func (ui *Ui) resizeGrid(width, height int) {
	termWidth, termHeight := width, height
	ui.grid.SetRect(calculateDimensions(termWidth, termHeight))
	ui.clearGrid()
	ui.renderGrid()
}

func (ui *Ui) drawHeading() {
	ui.heading.Title = "POMODORO"
	ui.heading.Text = `[The](fg:blue) [Pomodoro Technique](fg:bold,bg:red) [is a time management method developed by Francesco Cirillo in the late 1980s. The technique uses a timer to break down work into intervals, traditionally 25 minutes in length, separated by short breaks (5 minutes). Each interval is known as a pomodoro, from the Italian word for 'tomato', after the tomato-shaped kitchen timer that Cirillo used as a university student.](fg:blue)
		Press "r" to reset the timer.
		Press "Ctrl + c" or "q" to exit.`
	ui.heading.TitleStyle.Fg = termui.ColorGreen
	ui.heading.BorderStyle.Fg = termui.ColorGreen
}

func (ui *Ui) drawLoader(workDurationSeconds int) {
	ui.setLoaderStyle()
	ui.setLoaderData(workText, 0, formatLabel(ui.loader.Percent, workDurationSeconds, 0))
	ui.loader.SetRect(0, 0, 25, 5)
}

func (ui *Ui) setLoaderStyle() {
	ui.loader.BarColor = termui.ColorYellow
	ui.loader.TitleStyle.Fg = termui.ColorMagenta
	ui.loader.BorderStyle.Fg = termui.ColorMagenta
	ui.loader.LabelStyle = termui.NewStyle(termui.ColorCyan)
}

func (ui *Ui) setLoaderData(title string, percent int, label string) {
	ui.loader.Title = title
	ui.loader.Percent = percent
	ui.loader.Label = label
}

func (ui *Ui) clearGrid() {
	termui.Clear()
}

func (ui *Ui) renderGrid() {
	termui.Render(ui.grid)
}
