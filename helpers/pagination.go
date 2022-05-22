package helpers

import (
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Meta struct {
	Page *int `name=page" json:"page,omitempty"`
	// PerPage    *wrapperspb.Int32Value `name=per_page,json=perPage" json:"per_page,omitempty"`
	PerPage    *int `name=per_page,json=perPage" json:"per_page,omitempty"`
	PageCount  *int `name=page_count,json=pageCount" json:"page_count,omitempty"`
	TotalCount *int `name=total_count,json=totalCount" json:"total_count,omitempty"`
	// Links      map[string]string      `name=links" json:"links,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type (
	// FetchQuery fetch query options
	FetchQuery struct {
		Type  string
		Query string
		Args  []interface{}
	}

	// PagingConfig paging config
	PagingConfig struct {
		DB         *gorm.DB
		FetchQuery []*FetchQuery
		Page       int
		PerPage    int
		OrderBy    []*PagingOrderBy
		ShowSQL    bool
		All        bool
	}

	PagingOrderBy struct {
		SortBy   string
		SortType string
	}

	// Paginator paging response
	Paginator struct {
		TotalRecord int         `json:"total_record"`
		TotalPage   int         `json:"total_page"`
		Records     interface{} `json:"data"`
		Offset      int         `json:"offset"`
		PerPage     int         `json:"PerPage"`
		Page        int         `json:"page"`
		PrevPage    int         `json:"prev_page"`
		NextPage    int         `json:"next_page"`
	}
)

// Paging query data with paging
func Paging(p *PagingConfig, result interface{}) (*Paginator, error) {
	db := p.DB
	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = 10
	}

	var paginator Paginator
	var count int64
	var offset int

	if err := db.Model(result).Count(&count).Error; err != nil {
		return nil, err
	}

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.PerPage
	}

	if len(p.FetchQuery) > 0 {
		for _, v := range p.FetchQuery {
			switch queryType := v.Type; queryType {
			case "select":
				db = db.Select(v.Query, v.Args...)
			case "join":
				db = db.Joins(v.Query, v.Args...)
			case "order":
				db = db.Order(v.Query)
			default:
			}
		}
	}

	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(
				clause.OrderByColumn{
					Column: clause.Column{
						Name: strings.ReplaceAll(o.SortBy, "`", ""),
					},
					Desc: validateSortType(o.SortType),
				},
			)
		}
	}

	if p.All {
		if err := db.Find(result).Error; err != nil {
			return nil, err
		}
		totalRec := int(count)
		paginator := &Paginator{
			TotalRecord: totalRec,
			Records:     result,
			Page:        1,
			Offset:      0,
			PerPage:     totalRec,
			TotalPage:   1,
		}
		return paginator, nil
	}

	if err := db.Limit(p.PerPage).Offset(offset).Find(result).Error; err != nil {
		return nil, err
	}

	paginator.TotalRecord = int(count)
	paginator.Records = result
	paginator.Page = p.Page
	paginator.Offset = offset
	paginator.PerPage = p.PerPage
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.PerPage)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator, nil
}

func validateSortType(sortType string) bool {
	switch strings.ToLower(sortType) {
	case "desc":
		return true
	default:
		return false
	}
}

func GeneratePagingOrder(args ...string) *PagingOrderBy {
	var sortByDefault = "created_at"
	var sortTypeDefault = "desc"
	switch len(args) {
	case 4:
		sortTypeDefault = args[3]
		fallthrough
	case 3:
		sortByDefault = args[2]
		fallthrough
	case 2:
		if args[1] != "" {
			sortTypeDefault = args[1]
		}
		fallthrough
	case 1:
		if args[0] != "" {
			sortByDefault = args[0]
		}
	}
	return &PagingOrderBy{
		SortBy:   sortByDefault,
		SortType: sortTypeDefault,
	}
}

// ToProtobuf parse paging to meta response
func (h *Paginator) ToProtobuf() *Meta {

	return &Meta{
		Page:       &h.Page, // &wrappers.Int32Value{Value: int32(h.Page)},
		PerPage:    &h.PerPage,
		PageCount:  &h.TotalPage,
		TotalCount: &h.TotalRecord,
	}
}
