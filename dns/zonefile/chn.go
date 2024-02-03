package zonefile

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// 创建一个zone文件对象
type ChnZone struct {
	// zone文件List
	// 初始化默认zone文件content
	defaultZoneFileList *glist.List
	// 运行时zone文件content
	runtimeZoneFileList *glist.List
}

type dnsRecord struct {
	DomainName string `json:"domainName"`
	TTL        string `json:"ttl,omitempty"`
	IN         string `json:"-"`
	Type       string `json:"type"`
	Priority   string `json:"priority,omitempty"`
	Data       string `json:"data"`
}

// type domainRecord struct {
// 	DomainName string `json:"domainName"`
// 	TTL        string `json:"ttl"`
// 	IN         string `json:"-"`
// 	Type       string `json:"type"`
// 	Data       string `json:"data"`
// }

// type mxRecord struct {
// 	DomainName string `json:"domainName"`
// 	TTL        string `json:"ttl"`
// 	IN         string `json:"-"`
// 	Type       string `json:"type"`
// 	Priority   string `json:"priority"`
// 	Data       string `json:"data"`
// }

// type reqRecord struct {
// 	DomainName string `json:"domainName"`
// 	Type       string `json:"type"`
// }

// type resRecord struct {
// 	DomainName string `json:"domainName"`
// 	TTL        string `json:"ttl"`
// 	Type       string `json:"type"`
// 	Data       string `json:"data"`
// }

// type resMXRecord struct {
// 	DomainName string `json:"domainName"`
// 	TTL        string `json:"ttl"`
// 	Type       string `json:"type"`
// 	Priority   string `json:"priority"`
// 	Data       string `json:"data"`
// }

func (p *ChnZone) Init() {
	fmt.Println("init ChnZone...")
	p.defaultZoneFileList = glist.New()
	p.runtimeZoneFileList = glist.New()

	// 填充默认的zone文件defaultZoneFileList
	p.initDefaultZoneFileList(p.defaultZoneFileList)

	// 读取chn.zone文件，填充druntimeZoneFileList
	p.readZoneContentFromFile("/var/named/chn.zone", p.runtimeZoneFileList)
}

func (p *ChnZone) initDefaultZoneFileList(defaultZoneFileList *glist.List) {

	defaultZoneFileList.PushBack("$ORIGIN chn.\n")
	defaultZoneFileList.PushBack("$TTL 120\n")
	defaultZoneFileList.PushBack("@ IN SOA a.gtld-servers.chn. master.hostname.com. (\n")
	defaultZoneFileList.PushBack("		2024010101 ; serial\n")
	defaultZoneFileList.PushBack("		3600       ; refresh (1 hour)\n")
	defaultZoneFileList.PushBack("		600        ; retry (10 minutes)\n")
	defaultZoneFileList.PushBack("		604800     ; expire (1 week)\n")
	defaultZoneFileList.PushBack("		120      ; minimum (2 minutes)\n")
	defaultZoneFileList.PushBack("	)\n\n")

	defaultZoneFileList.PushBack("; Nameservers\n")
	defaultZoneFileList.PushBack("@ IN NS a.gtld-servers.chn.\n\n")

	defaultZoneFileList.PushBack("; Mailservers\n\n")
	defaultZoneFileList.PushBack("; Reverse DNS Records (PTR)\n\n")
	defaultZoneFileList.PushBack("; CNAME\n\n")
	defaultZoneFileList.PushBack("; TXT\n\n")
	defaultZoneFileList.PushBack("; HOST RECORDS\n\n")
}

func (p *ChnZone) readZoneContentFromFile(filePath string, runtimeZoneFileList *glist.List) error {

	// 读取chn.zone文件，填充defaultZoneFileList
	err := gfile.ReadLines(filePath, func(s string) error {
		runtimeZoneFileList.PushBack(s)
		return nil
	})
	return err
}

func (p *ChnZone) GetDefaultZoneFileList() *glist.List {
	return p.defaultZoneFileList
}

func (p *ChnZone) GetRuntimeZoneFileList() *glist.List {
	return p.runtimeZoneFileList
}

