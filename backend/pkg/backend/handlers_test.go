package backend

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegisterHandler(t *testing.T) {
	router := gin.Default()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	if err != nil {
		t.Fail()
	}

	err = RegisterHandler(router, db)

	if err != nil {
		t.Fail()
	}

}
