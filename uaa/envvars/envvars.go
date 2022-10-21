package envvars

type EnvironmentVariables int64

const (
	UaaAuthUrl EnvironmentVariables = iota
	UaaClientId
	UaaDebug
	UaaDialTimeout
	UaaClientSecret
	UaaLoginUrl
	UaaSkipSslValidation
	UaaTrace
)

func (s EnvironmentVariables) String() string {
	switch s {
	case UaaAuthUrl:
		return "UAA_AUTH_URL"
	case UaaClientId:
		return "UAA_CLIENT_ID"
	case UaaClientSecret:
		return "UAA_CLIENT_SECRET"
	case UaaDebug:
		return "UAA_DEBUG"
	case UaaDialTimeout:
		return "UAA_DIAL_TIMEOUT"
	case UaaLoginUrl:
		return "UAA_LOGIN_URL"
	case UaaSkipSslValidation:
		return "UAA_SKIP_SSL_VALIDATION"
	case UaaTrace:
		return "UAA_TRACE"
	}
	return "unknown"
}
