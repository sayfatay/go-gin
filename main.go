package main

import (
	"finance/car-finance/back-end/configs"
	"finance/car-finance/back-end/db"
	"finance/car-finance/back-end/entities"
	"finance/car-finance/back-end/routes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	gintemplate "github.com/foolin/gin-template"
	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
)

func init() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
	// conf, _ := config.NewConfig()
	// src := env.NewSource()
	// conf.Load(src)
}

func main() {
	// Setup Web Service
	// Create RESTful handler (using Gin)

	mysqlDefaultConn := db.CreateMySQLConnection(&configs.MySQLConn{
		Config: mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
				os.Getenv("MYSQL_USERNAME"),
				os.Getenv("MYSQL_PASSWORD"),
				os.Getenv("MYSQL_HOST"),
				os.Getenv("MYSQL_PORT"),
				os.Getenv("MYSQL_DBNAME"),
			),
		},
	})
	mysqlDB, _ := mysqlDefaultConn.DB()
	mysqlDefaultConn.AutoMigrate(entities.User{})
	defer mysqlDB.Close()

	router := gin.Default()
	router.Use(core())
	router.Static("/assets", "public")
	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "views",          //template root path
		Extension: ".tpl",           //file extension
		Master:    "layouts/master", //master layout file
		// Partials:  []string{"partials/head"}, //partial files
		Funcs: template.FuncMap{
			"time": func(t string) string {
				return time.Now().Format(t)
			},
			"debug": func() bool {
				return os.Getenv("APP_DEBUG") == "true"
			},
			"getEnv": func(v string) string {
				return os.Getenv(v)
			},
		},
		DisableCache: false, //if disable cache, auto reload template file for debug.
	})
	router.Use(ginCors.Default())
	router.MaxMultipartMemory = 10 << 20 // 8 MiB

	// router.Handle("/", apmhttp.Wrap(routes.WebRouter("/", router)))
	// router.Handle("/", apmhttp.Wrap(routes.WebRouter("/", router)))
	router = routes.SetupRouter(router)
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// Run service
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func core() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if v := os.Getenv("CORS_ALLOW_ORIGIN_HEADER"); strings.TrimSpace(v) != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", fmt.Sprintf("%s, %s", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With", v))
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		// c.Request.Header.Del("Origin")
		c.Next()
	}
}
