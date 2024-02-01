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

type domainRecord struct {
	DomainName string `json:"domainName"`
	TTL        string `json:"ttl"`
	IN         string `json:"-"`
	Type       string `json:"type"`
	Data       string `json:"data"`
}

type reqRecord struct {
	DomainName string `json:"domainName"`
	Type       string `json:"type"`
}

type resRecord struct {
	DomainName string `json:"domainName"`
	TTL        string `json:"ttl"`
	Type       string `json:"type"`
	Data       string `json:"data"`
}

func (p *ChnZone) Init() {
	fmt.Println("init ChnZone...")
	p.defaultZoneFileList = glist.New()
	p.runtimeZoneFileList = glist.New()

	// 填充默认的zone文件defaultZoneFileList
	p.initDefaultZoneFileList(p.defaultZoneFileList)

	// 读取chn.zone文件，填充druntimeZoneFileList
	p.readZoneContentFromFile("chn.zone", p.runtimeZoneFileList)
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

func (p *ChnZone) AddRecord(jsonRecord string) error {
	if err := p.checkJaonRecord(jsonRecord); err != nil {
		return err
	}
	// 反序列化jsonRecord
	var record domainRecord
	err := json.Unmarshal([]byte(jsonRecord), &record)
	if err != nil {
		fmt.Println("Error unmarshal jsonRecord:", err)
		return err
	}

	// 填充record.IN
	record.IN = "IN"

	err = p.findRecord(record)
	if err != nil {
		return err
	}

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
		fmt.Println("Error record type:", record.Type)
	}

	if err != nil {
		fmt.Println("Error add record:", err)
		return err
	}

	// 递增serial
	err = p.incrementSerial()
	if err != nil {
		fmt.Println("Error increment serial:", err)
		return err
	}

	// 把runtimeZoneFileList写入文件
	err = p.writeRuntimeZoneFile()
	return err
}

func (p *ChnZone) DelRecord(jsonRecord string) error {
	if err := p.checkJaonRecord(jsonRecord); err != nil {
		return err
	}
	// 反序列化jsonRecord
	var record domainRecord
	err := json.Unmarshal([]byte(jsonRecord), &record)
	if err != nil {
		fmt.Println("Error unmarshal jsonRecord:", err)
		return err
	}

	// 填充record.IN
	record.IN = "IN"

	// 匹配record与runtimeZoneFileList中的记录，找到后删除
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		items := gstr.Split(e.Value.(string), " ")
		if gstr.Trim(items[0]) == record.DomainName {
			if gstr.Trim(items[3]) == record.Type {
				var data = record.Data
				if record.Type == "TXT" {
					data = strconv.Quote(data)
				}
				if record.Type == "NS" {
					data = data + "."
				}
				if data == gstr.Trim(items[4]) {
					p.runtimeZoneFileList.Remove(e)
					// 递增serial
					err = p.incrementSerial()
					if err != nil {
						fmt.Println("Error increment serial:", err)
						return err
					}
					// 把runtimeZoneFileList写入文件
					err = p.writeRuntimeZoneFile()
					return err
				}
			}
		}
	}
	return fmt.Errorf("DelRecord... not found record")
}

func (p *ChnZone) ModifyRecord(jsonRecord string) error {
	if err := p.checkJaonRecord(jsonRecord); err != nil {
		return err
	}
	var strContent string
	// 反序列化jsonRecord
	var record domainRecord
	err := json.Unmarshal([]byte(jsonRecord), &record)
	if err != nil {
		fmt.Println("Error unmarshal jsonRecord:", err)
		return err
	}

	// 填充record.IN
	record.IN = "IN"

	// 匹配record.type, recrod.domainName与runtimeZoneFileList中的记录，找到后修改
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		items := gstr.Split(e.Value.(string), " ")
		if gstr.Trim(items[0]) == record.DomainName {
			if gstr.Trim(items[3]) == record.Type {
				if record.Type == "TXT" {
					strContent = record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + strconv.Quote(record.Data)
				} else {
					strContent = record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data
				}
				e.Value = strContent
				// 递增serial
				err = p.incrementSerial()
				if err != nil {
					fmt.Println("Error increment serial:", err)
					return err
				}
				// 把runtimeZoneFileList写入文件
				err = p.writeRuntimeZoneFile()
				return err
			}
		}
	}
	return fmt.Errorf("ModifyRecord... not found record")
}

