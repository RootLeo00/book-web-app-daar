package backend

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title     string
	Author    string
	Language  string
	Text      string `gorm:"type:text"`
	ImageURL  string
	CRank     float64
	Occurance uint
}

type IndexedBook struct {
	gorm.Model
	Title               string
	WorldOccurancesJSON string `gorm:"type:text"`
}

type JaccardNeighbors struct {
	gorm.Model
	NeighborsJSON string `gorm:"type:text"`
}