func (p *ChnZone) PrintDefaultZoneFileList() {
	for e := p.defaultZoneFileList.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func (p *ChnZone) PrintRuntimeZoneFileList() {
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func (p *ChnZone) WriteDefaultZoneFile() {
	strContent := ""
	for e := p.defaultZoneFileList.Front(); e != nil; e = e.Next() {
		strContent += e.Value.(string)
	}
	fmt.Println(strContent)
	// 把strContent写入文件
	err := os.WriteFile("testFile.zone", []byte(strContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Default zone file written successfully.")
}

func (p *ChnZone) AddDNSRecord(jsonRecord string) error {
	// 反序列化jsonRecord
	var record dnsRecord
	err := json.Unmarshal([]byte(jsonRecord), &record)
	if err != nil {
		fmt.Println("Error unmarshal jsonRecord:", err)
		return err
	}
	//增加DNS记录时，需要输入的信息包括：域名、TTL、类型、优先级、数据，其中优先级是MX记录特有的
	//检查输入的数据是否合法
	if gstr.Trim(record.DomainName) == "" || gstr.Trim(record.TTL) == "" || gstr.Trim(record.Type) == "" || gstr.Trim(record.Data) == "" {
		return fmt.Errorf("域名、TTL、类型、数据不能为空")
	}
	if record.Type == "MX" && record.Priority == "" {
		return fmt.Errorf("MX记录必须指定优先级")
	}
	//检查输入的数据是否合法
	if record.Type != "MX" && record.Type != "A" && record.Type != "A9" && record.Type != "NS" && record.Type != "PTR" && record.Type != "CNAME" && record.Type != "TXT" {
		return fmt.Errorf("不支持的类型")
	}
	err = p.checkTTL(record.TTL)
	if err != nil {
		return err
	}
	if record.Type == "MX" {
		err = p.checkPriority(record.Priority)
		if err != nil {
			return err
		}
	}
	if record.Type == "A" {
		err = p.checkIPv4Address(record.Data)
		if err != nil {
			return err
		}
	}
	if record.Type == "A9" {
		err = p.checkIPv9Address(record.Data)
		if err != nil {
			return err
		}
	}

	//先检查是否已经存在相同的记录
	if p.findRecord(record) {
		return fmt.Errorf("已存在相同的记录")
	}

	//增加DNS记录
	switch record.Type {
	case "NS":
		err = p.addNSRecord(record)
	case "MX":
		err = p.addMXRecord(record)
	case "PTR":
		err = p.addPTRRecord(record)
	case "CNAME":
		err = p.addCNAMERecord(record)
	case "TXT":
		err = p.addTXTRecord(record)
	case "A":
		err = p.addDomainRecord(record)
	case "A9":
		err = p.addDomainRecord(record)
	default:
		err = fmt.Errorf("不支持的类型")
	}
	// 递增serial
	err = p.incrementSerial()
	if err != nil {
		fmt.Println("Error increment serial:", err)
		return err
	}
	err = p.WriteZoneFile()

	return err
}

// func (p *ChnZone) AddRecord(jsonRecord string) error {
// 	if err := p.checkJaonRecord(jsonRecord); err != nil {
// 		return err
// 	}
// 	// 反序列化jsonRecord
// 	var record domainRecord
// 	err := json.Unmarshal([]byte(jsonRecord), &record)
// 	if err != nil {
// 		fmt.Println("Error unmarshal jsonRecord:", err)
// 		return err
// 	}

// 	// 填充record.IN
// 	record.IN = "IN"

// 	err = p.findRecord(record)
// 	if err != nil {
// 		return err
// 	}

// 	switch record.Type {
// 	case "NS":
// 		err = p.addNSRecord(record)
// 	case "MX":
// 		err = fmt.Errorf("请使用AddMXRecord")
// 	case "PTR":
// 		err = p.addPTRRecord(record)
// 	case "CNAME":
// 		err = p.addCNAMERecord(record)
// 	case "TXT":
// 		err = p.addTXTRecord(record)
// 	case "A":
// 		err = p.addDomainRecord(record)
// 	case "A9":
// 		err = p.addDomainRecord(record)
// 	default:
// 		err = fmt.Errorf("不支持的类型")
// 	}

// 	if err != nil {
// 		//fmt.Println("Error add record:", err)
// 		return err
// 	}

// 	// 递增serial
// 	err = p.incrementSerial()
// 	if err != nil {
// 		fmt.Println("Error increment serial:", err)
// 		return err
// 	}
// 	err = p.WriteZoneFile()
// 	return err
// }

func (p *ChnZone) DelDNSRecord(jsonRecord string) error {
	// 反序列化jsonRecord
	var record dnsRecord
	err := json.Unmarshal([]byte(jsonRecord), &record)
	if err != nil {
		fmt.Println("Error unmarshal jsonRecord:", err)
		return err
	}
	//删除DNS记录时，需要输入的信息包括：域名、类型、数据
	//检查输入的数据是否合法
	if gstr.Trim(record.DomainName) == "" || gstr.Trim(record.Type) == "" || gstr.Trim(record.Data) == "" {
		return fmt.Errorf("域名、类型、数据不能为空")
	}

	//检查输入的数据是否合法
	if record.Type != "MX" && record.Type != "A" && record.Type != "A9" && record.Type != "NS" && record.Type != "PTR" && record.Type != "CNAME" && record.Type != "TXT" {
		return fmt.Errorf("不支持的类型")
	}

	err = p.findDNSRecordAndDelete(record)
	if err != nil {
		return err
	}
	// 递增serial
	err = p.incrementSerial()
	if err != nil {
		fmt.Println("Error increment serial:", err)
		return err
	}
	err = p.WriteZoneFile()
	return err
}

// func (p *ChnZone) ModifyRecord(jsonRecord string) error {
// 	if err := p.checkJaonRecord(jsonRecord); err != nil {
// 		return err
// 	}
// 	var strContent string
// 	// 反序列化jsonRecord
// 	var record domainRecord
// 	err := json.Unmarshal([]byte(jsonRecord), &record)
// 	if err != nil {
// 		fmt.Println("Error unmarshal jsonRecord:", err)
// 		return err
// 	}

// 	// 填充record.IN
// 	record.IN = "IN"

// 	// 匹配record.type, recrod.domainName与runtimeZoneFileList中的记录，找到后修改
// 	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 		items := gstr.Split(e.Value.(string), " ")
// 		if gstr.Trim(items[0]) == record.DomainName {
// 			if gstr.Trim(items[3]) == record.Type {
// 				if record.Type == "TXT" {
// 					strContent = record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + strconv.Quote(record.Data)
// 				} else {
// 					strContent = record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data
// 				}
// 				e.Value = strContent
// 				// 递增serial
// 				err = p.incrementSerial()
// 				if err != nil {
// 					fmt.Println("Error increment serial:", err)
// 					return err
// 				}
// 				err = p.WriteZoneFile()
// 				return err
// 			}
// 		}
// 	}
// 	return fmt.Errorf("ModifyRecord... not found record")
// }

// func (p *ChnZone) ModifyMXRecord(jsonRecord string) error {
// 	if err := p.checkJaonRecord(jsonRecord); err != nil {
// 		return err
// 	}
// 	var strContent string
// 	// 反序列化jsonRecord
// 	var record mxRecord
// 	err := json.Unmarshal([]byte(jsonRecord), &record)
// 	if err != nil {
// 		fmt.Println("Error unmarshal jsonRecord:", err)
// 		return err
// 	}

// 	if record.Type != "MX" {
// 		return fmt.Errorf("此方法仅支持MX类型记录")
// 	}

// 	if gstr.Trim(record.Priority) == "" {
// 		return fmt.Errorf("MX记录必须指定优先级")
// 	}

// 	// 填充record.IN
// 	record.IN = "IN"

// 	// 匹配record.type, recrod.domainName与runtimeZoneFileList中的记录，找到后修改
// 	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 		items := gstr.Split(e.Value.(string), " ")
// 		if gstr.Trim(items[0]) == record.DomainName {
// 			if gstr.Trim(items[3]) == record.Type {

// 				strContent = record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Priority + " " + record.Data

// 				e.Value = strContent
// 				// 递增serial
// 				err = p.incrementSerial()
// 				if err != nil {
// 					fmt.Println("Error increment serial:", err)
// 					return err
// 				}
// 				err = p.WriteZoneFile()
// 				return err
// 			}
// 		}
// 	}
// 	return fmt.Errorf("ModifyMXRecord... not found record")
// }

// 测试
func (p *ChnZone) QueryDNSRecord(jsonReq string) ([]dnsRecord, error) {
	//动态分配dnsRecords
	var dnsRecords []dnsRecord
	if jsonReq == "" {
		//jsonReq == "" 时，获取所有记录
		start := false
		for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
			if e.Value.(string) == "; Nameservers" {
				start = true
			}
			if start {
				items := gstr.SplitAndTrim(e.Value.(string), " ")
				if len(items) == 0 || items[0] == ";" {
					continue
				}
				switch items[3] {
				case "MX":
					//MX 记录特殊处理，多一项优先级
					dnsRecords = append(dnsRecords, dnsRecord{
						DomainName: items[0],
						TTL:        items[1],
						IN:         items[2],
						Type:       items[3],
						Priority:   items[4],
						Data:       items[5],
					})
				case "TXT":
					//TXT 记录特殊处理，最后的data有引号
					data := gstr.SplitAndTrim(e.Value.(string), `"`)[1]
					dnsRecords = append(dnsRecords, dnsRecord{
						DomainName: items[0],
						TTL:        items[1],
						IN:         items[2],
						Type:       items[3],
						Data:       data,
					})
				default:
					dnsRecords = append(dnsRecords, dnsRecord{
						DomainName: items[0],
						TTL:        items[1],
						IN:         items[2],
						Type:       items[3],
						Data:       items[4],
					})
				}
			}
		}
		return dnsRecords, nil
	} else {
		var req dnsRecord
		err := json.Unmarshal([]byte(jsonReq), &req)
		if err != nil {
			fmt.Println("Error unmarshal jsonRecord:", err)
			return nil, err
		}
		//获取指定域名的记录，有多种组合查询模式
		// 1.查询指定Type所有记录
		if req.DomainName == "" && req.Type != "" && req.Data == "" {
			start := false
			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
				if e.Value.(string) == "; Nameservers" {
					start = true
				}
				if start {
					items := gstr.SplitAndTrim(e.Value.(string), " ")
					if len(items) == 0 || items[0] == ";" || items[3] != req.Type {
						continue
					}

					switch items[3] {
					case "MX":
						//MX 记录特殊处理，多一项优先级
						dnsRecords = append(dnsRecords, dnsRecord{
							DomainName: items[0],
							TTL:        items[1],
							IN:         items[2],
							Type:       items[3],
							Priority:   items[4],
							Data:       items[5],
						})
					case "TXT":
						//TXT 记录特殊处理，最后的data有引号
						data := gstr.SplitAndTrim(e.Value.(string), `"`)[1]
						dnsRecords = append(dnsRecords, dnsRecord{
							DomainName: items[0],
							TTL:        items[1],
							IN:         items[2],
							Type:       items[3],
							Data:       data,
						})
					default:
						dnsRecords = append(dnsRecords, dnsRecord{
							DomainName: items[0],
							TTL:        items[1],
							IN:         items[2],
							Type:       items[3],
							Data:       items[4],
						})
					}
				}
			}
			return dnsRecords, nil
		}

		// 2.查询指定Type + DomainName 记录
		if req.DomainName != "" && req.Type != "" && req.Data == "" {
			start := false
			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
				if e.Value.(string) == "; Nameservers" {
					start = true
				}
				if start {
					items := gstr.SplitAndTrim(e.Value.(string), " ")
					if len(items) == 0 || items[0] == ";" || items[3] != req.Type || items[0] != req.DomainName {
						continue
					}

					switch items[3] {
					case "MX":
						//MX 记录特殊处理，多一项优先级
						dnsRecords = append(dnsRecords, dnsRecord{
							DomainName: items[0],
							TTL:        items[1],
							IN:         items[2],
							Type:       items[3],
							Priority:   items[4],
							Data:       items[5],
						})
					case "TXT":
						//TXT 记录特殊处理，最后的data有引号
						data := gstr.SplitAndTrim(e.Value.(string), `"`)[1]
						dnsRecords = append(dnsRecords, dnsRecord{
							DomainName: items[0],
							TTL:        items[1],
							IN:         items[2],
							Type:       items[3],
							Data:       data,
						})
					default:
						dnsRecords = append(dnsRecords, dnsRecord{
							DomainName: items[0],
							TTL:        items[1],
							IN:         items[2],
							Type:       items[3],
							Data:       items[4],
						})
					}
				}
			}
			return dnsRecords, nil
		}

		// 3.查询指定Type + DomainName + data 记录
		if req.DomainName != "" && req.Type != "" && req.Data != "" {
			start := false
			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
				if e.Value.(string) == "; Nameservers" {
					start = true
				}
				if start {
					items := gstr.SplitAndTrim(e.Value.(string), " ")
					if len(items) == 0 || items[0] == ";" || items[3] != req.Type || items[0] != req.DomainName {
						continue
					}

					switch items[3] {
					case "MX":
						//MX 记录特殊处理，多一项优先级
						if req.Data == items[5] {
							dnsRecords = append(dnsRecords, dnsRecord{
								DomainName: items[0],
								TTL:        items[1],
								IN:         items[2],
								Type:       items[3],
								Priority:   items[4],
								Data:       items[5],
							})
						}
					case "TXT":
						//TXT 记录特殊处理，最后的data有引号
						data := gstr.SplitAndTrim(e.Value.(string), `"`)[1]
						if req.Data == data {
							dnsRecords = append(dnsRecords, dnsRecord{
								DomainName: items[0],
								TTL:        items[1],
								IN:         items[2],
								Type:       items[3],
								Data:       data,
							})
						}
					default:
						if req.Data == items[4] {
							dnsRecords = append(dnsRecords, dnsRecord{
								DomainName: items[0],
								TTL:        items[1],
								IN:         items[2],
								Type:       items[3],
								Data:       items[4],
							})
						}
					}
				}
			}
			return dnsRecords, nil
		}
	}
	return dnsRecords, nil
}

