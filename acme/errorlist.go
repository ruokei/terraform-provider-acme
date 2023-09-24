package acme

import "strings"

const (
	ERR_ACME_ZEROSSL_NEW_ACCOUNT_TOO_MANY_REQUESTS   = "429 ::POST :: https://acme.zerossl.com/v2/DV90/newAccount"
	ERR_ACME_ZEROSSL_NEW_NONCE_TOO_MANY_REQUESTS     = "429 ::HEAD :: https://acme.zerossl.com/v2/DV90/newNonce"
	ERR_ACME_ZEROSSL_GET_DIRECTORY_TOO_MANY_REQUESTS = "429 ::GET :: https://acme.zerossl.com/v2/DV90"
	ERR_ACME_ZEROSSL_AUTHZ_TOO_MANY_REQUESTS         = "429 ::POST :: https://acme.zerossl.com/v2/DV90/authz"
	ERR_ACME_ZEROSSL_NEW_ORDER_TOO_MANY_REQUESTS     = "429 ::POST :: https://acme.zerossl.com/v2/DV90/newOrder"
	ERR_ACME_ZEROSSL_ACCOUNT_TOO_MANY_REQUESTS       = "429 ::POST :: https://acme.zerossl.com/v2/DV90/account"
	ERR_ACME_ZEROSSL_REVOKE_CERT_TOO_MANY_REQUESTS   = "429 ::POST :: https://acme.zerossl.com/v2/DV90/revokeCert"
	ERR_ACME_LETSENCRYPT_RATE_LIMITED                = "429 :: POST :: https://acme-staging-v02.api.letsencrypt.org/acme/new-acct"
	ERR_ACME_TIME_LIMIT_EXCEEDED                     = "time limit exceeded"
	ERR_ACME_DOMAINS_HAD_A_PROBLEM                   = "error: one or more domains had a problem "
)

func isAbleToRetry(errCode string) bool {
	return strings.Contains(errCode, ERR_ACME_ZEROSSL_NEW_ACCOUNT_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ERR_ACME_ZEROSSL_NEW_NONCE_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ERR_ACME_ZEROSSL_GET_DIRECTORY_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ERR_ACME_ZEROSSL_AUTHZ_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ERR_ACME_ZEROSSL_NEW_ORDER_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ERR_ACME_ZEROSSL_ACCOUNT_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ERR_ACME_ZEROSSL_REVOKE_CERT_TOO_MANY_REQUESTS) ||
		strings.Contains(errCode, ERR_ACME_LETSENCRYPT_RATE_LIMITED) ||
		strings.Contains(errCode, ERR_ACME_TIME_LIMIT_EXCEEDED)
}
