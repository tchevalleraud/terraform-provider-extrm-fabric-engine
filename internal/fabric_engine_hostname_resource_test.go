// internal/provider/fabric_engine_hostname_resource_test.go
package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestExampleFunction_Known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::scaffolding::example("testvalue")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact("testvalue"),
					),
				},
			},
		},
	})
}

func TestAccFabricEngineHostnameResource(t *testing.T) {
	host := os.Getenv("EXTRM_FE_HOST")
	port := os.Getenv("EXTRM_FE_PORT")
	user := os.Getenv("EXTRM_FE_USERNAME")
	pass := os.Getenv("EXTRM_FE_PASSWORD")
	if host == "" || user == "" || pass == "" {
		t.Skip("Les variables d’environnement EXTRM_FE_HOST, EXTRM_FE_USERNAME et EXTRM_FE_PASSWORD doivent être définies")
	}
	if port == "" {
		port = "22"
	}

	configCreate := fmt.Sprintf(`
provider "extrm_fabric_engine" {
  host     = "%s"
  port     = %s
  username = "%s"
  password = "%s"
}

resource "extrm_fabric_engine_hostname" "test" {
  hostname = "LAB-VOSS01"
}
`, host, port, user, pass)

	configUpdate := fmt.Sprintf(`
provider "extrm_fabric_engine" {
  host     = "%s"
  port     = %s
  username = "%s"
  password = "%s"
}

resource "extrm_fabric_engine_hostname" "test" {
  hostname = "LAB-VOSS02"
}
`, host, port, user, pass)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		// Les tests d’acceptation doivent être exécutés avec TF_ACC=1.
		Steps: []resource.TestStep{
			{
				// Étape 1 : création initiale du hostname
				Config: configCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("extrm_fabric_engine_hostname.test", "hostname", "LAB-VOSS01"),
					resource.TestCheckResourceAttrSet("extrm_fabric_engine_hostname.test", "id"),
				),
			},
			{
				// Étape 2 : mise à jour du hostname
				Config: configUpdate,
				Check:  resource.TestCheckResourceAttr("extrm_fabric_engine_hostname.test", "hostname", "LAB-VOSS02"),
			},
			{
				// Étape 3 : import pour valider l’état
				ResourceName:      "extrm_fabric_engine_hostname.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
