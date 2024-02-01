package main

import (
	_ "newCHNTLDManager/internal/packed"

	"newCHNTLDManager/dns/zonefile"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	//"github.com/gogf/gf/v2/os/gctx"
	//"newCHNTLDManager/internal/cmd"
)

// func main() {
// 	cmd.Main.Run(gctx.GetInitCtx())
// }

func main() {
	chnZone := new(zonefile.ChnZone)
	chnZone.Init()
	s := g.Server()
	s.BindHandler("/GetAllDomainRecord", func(r *ghttp.Request) {
		res, err := chnZone.GetAllDomainRecord()
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"success": false,
				"msg":     err.Error(),
			})
		} else {
			r.Response.WriteJsonExit(g.Map{
				"success":        true,
				"msg":            "ok",
				"totalCount":     len(res),
				"recordListJson": res,
			})
		}
	})
	s.BindHandler("/AddDomainRecord", func(r *ghttp.Request) {
		err := chnZone.AddRecord(r.GetBodyString())
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"success": false,
				"msg":     err.Error(),
			})
		} else {
			r.Response.WriteJsonExit(g.Map{
				"success": true,
				"msg":     "ok",
			})
		}
	})
	s.BindHandler("/DelDomainRecord", func(r *ghttp.Request) {
		err := chnZone.DelRecord(r.GetBodyString())
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"success": false,
				"msg":     err.Error(),
			})
		} else {
			r.Response.WriteJsonExit(g.Map{
				"success": true,
				"msg":     "ok",
			})
		}
	})
	s.BindHandler("/ModifyDomainRecord", func(r *ghttp.Request) {
		err := chnZone.ModifyRecord(r.GetBodyString())
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"success": false,
				"msg":     err.Error(),
			})
		} else {
			r.Response.WriteJsonExit(g.Map{
				"success": true,
				"msg":     "ok",
			})
		}
	})
	s.SetPort(80)
	s.Run()
}
