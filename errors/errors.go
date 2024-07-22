package errors

import "errors"

var (
	// データベースエラー
	ErrDb = errors.New("db error")

	// リソースが見つからなかったことを示すエラー
	ErrNotFound = errors.New("not found")
)
