package main

import (
	"github.com/gin-gonic/gin"
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
	//viper.SetConfigFile("./scarlet.backend.yaml")
	//err := viper.ReadInConfig()
	//if err != nil {
	//	tool.GetLogger().Fatal("缺少配置文件./scarlet.backend.toml")
	//}
	//viper.Set("mysql.addr","127.0.0.1:3306")
	//viper.Set("mysql.user","scarlet")
	//viper.Set("mysql.password","scarlet")
	//viper.Set("mysql.database","scarlet")
	//viper.Set("redis.addr","127.0.0.1:3306")
	//viper.Set("redis.db",0)
	r := gin.Default()
	r.Use(controller.JWT())
	r.POST("/user", controller.AddUser)
	r.POST("/login", controller.UserLogin)
	//r.POST("/user/password",controller.LoginRequired(),controller.UpdatePassword)
	// 按照规定 delete 不需要读消息体... 不知道GO是否会读

	userGroup := r.Group("/user")
	userGroup.Use(controller.LoginRequired()) // 一个基于JWT的SESSION管理中间件
	userGroup.PUT("/user/server", controller.UpdateServer)
	userGroup.GET("/user/server", controller.GetServers)
	userGroup.DELETE("/user/server", controller.DeleteServer)
	userGroup.POST("/user/uri", controller.AddURI)
	userGroup.DELETE("/user/uri", controller.DeleteURI)
	userGroup.POST("/user/rule", controller.GetRules)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := r.Run(":8080"); err != nil {
		tool.GetLogger().Fatal("Addres Already Used")
	}

}
