package main

import (
	_ "newCHNTLDManager/internal/packed"
	"sync"

	"newCHNTLDManager/dns/service"
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
	var mLock = new(sync.Mutex)
	chnZone := new(zonefile.ChnZone)
	chnZone.Init()
	s := g.Server()

	//测试
	s.BindHandler("/QueryDNSRecord", func(r *ghttp.Request) {
		mLock.Lock()
		res, err := chnZone.QueryDNSRecord(r.GetBodyString())
		defer mLock.Unlock()
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

	s.BindHandler("/AddDNSRecord", func(r *ghttp.Request) {
		mLock.Lock()
		err := chnZone.AddDNSRecord(r.GetBodyString())
		defer mLock.Unlock()
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

	s.BindHandler("/DelDNSRecord", func(r *ghttp.Request) {
		mLock.Lock()
		err := chnZone.DelDNSRecord(r.GetBodyString())
		defer mLock.Unlock()
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

	s.BindHandler("/ReloadZone", func(r *ghttp.Request) {
		out, err := service.ReloadZone()
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"success": false,
				"msg":     err.Error(),
			})
		} else {
			r.Response.WriteJsonExit(g.Map{
				"success": true,
				"msg":     out,
			})
		}
	})

	s.BindHandler("/QueryDnsServiceStatus", func(r *ghttp.Request) {
		out, err := service.DnsServiceStatus()
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"success": false,
				"msg":     err.Error(),
			})
		} else {
			r.Response.WriteJsonExit(g.Map{
				"success": true,
				"msg":     out,
			})
		}
	})

	s.BindHandler("/RestartDnsService", func(r *ghttp.Request) {
		out, err := service.RestartDnsService()
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"success": false,
				"msg":     err.Error(),
			})
		} else {
			r.Response.WriteJsonExit(g.Map{
				"success": true,
				"msg":     out,
			})
		}
	})
	// s.BindHandler("/GetDomainRecord", func(r *ghttp.Request) {
	// 	mLock.Lock()
	// 	res, err := chnZone.GetDomainRecord(r.GetBodyString())
	// 	defer mLock.Unlock()
	// 	if err != nil {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": false,
	// 			"msg":     err.Error(),
	// 		})
	// 	} else {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success":        true,
	// 			"msg":            "ok",
	// 			"totalCount":     len(res),
	// 			"recordListJson": res,
	// 		})
	// 	}

	// })
	// s.BindHandler("/GetDomainMXRecord", func(r *ghttp.Request) {
	// 	mLock.Lock()
	// 	res, err := chnZone.GetDomainMXRecord(r.GetBodyString())
	// 	defer mLock.Unlock()
	// 	if err != nil {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": false,
	// 			"msg":     err.Error(),
	// 		})
	// 	} else {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success":        true,
	// 			"msg":            "ok",
	// 			"totalCount":     len(res),
	// 			"recordListJson": res,
	// 		})
	// 	}
	// })

	// s.BindHandler("/AddDomainRecord", func(r *ghttp.Request) {
	// 	mLock.Lock()
	// 	err := chnZone.AddRecord(r.GetBodyString())
	// 	defer mLock.Unlock()
	// 	if err != nil {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": false,
	// 			"msg":     err.Error(),
	// 		})
	// 	} else {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": true,
	// 			"msg":     "ok",
	// 		})
	// 	}
	// })

	// s.BindHandler("/AddDomainMXRecord", func(r *ghttp.Request) {
	// 	mLock.Lock()
	// 	err := chnZone.AddMXRecord(r.GetBodyString())
	// 	defer mLock.Unlock()
	// 	if err != nil {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": false,
	// 			"msg":     err.Error(),
	// 		})
	// 	} else {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": true,
	// 			"msg":     "ok",
	// 		})
	// 	}
	// })

	// s.BindHandler("/DelDomainRecord", func(r *ghttp.Request) {
	// 	mLock.Lock()
	// 	err := chnZone.DelRecord(r.GetBodyString())
	// 	defer mLock.Unlock()
	// 	if err != nil {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": false,
	// 			"msg":     err.Error(),
	// 		})
	// 	} else {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": true,
	// 			"msg":     "ok",
	// 		})
	// 	}
	// })
	// s.BindHandler("/ModifyDomainRecord", func(r *ghttp.Request) {
	// 	mLock.Lock()
	// 	err := chnZone.ModifyRecord(r.GetBodyString())
	// 	defer mLock.Unlock()
	// 	if err != nil {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": false,
	// 			"msg":     err.Error(),
	// 		})
	// 	} else {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": true,
	// 			"msg":     "ok",
	// 		})

	// 	}
	// })
	// s.BindHandler("/ModifyDomainMXRecord", func(r *ghttp.Request) {
	// 	mLock.Lock()
	// 	err := chnZone.ModifyMXRecord(r.GetBodyString())
	// 	defer mLock.Unlock()
	// 	if err != nil {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": false,
	// 			"msg":     err.Error(),
	// 		})
	// 	} else {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"success": true,
	// 			"msg":     "ok",
	// 		})

	// 	}
	// })

	s.SetPort(80)
	s.Run()
}
