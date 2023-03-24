package types

type User struct {
	Node
	UserName string
	Email    string
	Password string
	Person   string
}

// Context is the struct passed to templates
type Context struct {
	User //using
	Person
	Persons   []Person
	CSRFToken string
	Page
}

type Page struct {
	IsDisabled bool
	IsOwner    bool
}
