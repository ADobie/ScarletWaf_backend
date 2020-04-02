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
