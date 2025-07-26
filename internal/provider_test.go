package provider

import (
	"os"
	"testing"
)

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
