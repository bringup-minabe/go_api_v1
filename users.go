package main

import (
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type User struct {
	Id        uint      `json:"id"`
	Uuid      string    `json:"uuid"`
	Username  string    `json:"username"`
	Created   time.Time `json:"created"`
	LastName  string    `json:"last_name"`
	FirstName string    `json:"first_name"`
}

type Res struct {
	Data     []User
	Paginate Paginate
}

/**
 * GetUsers
 */
func GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {

		var Db, _ = DbConnect()

		var paginate Paginate

		// set paginate params
		SetPaginateParams(&paginate, c)

		/**
		 * Select
		 * 並列処理で通常検索とカウント用検索をおこなう
		 */

		// set Where
		where := map[string]interface{}{
			"is_deleted": 0,
		}

		// CPU数取得
		cpus := runtime.NumCPU()

		//Find
		uf := make(chan []User, cpus)
		go func() {
			var users []User
			Db.
				Where(where).
				Limit(paginate.Limit).
				Offset(paginate.Offset).
				Order("id desc").
				Find(&users)
			uf <- users
		}()

		//Count
		uc := make(chan int, cpus)
		go func() {
			count := 0
			Db.
				Model(User{}).
				Where(where).
				Count(&count)
			uc <- count
		}()

		users := <-uf
		users_count := <-uc
		paginate.Count = users_count

		// set paginate params
		SetPaginateParams(&paginate, c)

		/**
		 * set Responce
		 */
		var res Res
		res.Data = users
		res.Paginate = paginate

		return c.JSON(http.StatusOK, res)
	}
}

/**
 * GetUser
 */
func GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		var Db, _ = DbConnect()

		id, _ := strconv.Atoi(c.Param("id"))

		/**
		 * Select
		 */
		var user []User
		if err := Db.Where("id = ?", id).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		} else {
			if len(user) == 0 {
				return echo.NewHTTPError(http.StatusNotFound)
			} else {
				return c.JSON(http.StatusOK, user)
			}
		}
	}
}

/**
 * AddUser
 */
func AddUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		var Db, _ = DbConnect()

		user := User{}

		user.Username = c.FormValue("name")
		user.LastName = c.FormValue("last_name")
		user.FirstName = c.FormValue("first_name")
		user.Created = time.Now()
		user.Uuid = uuid.NewV4().String()

		// Db.NewRecord(user)
		Db.Create(&user)

		jsonMap := map[string]string{}

		return c.JSON(http.StatusOK, jsonMap)
	}
}
