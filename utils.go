package pomodoro

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func newPomodoro(flags Flags, ui *Ui, notifier *Notifier) *Pomodoro {
	return &Pomodoro{
		flags:        flags,
		ui:           ui,
		notifier:     notifier,
		secondTicker: time.NewTicker(time.Second),
		status: PomodoroStatus{
			isWork: true,
		},
	}
}

func newFlags(workDuration, shortBreakDuration, longBreakDuration int) Flags {
	return Flags{
		workDuration:       workDuration,
		shortBreakDuration: shortBreakDuration,
		longBreakDuration:  longBreakDuration,
	}
}

func newUi(grid *termui.Grid, heading *widgets.Paragraph, loader *widgets.Gauge) *Ui {
	return &Ui{
		grid:    grid,
		heading: heading,
		loader:  loader,
	}
}

func newNotifier(sound beep.StreamSeeker, buffer *beep.Buffer) *Notifier {
	return &Notifier{
		sound:  sound,
		buffer: buffer,
	}
}

func initUi() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termtermui: %v", err)
	}
}

func setFlags() Flags {
	workDurationPtr := flag.Int("work", 25, "Duration of work interval in minutes")
	shortBreakPtr := flag.Int("shortBreak", 5, "Duration of short break interval in minutes")
	longBreakPtr := flag.Int("longBreak", 5, "Duration of long break interval in minutes")

	flag.Parse()

	workDuration := *workDurationPtr * 60
	shortBreakDuration := *shortBreakPtr * 60
	longBreakDuration := *longBreakPtr * 60

	return newFlags(workDuration, shortBreakDuration, longBreakDuration)
}

func initSoundBuffer() (beep.StreamSeeker, *beep.Buffer) {
	f, err := os.Open("../notification.wav")
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

	sound := buffer.Streamer(0, buffer.Len())

	return sound, buffer
}

func calculateDimensions(termWidth, termHeight int) (x1, y1, x2, y2 int) {
	x1, y1 = termWidth/2-30, termHeight/2-10
	x2, y2 = termWidth/2+30, termHeight/2+10
	return
}

func isLongBreak(pomodoroCount int) bool {
	return pomodoroCount%4 == 0
}

func isWorkOver(secondsPassed, workDuration int) bool {
	return secondsPassed == workDuration
}

func isBreakOver(secondsPassed, breakDuration int) bool {
	return secondsPassed == breakDuration
}

func calculatePercentage(secondsPassed, totalSeconds int) int {
	return (secondsPassed * 100) / totalSeconds
}

func formatLabel(percent, duration, secondsPassed int) string {
	return fmt.Sprintf(
		"%v%% (%vm %vs left)",
		percent,
		(duration-secondsPassed)/60,
		(duration-secondsPassed)%60,
	)
}
