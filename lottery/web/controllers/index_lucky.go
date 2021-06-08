package controllers

import (
	"fmt"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/web/utils"
)

// http://localhost:8080/lucky
func (c *IndexController) GetLucky() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	// 1. 验证登录用户
	loginuser := comm.GetLoginUser(c.Ctx.Request())
	if loginuser == nil || loginuser.Uid < 1 {
		rs["code"] = 101
		rs["msg"] = "请先登录，再来抽奖"
		return rs
	}

	// 2.用户抽奖分布式锁定
	ok := utils.LockLucky(loginuser.Uid)
	if ok {
		defer utils.UnlockLucky(loginuser.Uid)
	} else {
		rs["code"] = 102
		rs["msg"] = "正在抽奖中，请稍后重试"
		return rs
	}
	// 3.验证用户今日参与次数
	ok = c.checkUserday(loginuser.Uid)
	if !ok {
		rs["code"] = 103
		rs["msg"] = "今日的抽奖次数已用完，明天再来吧"
		return rs
	}
	// 4.验证IP今日的参与次数
	ip := comm.ClientIP(c.Ctx.Request())
	ipDayNum := utils.IncrIpLuckyNum(ip)
	if ipDayNum > conf.IpLimitMax {
		rs["code"] = 104
		rs["msg"] = "相同IP参与次数过多，明天再来参与吧"
		return rs
	}
	limitBlack := false // 黑名单
	if ipDayNum > conf.IpPrizeMax {
		limitBlack = true
	}
	// 5.验证IP白名单
	var blackipInfo *models.LtBlackip
	if !limitBlack {
		ok, blackipInfo = c.checkBlackip(ip)
		if !ok {
			fmt.Println("黑名单中的IP", ip, limitBlack)
			limitBlack = true
		}
	}
	// 6.验证用户黑名单
	var userInfo *models.LtUser
	if !limitBlack {
		ok, userInfo = c.checkBlackUser(loginuser.Uid)
		if !ok {
			fmt.Println("黑名单中的用户", loginuser.Uid, limitBlack)
			limitBlack = true
		}
	}
	// 7.获得抽奖编码
	prizeCode := comm.Random(10000)
	// 8.匹配奖品是否中奖
	prizeGift := c.prize(prizeCode, limitBlack)
	if prizeGift == nil || prizeGift.PrizeNum < 0 || (prizeGift.PrizeNum > 0 && prizeGift.LeftNum <= 0) {
		rs["code"] = 205
		rs["msg"] = "很遗憾，没有中奖，请下次再试"
		return rs
	}
	// 9.有限制奖品发放
	// 10.
	return rs
}
