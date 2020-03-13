package api

import (
	"github.com/gin-gonic/gin"
	"huage.tech/mini/app/bean"
	"huage.tech/mini/app/dao"
	"huage.tech/mini/app/util"
	"strconv"
)

func EntryCreate(c *gin.Context) {
	var err error
	var form bean.Entry

	err = c.ShouldBind(&form)
	if err != nil {
		util.AbortNewResultErrorOfClient(c, err)
		return
	}

	r, err := dao.EntryCreate(form)
	if err != nil || r.ID == 0 {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}

	c.JSON(200, util.NewResultOKofWrite(r, 1))
}

func EntryList(c *gin.Context) {
	r, err := dao.EntryList()
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	c.JSON(200, util.NewResultOKofRead(r, len(r)))
	return
}

func EntryDelete(c *gin.Context) {
	EntryId := c.Param("id")
	id, err := strconv.ParseInt(EntryId, 10, 64)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	err = dao.EntryDelete(id)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	c.JSON(200, util.NewResultOKofWrite(nil, 1))
}

func EntryRead(c *gin.Context) {
	EntryId := c.Param("id")
	id, err := strconv.ParseInt(EntryId, 10, 64)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	r, err := dao.EntryRead(id)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	c.JSON(200, util.NewResultOKofRead(r, 1))
}

func EntryUpdate(c *gin.Context) {
	var err error
	var form bean.Entry
	EntryId := c.Param("id")
	id, err := strconv.ParseInt(EntryId, 10, 64)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}

	err = c.ShouldBind(&form)
	if err != nil {
		util.AbortNewResultErrorOfClient(c, err)
		return
	}
	form.ID = id
	r, err := dao.EntryUpdate(form)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}

	c.JSON(200, util.NewResultOKofWrite(r, 1))
}

func Entries(c *gin.Context) {
	v, err := dao.Entries()
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	tree := util.MenuTree(v)
	c.JSON(200, util.NewResultOKofRead(tree, 1))
}
