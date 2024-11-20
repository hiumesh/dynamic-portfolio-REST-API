package api

type ErrorCode = string

const (
	// ErrorCodeUnknown should not be used directly, it only indicates a failure in the error handling system in such a way that an error code was not assigned properly.
	ErrorCodeUnknown ErrorCode = "unknown"

	// ErrorCodeUnexpectedFailure signals an unexpected failure such as a 500 Internal Server Error.
	ErrorCodeUnexpectedFailure ErrorCode = "unexpected_failure"

	ErrorCodeValidationFailed       ErrorCode = "validation_failed"
	ErrorCodeBadJSON                ErrorCode = "bad_json"
	ErrorCodeEmailExists            ErrorCode = "email_exists"
	ErrorCodeBadJWT                 ErrorCode = "bad_jwt"
	ErrorCodeNotAdmin               ErrorCode = "not_admin"
	ErrorCodeNoAuthorization        ErrorCode = "no_authorization"
	ErrorCodeUserNotFound           ErrorCode = "user_not_found"
	ErrorCodeUserBanned             ErrorCode = "user_banned"
	ErrorCodeInviteNotFound         ErrorCode = "invite_not_found"
	ErrorCodeUnexpectedAudience     ErrorCode = "unexpected_audience"
	ErrorCodeIdentityAlreadyExists  ErrorCode = "identity_already_exists"
	ErrorCodeCaptchaFailed          ErrorCode = "captcha_failed"
	ErrorCodeSMSSendFailed          ErrorCode = "sms_send_failed"
	ErrorCodeEmailNotConfirmed      ErrorCode = "email_not_confirmed"
	ErrorCodeUserAlreadyExists      ErrorCode = "user_already_exists"
	ErrorCodeConflict               ErrorCode = "conflict"
	ErrorCodeIdentityNotFound       ErrorCode = "identity_not_found"
	ErrorCodeOverRequestRateLimit   ErrorCode = "over_request_rate_limit"
	ErrorCodeOverEmailSendRateLimit ErrorCode = "over_email_send_rate_limit"
	ErrorCodeOverSMSSendRateLimit   ErrorCode = "over_sms_send_rate_limit"
	ErrorCodeRequestTimeout         ErrorCode = "request_timeout"
)
