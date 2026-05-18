package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const DateLayout = "2006-01-02"

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%s"`, d.Format(DateLayout)), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("date: expected quoted string, got %s", data)
	}
	t, err := time.Parse(DateLayout, string(data[1:len(data)-1]))
	if err != nil {
		return fmt.Errorf("date: %w", err)
	}
	d.Time = t
	return nil
}

func (d *Date) Scan(src any) error {
	switch v := src.(type) {
	case time.Time:
		d.Time = v
		return nil
	case nil:
		d.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("date: cannot scan %T", src)
	}
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}
