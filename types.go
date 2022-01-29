package pomodoro

import (
	"time"

	"github.com/faiface/beep"
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Flags struct {
	workDuration, shortBreakDuration, longBreakDuration int
}

type Ui struct {
	grid    *termui.Grid
	heading *widgets.Paragraph
	loader  *widgets.Gauge
}

type Notifier struct {
	sound  beep.StreamSeeker
	buffer *beep.Buffer
}

type PomodoroStatus struct {
	isWork, isShortBreak, isLongBreak bool
	pomodoroCount, tickerCount        int
}

type Pomodoro struct {
	flags        Flags
	ui           *Ui
	notifier     *Notifier
	secondTicker *time.Ticker
	status       PomodoroStatus
}
