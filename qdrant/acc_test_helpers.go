package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testCheckHasListAttr verifies a list attribute has at least one element.
func testCheckHasListAttr(resourceName, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if n, ok := rs.Primary.Attributes[attr+".#"]; !ok || n == "0" {
			return fmt.Errorf("expected %s to have at least one element, got %q", attr, n)
		}
		return nil
	}
}

// testAccCheckListNonEmpty returns a TestCheckFunc that asserts the given resource list attribute is present and non-empty in state.
func testAccCheckListNonEmpty(name, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", name)
		}
		// Terraform stores list length at "<attr>.#"
		if n := rs.Primary.Attributes[attr+".#"]; n == "" || n == "0" {
			return fmt.Errorf("expected %s to be non-empty, got length=%q", attr, n)
		}
		return nil
	}
}

// testAccCheckAttrEqual asserts name1.attr1 == name2.attr2 (string compare).
func testAccCheckAttrEqual(name1, attr1, name2, attr2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, ok := s.RootModule().Resources[name1]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", name1)
		}
		rs2, ok := s.RootModule().Resources[name2]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", name2)
		}
		v1, ok := rs1.Primary.Attributes[attr1]
		if !ok {
			return fmt.Errorf("attribute %s not found on %s", attr1, name1)
		}
		v2, ok := rs2.Primary.Attributes[attr2]
		if !ok {
			return fmt.Errorf("attribute %s not found on %s", attr2, name2)
		}
		if v1 != v2 {
			return fmt.Errorf("expected %s.%s (%q) to equal %s.%s (%q)", name1, attr1, v1, name2, attr2, v2)
		}
		return nil
	}
}
