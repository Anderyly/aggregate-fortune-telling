package main

import (
	"gin/ay"
	"gin/routers"
	"gin/service"
	"gin/task"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"net/http"
)

var (
	r *gin.Engine
)

func main() {

	ay.Yaml = ay.InitConfig()
	ay.Sql()
	go ay.WatchConf()

	//models.S()
	//return

	// 开启定时任务
	c := cron.New()
	//c.AddFunc("45 2 * * *", task.InviteAmount)
	//c.AddFunc("@every 1s", task.InviteAmount)
	c.AddFunc("@every 2h10m10s", task.InviteAmount)
	c.Start()

	r = gin.Default()
	gin.SetMode(gin.DebugMode)
	r = service.Set(r)
	r.StaticFS("/static/", http.Dir("./static"))
	r = routers.GinRouter(r)

	err := r.Run(":8080")
	if err != nil {
		panic(err.Error())
	}
	select {}
}
