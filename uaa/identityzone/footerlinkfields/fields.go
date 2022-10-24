package footerlinkfields

type FooterLinkField int64

const (
	Name FooterLinkField = iota
	Url
)

func (s FooterLinkField) String() string {
	switch s {
	case Name:
		return "name"
	case Url:
		return "url"
	}
	return "unknown"
}
