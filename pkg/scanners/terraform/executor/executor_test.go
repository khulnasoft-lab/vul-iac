package executor

import (
	"context"
	"testing"

	"github.com/khulnasoft-lab/defsec/pkg/providers"
	"github.com/khulnasoft-lab/defsec/pkg/scan"
	"github.com/khulnasoft-lab/defsec/pkg/severity"
	"github.com/khulnasoft-lab/defsec/pkg/terraform"
	"github.com/khulnasoft-lab/vul-iac/pkg/scanners/terraform/parser"
	"github.com/khulnasoft-lab/vul-iac/test/testutil"
	"github.com/khulnasoft-lab/vul-policies/pkg/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var panicRule = scan.Rule{
	Provider:  providers.AWSProvider,
	Service:   "service",
	ShortCode: "abc",
	Severity:  severity.High,
	CustomChecks: scan.CustomChecks{
		Terraform: &scan.TerraformCustomCheck{
			RequiredTypes:  []string{"resource"},
			RequiredLabels: []string{"problem"},
			Check: func(resourceBlock *terraform.Block, _ *terraform.Module) (results scan.Results) {
				if resourceBlock.GetAttribute("panic").IsTrue() {
					panic("This is fine")
				}
				return
			},
		},
	},
}

func Test_PanicInCheckNotAllowed(t *testing.T) {

	reg := rules.Register(panicRule, nil)
	defer rules.Deregister(reg)

	fs := testutil.CreateFS(t, map[string]string{
		"project/main.tf": `
resource "problem" "this" {
	panic = true
}
`,
	})

	p := parser.New(fs, "", parser.OptionStopOnHCLError(true))
	err := p.ParseFS(context.TODO(), "project")
	require.NoError(t, err)
	modules, _, err := p.EvaluateAll(context.TODO())
	require.NoError(t, err)
	results, _, _ := New().Execute(modules)
	assert.Equal(t, len(results.GetFailed()), 0)
}

func Test_PanicInCheckAllowed(t *testing.T) {

	reg := rules.Register(panicRule, nil)
	defer rules.Deregister(reg)

	fs := testutil.CreateFS(t, map[string]string{
		"project/main.tf": `
resource "problem" "this" {
	panic = true
}
`,
	})

	p := parser.New(fs, "", parser.OptionStopOnHCLError(true))
	err := p.ParseFS(context.TODO(), "project")
	require.NoError(t, err)
	modules, _, err := p.EvaluateAll(context.TODO())
	require.NoError(t, err)
	_, _, err = New(OptionStopOnErrors(false)).Execute(modules)
	assert.Error(t, err)
}

func Test_PanicNotInCheckNotIncludePassed(t *testing.T) {

	reg := rules.Register(panicRule, nil)
	defer rules.Deregister(reg)

	fs := testutil.CreateFS(t, map[string]string{
		"project/main.tf": `
resource "problem" "this" {
	panic = true
}
`,
	})

	p := parser.New(fs, "", parser.OptionStopOnHCLError(true))
	err := p.ParseFS(context.TODO(), "project")
	require.NoError(t, err)
	modules, _, err := p.EvaluateAll(context.TODO())
	require.NoError(t, err)
	results, _, _ := New().Execute(modules)
	assert.Equal(t, len(results.GetFailed()), 0)
}

func Test_PanicNotInCheckNotIncludePassedStopOnError(t *testing.T) {

	reg := rules.Register(panicRule, nil)
	defer rules.Deregister(reg)

	fs := testutil.CreateFS(t, map[string]string{
		"project/main.tf": `
resource "problem" "this" {
	panic = true
}
`,
	})

	p := parser.New(fs, "", parser.OptionStopOnHCLError(true))
	err := p.ParseFS(context.TODO(), "project")
	require.NoError(t, err)
	modules, _, err := p.EvaluateAll(context.TODO())
	require.NoError(t, err)

	_, _, err = New(OptionStopOnErrors(false)).Execute(modules)
	assert.Error(t, err)
}