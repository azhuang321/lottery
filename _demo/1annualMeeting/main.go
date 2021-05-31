/**
 * 年会抽奖实现
 * curl http://localhost:8080
 * curl --data "users=zhangsan,lisi,zhaowu,wangmazi" http://localhost:8080/import
 * curl http://localhost:8080/lucky
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

var userList []string

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
	//userList = make([]string, 0)
	userList = []string{}
	app.Run(iris.Addr(":8080"))
}

func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("当前总共参加抽奖人数：%d\n", count)
}

// POST http://localhost:8080/import
// params: users string
func (c *lotteryController) PostImport() string {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	count1 := len(users)
	for _, u := range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	return fmt.Sprintf("需导入用户数：%d；成功导入用户数：%d\n", count1, count2)

}

//GET http://localhost/lucky
func (c *lotteryController) GetLucky() string {
	count := len(userList)
	fmt.Println(count)
	if count == 0 {
		return fmt.Sprintf("已无参与用户，请先导入")
	} else if count == 1 {
		user := userList[0]
		userList = []string{}
		return fmt.Sprintf("当前中奖用户：%s；剩余用户：0；", user)
	} else {
		seed := time.Now().UnixNano()                                //随机数
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count)) //产生0~count的随机int32数
		user := userList[index]
		userList = append(userList[0:index], userList[index+1:]...)
		return fmt.Sprintf("当前中奖用户：%s；剩余用户数：%d\n", user, count-1)
	}
}
