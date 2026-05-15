package models

type FriendSource int64

const (
	Internal FriendSource = 1 + iota
	Steam
)

type UserData struct {
	Id      int64
	Name    string
	SteamId *int64
	EosId   *string
}

type EosData struct {
	EosId string
	Name  string
}

type Friend struct {
	Data   *UserData
	Source FriendSource
}

type PageKey struct {
	LastUserId int64 `json:"lastUserId"`
}

type AddFriendResult int

const (
	AddFriendNoop AddFriendResult = iota
	AddFriendRequestSent
	AddFriendAccepted
)

type RemoveFriendResult int

const (
	RemoveFriendNoop RemoveFriendResult = iota
	RemoveFriendRequestCancelled
	RemoveFriendRequestDeclined
	RemoveFriendFriendRemoved
)

type SteamData struct {
	SteamId string `json:"steamId"`
	Name    string `json:"name"`
}
