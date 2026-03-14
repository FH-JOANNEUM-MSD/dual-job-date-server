package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"math/big"
)

// GetSupabasePublicKey konvertiert die JWK-Koordinaten in einen ECDSA Public Key.
// Da der Key public ist, können wir ihn hier direkt hinterlegen oder später aus der .env laden.
func GetSupabasePublicKey() *ecdsa.PublicKey {
	// Deine X und Y Werte von Supabase
	xStr := "5Unoc0jPrbl4FtJevJ-AR05eqT2DdOSuwqHWd31SgH0"
	yStr := "4HA_AIR9XvfqHdswYXkRPd3hIkiQ-42DcimpF7rg4hc"

	xBytes, _ := base64.RawURLEncoding.DecodeString(xStr)
	yBytes, _ := base64.RawURLEncoding.DecodeString(yStr)

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}
}
