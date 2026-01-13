package internal

import "errors"

var (
	ErrImageNotSupported = errors.New("err_image_not_supported")
	ErrImageTooSmall     = errors.New("err_image_too_small")
	ErrDataNotFound      = errors.New("err_data_not_found")
	ErrDecryptionFailed  = errors.New("err_decryption_failed")
	ErrExtensionTooLong  = errors.New("err_extension_too_long")
	
	ErrInvalidAPIKey = errors.New("err_invalid_api_key")
	ErrNetworkIssue  = errors.New("err_network_issue")
	ErrRateLimited   = errors.New("err_rate_limited")
	
	ErrInternal = errors.New("err_internal")
	
	ErrNoCarrier     = errors.New("err_no_carrier")
	ErrNoSource      = errors.New("err_no_source")
	ErrNoText        = errors.New("err_no_text")
	ErrNoFile        = errors.New("err_no_file")
	ErrPasswordShort = errors.New("err_password_short")
	ErrSessionExpired = errors.New("err_session_expired")
)
