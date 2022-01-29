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
	uiEvents := termui.PollEvents()
	p.ui.renderGrid()
	p.secondTicker = time.NewTicker(time.Second)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				// Quit the program
				return
			case "r":
				// Reset the Pomodoro
				p.resetStatus()
				p.secondTicker.Stop()
				p.ui.setLoaderData(workText, 0, formatLabel(p.ui.loader.Percent, p.flags.workDuration, 0))
				p.ui.renderGrid()
				p.notifier.play()
				p.secondTicker.Reset(time.Second)
			case "<Resize>":
				// Resize UI
				payload := e.Payload.(termui.Resize)
				p.ui.resizeGrid(payload.Width, payload.Height)
			}
		case <-p.secondTicker.C:
			// Updates the status after every second
			p.status.tickerCount++
			p.setStatus()
			p.updateUi()
			p.ui.renderGrid()
		}
	}
}

func (p *Pomodoro) setStatus() {
	switch {
	case p.status.isWork:
		if isWorkOver(p.status.tickerCount, p.flags.workDuration) {
			p.status.isWork = false
			p.status.pomodoroCount++
			p.endPomodoro()

			if isLongBreak(p.status.pomodoroCount) {
				p.status.isLongBreak = true
				p.status.pomodoroCount = 0
			} else {
				// Short Break
				p.status.isShortBreak = true
			}
		}
	case p.status.isShortBreak:
		if isBreakOver(p.status.tickerCount, p.flags.shortBreakDuration) {
			p.status.isShortBreak = false
			p.status.isWork = true
			p.endPomodoro()
		}
	case p.status.isLongBreak:
		if isBreakOver(p.status.tickerCount, p.flags.longBreakDuration) {
			p.status.isLongBreak = false
			p.status.isWork = true
			p.endPomodoro()
		}
	}
}

func (p *Pomodoro) updateUi() {
	switch {
	case p.status.isWork:
		percent := calculatePercentage(p.status.tickerCount, p.flags.workDuration)
		label := formatLabel(percent, p.flags.workDuration, p.status.tickerCount)
		p.ui.setLoaderData(workText, percent, label)
	case p.status.isShortBreak:
		percent := calculatePercentage(p.status.tickerCount, p.flags.shortBreakDuration)
		label := formatLabel(percent, p.flags.shortBreakDuration, p.status.tickerCount)
		p.ui.setLoaderData(shortBreakText, percent, label)
	case p.status.isLongBreak:
		percent := calculatePercentage(p.status.tickerCount, p.flags.longBreakDuration)
		label := formatLabel(percent, p.flags.longBreakDuration, p.status.tickerCount)
		p.ui.setLoaderData(longBreakText, percent, label)
	}
}

func (p *Pomodoro) endPomodoro() {
	p.notifier.play()
	p.status.tickerCount = 0
}

func (p *Pomodoro) resetStatus() {
	p.status.isWork = true
	p.status.isLongBreak = false
	p.status.isShortBreak = false
	p.status.pomodoroCount = 0
	p.status.tickerCount = 0
}
