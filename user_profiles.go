package main

import (
	"github.com/jinzhu/gorm"
)

type UserProfile struct {
	gorm.Model
	UserID    uint
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}
