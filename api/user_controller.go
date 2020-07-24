package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/zackartz/blog-api/auth"
	"github.com/zackartz/blog-api/models"
	"github.com/zackartz/blog-api/utils"
	"github.com/zackartz/blog-api/utils/formaterror"
	"net/http"
	"strconv"
)

func (s *Server) GetAllUsers(ctx *fiber.Ctx) {
	user := models.User{}
	users, err := user.GetAllUsers(s.DB)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	utils.JSON(ctx, http.StatusOK, users)
}

func (s *Server) CreateUser(ctx *fiber.Ctx) {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Printf("%s", user.Password)
	err := user.Validate("")
	if err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	pwBytes, err := models.Hash(user.Password)
	user.Password = string(pwBytes)
	if err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.CreateUser(s.DB)
	fmt.Printf("%s", user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		utils.Error(ctx, http.StatusInternalServerError, formattedError)
		return
	}
	ctx.Set("Location", fmt.Sprintf("%s%s/%d", ctx.Hostname(), ctx.OriginalURL(), userCreated.ID))
	utils.JSON(ctx, http.StatusCreated, userCreated)
}

func (s *Server) GetUser(ctx *fiber.Ctx) {
	id := ctx.Params("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.JSON(ctx, http.StatusBadRequest, err)
		return
	}
	user := &models.User{}
	user, err = user.GetUserByID(s.DB, uint32(uid))
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	Return(ctx, user)
}

func (s *Server) GetMe(ctx *fiber.Ctx) {
	id, err := auth.ExtractTokenID(ctx)
	if err != nil {
		utils.JSON(ctx, http.StatusUnauthorized, err)
		return
	}
	user := &models.User{}
	user, err = user.GetUserByID(s.DB, id)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	Return(ctx, user)
}

func (s *Server) UpdateUser(ctx *fiber.Ctx) {
	uid, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	tokenId, err := auth.ExtractTokenID(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	if tokenId != uint32(uid) {
		utils.Error(ctx, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateUser(s.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		utils.Error(ctx, http.StatusInternalServerError, formattedError)
		return
	}
	utils.JSON(ctx, http.StatusOK, updatedUser)
}

func (s *Server) DeleteUser(ctx *fiber.Ctx) {
	user := new(models.User)

	uid, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	tokenId, err := auth.ExtractTokenID(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	if tokenId != 0 && tokenId != uint32(uid) {
		utils.Error(ctx, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteUserByID(s.DB, uint32(uid))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Set("Entity", fmt.Sprintf("%d", uid))
	utils.JSON(ctx, http.StatusNoContent, "")
}

func Return(ctx *fiber.Ctx, user *models.User) {
	user.Password = ""
	if user.ShowEmail == false {
		user.Email = ""
	}
	utils.JSON(ctx, http.StatusOK, user)
}
