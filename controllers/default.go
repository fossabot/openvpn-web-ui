package controllers

import (
	"github.com/vuonglequoc/openvpn-web-ui/lib"
	"github.com/vuonglequoc/openvpn-web-ui/models"
	"github.com/beego/beego"

	mi "github.com/vuonglequoc/go-openvpn/server/mi"
)

type MainController struct {
	BaseController
}

func (c *MainController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Status",
	}
}

func (c *MainController) Get() {
	c.Data["sysinfo"] = lib.GetSystemInfo()
	lib.Dump(lib.GetSystemInfo())
	client := mi.NewClient(models.GlobalCfg.MINetwork, models.GlobalCfg.MIAddress)
	status, err := client.GetStatus()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovstatus"] = status
	}
	lib.Dump(status)

	version, err := client.GetVersion()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovversion"] = version.OpenVPN
	}
	lib.Dump(version)

	pid, err := client.GetPid()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovpid"] = pid
	}
	lib.Dump(pid)

	loadStats, err := client.GetLoadStats()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovstats"] = loadStats
	}
	lib.Dump(loadStats)

	c.TplName = "index.html"
}
