package main

import (
	"context"
	"fmt"
	conf "github.com/lirchis/goblog/configs"
	"github.com/lirchis/goblog/internal/model"
	"github.com/lirchis/goblog/internal/repo"
	"github.com/zxysilent/logs"
)

func main() {
	fmt.Printf("hello world\n")
	ctx := logs.TraceCtx(context.Background())
	defer logs.Close()
	conf.Init()
	repo.Init(ctx)
	defer repo.Close()
	ipt := &model.IptId{Id: 102}
	mod := &model.Post{Id: ipt.Id}
	err := repo.PostGet(mod)
	if err != nil {
		fmt.Printf("失败")
	}
	fmt.Printf(mod.Markdown)
}
