package model

import (
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%s"`, d.Format(dateLayout)), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("date: expected quoted string, got %s", data)
	}
	t, err := time.Parse(dateLayout, string(data[1:len(data)-1]))
	if err != nil {
		return fmt.Errorf("date: %w", err)
	}
	d.Time = t
	return nil
}
