package timerimpl

type Config struct {
	TimerGroup map[string]TimerConfig `json:"timer_group"`
}

type TimerConfig struct {
	TriggerHour   int `json:"trigger_hour"`
	TriggerMin    int `json:"trigger_min"`
	TriggerSecond int `json:"trigger_second"`
}
