package service

import (
	"scarlet/common"
)

type UriService struct{}

func NewURIService() *UriService {
	return new(UriService)
}

func (u *UriService) Add(uri common.URI) {
	mysqlClient.Create(&uri)
	return
}

func (u *UriService) Own(uriId uint, serverId uint) bool {
	var num int
	mysqlClient.Table("uris").Where("id = ? and server_id = ?", uriId, serverId).Count(&num)
	if num >= 1 {
		return true
	} else {
		return false
	}
}

func (u *UriService) Delete(uri common.URI) {
	mysqlClient.Delete(uri)
	return
}

func (u *UriService) Get(uriId uint) common.URI {
	var uri common.URI
	uri.ID = uriId
	mysqlClient.First(&uri)
	return uri
}

func (u *UriService) GetByServerID(serverID uint) (uris []common.URI) {
	mysqlClient.Where("server_id = ?", serverID).Find(&uris)
	return
}

// 通过serverId 和 path 检查是否已存在
func (u *UriService) Exist(uri common.URI) bool {
	var num int
	mysqlClient.Debug().Table("uris").Where("server_id = ? and path = ?", uri.ServerID, uri.Path).Count(&num)
	if num >= 1 {
		return true
	} else {
		return false
	}
}
