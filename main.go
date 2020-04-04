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

	// 按照规定 delete 不需要读消息体... 不知道GO是否会读
	// rule不存在主键，或者说需要多个条件。不太适合用纯Restful 修改API
	userGroup := r.Group("/user")
	userGroup.Use(controller.LoginRequired()) // 一个基于JWT的SESSION管理中间件

	userGroup.POST("/server/add", controller.AddServer)
	userGroup.POST("/server/get", controller.GetServers)
	userGroup.POST("/server/delete", controller.DeleteServer)

	userGroup.POST("/uri/add", controller.AddURI)
	userGroup.POST("/uri/delete", controller.DeleteURI)
	userGroup.POST("/uri/get", controller.GetURI)

	userGroup.POST("/rule/get", controller.GetRules)
	userGroup.POST("/rule/delete", controller.DeleteRule)
	userGroup.POST("/rule/add", controller.AddRule)

	userGroup.POST("/switch/waf", controller.WafStatus)
	userGroup.POST("/switch/change", controller.ChangeSwitch)
	userGroup.POST("/switch/get", controller.GetSwitch)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := r.Run(":8080"); err != nil {
		tool.GetLogger().Fatal("Address Already Used")
	}

}
