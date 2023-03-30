package huaweidns

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	"go.uber.org/zap"
)

var (
	_ caddyfile.Unmarshaler = (*Module)(nil)
)

func init() {
	caddy.RegisterModule(&Module{})
}

// Module wraps dnspodcn.Provider.
type Module struct {
	*Provider
	logger *zap.SugaredLogger
}

// CaddyModule returns the Caddy module information.
func (m *Module) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "dns.providers.huaweidns",
		New: func() caddy.Module {
			return &Module{
				Provider: &Provider{},
			}
		},
	}
}

func (m *Module) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger().Sugar()
	m.logger.Infof("AccessKey:%s, SecretAccessKey:%s, RegionID:%s, EndPoint:%s, ZoneId:%s",
		m.AccessKey, m.SecretAccessKey, m.RegionID, m.EndPoint, m.ZoneId)
	m.initClient()
	return nil
}

func (m *Module) Validate() error {
	_, err := m.dnsClient.ListApiVersions(&model.ListApiVersionsRequest{})
	if err != nil {
		m.logger.Infof("huawei dns privider validate failed:%v", err)
	}
	return err
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	dnspodcn [<app_id> <app_token>] {
//	    app_id <app_id>
//	    app_token <app_token>
//	}
func (m *Module) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	repl := caddy.NewReplacer()

	for d.Next() {
		//插件名后面的参数
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "access_key":
				if d.NextArg() {
					m.AccessKey = repl.ReplaceAll(d.Val(), "")
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "secret_access_key":
				if d.NextArg() {
					m.SecretAccessKey = repl.ReplaceAll(d.Val(), "")
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "region_id":
				if d.NextArg() {
					m.RegionID = repl.ReplaceAll(d.Val(), "")
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "end_point":
				if d.NextArg() {
					m.EndPoint = repl.ReplaceAll(d.Val(), "")
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "zone_id":
				if d.NextArg() {
					m.ZoneId = repl.ReplaceAll(d.Val(), "")
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if m.AccessKey == "" || m.SecretAccessKey == "" {
		return d.Err("missing Access Key or Secret Access Key")
	}
	return nil
}
