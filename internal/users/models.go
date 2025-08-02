package users

type User struct {
	Id           int
	Email        string
	HashPass     string
	RefreshToken string
	Role         string
	IsActive     bool
}

type AdminConfig struct {
	AdminEmail string
	HashedPass string
	Role       string
	Is_Active  bool
}
