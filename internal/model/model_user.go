package model

import (
	"github.com/Zhoangp/Auth-Service/pkg/common"
	"github.com/Zhoangp/Auth-Service/pkg/utils"
	"time"
)

type Users struct {
	common.SQLModel
	Email        string        `gorm:"column:email"`
	Password     string        `gorm:"column:password"`
	FirstName    string        `gorm:"column:firstName"`
	LastName     string        `gorm:"column:lastName"`
	Phone        string        `gorm:"column:phoneNumber"`
	Role         string        `gorm:"column:role"`
	Address      string        `gorm:"column:address"`
	IsInstructor bool          `gorm:"column:is_instructor"`
	Verified     bool          `gorm:"column:verified"`
	Avatar       *common.Image `gorm:"column:picture"`
	LastLogin    time.Time     `gorm:"column:lastLogin"`
}
type UserRegister struct {
	FirstName string `gorm:"column:firstName"`
	LastName  string `gorm:"column:lastName"`
	Phone     string `gorm:"column:phoneNumber"`
	Role      string `gorm:"column:role"`
	Address   string `gorm:"column:address"`
}
type UserLogin struct {
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (Users) TableName() string {
	return "Users"
}

func (u *Users) GetUserId() int {
	return u.Id
}

func (u *Users) GetUserEmail() string {
	return u.Email
}

func (u *Users) GetUserRole() string {
	return u.Role
}
func (u *Users) PrepareCreate() error {

	passHashed, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = passHashed
	u.Role = "guest"
	return nil
}
