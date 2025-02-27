// Package routers defines application routes
// @APIVersion 1.0.0
// @Title OpenVPN API
// @Description REST API allows you to control and monitor your OpenVPN server
// @Contact adam.walach@gmail.com
// License Apache 2.0
// LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/vuonglequoc/openvpn-web-ui/controllers"
	"github.com/beego/beego"
)

func init() {
	beego.SetStaticPath("/swagger", "swagger")
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{}, "get,post:Login")
	beego.Router("/logout", &controllers.LoginController{}, "get:Logout")
	beego.Router("/profile", &controllers.ProfileController{})
	beego.Router("/settings", &controllers.SettingsController{})
	beego.Router("/ov/config", &controllers.OVConfigController{})
	beego.Router("/logs", &controllers.LogsController{})

	beego.Include(&controllers.CertificatesController{})

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/session",
			beego.NSInclude(
				&controllers.APISessionController{},
			),
		),
		beego.NSNamespace("/sysload",
			beego.NSInclude(
				&controllers.APISysloadController{},
			),
		),
		beego.NSNamespace("/signal",
			beego.NSInclude(
				&controllers.APISignalController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
