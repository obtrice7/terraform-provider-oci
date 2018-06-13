// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"
	"testing"

	"os"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	CompartmentRequiredOnlyResource = CompartmentResourceDependencies + `
resource "oci_identity_compartment" "test_compartment" {
	#Required
	compartment_id = "${var.compartment_id}"
	description = "${var.compartment_description}"
	name = "${var.compartment_name}"
}
`

	CompartmentResourceConfig = CompartmentResourceDependencies + `
resource "oci_identity_compartment" "test_compartment" {
	#Required
	compartment_id = "${var.compartment_id}"
	description = "${var.compartment_description}"
	name = "${var.compartment_name}"

	#Optional
	defined_tags = "${var.compartment_defined_tags}"
	freeform_tags = "${var.compartment_freeform_tags}"
}
`
	CompartmentPropertyVariables = `
variable "compartment_defined_tags" { default = {"example-tag-namespace.example-tag"= "value"} }
variable "compartment_description" { default = "For network components" }
variable "compartment_freeform_tags" { default = {"Department"= "Finance"} }
variable "compartment_name" { default = "Network" }

`
	CompartmentResourceDependencies = DefinedTagsDependencies
)

func TestIdentityCompartmentResource_basic(t *testing.T) {
	provider := testAccProvider
	config := testProviderConfig()

	os.Setenv("TF_VAR_tag_namespace_compartment", getRequiredEnvSetting("compartment_id_for_create"))
	compartmentId := getRequiredEnvSetting("tenancy_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_identity_compartment.test_compartment"
	datasourceName := "data.oci_identity_compartments.test_compartments"

	var resId, resId2 string

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify create
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            config + CompartmentPropertyVariables + compartmentIdVariableStr + CompartmentRequiredOnlyResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "description", "For network components"),
					resource.TestCheckResourceAttr(resourceName, "name", "Network"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},
			// verify create with optionals
			{
				Config: config + CompartmentPropertyVariables + compartmentIdVariableStr + CompartmentResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "For network components"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "Network"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// verify updates to updatable parameters except name.
			// TODO add name updatability when we compartment delete becomes available
			{
				Config: config + `
variable "compartment_defined_tags" { default = {"example-tag-namespace.example-tag"= "updatedValue"} }
variable "compartment_description" { default = "description2" }
variable "compartment_freeform_tags" { default = {"Department"= "Accounting"} }
variable "compartment_name" { default = "Network" }

                ` + compartmentIdVariableStr + CompartmentResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "Network"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config + `
variable "compartment_defined_tags" { default = {"example-tag-namespace.example-tag"= "updatedValue"} }
variable "compartment_description" { default = "description2" }
variable "compartment_freeform_tags" { default = {"Department"= "Accounting"} }
variable "compartment_name" { default = "Network" }

data "oci_identity_compartments" "test_compartments" {
	#Required
	compartment_id = "${var.compartment_id}"

    filter {
    	name = "id"
    	values = ["${oci_identity_compartment.test_compartment.id}"]
    }
}
                ` + compartmentIdVariableStr + CompartmentResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),

					resource.TestCheckResourceAttr(datasourceName, "compartments.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "compartments.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "compartments.0.defined_tags.%", "1"),
					resource.TestCheckResourceAttr(datasourceName, "compartments.0.description", "description2"),
					resource.TestCheckResourceAttr(datasourceName, "compartments.0.freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "compartments.0.id"),
					resource.TestCheckResourceAttr(datasourceName, "compartments.0.name", "Network"),
					resource.TestCheckResourceAttrSet(datasourceName, "compartments.0.state"),
					resource.TestCheckResourceAttrSet(datasourceName, "compartments.0.time_created"),
				),
			},
		},
	})
}
