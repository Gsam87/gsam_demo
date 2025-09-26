package entity

type AddUser struct {
	Email    string
	Password string
}

type ModifyUserPassword struct {
	Email       string
	OldPassword string
	NewPassword string
}

type GetUser struct {
	Email    string
	Password string
}