// 获取除MX记录外的其他类型记录
// func (p *ChnZone) GetDomainRecord(jsonReq string) ([]resRecord, error) {
// 	//动态分配resRecords
// 	var resRecords []resRecord
// 	var req reqRecord
// 	err := json.Unmarshal([]byte(jsonReq), &req)
// 	if err != nil {
// 		fmt.Println("Error unmarshal jsonRecord:", err)
// 		return nil, err
// 	}

// 	if req.Type == "MX" {
// 		return nil, fmt.Errorf("MX记录请调用GetDomainMXRecord")
// 	}

// 	if req.DomainName == "*" {
// 		//获取对应type所有记录
// 		switch req.Type {
// 		case "NS":
// 			start := false
// 			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 				if e.Value.(string) == "; Nameservers" {
// 					start = true
// 				}
// 				if e.Value.(string) == "; Mailservers" {
// 					return resRecords, nil
// 				}
// 				if start {
// 					items := gstr.SplitAndTrim(e.Value.(string), " ")
// 					if len(items) == 0 {
// 						continue
// 					}
// 					if items[0] != ";" && items[0] != "" {
// 						resRecords = append(resRecords, resRecord{
// 							DomainName: items[0],
// 							TTL:        items[1],
// 							Type:       items[3],
// 							Data:       items[4],
// 						})
// 					}
// 				}
// 			}
// 		case "PTR":
// 			start := false
// 			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 				if e.Value.(string) == "; Reverse DNS Records (PTR)" {
// 					start = true
// 				}
// 				if e.Value.(string) == "; TXT" {
// 					return resRecords, nil
// 				}
// 				if start {
// 					items := gstr.SplitAndTrim(e.Value.(string), " ")
// 					if len(items) == 0 {
// 						continue
// 					}
// 					if items[0] != ";" && items[0] != "" {
// 						resRecords = append(resRecords, resRecord{
// 							DomainName: items[0],
// 							TTL:        items[1],
// 							Type:       items[3],
// 							Data:       items[4],
// 						})
// 					}
// 				}
// 			}
// 		case "CNAME":
// 			start := false
// 			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 				if e.Value.(string) == "; CNAME" {
// 					start = true
// 				}
// 				if e.Value.(string) == "; HOST RECORDS" {
// 					return resRecords, nil
// 				}
// 				if start {
// 					items := gstr.SplitAndTrim(e.Value.(string), " ")
// 					if len(items) == 0 {
// 						continue
// 					}
// 					if items[0] != ";" && items[0] != "" {
// 						resRecords = append(resRecords, resRecord{
// 							DomainName: items[0],
// 							TTL:        items[1],
// 							Type:       items[3],
// 							Data:       items[4],
// 						})
// 					}
// 				}
// 			}
// 		case "TXT":
// 			start := false
// 			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 				if e.Value.(string) == "; TXT" {
// 					start = true
// 				}
// 				if e.Value.(string) == "; CNAME" {
// 					return resRecords, nil
// 				}
// 				if start {
// 					items := gstr.SplitAndTrim(e.Value.(string), " ")
// 					//发现空行，跳过继续解析next
// 					if len(items) == 0 {
// 						continue
// 					}
// 					if items[0] != ";" && items[0] != "" {
// 						resRecords = append(resRecords, resRecord{
// 							DomainName: items[0],
// 							TTL:        items[1],
// 							Type:       items[3],
// 							Data:       gstr.Str(e.Value.(string), `"`),
// 							//Data:       items[4],
// 						})
// 					}
// 				}
// 			}
// 		case "A":
// 			start := false
// 			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 				if e.Value.(string) == "; HOST RECORDS" {
// 					start = true
// 				}
// 				if e.Next() == nil {
// 					return resRecords, nil
// 				}
// 				if start {
// 					items := gstr.SplitAndTrim(e.Value.(string), " ")
// 					if len(items) == 0 {
// 						continue
// 					}
// 					if items[0] != ";" && items[0] != "" {
// 						if items[3] != "A" {
// 							continue
// 						}
// 						resRecords = append(resRecords, resRecord{
// 							DomainName: items[0],
// 							TTL:        items[1],
// 							Type:       items[3],
// 							Data:       items[4],
// 						})
// 					}
// 				}
// 			}
// 		case "A9":
// 			start := false
// 			for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 				if e.Value.(string) == "; HOST RECORDS" {
// 					start = true
// 				}
// 				if e.Next() == nil {
// 					return resRecords, nil
// 				}
// 				if start {
// 					items := gstr.SplitAndTrim(e.Value.(string), " ")
// 					if len(items) == 0 {
// 						continue
// 					}
// 					if items[0] != ";" && items[0] != "" {
// 						if items[3] != "A9" {
// 							continue
// 						}
// 						resRecords = append(resRecords, resRecord{
// 							DomainName: items[0],
// 							TTL:        items[1],
// 							Type:       items[3],
// 							Data:       items[4],
// 						})
// 					}
// 				}
// 			}

