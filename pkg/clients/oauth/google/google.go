package google

const (
	URL string = ""
)

type Provider struct {
}

func (p *Provider) Authenticate(identifier, token, onetimecode string) string {
	return ""
}
