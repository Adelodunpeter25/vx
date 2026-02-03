package terminalpane

import "time"

type Config struct {
	Shell       string
	Dir         string
	Env         []string
	Scrollback  int
	MaxFPS      int
	MinRedraw   time.Duration
	EnableMouse bool
}

func DefaultConfig() Config {
	return Config{
		Scrollback: 2000,
		MaxFPS:     60,
	}
}
