package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"xtrm-fabric-engine": providerserver.NewProtocol6WithError(New("test")()),
}

var testAccProtoV6ProviderFactoriesWithEcho = map[string]func() (tfprotov6.ProviderServer, error){
	"xtrm-fabric-engine": providerserver.NewProtocol6WithError(New("test")()),
	"echo":               echoprovider.NewProviderServer(),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("EXTRM_FE_HOST"); v == "" {
		t.Fatal("EXTRM_FE_HOST env variable must be set for acceptance tests")
	}
	if v := os.Getenv("EXTRM_FE_USERNAME"); v == "" {
		t.Fatal("EXTRM_FE_USERNAME env variable must be set for acceptance tests")
	}
	if v := os.Getenv("EXTRM_FE_PASSWORD"); v == "" {
		t.Fatal("EXTRM_FE_PASSWORD env variable must be set for acceptance tests")
	}
}
