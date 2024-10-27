package domain

type User struct {
	Id       int64  `gorm:"primaryKey"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}
