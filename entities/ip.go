package entities

type IPData struct {
	ID            uint64  `json:"id" gorm:"primary_key"`
	IP            string  `json:"ip" gorm:"not null"`
	IPDate        string  `json:"ip_date" gorm:"unique;not null"`
	Date          string  `json:"date"` // 2020-02-20
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	TimeZone      string  `json:"timezone"`
	Currency      string  `json:"currency"`
	ISP           string  `json:"isp"`
	ORG           string  `json:"org"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
	Register      int64   `json:"register"`
}