// 		default:
// 			return nil, fmt.Errorf("不支持的类型")

// 		}
// 	}

// 	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 		items := gstr.SplitAndTrim(e.Value.(string), " ")
// 		if len(items) == 0 {
// 			continue
// 		}
// 		if gstr.Trim(items[0]) == req.DomainName {
// 			if gstr.Trim(items[3]) == req.Type {
// 				//发现记录
// 				resRecords = append(resRecords, resRecord{
// 					DomainName: items[0],
// 					TTL:        items[1],
// 					Type:       items[3],
// 					Data:       items[4],
// 				})
// 			}
// 		}
// 	}

// 	return resRecords, nil
// }

// func (p *ChnZone) GetDomainMXRecord(jsonReq string) ([]resMXRecord, error) {
// 	//MX 记录特殊处理
// 	var resMXRecords []resMXRecord
// 	var req reqRecord
// 	err := json.Unmarshal([]byte(jsonReq), &req)
// 	if err != nil {
// 		fmt.Println("Error unmarshal jsonRecord:", err)
// 		return nil, err
// 	}

// 	if req.Type != "MX" {
// 		return nil, fmt.Errorf("此方法仅支持MX类型记录")
// 	}
// 	if req.DomainName == "*" {
// 		//获取对应type所有记录
// 		start := false
// 		for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 			if e.Value.(string) == "; Mailservers" {
// 				start = true
// 			}
// 			if e.Value.(string) == "; Reverse DNS Records (PTR)" {
// 				return resMXRecords, nil
// 			}
// 			if start {
// 				items := gstr.SplitAndTrim(e.Value.(string), " ")
// 				//发现空行，跳过继续解析next
// 				if len(items) == 0 {
// 					continue
// 				}
// 				if items[0] != ";" && items[0] != "" {
// 					resMXRecords = append(resMXRecords, resMXRecord{
// 						DomainName: items[0],
// 						TTL:        items[1],
// 						Type:       items[3],
// 						Priority:   items[4],
// 						Data:       items[5],
// 					})
// 				}
// 			}
// 		}
// 	}
// 	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 		items := gstr.SplitAndTrim(e.Value.(string), " ")
// 		if len(items) == 0 {
// 			continue
// 		}
// 		if gstr.Trim(items[0]) == req.DomainName {
// 			if gstr.Trim(items[3]) == req.Type {
// 				//发现记录
// 				resMXRecords = append(resMXRecords, resMXRecord{
// 					DomainName: items[0],
// 					TTL:        items[1],
// 					Type:       items[3],
// 					Priority:   items[4],
// 					Data:       items[5],
// 				})
// 			}
// 		}
// 	}
// 	return resMXRecords, nil
// }

