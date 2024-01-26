package backend

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Book struct {
	Model
	BookID    uint    `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Language  string  `json:"language"`
	Text      string  `gorm:"type:text" json:"text"`
	ImageURL  string  `json:"imageBook"`
	CRank     float64 `json:"crank"`
	Occurance uint    `json:"occurence"`
}

type IndexedBook struct {
	Model
	BookID              uint   `json:"id"`
	Title               string `json:"title"`
	WorldOccurancesJSON string `gorm:"type:text" json:"wordOcc"`
}

type JaccardNeighbors struct {
	Model
	BookID        uint   `json:"id"`
	NeighborsJSON string `gorm:"type:text" json:"neightbors"`
}
