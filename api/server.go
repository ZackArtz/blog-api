package api

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zackartz/blog-api/models"
	"log"
	"os"
)

type Server struct {
	DB     *gorm.DB
	Router *fiber.App
}

func (s *Server) Initialize() {
	var err error

	DBURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	s.DB, err = gorm.Open("mysql", DBURI)
	if err != nil {
		log.Panicf("could not connect to DB %v", err)
	}

	s.DB.AutoMigrate(&models.User{}, &models.Article{})

	s.initializeRoutes()

	log.Fatal(s.Router.Listen(1337))
}
