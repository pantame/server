package entities

type Session struct {
	ID         uint64 `json:"id" gorm:"primary_key"`
	UserID     uint64 `json:"user_id" gorm:"not null"`
	Token      string `json:"token,omitempty" gorm:"unique;not null"`
	AccessPass string `json:"access_pass"`
	Ip         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Active     bool   `json:"active" gorm:"default:true;not null"`
	Register   int64  `json:"register"`
	Change     int64  `json:"change"`
}
