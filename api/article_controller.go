package api

import (
	"github.com/gofiber/fiber"
	"github.com/zackartz/blog-api/models"
	"github.com/zackartz/blog-api/utils"
	"net/http"
	"strconv"
)

type ArticleResponse struct {
	models.Article

	*models.User `json:"author,omitempty"`
}

func (s *Server) NewArticleResponse(article models.Article) ArticleResponse {
	user := &models.User{}
	user, err := user.GetUserByID(s.DB, article.AuthorID)
	if err != nil {
		return ArticleResponse{}
	}
	article.AuthorID = 0
	resp := ArticleResponse{Article: article}
	resp.User = user
	return resp
}

func (s *Server) GetArticles(ctx *fiber.Ctx) {
	article := models.Article{}
	articles, err := article.GetAllArticles(s.DB)
	user := models.User{}
	if err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	var resp []ArticleResponse
	for _, article := range articles {
		res := ArticleResponse{Article: article}
		res.User, err = user.GetUserByID(s.DB, article.AuthorID)
		if err != nil {
			utils.Error(ctx, http.StatusUnprocessableEntity, err)
			return
		}
		resp = append(resp, res)
	}
	utils.JSON(ctx, http.StatusOK, resp)
}

func (s *Server) CreateArticle(ctx *fiber.Ctx) {
	article := &models.Article{}
	if err := ctx.BodyParser(article); err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	article, err := article.CreateArticle(s.DB)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	utils.JSON(ctx, http.StatusCreated, s.NewArticleResponse(*article))
}

func (s *Server) GetArticleByID(ctx *fiber.Ctx) {
	id := ctx.Params("id")
	aid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	article := &models.Article{}
	article, err = article.GetArticleByID(s.DB, uint32(aid))
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	utils.JSON(ctx, http.StatusCreated, s.NewArticleResponse(*article))
}

func (s *Server) GetArticleBySlug(ctx *fiber.Ctx) {
	slug := ctx.Params("slug")
	article := &models.Article{}
	article, err := article.GetArticleBySlug(s.DB, slug)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	utils.JSON(ctx, http.StatusCreated, s.NewArticleResponse(*article))
}

func (s *Server) UpdateArticle(ctx *fiber.Ctx) {
	id := ctx.Params("id")
	aid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	article := &models.Article{}
	article, err = article.GetArticleByID(s.DB, uint32(aid))
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	data := &models.Article{}
	if err := ctx.BodyParser(data); err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	updateArticle, err := data.UpdateArticle(s.DB, article.ID)
	if err != nil {
		utils.Error(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	utils.JSON(ctx, http.StatusOK, s.NewArticleResponse(updateArticle))
}

func (s *Server) DeleteArticle(ctx *fiber.Ctx) {
	var err error
	id := ctx.Params("id")
	aid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	article := &models.Article{}
	article, err = article.GetArticleByID(s.DB, uint32(aid))
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err)
		return
	}
	res, err := article.DeleteArticle(s.DB, article.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(ctx, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: string(res),
	})
}
