package valueObject

type Url struct {
	url string 
}

func (u *Url) Url() string {
	return u.url
}

func NewUrl(url string) {
	
}