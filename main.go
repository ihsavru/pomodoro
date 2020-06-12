package main

import (
	"fmt"
	"flag"
	"log"
	"time"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep"
	"os"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func isLongBreak(pomodoroCount int) bool {
	return pomodoroCount % 4 == 0
}

func isWorkOver(secondsPassed , workDuration int) bool {
	return secondsPassed == workDuration
}

func isBreakOver(secondsPassed , breakDuration int) bool {
	return secondsPassed == breakDuration
}

func formatLabel(percent , duration , secondsPassed int) string {
	return fmt.Sprintf(
		"%v%% (%vm %vs left)",
		percent,
		(duration - secondsPassed ) / 60,
		(duration - secondsPassed ) % 60,
	)
}

func calculateDimensions(termWidth int, termHeight int) (int, int, int, int) {
	var x1, y1 int =  termWidth / 2 - 30, termHeight / 2 - 10
	var x2, y2 int = termWidth / 2 + 30, termHeight / 2 + 10
	return x1, y1, x2, y2
}

func main() {
	workDurationPtr := flag.Int("work", 25, "Duration of work interval in minutes")
	shortBreakPtr := flag.Int("shortBreak", 5, "Duration of short break interval in minutes")
	longBreakPtr := flag.Int("longBreak", 5, "Duration of long break interval in minutes")

	flag.Parse()

	workDurationMinutes := *workDurationPtr
	workDurationSeconds := workDurationMinutes * 60

	shortBreakDurationMinutes := *shortBreakPtr
	shortBreakDurationSeconds := shortBreakDurationMinutes * 60

	longBreakDurationMinutes := *longBreakPtr
	longBreakDurationSeconds := longBreakDurationMinutes * 60

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(calculateDimensions(termWidth, termHeight))

	workText := "WORK"
	shortBreakText := "SHORT BREAK"
	longBreakText := "LONG BREAK"

	heading := widgets.NewParagraph()
	heading.Title = "POMODORO"
	heading.Text = `[The](fg:blue) [Pomodoro Technique](fg:bold,bg:red) [is a time management method developed by Francesco Cirillo in the late 1980s. The technique uses a timer to break down work into intervals, traditionally 25 minutes in length, separated by short breaks (5 minutes). Each interval is known as a pomodoro, from the Italian word for 'tomato', after the tomato-shaped kitchen timer that Cirillo used as a university student.](fg:blue)`
	heading.TitleStyle.Fg = ui.ColorGreen
	heading.BorderStyle.Fg = ui.ColorGreen

	loader := widgets.NewGauge()
	loader.Title = workText
	loader.Label = formatLabel(loader.Percent, workDurationSeconds, 0)
	loader.BarColor = ui.ColorYellow
	loader.TitleStyle.Fg = ui.ColorMagenta
	loader.BorderStyle.Fg = ui.ColorMagenta
	loader.LabelStyle = ui.NewStyle(ui.ColorCyan)
	loader.Percent = 0
	loader.SetRect(0, 0, 25, 5)

	grid.Set(
		ui.NewRow(2.0 / 3,
			ui.NewCol(1, heading),
		),
		ui.NewRow(1.0 / 3,
			ui.NewCol(1, loader),
		),
	)

	f, err := os.Open(os.Getenv("GOPATH") + "/src/github.com/ihsavru/pomodoro-cli/notification.wav")
	if err != nil {
		log.Fatalf("failed to load audio file: %v", err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatalf("failed to decode notification.wav: %v", err)
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	isWork := true

	sound := buffer.Streamer(0, buffer.Len())
	speaker.Play(sound)

	tickerCount := 0
	pomodoroCount := 0
	uiEvents := ui.PollEvents()
	ui.Render(grid)
	secondTicker := time.NewTicker(time.Second)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				termWidth, termHeight = payload.Width, payload.Height
				grid.SetRect(calculateDimensions(termWidth, termHeight))
				ui.Clear()
				ui.Render(grid)
			}
		case <- secondTicker.C:
			tickerCount++
			if isWork {
				loader.Percent = ( tickerCount * 100 ) / workDurationSeconds
				loader.Label = formatLabel(loader.Percent, workDurationSeconds, tickerCount)

				if isWorkOver(tickerCount, workDurationSeconds) {
					isWork = false
					tickerCount = 0
					pomodoroCount++
					loader.Title = shortBreakText
					if isLongBreak(pomodoroCount) {
						loader.Title = longBreakText
					}
					sound := buffer.Streamer(0, buffer.Len())
					speaker.Play(sound)
				}
			} else {
				breakDurationSeconds := shortBreakDurationSeconds
				if isLongBreak(pomodoroCount) {
					breakDurationSeconds = longBreakDurationSeconds
				}
				loader.Percent = ( tickerCount * 100 ) / breakDurationSeconds
				loader.Label = formatLabel(loader.Percent, breakDurationSeconds, tickerCount)

				if isBreakOver(tickerCount, breakDurationSeconds) {
					isWork = true
					tickerCount = 0
					loader.Title = workText
					sound := buffer.Streamer(0, buffer.Len())
					speaker.Play(sound)
				}
			}
			ui.Render(grid)
		}
	}
}
