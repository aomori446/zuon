package core

import (
	"errors"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/aomori446/zuon/front/i18n"
	"github.com/aomori446/zuon/internal"
)

func ShowLocalizedError(err error, parent fyne.Window) {
	if err == nil {
		return
	}
	
	var msg string
	switch {
	case errors.Is(err, internal.ErrImageNotSupported):
		msg = i18n.T("err_image_not_supported")
	case errors.Is(err, internal.ErrImageTooSmall):
		msg = i18n.T("err_image_too_small")
	case errors.Is(err, internal.ErrDataNotFound):
		msg = i18n.T("err_data_not_found")
	case errors.Is(err, internal.ErrDecryptionFailed):
		msg = i18n.T("err_decryption_failed")
	case errors.Is(err, internal.ErrInvalidAPIKey):
		msg = i18n.T("err_invalid_api_key")
	case errors.Is(err, internal.ErrNetworkIssue):
		msg = i18n.T("err_network_issue")
	case errors.Is(err, internal.ErrRateLimited):
		msg = i18n.T("err_rate_limited")
	case errors.Is(err, internal.ErrExtensionTooLong):
		msg = i18n.T("err_extension_too_long")
	case errors.Is(err, internal.ErrInternal):
		msg = i18n.T("err_internal")
	case errors.Is(err, internal.ErrNoCarrier):
		msg = i18n.T("err_no_carrier")
	case errors.Is(err, internal.ErrNoSource):
		msg = i18n.T("err_no_extract_source")
	case errors.Is(err, internal.ErrNoText):
		msg = i18n.T("err_no_text")
	case errors.Is(err, internal.ErrNoFile):
		msg = i18n.T("err_no_file")
	case errors.Is(err, internal.ErrPasswordShort):
		msg = i18n.T("err_password_short")
	case errors.Is(err, internal.ErrSessionExpired):
		msg = i18n.T("err_session_expired")
	default:
		
		msg = i18n.T("err_internal") + "\n\nDetails: " + err.Error()
	}
	
	d := dialog.NewInformation(i18n.T("err_title"), msg, parent)
	d.Show()
}