// func (p *ChnZone) AddMXRecord(jsonRecord string) error {
// 	if err := p.checkJaonRecord(jsonRecord); err != nil {
// 		return err
// 	}
// 	// 反序列化jsonRecord
// 	var record mxRecord
// 	err := json.Unmarshal([]byte(jsonRecord), &record)
// 	if err != nil {
// 		fmt.Println("Error unmarshal jsonRecord:", err)
// 		return err
// 	}

// 	if gstr.Trim(record.Priority) == "" {
// 		record.Priority = "10"
// 	}

// 	// 填充record.IN
// 	record.IN = "IN"

// 	err = p.findMXRecord(record)
// 	if err != nil {
// 		return err
// 	}
// 	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Priority + " " + record.Data
// 	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 		if e.Value.(string) == "; Mailservers" {
// 			p.runtimeZoneFileList.InsertAfter(e, strRecord)
// 			break
// 		}
// 	}
// 	// 递增serial
// 	err = p.incrementSerial()
// 	if err != nil {
// 		fmt.Println("Error increment serial:", err)
// 		return err
// 	}
// 	err = p.WriteZoneFile()
// 	return err
// }

func (p *ChnZone) addPTRRecord(record dnsRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Reverse DNS Records (PTR)" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found Reverse DNS Records (PTR) area")
}

