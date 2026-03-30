package k8s

import "crypto"

type JwksFetcher interface {
	getKeyFromJwks(kid string) crypto.PublicKey
}
