package model

import (
	"errors"
	"github.com/Zhoangp/Auth-Service/pkg/common"
)

// Định nghĩa các error cho riêng phần User-Service
var (
	ErrEmailOrPasswordInvalid = common.NewCustomError(
		errors.New("email or password invalid"),
		"email or password invalid",
	)

	ErrEmailExisted = common.NewCustomError(
		errors.New("email has already existed"),
		"email has already existed",
	)

	ErrCannotCreateAccount = common.NewCustomError(
		errors.New("can not create your account"),
		"can not create your account",
	)
	ErrCannotUpdateUser = common.NewCustomError(
		errors.New("can not update your account"),
		"can not update your account",
	)
	ErrAccountNotVerified = common.NewCustomError(
		errors.New("This account has not been verified!"),
		"This account has not been verified!")
)
