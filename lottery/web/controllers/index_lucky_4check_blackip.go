package controllers

import (
	"lottery/models"
	"time"
)

func (c *IndexController) checkBlackip(ip string) (bool, *models.LtBlackip) {
	info := c.ServiceBlackip.GetByIp(ip)
	if info == nil || info.Ip == "" {
		return true, nil
	}
	if info.Blacktime > int(time.Now().Unix()) {
		return false, info //ip黑名单存在，并且还在黑名单有效期内
	}
	return true, info
}
