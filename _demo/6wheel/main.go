/**
 * 大转盘活动
 * 每一次转动抽奖。后端计算出这次的抽奖中奖情况。并返回对应的奖品信息
 * 线程不安全，因为获奖概率低，并发更新库存冲突很少能出现，不容易出现线程安全问题
 * 压力测试：wrk -t10 -c100 -d5 "http://localhost:8080/prize"
 */
package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex = sync.Mutex{}

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

type Prate struct {
	Rate  int // 万分之N的中奖概率
	Total int // 总数量限制,0 表示无限数量
	CodeA int //中奖概率起始编码(包含)
	CodeB int //中奖概率结束编码(包含)
	Left  int //剩余数
}

//奖品列表
var prizeList []string = []string{
	"一等奖，火星单程票",
	"二等奖，南极之旅",
	"三等奖，iphone",
	"", //无中奖
}

// 奖品的中奖概率设置，与上面的prizeList 对应设置
var rateList []Prate = []Prate{
	{
		Rate:  1,
		Total: 1,
		CodeA: 0,
		CodeB: 0,
		Left:  1,
	},
	{
		Rate:  2,
		Total: 2,
		CodeA: 1,
		CodeB: 2,
		Left:  2,
	},
	{
		Rate:  5,
		Total: 10,
		CodeA: 3,
		CodeB: 5,
		Left:  10,
	},
	{
		Rate:  100,
		Total: 0,
		CodeA: 0,
		CodeB: 9999,
		Left:  0,
	},
}

// GET http://localhost:8080/
func (c lotteryController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return fmt.Sprintf("大转盘奖品列表：<br/> %s", strings.Join(prizeList, "<br/>\n"))
}

func (c *lotteryController) GetDebug() string {
	return fmt.Sprintf("获奖概率:%v\n", rateList)
}

func (c *lotteryController) GetPrize() string {
	//第一步 抽奖 根据随机数匹配奖品
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	code := r.Intn(10000)

	var myprize string
	var prizeRate *Prate
	//  从奖品列表中匹配是否中奖
	for i, prize := range prizeList {
		rate := &rateList[i]
		if code >= rate.CodeA && code <= rate.CodeB {
			myprize = prize
			prizeRate = rate
			break
		}
	}
	if myprize == "" {
		myprize = "很遗憾，再来一次吧"
		return myprize
	}
	// 第二步 中奖后 开始发奖
	if prizeRate.Total == 0 {
		return myprize
	} else if prizeRate.Left > 0 {
		mu.Lock()
		prizeRate.Left -= 1
		mu.Unlock()
		return myprize
	} else {
		myprize = "很遗憾，再来一次吧"
		return myprize
	}
}
