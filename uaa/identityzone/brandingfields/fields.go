package brandingfields

type BrandingField int64

const (
	CompanyLogo BrandingField = iota
	CompanyName
	Favicon
	FooterLinks
	FooterText
)

func (s BrandingField) String() string {
	switch s {
	case CompanyLogo:
		return "company_logo"
	case CompanyName:
		return "company_name"
	case Favicon:
		return "favicon"
	case FooterLinks:
		return "footer_link"
	case FooterText:
		return "footer_text"
	}
	return "unknown"
}
