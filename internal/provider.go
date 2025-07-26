package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ExtrmFabricEngineProvider struct {
	version string
}

type ExtrmFabricEngineModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.string `tfsdk:"username"`
	Password types.string `tfsdk:"password"`
	Port     types.Int32  `tfsdk:"port"`
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
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for the Fabric Engine device.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for the Fabric Engine device.",
				Required:            true,
			},
			"port": schema.Int32Attribute{
				MarkdownDescription: "Port for the Fabric Engine device.",
				Optional:            true,
				Validators: []validator.Int32{
					int32validator.Between(0, 65535),
				},
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
}
