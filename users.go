package main

import (
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Uuid      string     `json:"uuid"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	LastName  string     `json:"last_name"`
	FirstName string     `json:"first_name"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
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

		// CPU数取得
		cpus := runtime.NumCPU()

		//Find
		uf := make(chan []User, cpus)
		go func() {
			var users []User
			Db.
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
		user := User{}
		if err := Db.Where("id = ?", id).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return c.JSON(http.StatusOK, user)
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

		user.Username = c.FormValue("username")
		user.LastName = c.FormValue("last_name")
		user.FirstName = c.FormValue("first_name")
		user.Uuid = uuid.NewV4().String()

		//パスワード
		password := c.FormValue("password")
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user.Password = string(hash)

		// Db.NewRecord(user)
		Db.Create(&user)

		return c.JSON(http.StatusOK, user)
	}
}

/**
 * EditUser
 */
func EditUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		var Db, _ = DbConnect()

		id, _ := strconv.Atoi(c.Param("id"))

		user := User{}
		if err := Db.Where("id = ?", id).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		user.Username = c.FormValue("username")
		user.LastName = c.FormValue("last_name")
		user.FirstName = c.FormValue("first_name")

		//パスワード
		password := c.FormValue("password")
		if password != "" {
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			user.Password = string(hash)
		}

		Db.Save(&user)

		return c.JSON(http.StatusOK, user)
	}
}
