package main

import (
	"github.com/labstack/echo"
	"strconv"
)

type Paginate struct {
	Limit     int  `json:"limit"`
	Page      int  `json:"page"`
	Offset    int  `json:"-"`
	Count     int  `json:"count"`
	PageCount int  `json:"page_count"`
	NextPage  bool `json:"nextPage"`
	PrevPage  bool `json:"prevPage"`
}

func SetPaginateParams(paginate *Paginate, c echo.Context) {

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}

	if paginate.Count > 0 {

		//page_count
		var page_count int = paginate.Count / limit
		var page_count_m = paginate.Count % limit
		if page_count_m > 0 {
			page_count = page_count + 1
		}
		paginate.PageCount = page_count

		//NextPage
		if page < page_count {
			paginate.NextPage = true
		}

		//PrevPage
		if page > 1 {
			paginate.PrevPage = true
		}
	}

	paginate.Limit = limit
	paginate.Page = page
	paginate.Offset = offset

}
