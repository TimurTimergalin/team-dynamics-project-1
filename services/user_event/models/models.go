package models

type PlayerStatus int64

const (
	Offline PlayerStatus = 1 + iota
	Online
	Busy
)

type PlayerUserData struct {
	Id     int64
	Name   string
	Rating int64
}
