package cmd

//"context"

//"github.com/gogf/gf/v2/os/gcmd"
//"newCHNTLDManager/dns/zonefile"
//"newCHNTLDManager/internal/controller/hello"

// func main() {
// 	s := g.Server()
// 	s.BindHandler("/GetAllDomainRecord", func(r *ghttp.Request) {
// 		r.Response.Write("哈喽世界！")
// 	})
// 	s.SetPort(80)
// 	s.Run()
// }

// var (
// 	Main = gcmd.Command{
// 		Name:  "main",
// 		Usage: "main",
// 		Brief: "start http server",
// 		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
// 			s := g.Server()
// 			s.Group("/", func(group *ghttp.RouterGroup) {
// 				group.Middleware(ghttp.MiddlewareHandlerResponse)
// 				group.Bind(
// 					hello.NewV1(),
// 				)
// 			})
// 			s.Run()
// 			println("start my app")
// 			// p := dnsFile.Person{Name: "zhangsan", Age: 22}
// 			// p.Run()

// 			// // 测试对象p的类型
// 			// fmt.Printf("p type is %T\n", p)
// 			// // 反射使用p的方法
// 			// v := reflect.ValueOf(p)
// 			// // 获取方法的数量
// 			// fmt.Printf("p has %d methods\n", v.NumMethod())
// 			// // 获取方法的名称
// 			// fmt.Printf("p's method 0 is: %s\n", v.Method(0).Type())
// 			// // 获取方法的参数和返回值
// 			// method := v.Method(0)
// 			// fmt.Printf("method name is %s\n", method.Type().Name())
// 			// fmt.Printf("method is %s\n", method.Type())
// 			// fmt.Printf("method's first parameter is %s\n", method.Type().In(0))
// 			// fmt.Printf("method's first output parameter is %s\n", method.Type().Out(0))

// 			chnZone := new(zonefile.ChnZone)
// 			chnZone.Init()
// 			//chnZone.PrintDefaultZoneFileList()
// 			//chnZone.AddRecord(`{"DomainName":"polo","TTL":"120","IN":"IN","Type":"A","Data":"192.168.1.1"}`)
// 			err = chnZone.AddRecord(`{"domainName":"ak47","TTL":"600","type":"A","data":"192.168.15.38"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}
// 			err = chnZone.AddRecord(`{"domainName":"polo","TTL":"600","type":"A9","data":"32768[86[21[4]111"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}
// 			err = chnZone.AddRecord(`{"domainName":"polo","TTL":"600","type":"MX","data":"abc.com"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}

// 			err = chnZone.AddRecord(`{"domainName":"poloTXT","TTL":"600","type":"TXT","data":"hello world"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}

// 			// err = chnZone.AddRecord(`{"domainName":"@","TTL":"600","type":"NS","data":"polo.gtld-servers.chn"}`)
// 			// if err != nil {
// 			// 	println(err.Error())
// 			// }

// 			err = chnZone.AddRecord(`{"domainName":"admin.abc","TTL":"600","type":"A","data":"192.168.1.100"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}

// 			err = chnZone.AddRecord(`{"domainName":"iccnea.xatu","TTL":"7200","type":"TXT","data":"cs.xatu.cn/iccnea-19/index.htm"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}

// 			err = chnZone.ModifyRecord(`{"domainName":"aa","TTL":"1000","type":"dfdf","data":"polo.com"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}

// 			// err = chnZone.DelRecord(`{"domainName":"@","TTL":"600","type":"NS","data":"polo.gtld-servers.chn"}`)
// 			// if err != nil {
// 			// 	println(err.Error())
// 			// }

// 			err = chnZone.DelRecord(`{"domainName":"*.boxes.shvpn","TTL":"600","type":"A","data":"202.170.218.74"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}

// 			req, err := chnZone.GetDomainRecord(`{"domainName":"polo","type":"A9"}`)
// 			if err != nil {
// 				println(err.Error())
// 			}
// 			println(req)
// 			println(len(req))

// 			allRes, err := chnZone.GetAllDomainRecord()
// 			if err != nil {
// 				println(err.Error())
// 			}
// 			println(allRes)
// 			println(len(allRes))
// 			chnZone.PrintRuntimeZoneFileList()
// 			//dnsFile.WriteDefaultZoneFile()
// 			return nil
// 		},
// 	}
// )
