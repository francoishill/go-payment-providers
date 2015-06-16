package PayfastUIConfig

type authorizedContext struct {
	FirstName string
	LastName  string
	Email     string
}

func CreateAuthorizedContext(firstName, lastName, email string) *authorizedContext {
	return &authorizedContext{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
}
