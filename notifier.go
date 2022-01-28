package pomodoro

import (
	"github.com/faiface/beep/speaker"
)

func (n *Notifier) play() {
	n.sound = n.buffer.Streamer(0, n.buffer.Len())
	speaker.Play(n.sound)
}
