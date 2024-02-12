package components

type Page struct {
	Path  string
	Title string
}

type Base struct {
	Page  *Page
	Pages []*Page
}
