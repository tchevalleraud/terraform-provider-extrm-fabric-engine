package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	provider2 "github.com/tchevalleraud/extrm-fabric-engine/internal"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"xtrm-fabric-engine": providerserver.NewProtocol6WithError(provider2.New("test")()),
}
