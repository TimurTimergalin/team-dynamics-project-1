package models

import "crypto/rsa"

type UserId struct {
	SteamId  *int64
	PlayerId *int64
}

type TokenPayload struct {
	UserId  UserId
	TokenId *string
}

type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type TokenWithExp struct {
	Token string
	ExpMs int64
}
