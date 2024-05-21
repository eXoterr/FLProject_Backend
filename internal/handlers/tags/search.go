package tags

import (
	"fmt"
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

const TAG_RESULTS_PER_PAGE = 20

func Search(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
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

		query := q.Get("q")
		if len(query) == 0 {
			render.JSON(w, r, utils.Error(w, "query is unspecified", 400))
			return
		}

		tags := []models.Tag{}
		offset := (pageNum - 1) * TAG_RESULTS_PER_PAGE
		result := store.Offset(offset).Where("LOWER(name) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", query)).Find(&tags)

		if result.Error != nil {
			logger.Error(result.Error.Error())
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		tagsResp := make([]payloads.Tag, len(tags))
		for idx, tag := range tags {
			tagsResp[idx] = payloads.Tag{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}

		render.JSON(w, r, utils.Success(w, tagsResp, 200))
	}
}
