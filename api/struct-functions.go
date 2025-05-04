package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func (s Status) String() string {
	return [...]string{
		"grey",
		"red",
		"amber",
		"green",
		"noUpdate",
	}[s]
}

func (s *Status) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	switch str {
	case "green":
		s = Ptr(Green)
	case "grey":
		s = Ptr(Grey)
	case "gray":
		s = Ptr(Grey)
	case "noUpdate":
		s = Ptr(NoUpdate)
	case "red":
		s = Ptr(Red)
	case "amber":
		s = Ptr(Amber)
	default:
		return fmt.Errorf("invalid status")
	}

	return nil
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (bs BoxSize) String() string {
	return [...]string{
		"dot",
		"micro",
		"dmicro",
		"small",
		"dsmall",
		"medium",
		"dmedium",
		"large",
		"dlarge",
		"xlarge",
	}[bs]
}

func (bs *BoxSize) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	switch str {
	case "dot":
		bs = Ptr(Dot)
	case "micro":
		bs = Ptr(Micro)
	case "dmicro":
		bs = Ptr(Dmicro)
	case "small":
		bs = Ptr(Small)
	case "dsmall":
		bs = Ptr(Dsmall)
	case "medium":
		bs = Ptr(Medium)
	case "dmedium":
		bs = Ptr(Dmedium)
	case "large":
		bs = Ptr(Large)
	case "dlarge":
		bs = Ptr(Dlarge)
	case "xlarge":
		bs = Ptr(Xlarge)
	default:
		return fmt.Errorf("invalid box size")
	}

	return nil
}

func (bs BoxSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(bs.String())
}

func DurationFromString(ds *DurationString) *time.Duration {
	if ds == nil || *ds == "" {
		return nil
	}
	str := string(*ds)

	// Try parsing as an integer (interpreted as seconds)
	if seconds, err := strconv.Atoi(str); err == nil {
		d := time.Duration(seconds) * time.Second
		return &d
	}

	// Try parsing as full duration string (e.g., "10s", "5m")
	if d, err := time.ParseDuration(str); err == nil {
		return &d
	} else {
		return nil
	}
}

func DurationStringFromDuration(d *time.Duration) *DurationString {
	if d == nil {
		return nil
	}
	s := DurationString(d.String())
	return &s
}

func Ptr[T any](v T) *T {
	return &v
}
