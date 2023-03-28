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
	User      User //using
	Person    Person
	Persons   []Person
	CSRFToken string
	Page      Page
	Error     error
}

type Page struct {
	IsDisabled bool
	IsOwner    bool
}
