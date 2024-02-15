package components

type Page struct {
	Path       string
	Title      string
	Breadcrumb string
}

type Base struct {
	Navbar       []*Page
	NavbarActive *Page
}
