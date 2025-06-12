package common

import (
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
)

type PageResponse struct {
	Limit      int         `query:"limit" json:"limit"`
	Page       int         `query:"page" json:"page"`
	Sort       string      `query:"sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Items      interface{} `json:"items"`
}

func (p *PageResponse) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *PageResponse) GetLimit() int {
	if p.Limit > 100 {
		p.Limit = 100
	} else if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *PageResponse) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func NewPageResponse(model interface{}, request *http.Request, db *gorm.DB) *PageResponse {
	var pageResponse PageResponse
	query := request.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	var totalRows int64
	db.Model(model).Count(&totalRows)
	pageResponse.Limit = limit
	pageResponse.Page = page
	pageResponse.TotalRows = totalRows

	totalPages := int(math.Ceil(float64(totalRows) / float64(pageResponse.GetLimit())))
	pageResponse.TotalPages = totalPages

	return &pageResponse
}

func (p *PageResponse) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit())
	}
}
