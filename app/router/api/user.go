package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"huage.tech/mini/app/bean"
	"huage.tech/mini/app/config"
	"huage.tech/mini/app/dao"
	"huage.tech/mini/app/util"
	"strconv"
)

func UserCreate(c *gin.Context) {
	var err error
	var form bean.User

	err = c.ShouldBind(&form)
	if err != nil {
		util.AbortNewResultErrorOfClient(c, err)
		return
	}
	TenantId, _ := c.MustGet("TENANT_ID").(int64)
	form.TenantId = TenantId

	form.Password = util.Md5(config.JwtSecret + form.Password)
	r, err := dao.UserCreate(form)
	if err != nil || r.ID == 0 {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}

	c.JSON(200, util.NewResultOKofWrite(r, 1))
}

func UserList(c *gin.Context) {
	account := c.DefaultQuery("account", "")

	roleId, err := strconv.ParseInt(c.DefaultQuery("role", ""), 10, 64)

	if err != nil {
		roleId = 0
	}

	TenantId, _ := c.MustGet("TENANT_ID").(int64)
	orgId, err := strconv.ParseInt(c.DefaultQuery("org", ""), 10, 64)
	if err != nil {
		orgId = 0
	}

	pag := c.DefaultQuery("page", "1")
	lim := c.DefaultQuery("limit", "20")
	offset, limit, err := util.PageLimit(pag, lim)
	if err != nil {
		util.AbortNewResultErrorOfClient(c, err)
		return
	}

	r, count, err := dao.UserList(TenantId, account, roleId, orgId, offset, limit)

	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	c.JSON(200, util.NewResultOKofRead(r, count))
	return
}

func UserDelete(c *gin.Context) {
	UserId := c.Param("id")
	id, err := strconv.ParseInt(UserId, 10, 64)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	tokenUid, _ := c.MustGet("UID").(int64)
	if tokenUid == id {
		util.AbortNewResultErrorOfClient(c, errors.New("不允许删除自己"))
		return
	}
	TenantId, _ := c.MustGet("TENANT_ID").(int64)
	err = dao.UserDelete(TenantId, id)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	c.JSON(200, util.NewResultOKofWrite(nil, 1))
}

func UserRead(c *gin.Context) {
	UserId := c.Param("id")
	id, err := strconv.ParseInt(UserId, 10, 64)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	TenantId, _ := c.MustGet("TENANT_ID").(int64)
	r, err := dao.UserRead(id, TenantId)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}
	c.JSON(200, util.NewResultOKofRead(r, 1))
}

func UserUpdate(c *gin.Context) {
	var err error
	var form bean.User
	UserId := c.Param("id")
	id, err := strconv.ParseInt(UserId, 10, 64)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}

	err = c.ShouldBind(&form)
	if err != nil {
		util.AbortNewResultErrorOfClient(c, err)
		return
	}
	TenantId, _ := c.MustGet("TENANT_ID").(int64)
	form.TenantId = TenantId
	form.ID = id
	if len(form.Password) > 0 {
		form.Password = util.Md5(config.JwtSecret + form.Password)
	}

	r, err := dao.UserUpdate(form)
	if err != nil {
		util.AbortNewResultErrorOfServer(c, err)
		return
	}

	c.JSON(200, util.NewResultOKofWrite(r, 1))
}

func ChangPassword(c *gin.Context) {
	old := c.DefaultPostForm("old", "")
	new := c.DefaultPostForm("new", "")
	new1 := c.DefaultPostForm("repeat", "")
	if old == "" {
		util.AbortNewResultErrorOfClient(
			c, errors.New("请输入旧密码"))
		return
	}
	if new != new1 {
		util.AbortNewResultErrorOfClient(
			c, errors.New("两次新密码输入不一致"))
		return
	}
	TenantId, _ := c.MustGet("TENANT_ID").(int64)
	if uId, ok := c.Get("UID"); ok {
		if uid, ok := uId.(int64); ok {
			u, err := dao.UserRead(TenantId, uid)
			if err != nil {
				util.AbortNewResultErrorOfClient(
					nil, errors.New("不存在此用户"))
				return
			}
			if u.Password != util.Md5(config.JwtSecret+old) {
				util.AbortNewResultErrorOfClient(
					nil, errors.New("旧密码不正确"))
				return
			}
			u.Password = util.Md5(config.JwtSecret + new)
			r, err := dao.UserUpdate(u)
			if err != nil {
				util.AbortNewResultErrorOfServer(c, err)
				return
			}

			c.JSON(200, util.NewResultOKofWrite(r, 1))
			return
		}

	}
}
