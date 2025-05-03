package model

import "github.com/pkg/errors"

var (
	ErrorUserNotFound = errors.New("user not found")
	ErrorBuildQuery   = errors.New("failed Build Query")
	ErrorScanRows     = errors.New("failed Scan Rows")
	ErrorExecuteQuery = errors.New("failed Execute Query")
)
