package tables

import "gorm.io/gorm"

type Policy struct {
	gorm.Model

	Role       string `gorm:"type:text;not null;"`
	Resource   string `gorm:"type:text;not null;"`
	Permission string `gorm:"type:text;not null;check:Permissions in ('write','read','update','delete')"`
}
