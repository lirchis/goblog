package main

import (
	"context"
	conf "github.com/lirchis/goblog/configs"
	"github.com/lirchis/goblog/internal/repo"
	"github.com/lirchis/goblog/internal/router"
	"github.com/lirchis/goblog/internal/utils"
	"github.com/zxysilent/logs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	func main() {
//		fmt.Printf("hello world\n")
//		ctx := logs.TraceCtx(context.Background())
//		defer logs.Close()
//		conf.Init()
//		repo.Init(ctx)
//		defer repo.Close()
//		ipt := &model.IptId{Id: 102}
//		mod := &model.Post{Id: ipt.Id}
//		err := repo.PostGet(mod)
//		if err != nil {
//			fmt.Printf("失败")
//		}
//		fmt.Printf(mod.Markdown)
//	}
func main() {
	ctx := logs.TraceCtx(context.Background())
	defer logs.Close()
	logs.SetText()
	logs.Ctx(ctx).Info("app initializing")
	conf.Init()
	repo.Init(ctx)
	defer repo.Close()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logs.Info("app running")
	if !utils.FileExist("./static") {
		os.MkdirAll("./static", 0755)
	}
	if !utils.FileExist("./dist") {
		os.MkdirAll("./dist", 0755)
	}
	app := router.Init()
	// json.NewEncoder(os.Stdout).Encode(app.Routes())
	go func() {
		if err := app.Start(conf.Addr()); err != nil {
			if err != http.ErrServerClosed { //Normal exit
				logs.Ctx(ctx).Error("Server Shutdown:", err)
			}
		}
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		if err != http.ErrServerClosed { //Normal exit
			logs.Ctx(ctx).Error("Server Shutdown:", err)
		}
	}
	logs.Ctx(ctx).Info("app quitted")
}
