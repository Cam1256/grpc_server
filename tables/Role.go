package tables

import "gorm.io/gorm"

type RoleDef struct {
	gorm.Model

	User   string `gorm:"type:text;not null;"`
	Name   string `gorm:"type:text;not null;"`
	Tenant string `gorm:"type:text;not null;"`
}
