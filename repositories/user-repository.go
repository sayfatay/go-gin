package repositories

import (
	"context"
	"finance/car-finance/back-end/db"
	"finance/car-finance/back-end/entities"
	"finance/car-finance/back-end/helpers"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	mysql db.MySQL
}

func (r *UserRepository) QueryBinding(ct context.Context, req *entities.RequestUser) (*gorm.DB, error) {
	qr := r.mysql.DB().Model(new(entities.RequestUser))
	if v := req.ID; v != 0 {
		qr = qr.Where("id = ?", v)
	}

	if v := req.Keyword; v != "" {
		qr = qr.Where("lower(first_name) like ? or lower(last_name) like ?", fmt.Sprintf("%%%s%%", v), fmt.Sprintf("%%%s%%", v))
	}

	return qr, nil
}

func (r *UserRepository) GetUsers(ct context.Context, req *entities.RequestUser) ([]*entities.User, *helpers.Paginator, error) {
	var models []*entities.User
	qr, errQuery := r.QueryBinding(ct, req)
	if errQuery != nil {
		// bootstrap.Logger(ct).Error(errQuery)
		return nil, nil, errQuery
	}

	res, err := helpers.Paging(&helpers.PagingConfig{
		DB:      qr,
		Page:    int(req.Page),
		PerPage: int(req.PerPage),
		OrderBy: []*helpers.PagingOrderBy{helpers.GeneratePagingOrder(req.SortBy, req.SortType)},
	}, &models)
	if err != nil {
		return nil, nil, err
	}

	return models, res, nil
}
func (r *UserRepository) GetUser(ct context.Context) ([]string, error) {
	model := []*entities.User{}
	list := []string{}
	if err := r.mysql.DB().Model(new(entities.User)).Find(&model).Error; err != nil {
		// .Where("acked = FALSE and age(now(), created_at) < '30 minutes'")
		return nil, err
	}
	for _, v := range model {
		list = append(list, v.FirstName)
	}
	return list, nil
}

func (r *UserRepository) Create(ct context.Context, data *entities.User) (*entities.User, error) {
	tx := r.mysql.DB().Begin()
	if err := tx.Create(&data).Error; err != nil {
		// bootstrap.Logger(ct).Error(err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		// bootstrap.Logger(ct).Error(err)
		tx.Rollback()
		return nil, err
	}
	return data, nil
}
