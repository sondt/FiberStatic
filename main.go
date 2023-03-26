package main

import (
	"FiberStatic/services"
	imageHelper "FiberStatic/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

var publicFolder = "public"

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
	})

	//config cache client
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("no-cache") == "true"
		},
		Expiration:   10 * 60 * 1000, // 10 minutes
		CacheControl: true,
	}))

	//config logger
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${locals:requestid} ${status} - ${method} ${path} - ${ip} - ${latency}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Ho_Chi_Minh",
	}))

	//config rewrite image
	app.Get("/:type/images/upload/*", func(c *fiber.Ctx) error {
		return imageHelper.ProcessRewriteImage(c)
	})

	//weather api
	app.Get("/weather/:cityCode", func(c *fiber.Ctx) error {
		content := services.GetWeather(c.Params("cityCode"))
		return c.JSON(content)
		//return c.SendString(content)
	})

	//config static folder
	app.Static("/", "./"+publicFolder, fiber.Static{
		Compress:  true,
		ByteRange: true,
		Index:     "index.html",
		Browse:    true,
		MaxAge:    3600,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! 👋!")
	})
	port := "80"
	if os.Getenv("ASPNETCORE_PORT") != "" { // get enviroment variable that set by ACNM
		port = os.Getenv("ASPNETCORE_PORT")
	}
	//103.27.237.189

	err := app.Listen(":" + port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
