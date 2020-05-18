package models

// Motion is a basic model of an item for SIU
type Motion struct {
	ID       string
	Name     string
	URL      string
	Shortcut string
	Usage    int32
}
