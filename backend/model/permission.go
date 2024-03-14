package model

type Permission struct {
	PermissionId uint64 `json:"permission-id" gorm:"type=int;primary_key"`
	Privilege    string `json:"privilege" gorm:"varchar(32); not null"`
}
