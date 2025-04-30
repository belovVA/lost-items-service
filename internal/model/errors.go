package model

import "github.com/pkg/errors"

var (
	ErrorUserNotFound     = errors.New("user not found")
	ErrorFailedBuildQuery = errors.New("failed Build Query")
	ErrorFailedCreateUser = errors.New("failed to Create User")
)
