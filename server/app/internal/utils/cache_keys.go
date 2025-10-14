package utils

import (
	"fmt"
)

func GetRefreshTokenKey(token string) string {
	return fmt.Sprintf("refresh_token:%s", token)
}

func GetEmailVerificationTokenKey(email string) string {
	return fmt.Sprintf("email_verification_token:%s", email)
}

func GetPasswordResetTokenKey(email string) string {
	return fmt.Sprintf("password_reset_token:%s", email)
}