func (p *ChnZone) GetDomainRecord(jsonReq string) ([]resRecord, error) {
	//var resRecords []resRecord = make([]resRecord, 0)
	//动态分配resRecords
	var resRecords []resRecord
	var req reqRecord
	err := json.Unmarshal([]byte(jsonReq), &req)
	if err != nil {
		fmt.Println("Error unmarshal jsonRecord:", err)
		return nil, err
	}

	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		items := gstr.Split(e.Value.(string), " ")
		if gstr.Trim(items[0]) == req.DomainName {
			if gstr.Trim(items[3]) == req.Type {
				//发现记录
				resRecords = append(resRecords, resRecord{
					DomainName: items[0],
					TTL:        items[1],
					Type:       items[3],
					Data:       items[4],
				})
			}
		}
	}

	return resRecords, nil
}

func (p *ChnZone) GetAllDomainRecord() ([]resRecord, error) {
	var resRecords []resRecord
	var start bool = false
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Nameservers" {
			start = true
		}
		if start {
			items := gstr.Split(e.Value.(string), " ")
			if items[0] != ";" && items[0] != "" {
				resRecords = append(resRecords, resRecord{
					DomainName: items[0],
					TTL:        items[1],
					Type:       items[3],
					Data:       items[4],
				})
			}
		}

	}

	// resJson, err := gjson.Marshal(resRecords)
	// if err != nil {
	// 	fmt.Println("Error marshal resRecords:", err)
	// 	return nil, err
	// }
	// fmt.Println(string(resJson))
	return resRecords, nil
}

func (p *ChnZone) addMXRecord(record domainRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Mailservers" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found Mailservers area")
}

func (p *ChnZone) addPTRRecord(record domainRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Reverse DNS Records (PTR)" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found Reverse DNS Records (PTR) area")
}

func (p *ChnZone) addCNAMERecord(record domainRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data + "."
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; CNAME" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found CNAME area")
}

func (p *ChnZone) addTXTRecord(record domainRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + strconv.Quote(record.Data)
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; TXT" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found TXT area")
}

func (p *ChnZone) addNSRecord(record domainRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data + "."
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; Nameservers" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found Nameservers area")
}

func (p *ChnZone) addDomainRecord(record domainRecord) error {
	strRecord := record.DomainName + " " + record.TTL + " " + "IN" + " " + record.Type + " " + record.Data
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == "; HOST RECORDS" {
			p.runtimeZoneFileList.InsertAfter(e, strRecord)
			return nil
		}
	}
	return fmt.Errorf("not found HOST RECORDS area")
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

func (p *ChnZone) writeRuntimeZoneFile() error {
	strContent := ""
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		strContent += e.Value.(string) + "\n"
	}
	// 把strContent写入文件
	err := os.WriteFile("runtime.zone", []byte(strContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}

func (p *ChnZone) checkJaonRecord(jsonRecord string) error {
	// 检查jsonRecord是否符合规范
	// 反序列化jsonRecord
	var record domainRecord
	var err error
	err = json.Unmarshal([]byte(jsonRecord), &record)
	if err != nil {
		fmt.Println("Error unmarshal jsonRecord:", err)
		return err
	}
	if record.DomainName == "" || record.TTL == "" || record.Type == "" || record.Data == "" {
		return fmt.Errorf("fields is empty")
	}
	if gstr.Trim(record.DomainName) == "" || gstr.Trim(record.TTL) == "" || gstr.Trim(record.Type) == "" || gstr.Trim(record.Data) == "" {
		return fmt.Errorf("DomainName is empty")
	}
	ttl, err := strconv.Atoi(record.TTL)
	if err != nil {
		return err
	}
	if ttl <= 0 {
		return fmt.Errorf("TTL 不能是负数或0")
	}
	if record.Type == "A" {
		//判断是否是合法的IP地址
		address := net.ParseIP(record.Data)
		if address == nil {
			return fmt.Errorf("data is not a valid IPv4 address")
		}
	}

	if record.Type == "A9" {
		//判断是否是合法的IPv9地址
		if gstr.Contains(record.Data, " ") {
			return fmt.Errorf("data is not a valid IPv9 address，不能包含空格")
		}
		if !gstr.Contains(record.Data, "[") {
			return fmt.Errorf("data is not a valid IPv9 address")
		}
		items := gstr.Split(record.Data, "[")
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
			if !gstr.Contains(record.Data, "]") {
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

	}
	return nil
}

func (p *ChnZone) findRecord(record domainRecord) error {
	for e := p.runtimeZoneFileList.Front(); e != nil; e = e.Next() {
		items := gstr.Split(e.Value.(string), " ")
		if gstr.Trim(items[0]) == record.DomainName {
			if gstr.Trim(items[3]) == record.Type {
				//发现同样的记录
				return fmt.Errorf("find same record")
			}
		}
	}
	return nil
}
