package auth

import "time"

const (
	JWTTokenKey = "JWTtoken"

	tokenExpiretime = time.Minute * 5

	testPrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIFfa14vgswEH/ySQmOaZ+padFPNs2db03TMDG0SzF/1ZoAoGCCqGSM49
AwEHoUQDQgAEf8SFq4YvKfUWFnZWed4ULWS5j5ufpYJ/rzKX98nNtU8OVlfeUQ4b
PYeiaEpP5sjPiMI++w/2OGhHxkiLl1vKRQ==
-----END EC PRIVATE KEY-----`

	testPubliKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEf8SFq4YvKfUWFnZWed4ULWS5j5uf
pYJ/rzKX98nNtU8OVlfeUQ4bPYeiaEpP5sjPiMI++w/2OGhHxkiLl1vKRQ==
-----END PUBLIC KEY-----`
)

/*
const testPrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIHIxGwJwHqwhOjUFSsObAtmG7wAxZ4znLjkK7c/fXXj9oAoGCCqGSM49
AwEHoUQDQgAEf8SFq4YvKfUWFnZWed4ULWS5j5ufpYJ/rzKX98nNtU8OVlfeUQ4b
PYeiaEpP5sjPiMI++w/2OGhHxkiLl1vKRQ==
-----END EC PRIVATE KEY-----`

const testPubliKey = `MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEf8SFq4YvKfUWFnZWed4ULWS5j5ufpYJ/rzKX98nNtU8OVlfeUQ4bPYeiaEpP5sjPiMI++w/2OGhHxkiLl1vKRQ==`
*/
