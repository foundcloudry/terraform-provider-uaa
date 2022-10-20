package apiheaders

type ApiHeader int64

const (
	Authorization ApiHeader = iota
	IfMatch
	ZoneId
)

func (s ApiHeader) String() string {
	switch s {
	case Authorization:
		return "Authorization"
	case IfMatch:
		return "If-Match"
	case ZoneId:
		return "X-Identity-Zone-Id"
	}
	return "unknown"
}
