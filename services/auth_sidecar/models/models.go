package models

type UserId struct {
	SteamId  *int64
	PlayerId *int64
}

// AuthorityMap maps a UserId to a list of roles. A zero-value UserId means "any user".
type AuthorityMap map[UserId][]string
