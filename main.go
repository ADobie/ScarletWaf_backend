package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"scarlet/controller"
	_ "scarlet/docs"
	"scarlet/tool"
)

// @title Scarlet Backend
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @query.collection.format multi
func main() {
	//viper.SetConfigName("scarlet.backend")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath(".")
	viper.SetConfigFile("./scarlet.backend.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		tool.GetLogger().Fatal("缺少配置文件./scarlet.backend.toml")
	}
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.POST("/user", controller.AddUser)
	r.POST("/user/server", controller.LoginRequired(), controller.UpdateServer)
	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "index.tpl", nil)
	})
	r.POST("/login", controller.UserLogin)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := r.Run(":8080"); err != nil {
		tool.GetLogger(err)
	}

}
