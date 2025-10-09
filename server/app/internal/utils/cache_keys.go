package utils

import "fmt"

func GetRefreshTokenKey(token string) string {
	return fmt.Sprintf("refresh_token:%s", token)
}
