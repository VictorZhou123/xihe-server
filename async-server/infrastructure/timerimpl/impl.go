package timerimpl

import (
	"time"

	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

// timer
type TimerGroup struct {
	group map[string]DTimer
}

type DTimer struct {
	name   string
	handle func(interface{}) error
	timer  *time.Timer
	cfg    *TimerConfig
}

func newTimer(hour, min, second int) *time.Timer {
	return time.NewTimer(
		utils.RemainTime(hour, min, second),
	)
}

func NewTimerGroup(
	cfg *Config, handles map[string]func(interface{}) error,
) (g TimerGroup) {

	m := make(map[string]DTimer, len(handles))

	for key := range handles {
		f := handles[key]
		if config, ok := cfg.TimerGroup[key]; !ok {
			logrus.Warnf("timer init error, cannot found timer name: %s", key)
		} else {
			m[key] = DTimer{
				name:   key,
				handle: f,
				timer:  newTimer(config.TriggerHour, config.TriggerMin, config.TriggerSecond),
				cfg:    &config,
			}
		}
	}

	g.group = m

	return
}

// TimerWatcher
type TimerWatcher struct {
	timers TimerGroup
	cfg    Config
}

func NewTimerWatcher(
	cfg Config,
	handles map[string]func(interface{}) error,
) *TimerWatcher {
	return &TimerWatcher{
		timers: NewTimerGroup(&cfg, handles),
		cfg:    cfg,
	}
}

func (t *TimerWatcher) startTimer() {
	logrus.Debug("start timer")

	for key := range t.timers.group {
		timer := t.timers.group[key]

		go func() {
			for {
				<-timer.timer.C
				logrus.Infof("timer: %s called in %s", timer.name, utils.Date())

				if err := timer.handle(nil); err != nil {
					logrus.Errorf("timer handle called error, err: %s", err.Error())
				}

				timer.timer.Reset(utils.RemainTime(
					timer.cfg.TriggerHour,
					timer.cfg.TriggerMin,
					timer.cfg.TriggerSecond,
				))
			}
		}()
	}
}

func (t *TimerWatcher) stop() {
	for key := range t.timers.group {
		timer := t.timers.group[key]
		timer.timer.Stop()
	}
}

func (t *TimerWatcher) Run() {

	t.startTimer()

}

func (t *TimerWatcher) Exit() {
	t.stop()
}
