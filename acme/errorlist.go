package acme

import "strings"

const (
	ZEROSSL_ERR_NEW_ACCOUNT_TOO_MANY_REQUESTS   = "429 ::POST :: https://acme.zerossl.com/v2/DV90/newAccount"
	ZEROSSL_ERR_NEW_NONCE_TOO_MANY_REQUESTS     = "429 ::HEAD :: https://acme.zerossl.com/v2/DV90/newNonce"
	ZEROSSL_ERR_GET_DIRECTORY_TOO_MANY_REQUESTS = "429 ::GET :: https://acme.zerossl.com/v2/DV90"
	ZEROSSL_ERR_AUTHZ_TOO_MANY_REQUESTS         = "429 ::POST :: https://acme.zerossl.com/v2/DV90/authz"
	ZEROSSL_ERR_NEW_ORDER_TOO_MANY_REQUESTS     = "429 ::POST :: https://acme.zerossl.com/v2/DV90/newOrder"
	ZEROSSL_ERR_ACCOUNT_TOO_MANY_REQUESTS       = "429 ::POST :: https://acme.zerossl.com/v2/DV90/account"
	ZEROSSL_ERR_REVOKE_CERT_TOO_MANY_REQUESTS   = "429 ::POST :: https://acme.zerossl.com/v2/DV90/revokeCert"
	LETSENCRYPT_ERR_RATE_LIMITED                = "429 :: POST :: https://acme-staging-v02.api.letsencrypt.org/acme/new-acct"
	ERR_TIME_LIMIT_EXCEEDED                     = "time limit exceeded"
)

func isAbleToRetry(errCode string) bool {
	return strings.Contains(errCode, ZEROSSL_ERR_NEW_ACCOUNT_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ZEROSSL_ERR_NEW_NONCE_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ZEROSSL_ERR_GET_DIRECTORY_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ZEROSSL_ERR_AUTHZ_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ZEROSSL_ERR_NEW_ORDER_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ZEROSSL_ERR_ACCOUNT_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ZEROSSL_ERR_REVOKE_CERT_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, LETSENCRYPT_ERR_RATE_LIMITED) ||
		strings.Contains(errCode, ERR_TIME_LIMIT_EXCEEDED)
}
