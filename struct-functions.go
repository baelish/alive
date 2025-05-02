package main

import (
	"encoding/json"
	"errors"
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
		*s = Green
	case "grey":
		*s = Grey
	case "gray":
		*s = Grey
	case "noUpdate":
		*s = NoUpdate
	case "red":
		*s = Red
	case "amber":
		*s = Amber
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
		*bs = Dot
	case "micro":
		*bs = Micro
	case "dmicro":
		*bs = Dmicro
	case "small":
		*bs = Small
	case "dsmall":
		*bs = Dsmall
	case "medium":
		*bs = Medium
	case "dmedium":
		*bs = Dmedium
	case "large":
		*bs = Large
	case "dlarge":
		*bs = Dlarge
	case "xlarge":
		*bs = Xlarge
	default:
		return fmt.Errorf("invalid box size")
	}

	return nil
}

func (bs BoxSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(bs.String())
}

type Duration struct {
	time.Duration
	bool
}

// This is a custom UnmarshalJSON() for a time.Duration that can have a null
// value ("") as well as being backward compatible with supplying a string with
// the number of seconds.
func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		d.bool = true
		return nil
	case string:
		var err error
		if value == "" {
			d.Duration = 0
			d.bool = false
			return nil
		}
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			i, err2 := strconv.Atoi(value)
			if err2 != nil {
				return errors.Join(err, err2)
			}
			d.Duration = time.Duration(i) * time.Second
			d.bool = true
			return nil
		} else {
			d.bool = true
			return nil
		}
	default:
		return fmt.Errorf("invalid type for duration (%T)", v)
	}
}

// This is a custom MarshalJSON() for a time.Duration that can have a null value
// ("")
func (d Duration) MarshalJSON() (b []byte, err error) {
	if d.bool {
		return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
	} else {
		return []byte(`""`), nil
	}
}
