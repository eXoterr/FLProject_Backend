package category

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/handlers/payloads"
	"github.com/eXoterr/FLProject/internal/storage/models"
	"github.com/eXoterr/FLProject/internal/utils"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

const CATEGORIES_PER_PAGE = 20

func GetList(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()
		page := q.Get("page")
		if len(page) == 0 {
			render.JSON(w, r, utils.Error(w, "page is unspecified", 400))
			return
		}

		pageNum, err := strconv.Atoi(page)
		if err != nil {
			render.JSON(w, r, utils.Error(w, "unable to parse page number", 400))
			return
		}

		categories := []models.Category{}
		offset := (pageNum - 1) * CATEGORIES_PER_PAGE
		result := store.Offset(offset).Find(&categories)

		if result.Error != nil {
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		categoriesResp := make([]payloads.Category, len(categories))
		for idx, category := range categories {
			categoriesResp[idx] = payloads.Category{
				ID:   category.ID,
				Name: category.Name,
			}
		}

		render.JSON(w, r, utils.Success(w, categoriesResp, 200))
	}
}
