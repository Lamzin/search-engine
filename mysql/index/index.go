package sqlindex

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Row struct {
	gorm.Model
	lexeme    string `gorm:"index:lexeme"`
	docID     uint32
	frequency uint32
}

type SQLIndex struct {
	dbPath string
	db     *gorm.DB
}

func NewSQLIndex(dbPath string) (*SQLIndex, error) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Row{})
	return &SQLIndex{dbPath, db}, nil
}

func (index *SQLIndex) Add(lexeme string, docID uint32, frequency uint32) error {
	index.db.Create(&Row{
		lexeme:    lexeme,
		docID:     docID,
		frequency: frequency,
	})
	return index.db.Error
}

func (index *SQLIndex) Close() error {
	return index.db.Close()
}
