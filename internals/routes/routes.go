package routes

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type fn func(filterType string, value string, filter string, isLast bool) string

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func DefineRoutes(router *gin.Engine, db *sql.DB) {

	router.GET("/albums", func(ctx *gin.Context) {
		albums := getAllAlbums(db)
		ctx.JSON(200, map[string]any{"data": albums})
	})

	router.GET("/albums/:title", func(ctx *gin.Context) {
		albums := getAlbumsFilter(db, ctx)
		ctx.JSON(200, map[string]any{"data": albums})
	})

	router.POST("/albums", func(ctx *gin.Context) {
		albums := getAdvancedAlbumsFilter(db, ctx)
		ctx.JSON(200, map[string]any{"data": albums})
	})

	router.Run("localhost:8080")
}

func getAllAlbums(db *sql.DB) []Album {
	var albums []Album
	rows, _ := db.Query("SELECT * FROM album")

	defer rows.Close()
	for rows.Next() {
		var alb Album
		rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		albums = append(albums, alb)
	}

	return albums
}

func getAlbumsFilter(db *sql.DB, ctx *gin.Context) []Album {
	var albums []Album
	name := ctx.Param("title")
	rows, _ := db.Query("SELECT * FROM album where artist = ?", name)

	defer rows.Close()
	for rows.Next() {
		var alb Album
		rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		albums = append(albums, alb)
	}

	return albums
}

func getAdvancedAlbumsFilter(db *sql.DB, ctx *gin.Context) []Album {
	var jsonData = make(map[string]string)
	if err := ctx.BindJSON(&jsonData); err != nil {
		panic(err)
	}

	filter := "SELECT * FROM album WHERE"

	filterMap := map[string]fn{
		"title":  addFilter,
		"artist": addFilter,
	}

	count := 0
	for key, value := range jsonData {
		filter = filterMap[key](key, value, filter, (count+1) == len(jsonData))

		count++
	}

	var albums []Album
	rows, err := db.Query(filter)

	if err != nil {
		fmt.Println(filter)
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var alb Album
		rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		albums = append(albums, alb)
	}

	return albums
}

func addFilter(filterType string, value string, filter string, isLast bool) string {
	filter += " " + filterType + " LIKE('%" + value + "%')"
	if !isLast {
		filter += " AND"
	}

	return filter
}
