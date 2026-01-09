package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// Error response for api calls
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Event struct is used to stream events to dashboard.
type Event struct {
	ID          string         `json:"id,omitempty"`
	After       string         `json:"after,omitempty"`
	Box         *Box           `json:"box,omitempty"`
	Status      Status         `json:"status,omitempty"`
	Message     string         `json:"lastMessage,omitempty"`
	ExpireAfter *time.Duration `json:"expireAfter,omitempty"`
	MaxTBU      *time.Duration `json:"maxTBU,omitempty"`
	Type        string         `json:"type"`
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
	ExpireAfter *time.Duration     `json:"expireAfter,omitempty"`
	MaxTBU      *time.Duration     `json:"maxTBU,omitempty"`
	Messages    []Message          `json:"messages"`
	LastUpdate  time.Time          `json:"lastUpdate"`
	LastMessage string             `json:"lastMessage"`
	Links       []Links            `json:"links"`
}
