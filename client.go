package pomodoro

import (
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func Start() {
	// Build Flags
	flags := setFlags()

	// Build Ui
	initUi()
	defer termui.Close()
	grid := termui.NewGrid()
	heading := widgets.NewParagraph()
	loader := widgets.NewGauge()
	ui := newUi(grid, heading, loader)

	// Build Notifier
	sound, buffer := initSoundBuffer()
	notifier := newNotifier(sound, buffer)

	// Build Pomodoro
	pomodoro := newPomodoro(flags, ui, notifier)
	pomodoro.drawUi()
	pomodoro.startTimer()
}
