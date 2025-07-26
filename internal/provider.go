package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"os"

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
	Port     types.Int32  `tfsdk:"port"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
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

	// HOST
	var host string
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Connot use unknown value as host")
		return
	}

	if config.Host.IsNull() {
		host = os.Getenv("EXTRM_FE_HOST")
	} else {
		host = config.Host.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddError(
			"Unable to find host",
			"Host cannot be an empty string")
		return
	}

	// PORT
	var port int32
	if config.Port.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Connot use unknown value as port")
		return
	}

	if config.Port.IsNull() {
		resp.Diagnostics.AddError(
			"Unable to find port",
			"Port cannot be an null integer")
		return
	} else {
		port = config.Port.ValueInt32()
	}

	if port == 0 {
		resp.Diagnostics.AddError(
			"Unable to find port",
			"Port cannot be an empty integer")
		return
	}

	// Username
	var username string
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Connot use unknown value as username")
		return
	}

	if config.Username.IsNull() {
		username = os.Getenv("EXTRM_FE_USERNAME")
	} else {
		username = config.Username.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddError(
			"Unable to find host",
			"Username cannot be an empty string")
		return
	}

	// Username
	var password string
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Connot use unknown value as password")
		return
	}

	if config.Password.IsNull() {
		password = os.Getenv("EXTRM_FE_PASSWORD")
	} else {
		password = config.Password.ValueString()
	}

	if password == "" {
		resp.Diagnostics.AddError(
			"Unable to find host",
			"Password cannot be an empty string")
		return
	}

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
