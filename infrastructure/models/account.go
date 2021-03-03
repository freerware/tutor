package models

type Account struct {
	Model

	GivenName string
	Surname   string
	Username  string `gorm:"column:PRIMARY_CREDENTIAL"`
}

func (a Account) TableName() string {
	return "ACCOUNT"
}