func (p *ChnZone) addCNAMERecord(record dnsRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data + "."
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; CNAME" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found CNAME area")
}

func (p *ChnZone) addTXTRecord(record dnsRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + strconv.Quote(record.Data)
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; TXT" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found TXT area")
}

func (p *ChnZone) addNSRecord(record dnsRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data + "."
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Nameservers" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found Nameservers area")
}

func (p *ChnZone) addDomainRecord(record dnsRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; HOST RECORDS" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found HOST RECORDS area")
}

func (p *ChnZone) addMXRecord(record dnsRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Priority + " " + record.Data
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Mailservers" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found Mailservers area")
}

func (p *ChnZone) incrementSerial() error {
	// 从runtimeZoneFileList中读取serial
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if gstr.Contains(e.Value.(string), "IN SOA") {
			//find "IN SOA" line,serial is the next line
			fmt.Println(e.Next().Value.(string))
			serial, err := strconv.Atoi(gstr.Trim(gstr.Split(e.Next().Value.(string), ";")[0]))
			if err != nil {
				fmt.Println("Error convert serial:", err)
				return err
			}
			serial++
			//在serial前面面加上/t

			e.Next().Value = "\t\t\t" + strconv.Itoa(serial) + " ; serial"
			return nil
		}
	}
	return nil
}

// func (p *ChnZone) writeRuntimeZoneFile() error {
// 	strContent := ""
// 	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 		strContent += e.Value.(string) + "\n"
// 	}
// 	// 把strContent写入文件
// 	err := os.WriteFile("runtime.zone", []byte(strContent), 0644)
// 	if err != nil {
// 		fmt.Println("Error writing file:", err)
// 		return err
// 	}
// 	return nil
// }

