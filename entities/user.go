package entities

type User struct {
	ID           uint64       `json:"id" gorm:"primary_key"`
	Username     string       `json:"username" gorm:"unique;size:120"`
	Name         string       `json:"name"`
	Key          string       `json:"key" gorm:"not null"`
	KeyVersion   uint64       `json:"key_version" gorm:"default:0;not null"`
	Level        uint64       `json:"level" gorm:"default:0;not null"`
	Limit        uint64       `json:"limit" gorm:"default:1000;not null"`   // Limite de armazenamento em bytes
	LimitUsed    uint64       `json:"limit_used" gorm:"default:0;not null"` // Limite usado em bytes
	Register     int64        `json:"register" gorm:"not null"`
	Change       int64        `json:"change"`
	AccessPasses []AccessPass `json:"access_pass,omitempty" gorm:"Foreignkey:UserID"`
}

type AccessPass struct {
	ID       uint64 `json:"id" gorm:"primary_key"`
	UserID   uint64 `json:"user_id" gorm:"not null"`
	Pass     string `json:"pass" gorm:"unique;not null"`
	Type     string `json:"type" gorm:"default:'mail';not null"`
	Register int64  `json:"register"`
}

func (u User) AvailableLimit() uint64 {
	return u.Limit - u.LimitUsed
}
