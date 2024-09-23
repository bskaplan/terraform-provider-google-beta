// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccComputeRegionNetworkFirewallPolicyWithRules_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckComputeRegionNetworkFirewallPolicyWithRulesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkFirewallPolicyWithRules_full(context),
			},
			{
				ResourceName:            "google_compute_region_network_firewall_policy_with_rules.region-network-firewall-policy-with-rules",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccComputeRegionNetworkFirewallPolicyWithRules_update(context),
			},
			{
				ResourceName:            "google_compute_region_network_firewall_policy_with_rules.region-network-firewall-policy-with-rules",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccComputeRegionNetworkFirewallPolicyWithRules_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  provider = google-beta
}

resource "google_compute_region_network_firewall_policy_with_rules" "region-network-firewall-policy-with-rules" {
  name        = "tf-test-tf-region-fw-policy-with-rules%{random_suffix}"
  region      = "us-west2"
  description = "Terraform test"
  provider    = google-beta

  rule {
    description    = "tcp rule"
    priority       = 1000
    enable_logging = true
    action         = "allow"
    direction      = "EGRESS"
    match {
      layer4_config {
        ip_protocol = "tcp"
        ports       = [8080, 7070]
      }
      dest_ip_ranges = ["11.100.0.1/32"]
      dest_fqdns = ["www.yyy.com", "www.zzz.com"]
      dest_region_codes = ["HK", "IN"]
      dest_threat_intelligences = ["iplist-search-engines-crawlers", "iplist-tor-exit-nodes"]
      dest_address_groups = [google_network_security_address_group.address_group_1.id]
    }
    target_secure_tag {
      name = "tagValues/${google_tags_tag_value.secure_tag_value_1.name}"
    }
  }
  rule {
      description    = "udp rule"
      rule_name      = "test-rule"
      priority       = 2000
      enable_logging = false
      action         = "deny"
      direction      = "INGRESS"
      match {
        layer4_config {
          ip_protocol = "udp"
        }
        src_ip_ranges = ["0.0.0.0/0"]
        src_fqdns = ["www.abc.com", "www.def.com"]
        src_region_codes = ["US", "CA"]
        src_threat_intelligences = ["iplist-known-malicious-ips", "iplist-public-clouds"]
        src_address_groups = [google_network_security_address_group.address_group_1.id]
        src_secure_tag {
          name = "tagValues/${google_tags_tag_value.secure_tag_value_1.name}"
        }
      }
      disabled = true
    }
}

resource "google_network_security_address_group" "address_group_1" {
  provider  = google-beta
  name        = "tf-test-tf-address-group%{random_suffix}"
  parent      = "projects/${data.google_project.project.name}"
  description = "Regional address group"
  location    = "us-west2"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_tags_tag_key" "secure_tag_key_1" {
  provider   = google-beta
  description = "Tag key"
  parent      = "projects/${data.google_project.project.name}"
  purpose     = "GCE_FIREWALL"
  short_name  = "tf-test-tf-tag-key%{random_suffix}"
  purpose_data = {
    network = "${data.google_project.project.name}/default"
  }
}

resource "google_tags_tag_value" "secure_tag_value_1" {
  provider   = google-beta
  description = "Tag value"
  parent      = "tagKeys/${google_tags_tag_key.secure_tag_key_1.name}"
  short_name  = "tf-test-tf-tag-value%{random_suffix}"
}
`, context)
}

func testAccComputeRegionNetworkFirewallPolicyWithRules_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  provider = google-beta
}

resource "google_compute_region_network_firewall_policy_with_rules" "region-network-firewall-policy-with-rules" {
  name = "tf-test-tf-fw-policy-with-rules%{random_suffix}"
  description = "Terraform test - update"
  region = "us-west2"
  provider = google-beta

  rule {
    description    = "tcp rule - changed"
    priority       = 1000
    enable_logging = false
    action         = "allow"
    direction      = "EGRESS"
    match {
      layer4_config {
        ip_protocol = "tcp"
        ports       = [8080, 7070]
      }
      dest_ip_ranges = ["11.100.0.1/32"]
    }
  }
  rule {
      description    = "new udp rule"
      priority       = 4000
      enable_logging = true
      action         = "deny"
      direction      = "INGRESS"
      match {
        layer4_config {
          ip_protocol = "udp"
        }
        src_ip_ranges = ["0.0.0.0/0"]
        src_fqdns = ["www.abc.com", "www.ghi.com"]
        src_region_codes = ["IT", "FR"]
        src_threat_intelligences = ["iplist-public-clouds"]
        src_address_groups = [google_network_security_address_group.address_group_1.id]
        src_secure_tag {
          name = "tagValues/${google_tags_tag_value.secure_tag_value_1.name}"
        }
      }
      disabled = false
    }
}

resource "google_network_security_address_group" "address_group_1" {
  provider  = google-beta
  name        = "tf-test-tf-address-group%{random_suffix}"
  parent      = "projects/${data.google_project.project.name}"
  description = "Regional address group"
  location    = "us-west2"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_tags_tag_key" "secure_tag_key_1" {
  provider   = google-beta
  description = "Tag key"
  parent      = "projects/${data.google_project.project.name}"
  purpose     = "GCE_FIREWALL"
  short_name  = "tf-test-tf-tag-key%{random_suffix}"
  purpose_data = {
    network = "${data.google_project.project.name}/default"
  }
}

resource "google_tags_tag_value" "secure_tag_value_1" {
  provider   = google-beta
  description = "Tag value"
  parent      = "tagKeys/${google_tags_tag_key.secure_tag_key_1.name}"
  short_name  = "tf-test-tf-tag-value%{random_suffix}"
}
`, context)
}