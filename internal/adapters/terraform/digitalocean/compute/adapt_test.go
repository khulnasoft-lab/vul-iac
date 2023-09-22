package compute

import (
	"testing"

	defsecTypes "github.com/khulnasoft-lab/defsec/pkg/types"

	"github.com/khulnasoft-lab/defsec/pkg/providers/digitalocean/compute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft-lab/vul-iac/internal/adapters/terraform/tftestutil"
	"github.com/khulnasoft-lab/vul-iac/test/testutil"
)

func Test_adaptDroplets(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  []compute.Droplet
	}{
		{
			name: "key as data reference",
			terraform: `
			data "digitalocean_ssh_key" "terraform" {
				name = "myKey"
			  }
			  
			resource "digitalocean_droplet" "example" {
				ssh_keys = [ data.digitalocean_ssh_key.myKey.id ]
			}
`,
			expected: []compute.Droplet{
				{
					Metadata: defsecTypes.NewTestMetadata(),
					SSHKeys: []defsecTypes.StringValue{
						defsecTypes.String("", defsecTypes.NewTestMetadata()),
					},
				},
			},
		},
		{
			name: "key as string",
			terraform: `
			data "digitalocean_ssh_key" "terraform" {
				name = "myKey"
			  }
			  
			resource "digitalocean_droplet" "example" {
				ssh_keys = [ "my-ssh-key" ]
			}
`,
			expected: []compute.Droplet{
				{
					Metadata: defsecTypes.NewTestMetadata(),
					SSHKeys: []defsecTypes.StringValue{
						defsecTypes.String("my-ssh-key", defsecTypes.NewTestMetadata()),
					},
				},
			},
		},
		{
			name: "defaults",
			terraform: `
			resource "digitalocean_droplet" "example" {
			}
`,
			expected: []compute.Droplet{
				{
					Metadata: defsecTypes.NewTestMetadata(),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := tftestutil.CreateModulesFromSource(t, test.terraform, ".tf")
			adapted := adaptDroplets(modules)
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func Test_adaptFirewalls(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  []compute.Firewall
	}{
		{
			name: "basic",
			terraform: `
			resource "digitalocean_firewall" "example" {
				droplet_ids = [digitalocean_droplet.web.id]

				outbound_rule {
				  protocol         = "tcp"
				  port_range       = "22"
				  destination_addresses = ["192.168.1.0/24"]
				}
						  
				inbound_rule {
					protocol         = "tcp"
					port_range       = "22"
					source_addresses = ["192.168.1.0/24", "fc00::/7"]
				}
			}
`,
			expected: []compute.Firewall{
				{
					Metadata: defsecTypes.NewTestMetadata(),
					OutboundRules: []compute.OutboundFirewallRule{
						{
							Metadata: defsecTypes.NewTestMetadata(),
							DestinationAddresses: []defsecTypes.StringValue{
								defsecTypes.String("192.168.1.0/24", defsecTypes.NewTestMetadata()),
							},
						},
					},
					InboundRules: []compute.InboundFirewallRule{
						{
							Metadata: defsecTypes.NewTestMetadata(),
							SourceAddresses: []defsecTypes.StringValue{
								defsecTypes.String("192.168.1.0/24", defsecTypes.NewTestMetadata()),
								defsecTypes.String("fc00::/7", defsecTypes.NewTestMetadata()),
							},
						},
					},
				},
			},
		},
		{
			name: "defaults",
			terraform: `
			resource "digitalocean_firewall" "example" {  
			}
`,
			expected: []compute.Firewall{
				{
					Metadata:      defsecTypes.NewTestMetadata(),
					OutboundRules: []compute.OutboundFirewallRule{},
					InboundRules:  []compute.InboundFirewallRule{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := tftestutil.CreateModulesFromSource(t, test.terraform, ".tf")
			adapted := adaptFirewalls(modules)
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func Test_adaptLoadBalancers(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  []compute.LoadBalancer
	}{
		{
			name: "basic",
			terraform: `
			resource "digitalocean_loadbalancer" "example" {

				forwarding_rule {
				  entry_port     = 443
				  entry_protocol = "https"
				
				  target_port     = 443
				  target_protocol = "https"
				}
			  }
`,
			expected: []compute.LoadBalancer{
				{
					Metadata: defsecTypes.NewTestMetadata(),
					ForwardingRules: []compute.ForwardingRule{
						{
							Metadata:      defsecTypes.NewTestMetadata(),
							EntryProtocol: defsecTypes.String("https", defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
		},
		{
			name: "defaults",
			terraform: `
			resource "digitalocean_loadbalancer" "example" {
			  }
`,
			expected: []compute.LoadBalancer{
				{
					Metadata:        defsecTypes.NewTestMetadata(),
					ForwardingRules: nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := tftestutil.CreateModulesFromSource(t, test.terraform, ".tf")
			adapted := adaptLoadBalancers(modules)
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func Test_adaptKubernetesClusters(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  []compute.KubernetesCluster
	}{
		{
			name: "basic",
			terraform: `
			resource "digitalocean_kubernetes_cluster" "example" {
				name   = "foo"
				region = "nyc1"
				version = "1.20.2-do.0"
				surge_upgrade = true
				auto_upgrade = true
			}
`,
			expected: []compute.KubernetesCluster{
				{
					Metadata:     defsecTypes.NewTestMetadata(),
					SurgeUpgrade: defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
					AutoUpgrade:  defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
				},
			},
		},
		{
			name: "defaults",
			terraform: `
			resource "digitalocean_kubernetes_cluster" "example" {
			}
`,
			expected: []compute.KubernetesCluster{
				{
					Metadata:     defsecTypes.NewTestMetadata(),
					SurgeUpgrade: defsecTypes.Bool(false, defsecTypes.NewTestMetadata()),
					AutoUpgrade:  defsecTypes.Bool(false, defsecTypes.NewTestMetadata()),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := tftestutil.CreateModulesFromSource(t, test.terraform, ".tf")
			adapted := adaptKubernetesClusters(modules)
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func TestLines(t *testing.T) {
	src := `
	data "digitalocean_ssh_key" "terraform" {
		name = "myKey"
	}
	  
	resource "digitalocean_droplet" "example" {
		ssh_keys = [ data.digitalocean_ssh_key.myKey.id ]
	}

	resource "digitalocean_firewall" "example" {

		outbound_rule {
		  destination_addresses = ["192.168.1.0/24"]
		}
				  
		inbound_rule {
			source_addresses = ["192.168.1.0/24", "fc00::/7"]
		}
	}

	resource "digitalocean_loadbalancer" "example" {

		forwarding_rule {
		  entry_port     = 443
		  entry_protocol = "https"
		}
	  }

	resource "digitalocean_kubernetes_cluster" "example" {
		name   = "foo"
		surge_upgrade = true
		auto_upgrade = true
	}`

	modules := tftestutil.CreateModulesFromSource(t, src, ".tf")
	adapted := Adapt(modules)

	require.Len(t, adapted.Droplets, 1)
	require.Len(t, adapted.Firewalls, 1)
	require.Len(t, adapted.KubernetesClusters, 1)
	require.Len(t, adapted.LoadBalancers, 1)

	droplet := adapted.Droplets[0]
	firewall := adapted.Firewalls[0]
	cluster := adapted.KubernetesClusters[0]
	loadBalancer := adapted.LoadBalancers[0]

	assert.Equal(t, 6, droplet.Metadata.Range().GetStartLine())
	assert.Equal(t, 8, droplet.Metadata.Range().GetEndLine())

	assert.Equal(t, 7, droplet.SSHKeys[0].GetMetadata().Range().GetStartLine())
	assert.Equal(t, 7, droplet.SSHKeys[0].GetMetadata().Range().GetEndLine())

	assert.Equal(t, 10, firewall.Metadata.Range().GetStartLine())
	assert.Equal(t, 19, firewall.Metadata.Range().GetEndLine())

	assert.Equal(t, 12, firewall.OutboundRules[0].Metadata.Range().GetStartLine())
	assert.Equal(t, 14, firewall.OutboundRules[0].Metadata.Range().GetEndLine())

	assert.Equal(t, 13, firewall.OutboundRules[0].DestinationAddresses[0].GetMetadata().Range().GetStartLine())
	assert.Equal(t, 13, firewall.OutboundRules[0].DestinationAddresses[0].GetMetadata().Range().GetEndLine())

	assert.Equal(t, 16, firewall.InboundRules[0].Metadata.Range().GetStartLine())
	assert.Equal(t, 18, firewall.InboundRules[0].Metadata.Range().GetEndLine())

	assert.Equal(t, 17, firewall.InboundRules[0].SourceAddresses[0].GetMetadata().Range().GetStartLine())
	assert.Equal(t, 17, firewall.InboundRules[0].SourceAddresses[0].GetMetadata().Range().GetEndLine())

	assert.Equal(t, 21, loadBalancer.Metadata.Range().GetStartLine())
	assert.Equal(t, 27, loadBalancer.Metadata.Range().GetEndLine())

	assert.Equal(t, 23, loadBalancer.ForwardingRules[0].Metadata.Range().GetStartLine())
	assert.Equal(t, 26, loadBalancer.ForwardingRules[0].Metadata.Range().GetEndLine())

	assert.Equal(t, 25, loadBalancer.ForwardingRules[0].EntryProtocol.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 25, loadBalancer.ForwardingRules[0].EntryProtocol.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 29, cluster.Metadata.Range().GetStartLine())
	assert.Equal(t, 33, cluster.Metadata.Range().GetEndLine())

	assert.Equal(t, 31, cluster.SurgeUpgrade.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 31, cluster.SurgeUpgrade.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 32, cluster.AutoUpgrade.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 32, cluster.AutoUpgrade.GetMetadata().Range().GetEndLine())
}