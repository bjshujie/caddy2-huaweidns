package huaweidns

import (
	"context"
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/impl"
	dns "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/region"
	"github.com/libdns/libdns"
	"strings"
	"time"
)

func (p *Provider) initClient() {
	endpoints := []string{p.EndPoint}
	credential := basic.NewCredentialsBuilder().WithAk(p.AccessKey).WithSk(p.SecretAccessKey).Build()
	httpConfig := config.DefaultHttpConfig().WithHttpHandler(httphandler.NewHttpHandler())
	hcClient := http_client.NewHcHttpClient(impl.NewDefaultHttpClient(httpConfig))
	hcClient = hcClient.WithCredential(credential).WithEndpoints(endpoints)
	client := dns.NewDnsClient(
		dns.DnsClientBuilder().
			WithRegion(region.ValueOf(p.RegionID)).
			WithCredential(credential).
			Build())
	p.dnsClient = client
}

func (p *Provider) getDomain(ctx context.Context, zone string) ([]libdns.Record, error) {
	request := &model.ListRecordSetsByZoneRequest{ZoneId: p.ZoneId}
	resp, err := p.dnsClient.ListRecordSetsByZone(request)
	if err != nil {
		return nil, err
	}
	dnsRecords := make([]libdns.Record, 0, len(*resp.Recordsets))
	for _, record := range *resp.Recordsets {
		dnsRecord := libdns.Record{}
		dnsRecord.ID = *record.Id
		dnsRecord.Name = *record.Name
		dnsRecord.Type = *record.Type
		dnsRecord.TTL = time.Duration(*record.Ttl) * time.Second
		dnsRecord.Value = strings.Join(*record.Records, ",")
		dnsRecord.Priority = 1
		dnsRecords = append(dnsRecords, dnsRecord)
	}
	return dnsRecords, nil
}

func (p *Provider) setRecord(ctx context.Context, zone string, record libdns.Record, clear bool) error {
	if clear {
		delRequest := &model.DeleteRecordSetRequest{ZoneId: zone, RecordsetId: record.ID}
		p.dnsClient.DeleteRecordSet(delRequest)
	} else {
		ttl := int32(record.TTL.Seconds())
		records := strings.Split(record.Value, ",")
		if strings.TrimSpace(record.ID) == "" {
			req := &model.CreateRecordSetReq{Name: record.Name,
				Description: nil, Type: record.Type, Status: nil, Ttl: &ttl, Records: records, Tags: nil}
			createRequest := &model.CreateRecordSetRequest{ZoneId: zone, Body: req}
			p.dnsClient.CreateRecordSet(createRequest)
		} else {
			req := &model.UpdateRecordSetReq{Name: record.Name,
				Description: nil, Type: record.Type, Ttl: &ttl, Records: &records}

			updateRequest := &model.UpdateRecordSetRequest{ZoneId: zone, RecordsetId: record.ID, Body: req}
			p.dnsClient.UpdateRecordSet(updateRequest)
		}
	}
	return nil
}
