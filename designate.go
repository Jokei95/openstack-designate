package openstack

import (
	designate "github.com/Jokei95/libdns-openstack-designate"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// Provider wraps the libdns OpenStack Designate provider as a Caddy DNS module.
type Provider struct{ *designate.Provider }

func init() { caddy.RegisterModule(Provider{}) }

func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.openstack-designate",
		New: func() caddy.Module { return &Provider{Provider: new(designate.Provider)} },
	}
}

// Provision expands environment variables in the provider credentials.
func (p *Provider) Provision(_ caddy.Context) error {
	if p.Provider == nil {
		p.Provider = new(designate.Provider)
	}
	repl := caddy.NewReplacer()
	a := &p.Provider.AuthOpenStack
	a.RegionName = repl.ReplaceAll(a.RegionName, "")
	a.TenantID = repl.ReplaceAll(a.TenantID, "")
	a.IdentityApiVersion = repl.ReplaceAll(a.IdentityApiVersion, "")
	a.ApplicationCredentialSecret = repl.ReplaceAll(a.ApplicationCredentialSecret, "")
	a.AuthURL = repl.ReplaceAll(a.AuthURL, "")
	a.ApplicationCredentialId = repl.ReplaceAll(a.ApplicationCredentialId, "")
	a.TenantName = repl.ReplaceAll(a.TenantName, "")
	a.EndpointType = repl.ReplaceAll(a.EndpointType, "")
	return nil
}

// UnmarshalCaddyfile accepts the legacy provider Caddyfile syntax:
//
//	openstack-designate {
//	    region_name REGION
//	    tenant_id TENANT_ID
//	    identity_api_version VERSION
//      application_credential_secret <string>
//	    auth_url AUTH_URL
//	    application_credential_id <string>
//	    tenant_name TENANT_NAME
//	    endpoint_type ENDPOINT_TYPE
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if p.Provider == nil {
		p.Provider = new(designate.Provider)
	}
	a := &p.Provider.AuthOpenStack

	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			var dst *string
			switch d.Val() {
			case "region_name":
				dst = &a.RegionName
			case "tenant_id":
				dst = &a.TenantID
			case "identity_api_version":
				dst = &a.IdentityApiVersion
			case "application_credential_secret":
				dst = &a.ApplicationCredentialSecret
			case "auth_url":
				dst = &a.AuthURL
			case "application_credential_id":
				dst = &a.ApplicationCredentialId
			case "tenant_name":
				dst = &a.TenantName
			case "endpoint_type":
				dst = &a.EndpointType
			default:
				return d.Errf("unrecognized subdirective %q", d.Val())
			}
			if !d.NextArg() {
				return d.ArgErr()
			}
			*dst = d.Val()
			if d.NextArg() {
				return d.ArgErr()
			}
		}
	}
	return nil
}

var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
