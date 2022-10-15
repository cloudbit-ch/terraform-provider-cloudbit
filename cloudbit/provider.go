package cloudbit

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/flowswiss/goclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ tfsdk.Provider = (*provider)(nil)

type Option func(p *provider)

func WithVersion(version string) Option {
	return func(p *provider) {
		p.version = version
	}
}

func WithDefaultEndpoint(endpoint string) Option {
	return func(p *provider) {
		p.defaultEndpoint = endpoint
	}
}

func New(opts ...Option) tfsdk.Provider {
	p := &provider{
		version:         "dev",
		defaultEndpoint: "https://api.cloudbit.ch/",
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

type provider struct {
	version         string
	defaultEndpoint string

	client     goclient.Client
	configured bool
}

type providerData struct {
	Token    types.String `tfsdk:"token"`
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"token": {
				Type:                types.StringType,
				MarkdownDescription: "authentication token for the cloudbit api",
				Optional:            true,
				Sensitive:           true,
			},
			"endpoint": {
				Type:                types.StringType,
				MarkdownDescription: "endpoint of the cloudbit api",
				Optional:            true,
			},
		},
	}, nil
}

func (p *provider) Configure(ctx context.Context, request tfsdk.ConfigureProviderRequest, response *tfsdk.ConfigureProviderResponse) {
	if p.configured {
		return
	}

	var data providerData
	diagnostics := request.Config.Get(ctx, &data)
	response.Diagnostics.Append(diagnostics...)
	if response.Diagnostics.HasError() {
		return
	}

	if data.Token.Null {
		if val, ok := os.LookupEnv("CLOUDBIT_TOKEN"); ok {
			data.Token = types.String{Value: val}
		} else {
			response.Diagnostics.AddError(
				"Missing Token",
				"The token is missing. Please set the token in the provider configuration or set the CLOUDBIT_TOKEN environment variable.",
			)
			return
		}
	}

	if data.Endpoint.Null {
		data.Endpoint = types.String{Value: p.defaultEndpoint}

		if val, ok := os.LookupEnv("CLOUDBIT_ENDPOINT"); ok {
			data.Endpoint = types.String{Value: val}
		}
	}

	p.client = goclient.NewClient(
		goclient.WithToken(data.Token.Value),
		goclient.WithBase(data.Endpoint.Value),
		goclient.WithUserAgent(fmt.Sprintf("terraform-provider-cloudbit/%s", p.version)),

		goclient.WithHTTPClientOption(func(c *http.Client) {
			c.Transport = logTransport{base: c.Transport}
		}),
	)

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"cloudbit_compute_certificate":                  computeCertificateResourceType{},
		"cloudbit_compute_elastic_ip":                   computeElasticIPResourceType{},
		"cloudbit_compute_elastic_ip_server_attachment": computeElasticIPServerAttachmentResourceType{},
		"cloudbit_compute_key_pair":                     computeKeyPairResourceType{},
		"cloudbit_compute_load_balancer":                computeLoadBalancerResourceType{},
		"cloudbit_compute_load_balancer_member":         computeLoadBalancerMemberResourceType{},
		"cloudbit_compute_load_balancer_pool":           computeLoadBalancerPoolResourceType{},
		"cloudbit_compute_network":                      computeNetworkResourceType{},
		"cloudbit_compute_network_interface":            computeNetworkInterfaceResourceType{},
		"cloudbit_compute_router":                       computeRouterResourceType{},
		"cloudbit_compute_router_interface":             computeRouterInterfaceResourceType{},
		"cloudbit_compute_router_route":                 computeRouterRouteResourceType{},
		"cloudbit_compute_security_group":               computeSecurityGroupResourceType{},
		"cloudbit_compute_security_group_rule":          computeSecurityGroupRuleResourceType{},
		"cloudbit_compute_server":                       computeServerResourceType{},
		"cloudbit_compute_volume":                       computeVolumeResourceType{},
		"cloudbit_compute_volume_attachment":            computeVolumeAttachmentResourceType{},

		"cloudbit_kubernetes_cluster": kubernetesClusterResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"cloudbit_location": locationDataSourceType{},
		"cloudbit_module":   moduleDataSourceType{},
		"cloudbit_product":  productDataSourceType{},

		"cloudbit_compute_certificate":                     computeCertificateDataSourceType{},
		"cloudbit_compute_elastic_ip":                      computeElasticIPDataSourceType{},
		"cloudbit_compute_image":                           computeImageDataSourceType{},
		"cloudbit_compute_key_pair":                        computeKeyPairDataSourceType{},
		"cloudbit_compute_load_balancer_algorithm":         computeLoadBalancerAlgorithmDataSourceType{},
		"cloudbit_compute_load_balancer_health_check_type": computeLoadBalancerHealthCheckTypeDataSourceType{},
		"cloudbit_compute_load_balancer_member":            computeLoadBalancerMemberDataSourceType{},
		"cloudbit_compute_load_balancer_pool":              computeLoadBalancerPoolDataSourceType{},
		"cloudbit_compute_load_balancer_protocol":          computeLoadBalancerProtocolDataSourceType{},
		"cloudbit_compute_network":                         computeNetworkDataSourceType{},
		"cloudbit_compute_network_interface":               computeNetworkInterfaceDataSourceType{},
		"cloudbit_compute_router":                          computeRouterDataSourceType{},
		"cloudbit_compute_router_interface":                computeRouterInterfaceDataSourceType{},
		"cloudbit_compute_router_route":                    computeRouterRouteDataSourceType{},
		"cloudbit_compute_security_group":                  computeSecurityGroupDataSourceType{},
		"cloudbit_compute_security_group_rule":             computeSecurityGroupRuleDataSourceType{},
		"cloudbit_compute_server":                          computeServerDataSourceType{},
		"cloudbit_compute_snapshot":                        computeSnapshotDataSourceType{},
		"cloudbit_compute_volume":                          computeVolumeDataSourceType{},

		"cloudbit_kubernetes_cluster":     kubernetesClusterDataSourceType{},
		"cloudbit_kubernetes_kube_config": kubernetesKubeConfigDataSourceType{},
	}, nil
}

func convertToLocalProviderType(p tfsdk.Provider) (prov *provider, diagnostics diag.Diagnostics) {
	prov, ok := p.(*provider)
	if !ok {
		diagnostics.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)

		return
	}

	return
}

func waitForCondition(ctx context.Context, check func(ctx context.Context) (bool, diag.Diagnostics)) (diagnostics diag.Diagnostics) {
	done, d := check(ctx)
	diagnostics.Append(d...)
	if done || diagnostics.HasError() {
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

		case <-ctx.Done():
			diagnostics.AddError("Timeout", "Timeout while waiting for condition")
			return
		}

		done, d = check(ctx)
		diagnostics.Append(d...)
		if done || diagnostics.HasError() {
			return
		}
	}
}

type logTransport struct {
	base http.RoundTripper
}

func (l logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	additionalContext := map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	}

	res, err := l.transport().RoundTrip(req)

	if err == nil {
		additionalContext["request_id"] = res.Header.Get("X-Request-ID")

		msg := fmt.Sprintf("request to `%s %s` resulted in `%s`", req.Method, req.URL.String(), res.Status)
		tflog.Trace(req.Context(), msg, additionalContext)
	} else {
		msg := fmt.Sprintf("request to `%s %s` resulted in `%s`", req.Method, req.URL.String(), err)
		tflog.Trace(req.Context(), msg, additionalContext)
	}

	return res, err
}

func (l logTransport) transport() http.RoundTripper {
	if l.base == nil {
		return http.DefaultTransport
	}

	return l.base
}
