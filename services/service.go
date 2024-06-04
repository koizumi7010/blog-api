package services

import "database/sql"

type MyAppService struct {
	db *sql.DB
}

// コンストラクタを定義
func NewMyAppService(db *sql.DB) *MyAppService {
	return &MyAppService{db: db}
}
