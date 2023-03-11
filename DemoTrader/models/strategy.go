package models

import (
	// "fmt"
	// "log"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Renko struct {
	gorm.Model
	Name 			string	`json:"name" db:"name"`
	Coin 			string	`json:"coin" db:"coin"`
	BrickSize  	float64	`json:"brick_size" db:"brick_size"`
}