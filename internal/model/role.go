package model

type Role struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name"`
}

type UserHasRole struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	UserID   uint   `json:"user_id" gorm:"index"`
	RoleID   uint   `json:"role_id" gorm:"index"`
	RoleName string `json:"role_name"`
}
