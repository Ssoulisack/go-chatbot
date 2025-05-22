package repositories

import "gorm.io/gorm"

type MessageRepo interface {
}

type MessageRepoImpl struct {
	db *gorm.DB
}


func NewMessageRepo(db *gorm.DB) MessageRepo {
	return &MessageRepoImpl{
		db: db,
	}
}