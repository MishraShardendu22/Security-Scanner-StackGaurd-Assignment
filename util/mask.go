package util

func MaskSecret(secret string) string {

	if len(secret) <= 8 {
		return "****"
	}

	return secret[:4] + "********" + secret[len(secret)-4:]
}
