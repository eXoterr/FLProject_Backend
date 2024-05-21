package order

import (
	"encoding/json"
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

const ORDERS_PER_PAGE = 20

func Search(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()
		page := q.Get("page")
		if len(page) == 0 {
			page = "0"
		}

		pageNum, err := strconv.Atoi(page)
		if err != nil {
			render.JSON(w, r, utils.Error(w, "unable to parse page number", 400))
			return
		}

		query := q.Get("q")

		queryTags := q.Get("tags")
		tags := payloads.Tags{}
		if len(queryTags) > 0 {
			err = json.Unmarshal([]byte(queryTags), &tags)
			if err != nil {
				logger.Error(err.Error())
				render.JSON(w, r, utils.InternalError(w))
				return
			}
		}

		categoryId := q.Get("cat")
		if len(categoryId) == 0 {
			categoryId = "0"
		}
		categoryNum, err := strconv.Atoi(categoryId)
		if err != nil {
			render.JSON(w, r, utils.Error(w, "unable to parse category id", 400))
			return
		}

		orders := []models.Order{}

		offset := (pageNum - 1) * ORDERS_PER_PAGE
		results := store.Offset(offset).Where(
			"LOWER(title) LIKE LOWER(?)",
			fmt.Sprintf("%%%s%%", query),
		)
		if results.Error != nil {
			logger.Error("aboba")
			logger.Error(results.Error.Error())
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		if categoryNum > 0 {
			results = results.Find(&orders, "category_id = ?", categoryNum)
		} else {
			results = results.Find(&orders)
		}

		/// PUPUPU

		if err != nil {
			logger.Error(results.Error.Error())
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		category := &models.Category{}
		result := store.First(category, "id = ?", categoryNum)
		if result.Error != nil {
			logger.Error(result.Error.Error())
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		ordersResp := make([]payloads.Order, len(orders))
		for idx, order := range orders {
			ordersResp[idx] = payloads.Order{
				Title:       order.Title,
				CategoryID:  categoryNum,
				Category:    category.Name,
				Description: order.Description,
				BudgetMin:   order.BudgetMin,
				BudgetMax:   order.BudgetMax,
				Deadline:    order.Deadline,
				Tags:        tags,
			}
		}

		render.JSON(w, r, utils.Success(w, ordersResp, 200))
	}
}
