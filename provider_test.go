package huaweidns

import (
	"context"
	"fmt"
	"testing"
)

func TestProvider_GetRecords(t *testing.T) {
	p := &Provider{
		AccessKey:       "15ECIPCQCLRKG5TI4FDH",
		SecretAccessKey: "XMB3UMuDXz2LKRYYfCOOqAqUI26DS3neL7YqS93f",
		RegionID:        "cn-north-4",
		EndPoint:        "https://dns.cn-north-4.myhuaweicloud.com/",
		ZoneId:          "8aace3ba66f3b68501672f3659a43053",
	}
	p.initClient()
	records, err := p.GetRecords(context.Background(), "bjshujie.com")
	if err != nil {
		fmt.Println(err)
	}

	for _, record := range records {
		fmt.Printf("Id:%s Name:%s Type:%s Value:%s TTL:%v Priority:%d", record.ID, record.Name, record.Type, record.Value, record.TTL, record.Priority)
		fmt.Println()
	}

}
