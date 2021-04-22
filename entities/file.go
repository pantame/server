package entities

type File struct {
	ID         uint64 `json:"id,omitempty" gorm:"primary_key"`
	UUID       string `json:"uuid" gorm:"unique;not null"`
	UserID     uint64 `json:"user_id" gorm:"not null"`
	Version    uint64 `json:"version" gorm:"default:1;not null"`
	Key        string `json:"key" gorm:"not null"`
	Metadata   string `json:"metadata" gorm:"not null"`
	Size       uint64 `json:"size" gorm:"not null"` // Tamanho do arquivo em bytes
	TotalParts uint64 `json:"total_parts" gorm:"not null"`
	PartsSent  uint64 `json:"parts_sent" gorm:"default:0;not null"`
	Register   int64  `json:"register" gorm:"not null"`
	Change     int64  `json:"change"`
}

func (f *File) OwnsAllParts() bool {
	if f.PartsSent == f.TotalParts {
		return true
	}
	return false
}

func (f *File) CalcTotalParts() uint64 {
	const partSizeLimit uint64 = 2500000

	total := f.Size / partSizeLimit
	re := f.Size % partSizeLimit

	if re != 0 {
		total += 1
	}

	return total
}
