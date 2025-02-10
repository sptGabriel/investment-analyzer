package rest

type Payload struct {
	Title string `json:"title"`
}

func (p Payload) WithTitle(title string) Payload {
	p.Title = title
	return p
}

var (
	InternalServerErrorDefault = Payload{Title: "Internal Server Error"}
)
