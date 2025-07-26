package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &ExtrmFabricEngineProvider{}
var _ provider.ProviderWithFunctions = &ExtrmFabricEngineProvider{}
var _ provider.ProviderWithEphemeralResources = &ExtrmFabricEngineProvider{}

type ExtrmFabricEngineProvider struct {
	version string
}

type ExtrmFabricEngineModel struct {
	Host     types.String `tfsdk:"host"`
	Port     types.Int    `tfsdk:"port"`
	Username types.string `tfsdk:"username"`
	Password types.string `tfsdk:"password"`
}

type ExtrmFabricEngineClient struct {
	Host     string
	Port     int32
	Username string
	Password string
}

func (p *ExtrmFabricEngineProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "extrm-fabric-engine"
	resp.Version = p.version
}

func (p *ExtrmFabricEngineProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Host of Fabric Engine device.",
				Optional:            true,
			},
			"port": schema.Int32Attribute{
				MarkdownDescription: "Port for the Fabric Engine device.",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for the Fabric Engine device.",
				Optional:            true,
				Sensitive:           true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for the Fabric Engine device.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ExtrmFabricEngineProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config ExtrmFabricEngineModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values @TODO

	client := &ExtrmFabricEngineClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ExtrmFabricEngineProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ExtrmFabricEngineProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *ExtrmFabricEngineProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *ExtrmFabricEngineProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ExtrmFabricEngineProvider{
			version: version,
		}
	}
}
