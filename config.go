package nrmg

import "fmt"

type Config struct {
	ZoneSizes []float64
	Cells     int
}

func (c *Config) normalizeSizes() {
	total := 0.
	for _, s := range c.ZoneSizes {
		total += s
	}
	for ix, s := range c.ZoneSizes {
		c.ZoneSizes[ix] = s / total
	}
}

func (c Config) String() string {
	var out string
	for ix, s := range c.ZoneSizes {
		out += fmt.Sprintf(" zone %d: %f |", ix, s)
	}
	return out
}
