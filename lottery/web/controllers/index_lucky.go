package controllers

import (
	"lottery/comm"
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
	// 4.验证IP今日的参与次数
	// 5.验证IP白名单
	// 6.验证用户黑名单
	// 7.获得抽奖编码
	// 8.匹配奖品是否中奖
	// 9.有限制奖品发放
	// 10.
	return rs
}
