package main

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	return fmt.Sprintf("%d:%02d:%02d", h, m, d/time.Second)
}
