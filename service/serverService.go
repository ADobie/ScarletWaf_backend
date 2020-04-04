package service

import (
	"github.com/jinzhu/gorm"
	"scarlet/common"
)

type ServerService struct{}

func NewServerService() *ServerService {
	return new(ServerService)
}

func (s *ServerService) GetByUserID(userId uint) (servers []common.Server) {
	mysqlClient.Where("user_id = ?", userId).Find(&servers)
	return
}

func (s *ServerService) Delete(serverID uint) {
	mysqlClient.Delete(common.Server{
		Model: gorm.Model{ID: serverID},
	})
	return
}

func (s *ServerService) Own(userId uint, serverID uint) bool {
	var count int
	mysqlClient.Table("servers").Where("user_id = ? and id =  ?", userId, serverID).Count(&count)
	if count >= 1 {
		return true
	} else {
		return false
	}
}

func (s *ServerService) Get(serverID uint) common.Server {
	var server common.Server
	server.ID = serverID
	mysqlClient.First(&server)
	return server
}
