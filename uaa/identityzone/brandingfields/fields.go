package brandingfields

type BrandingField int64

const (
	BannerBackgroundColor BrandingField = iota
	BannerLogo
	BannerText
	BannerTextColor
	BannerUrl
	CompanyLogo
	CompanyName
	Favicon
	FooterLinks
	FooterText
)

func (s BrandingField) String() string {
	switch s {
	case BannerBackgroundColor:
		return "banner_bg_color"
	case BannerLogo:
		return "banner_logo"
	case BannerText:
		return "banner_text"
	case BannerTextColor:
		return "banner_text_color"
	case BannerUrl:
		return "banner_url"
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
