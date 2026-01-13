package model

type User struct {
	ID              int
	Email           string
	HashedPassword  string
	Username        *string
	UserImageUri    *string
	CompanyName     *string
	CompanyImageUri *string
}
