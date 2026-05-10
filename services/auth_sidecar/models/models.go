package models

type UserId struct {
	SteamId  *int64
	PlayerId *int64
}

type Principal struct {
	UserId         *UserId
	ServiceAccount *string
}

// AuthorityMap maps a Principal to a list of roles.
// A zero-value Principal (all nil fields) means "any user" - does not apply to service accounts.
type AuthorityMap map[Principal][]string

func AnyUserPrincipal() Principal {
	return Principal{}
}

func UserPrincipal(userId UserId) Principal {
	return Principal{UserId: &userId}
}

func ServiceAccountPrincipal(sa string) Principal {
	return Principal{ServiceAccount: &sa}
}
