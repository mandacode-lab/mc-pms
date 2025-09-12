package utils

import (
	"strings"
)

func ExtractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func CreateBearerToken(token string) string {
	return "Bearer " + token
}