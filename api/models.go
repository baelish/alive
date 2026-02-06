package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Error response for api calls
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

func (d *Duration) Duration() time.Duration {
	if d == nil {
		return 0
	}
	return time.Duration(*d)
}

// Event struct is used to stream events to dashboard.
type Event struct {
	ID          string    `json:"id,omitempty"`
	After       string    `json:"after,omitempty"`
	Box         *Box      `json:"box,omitempty"`
	Status      Status    `json:"status,omitempty"`
	Message     string    `json:"lastMessage,omitempty"`
	ExpireAfter *Duration `json:"expireAfter"`
	MaxTBU      *Duration `json:"maxTBU"`
	Type        string    `json:"type"`
}

// Links describes a URL with a name.
type Links struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Message struct {
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	TimeStamp time.Time `json:"timeStamp"`
}

type Status int

const (
	Grey Status = iota
	Red
	Amber
	Green
	NoUpdate
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

type BoxSize int

const (
	Dot BoxSize = iota
	Micro
	Dmicro
	Small
	Dsmall
	Medium
	Dmedium
	Large
	Dlarge
	Xlarge
)

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
		// If it's a small number (< 1 million), treat as seconds (intuitive)
		// If it's huge (>= 1 million), treat as nanoseconds (Go-marshaled duration)
		if value == 0 {
			d = nil
		}
		if value < 1000000 {
			*d = Duration(time.Duration(value) * time.Second)
		} else {
			*d = Duration(time.Duration(value)) // nanoseconds
		}
		return nil
	case string:
		if value == "" {
			d = nil
			return nil
		}
		x, err := time.ParseDuration(value)
		*d = Duration(x)
		if err != nil {
			// Backward compatible: treat plain number string as seconds
			i, err2 := strconv.Atoi(value)
			if err2 != nil {
				return errors.Join(err, err2)
			}
			*d = Duration(time.Duration(i) * time.Second)
			return nil
		} else {
			return nil
		}
	default:
		return fmt.Errorf("invalid type for duration (%T)", v)
	}
}

// This is a custom MarshalJSON() for a time.Duration that can have a null value
// ("")
func (d Duration) MarshalJSON() (b []byte, err error) {
	return json.Marshal(d.String())
}

// Box represents a single item on our monitoring screen.
type Box struct {
	ID          string             `json:"id"`
	Description string             `json:"description,omitempty"`
	DisplayName string             `json:"displayName,omitempty"`
	Name        string             `json:"name"`
	Info        *map[string]string `json:"info,omitempty"`
	Parent      string             `json:"parent,omitempty"`
	Size        BoxSize            `json:"size"`
	Status      Status             `json:"status"`
	ExpireAfter *Duration          `json:"expireAfter,omitempty"`
	MaxTBU      *Duration          `json:"maxTBU,omitempty"`
	Messages    []Message          `json:"messages"`
	LastUpdate  time.Time          `json:"lastUpdate"`
	LastMessage string             `json:"lastMessage"`
	Links       []Links            `json:"links"`
}

func (b *Box) Sanitise() {
	if b.MaxTBU != nil && *b.MaxTBU == 0 {
		b.MaxTBU = nil
	}
	if b.ExpireAfter != nil && *b.ExpireAfter == 0 {
		b.ExpireAfter = nil
	}
}

func (b *Box) UnmarshalJSON(data []byte) error {
	type Alias Box
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	b.Sanitise()

	return nil
}
