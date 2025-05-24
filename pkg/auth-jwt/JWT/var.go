package jwt

import "time"

const (
	DefaultExpireTime = 5 * time.Minute
	DefaultMaxRefresh = 10 * time.Minute
	TokenExpireTime   = time.Minute * 5

	TestPrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIFfa14vgswEH/ySQmOaZ+padFPNs2db03TMDG0SzF/1ZoAoGCCqGSM49
AwEHoUQDQgAEf8SFq4YvKfUWFnZWed4ULWS5j5ufpYJ/rzKX98nNtU8OVlfeUQ4b
PYeiaEpP5sjPiMI++w/2OGhHxkiLl1vKRQ==
-----END EC PRIVATE KEY-----`

	TestPublicKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEf8SFq4YvKfUWFnZWed4ULWS5j5uf
pYJ/rzKX98nNtU8OVlfeUQ4bPYeiaEpP5sjPiMI++w/2OGhHxkiLl1vKRQ==
-----END PUBLIC KEY-----`
)
