package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/zackartz/blog-api/auth"
	"github.com/zackartz/blog-api/models"
	"github.com/zackartz/blog-api/utils"
	"github.com/zackartz/blog-api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (s *Server) Login(ctx *fiber.Ctx) {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Printf("%s", user.Password)
	user.Prepare()
	err := user.Validate("login")
	if err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := s.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		utils.Error(ctx, http.StatusUnprocessableEntity, formattedError)
		return
	}
	utils.JSON(ctx, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (s *Server) SignIn(email, password string) (string, error) {
	var err error
	user := &models.User{}
	err = s.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", errors.New("email")
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Printf("%s, %s : %v", password, user.Password, err)
		return "", errors.New("password")
	}
	return auth.CreateToken(user.ID)
}
