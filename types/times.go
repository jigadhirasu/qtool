package types

import (
	"encoding/json"
	"strings"
	"time"
)

func NewTime(value string) *Time {
	ti, _ := time.Parse(TimeFormat, value)
	return &Time{Time: ti}
}

func NewDate(value string) *Date {
	ti, _ := time.Parse(DateFormat, value)
	return &Date{Time: ti}
}

// 所有跟金額有關的整數需使用該倍率
const TimeFormat = "2006-01-02 15:04:05"

type Time struct {
	time.Time
}

func (ti Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(ti.Format(TimeFormat))
}

func (ti *Time) UnmarshalJSON(data []byte) error {
	source := strings.ReplaceAll(string(data), `"`, ``)
	t, err := time.Parse(TimeFormat, source)
	if err != nil {
		return err
	}

	ti.Time = t
	return nil
}

const DateFormat = "2006-01-02"

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(DateFormat))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	source := strings.ReplaceAll(string(data), `"`, ``)
	// source := strings.ReplaceAll(string(data), `\"`, ``)
	t, err := time.Parse(DateFormat, source)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

func (d Date) String() string {
	return ""
}
