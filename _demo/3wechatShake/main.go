/**
 * 微信摇一摇
 * 基础功能 ：
 * /lucky 抽奖接口
 * 压力测试：wrk -t10 -c10 -d5 http://localhost:8080/lucky(-t：线程数 -c：连接数 -d：持续时间)
 */
package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// 奖品类型

const (
	giftTypeCoin      = iota //虚拟币
	giftTypeCoupon           //不同券
	giftTypeCouponFix        //相同券
	giftTypeRealSmall        //实物小奖
	giftTypeRealLarge        //实物大奖
)

type gift struct {
	id       int      //奖品ID
	name     string   // 奖品名称
	pic      string   //奖品图片
	link     string   //奖品连接
	gType    int      //奖品类型
	data     string   //奖品数据 （特定的配置信息）
	dataList []string //奖品数据集合（不同的优惠券编码）
	total    int      //总数 0 不限量
	left     int      // 剩余数量
	inuse    bool     //是否使用中
	rate     int      //中间概率 万分之N，0-9999
	rateMin  int      //小于等于最小中奖编码
	rateMax  int      //小于中奖编码
}

// 最大中间号码
const rateMax = 10000

var logger *log.Logger

// 奖品列表
var giftList []*gift

type lotteryController struct {
	Ctx iris.Context
}

// 初始化日志
func initLog() {
	f, _ := os.Create("/var/log/lottery_demo.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

// 初始化奖品列表
func initGift() {
	giftList = make([]*gift, 5)
	g1 := gift{
		id:       1,
		name:     "手机大奖",
		pic:      "",
		link:     "",
		gType:    giftTypeRealLarge,
		data:     "",
		dataList: nil,
		total:    20000,
		left:     20000,
		inuse:    true,
		rate:     10000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[0] = &g1
	g2 := gift{
		id:       2,
		name:     "充电器",
		pic:      "",
		link:     "",
		gType:    giftTypeRealSmall,
		data:     "",
		dataList: nil,
		total:    0,
		left:     0,
		inuse:    false,
		rate:     10,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[1] = &g2
	g3 := gift{
		id:       3,
		name:     "优惠券满200减50",
		pic:      "",
		link:     "",
		gType:    giftTypeCouponFix,
		data:     "",
		dataList: nil,
		total:    0,
		left:     0,
		inuse:    false,
		rate:     500,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[2] = &g3
	g4 := gift{
		id:       4,
		name:     "直降50元",
		pic:      "",
		link:     "",
		gType:    giftTypeCoupon,
		data:     "",
		dataList: []string{"c01", "c02", "c03", "c04", "c05"},
		total:    0,
		left:     0,
		inuse:    false,
		rate:     100,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[3] = &g4
	g5 := gift{
		id:       5,
		name:     "金币",
		pic:      "",
		link:     "",
		gType:    giftTypeCoin,
		data:     "",
		dataList: nil,
		total:    0,
		left:     0,
		inuse:    false,
		rate:     5000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[4] = &g5
	// 整理数据 中奖区间数据
	rateStart := 0
	for _, data := range giftList {
		if !data.inuse {
			continue
		}
		data.rateMin = rateStart
		data.rateMax = rateStart + data.rate
		if data.rateMax >= rateMax {
			data.rateMax = rateMax
			rateStart = 0
		} else {
			rateStart += data.rate
		}
	}
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	initLog()
	initGift()
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

// 奖品数量信息 GET http://localhost:8080
func (c *lotteryController) Get() string {
	count, total := 0, 0
	for _, data := range giftList {
		if data.inuse && (data.total == 0 || (data.total > 0 && data.left > 0)) {
			count++
			total += data.left
		}
	}
	return fmt.Sprintf("当前有效奖品种类：%d；限量奖品数：%d", count, total)
}

func luckyCode() int32 {
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Int31n(int32(rateMax))
	return code
}

// 抽奖 GET http://localhost:8080/lucky
func (c *lotteryController) GetLucky() map[string]interface{} {
	code := luckyCode()
	ok := false
	result := make(map[string]interface{})
	result["success"] = ok
	for _, data := range giftList {
		if !data.inuse || (data.total > 0 && data.left <= 0) {
			continue
		}
		if data.rateMin <= int(code) && data.rateMax > int(code) {
			// 中奖了，抽奖编码在奖品范围内
			//开奖
			sendData := ""
			switch data.gType {
			case giftTypeCoin:
				ok, sendData = sendCoin(data)
			case giftTypeCoupon:
				ok, sendData = sendCoupon(data)
			case giftTypeCouponFix:
				ok, sendData = sendCouponFix(data)
			case giftTypeRealSmall:
				ok, sendData = sendRealSmall(data)
			case giftTypeRealLarge:
				ok, sendData = sendRealLarge(data)
			}
			if ok { //中奖
				//生成中奖记录
				saveLuckyData(code, data.id, data.name, data.link, sendData, data.left)
				result["success"] = ok
				result["id"] = data.id
				result["name"] = data.name
				result["link"] = data.link
				result["data"] = sendData
				break
			}
		}
	}
	return result
}

func sendCoin(data *gift) (bool, string) {
	if data.total == 0 { // 数量无限
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 不同值优惠券
func sendCoupon(data *gift) (bool, string) {
	if data.left > 0 {
		left := data.left - 1
		data.left = left
		return true, data.dataList[left]
	} else {
		return false, "奖品已发完"
	}
}

func sendCouponFix(data *gift) (bool, string) {
	if data.left > 0 {
		left := data.left - 1
		data.left = left
		return true, data.dataList[left]
	} else {
		return false, "奖品已发完"
	}
}

func sendRealSmall(data *gift) (bool, string) {
	if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

func sendRealLarge(data *gift) (bool, string) {
	if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 记录用户获奖信息
func saveLuckyData(code int32, id int, name, link, sendData string, left int) {
	logger.Printf("lucky,code=%d,gift=%d,name=%s,link=%s,data=%s,left=%d\n", code, id, name, link, sendData, left)
}
