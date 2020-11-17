package person

type User struct {
	Id string
	Username string
	Password string
	Email string
}

type UserHandler struct {
	sync.Mutex
	Store map[string]User
}