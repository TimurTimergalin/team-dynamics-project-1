package user_auth

import "crypto"

type PublicKeyFetcher interface {
	getPublicKeys() []crypto.PublicKey
}
