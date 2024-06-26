package kernel

import (
	"encoding/json"
	"github.com/lirchis/goblog/internal/model"
	"github.com/lirchis/goblog/internal/repo"
	"github.com/lirchis/goblog/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// DictGet doc
// @Tags kernel
// @Summary 通过id获取单条字典
// @Param key path string true "key"
// @Success 200 {object} utils.Reply{data=model.Dict} "返回数据"
// @Router /api/dict/:key [get]
func DictVal(ctx echo.Context) error {
	key := ctx.Param("key")
	mod, has := repo.DictGet(key)
	if !has {
		return ctx.JSON(utils.Fail("不存在"))
	}
	return ctx.JSON(utils.Succ("succ", json.RawMessage(mod.Value)))
}

// DictGet doc
// @Tags kernel
// @Summary 通过id获取单条字典
// @Param key path string true "key"
// @Success 200 {object} utils.Reply{data=model.Dict} "返回数据"
// @Router /api/dict/get/:key [get]
func DictGet(ctx echo.Context) error {
	key := ctx.Param("key")
	mod, has := repo.DictGet(key)
	if !has {
		return ctx.JSON(utils.Fail("不存在"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// DictPage doc
// @Auth mgmt
// @Tags kernel
// @Summary 获取字典分页
// @Param query query model.Page true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Dict} "返回数据"
// @Router /api/dict/page [get]
func DictPage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	count := repo.DictCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := repo.DictPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, count))
}

// DictAdd doc
// @Auth mgmt
// @Tags kernel
// @Summary 添加字典
// @Param token query string true "token"
// @Param body body model.Dict true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/dict/add [post]
func DictAdd(ctx echo.Context) error {
	ipt := &model.Dict{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Check(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Created = time.Now().UnixMilli()
	ipt.Updated = ipt.Created
	err = repo.DictAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// DictEdit doc
// @Auth mgmt
// @Tags kernel
// @Summary 修改字典
// @Param token query string true "token"
// @Param body body model.Dict true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/dict/edit [post]
func DictEdit(ctx echo.Context) error {
	ipt := &model.Dict{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Check(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Updated = time.Now().UnixMilli()
	err = repo.DictEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// DictDrop doc
// @Auth mgmt
// @Tags kernel
// @Summary 通过key删除单条字典
// @Param key query string true "key"
// @Param token query string true "token"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/dict/drop [post]
func DictDrop(ctx echo.Context) error {
	ipt := &struct {
		Key string `json:"key"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := repo.DictGet(ipt.Key)
	if !has {
		return ctx.JSON(utils.Fail("不存在"))
	}
	if mod.Inner {
		return ctx.JSON(utils.ErrOpt("内置数据无法删除"))
	}
	err = repo.DictDrop(ipt.Key)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
