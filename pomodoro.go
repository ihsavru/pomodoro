package pomodoro

import (
	"time"

	termui "github.com/gizak/termui/v3"
)

func (p *Pomodoro) drawUi() {
	p.ui.setGridRect()
	p.ui.drawHeading()
	p.ui.drawLoader(p.flags.workDuration)
	p.ui.setGridElements()
}

func (p *Pomodoro) startTimer() {
	isWork := true
	tickerCount, pomodoroCount := 0, 0

	uiEvents := termui.PollEvents()
	termui.Render(p.ui.grid)
	secondTicker := time.NewTicker(time.Second)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "r":
				tickerCount = 0
				pomodoroCount = 0
				secondTicker.Stop()
				p.ui.setLoaderData(workText, 0, formatLabel(p.ui.loader.Percent, p.flags.workDuration, 0))
				p.ui.renderGrid()
				p.notifier.play()
				secondTicker = time.NewTicker(time.Second)
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				p.ui.resizeGrid(payload.Width, payload.Height)
			}
		case <-secondTicker.C:
			tickerCount++
			if isWork {
				p.ui.loader.Percent = (tickerCount * 100) / p.flags.workDuration
				p.ui.loader.Label = formatLabel(p.ui.loader.Percent, p.flags.workDuration, tickerCount)

				if isWorkOver(tickerCount, p.flags.workDuration) {
					isWork = false
					tickerCount = 0
					pomodoroCount++
					p.ui.loader.Title = shortBreakText
					if isLongBreak(pomodoroCount) {
						p.ui.loader.Title = longBreakText
					}
					p.notifier.play()
				}
			} else {
				breakDurationSeconds := p.flags.shortBreakDuration
				if isLongBreak(pomodoroCount) {
					breakDurationSeconds = p.flags.longBreakDuration
				}
				p.ui.loader.Percent = (tickerCount * 100) / breakDurationSeconds
				p.ui.loader.Label = formatLabel(p.ui.loader.Percent, breakDurationSeconds, tickerCount)

				if isBreakOver(tickerCount, breakDurationSeconds) {
					isWork = true
					tickerCount = 0
					p.ui.loader.Title = workText
					p.notifier.play()
				}
			}
			p.ui.renderGrid()
		}
	}
}