func (p *ChnZone) WriteZoneFile() error {
	strContent := ""
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		strContent += e.Value.(string) + "\n"
	}
	// 把strContent写入文件
	err := os.WriteFile("/var/named/chn.zone", []byte(strContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}

func (p *ChnZone) checkTTL(ttl string) error {
	//检查ttl是否是数字
	ttlNum, err := strconv.Atoi(ttl)
	if err != nil {
		return err
	}
	if ttlNum <= 0 {
		return fmt.Errorf("TTL 不能是负数或0")
	}
	return nil
}

func (p *ChnZone) checkPriority(pri string) error {
	//检查ttl是否是数字
	priNum, err := strconv.Atoi(pri)
	if err != nil {
		return err
	}
	if priNum <= 0 {
		return fmt.Errorf("优先级 不能是负数或0")
	}
	return nil
}

func (p *ChnZone) checkIPv4Address(address string) error {
	//判断是否是合法的IP地址
	addr := net.ParseIP(address)
	if addr == nil {
		return fmt.Errorf("data is not a valid IPv4 address")
	}
	return nil
}

func (p *ChnZone) checkIPv9Address(address string) error {
	//判断是否是合法的IPv9地址
	if gstr.Contains(address, " ") {
		return fmt.Errorf("data is not a valid IPv9 address，不能包含空格")
	}
	if !gstr.Contains(address, "[") {
		return fmt.Errorf("data is not a valid IPv9 address")
	}
	items := gstr.Split(address, "[")
	if len(items) == 8 {
		for _, item := range items {
			address := net.ParseIP(item)
			//每个item必须是数字或者最后一个item是IPv4地址
			if !gstr.IsNumeric(item) && address == nil {
				return fmt.Errorf("data is not a valid IPv9 address")
			}
			if address == nil {
				decIP, err := strconv.Atoi(item)
				if err != nil {
					return err
				}
				if decIP < 0 {
					return fmt.Errorf("data is not a valid IPv9 address,单段IPv9地址必须大于等于0")
				}
			}
		}
		return nil
	}
	if len(items) > 8 {
		return fmt.Errorf("data is not a valid IPv9 address,IPv9地址段不能超过8段")
	}
	if len(items) < 8 {
		//考虑地址压缩情况
		if !gstr.Contains(address, "]") {
			return fmt.Errorf("data is not a valid IPv9 address,IPv9地址段小于8段")
		}
		compressLens := 0
		for _, item := range items {
			if gstr.Contains(item, "]") {
				//地址压缩
				compressItems := gstr.Split(item, "]")
				// if len(compressItems) != 3 {
				// 	return fmt.Errorf("data is not a valid IPv9 address,地址压缩错误")
				// }
				compressLen, err := strconv.Atoi(compressItems[0])
				if err != nil {
					return err
				}
				if compressLen <= 0 {
					return fmt.Errorf("data is not a valid IPv9 address,地址压缩错误,不能小于等于0")
				}
				if compressLen > 8 {
					return fmt.Errorf("data is not a valid IPv9 address,地址压缩错误,不能大于8")
				}
				compressLens += compressLen
				item = compressItems[1]
			}
			//fmt.Println("compressLens:", compressLens)
			//每个item必须是数字或者最后一个item是IPv4地址
			address := net.ParseIP(item)
			if !gstr.IsNumeric(item) && address == nil {
				return fmt.Errorf("data is not a valid IPv9 address")
			}
			if address == nil {
				decIP, err := strconv.Atoi(item)
				if err != nil {
					return err
				}
				if decIP < 0 {
					return fmt.Errorf("data is not a valid IPv9 address,单段IPv9地址必须大于等于0")
				}
			}
		}
		if len(items)+compressLens != 8 {
			return fmt.Errorf("data is not a valid IPv9 address,IPv9地址段不等于8段")
		}
	}
	return nil
}

// func (p *ChnZone) checkJaonRecord(jsonRecord string) error {
// 	// 检查jsonRecord是否符合规范
// 	// 反序列化jsonRecord
// 	var record domainRecord
// 	var err error
// 	err = json.Unmarshal([]byte(jsonRecord), &record)
// 	if err != nil {
// 		fmt.Println("Error unmarshal jsonRecord:", err)
// 		return err
// 	}
// 	// if record.DomainName == "" || record.TTL == "" || record.Type == "" || record.Data == "" {
// 	// 	return fmt.Errorf("fields is empty")
// 	// }
// 	if gstr.Trim(record.DomainName) == "" || gstr.Trim(record.TTL) == "" || gstr.Trim(record.Type) == "" || gstr.Trim(record.Data) == "" {
// 		return fmt.Errorf("fields is empty")
// 	}
// 	ttl, err := strconv.Atoi(record.TTL)
// 	if err != nil {
// 		return err
// 	}
// 	if ttl <= 0 {
// 		return fmt.Errorf("TTL 不能是负数或0")
// 	}
// 	if record.Type == "A" {
// 		//判断是否是合法的IP地址
// 		address := net.ParseIP(record.Data)
// 		if address == nil {
// 			return fmt.Errorf("data is not a valid IPv4 address")
// 		}
// 	}

// 	if record.Type == "A9" {
// 		//判断是否是合法的IPv9地址
// 		if gstr.Contains(record.Data, " ") {
// 			return fmt.Errorf("data is not a valid IPv9 address，不能包含空格")
// 		}
// 		if !gstr.Contains(record.Data, "[") {
// 			return fmt.Errorf("data is not a valid IPv9 address")
// 		}
// 		items := gstr.Split(record.Data, "[")
// 		if len(items) == 8 {
// 			for _, item := range items {
// 				address := net.ParseIP(item)
// 				//每个item必须是数字或者最后一个item是IPv4地址
// 				if !gstr.IsNumeric(item) && address == nil {
// 					return fmt.Errorf("data is not a valid IPv9 address")
// 				}
// 				if address == nil {
// 					decIP, err := strconv.Atoi(item)
// 					if err != nil {
// 						return err
// 					}
// 					if decIP < 0 {
// 						return fmt.Errorf("data is not a valid IPv9 address,单段IPv9地址必须大于等于0")
// 					}
// 				}
// 			}
// 			return nil
// 		}
// 		if len(items) > 8 {
// 			return fmt.Errorf("data is not a valid IPv9 address,IPv9地址段不能超过8段")
// 		}
// 		if len(items) < 8 {
// 			//考虑地址压缩情况
// 			if !gstr.Contains(record.Data, "]") {
// 				return fmt.Errorf("data is not a valid IPv9 address,IPv9地址段小于8段")
// 			}
// 			compressLens := 0
// 			for _, item := range items {
// 				if gstr.Contains(item, "]") {
// 					//地址压缩
// 					compressItems := gstr.Split(item, "]")
// 					// if len(compressItems) != 3 {
// 					// 	return fmt.Errorf("data is not a valid IPv9 address,地址压缩错误")
// 					// }
// 					compressLen, err := strconv.Atoi(compressItems[0])
// 					if err != nil {
// 						return err
// 					}
// 					if compressLen <= 0 {
// 						return fmt.Errorf("data is not a valid IPv9 address,地址压缩错误,不能小于等于0")
// 					}
// 					if compressLen > 8 {
// 						return fmt.Errorf("data is not a valid IPv9 address,地址压缩错误,不能大于8")
// 					}
// 					compressLens += compressLen
// 					item = compressItems[1]
// 				}
// 				//fmt.Println("compressLens:", compressLens)
// 				//每个item必须是数字或者最后一个item是IPv4地址
// 				address := net.ParseIP(item)
// 				if !gstr.IsNumeric(item) && address == nil {
// 					return fmt.Errorf("data is not a valid IPv9 address")
// 				}
// 				if address == nil {
// 					decIP, err := strconv.Atoi(item)
// 					if err != nil {
// 						return err
// 					}
// 					if decIP < 0 {
// 						return fmt.Errorf("data is not a valid IPv9 address,单段IPv9地址必须大于等于0")
// 					}
// 				}
// 			}
// 			if len(items)+compressLens != 8 {
// 				return fmt.Errorf("data is not a valid IPv9 address,IPv9地址段不等于8段")
// 			}
// 		}

// 	}
// 	return nil
// }

func (p *ChnZone) findRecord(record dnsRecord) bool {
	start := false
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Nameservers" {
			start = true
		}
		if start {
			items := gstr.SplitAndTrim(e.Value.(string), " ")
			if len(items) == 0 || items[0] == ";" {
				continue
			}
			if items[0] == record.DomainName {
				if items[3] == record.Type {
					if record.Type == "MX" {
						if items[5] == record.Data {
							return true
						}
					} else if record.Type == "TXT" {
						if gstr.SplitAndTrim(e.Value.(string), `"`)[1] == record.Data {
							return true
						}
					} else if record.Type == "NS" {
						if items[4] == record.Data+"." {
							return true
						}
					} else {
						if items[4] == record.Data {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (p *ChnZone) findDNSRecordAndDelete(record dnsRecord) error {
	start := false
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Nameservers" {
			start = true
		}
		if start {
			items := gstr.SplitAndTrim(e.Value.(string), " ")
			if len(items) == 0 || items[0] == ";" {
				continue
			}
			if items[0] == record.DomainName {
				if items[3] == record.Type {
					if record.Type == "MX" {
						if items[5] == record.Data {
							p.runtimeZoneFileList.Remove(e)
							return nil
						}
					} else if record.Type == "TXT" {
						if gstr.SplitAndTrim(e.Value.(string), `"`)[1] == record.Data {
							p.runtimeZoneFileList.Remove(e)
							return nil
						}
					} else if record.Type == "NS" {
						if items[4] == record.Data+"." {
							p.runtimeZoneFileList.Remove(e)
							return nil
						}
					} else {
						if items[4] == record.Data {
							p.runtimeZoneFileList.Remove(e)
							return nil
						}
					}
				}
			}
		}
	}
	return fmt.Errorf("not found record")
}

// func (p *ChnZone) findMXRecord(record mxRecord) error {
// 	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
// 		items := gstr.Split(e.Value.(string), " ")
// 		if gstr.Trim(items[0]) == record.DomainName {
// 			if gstr.Trim(items[3]) == record.Type {
// 				//发现同样的记录
// 				return fmt.Errorf("find same record")
// 			}
// 		}
// 	}
// 	return nil
// }
