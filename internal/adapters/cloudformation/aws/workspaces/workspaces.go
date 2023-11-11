package workspaces

import (
	"github.com/aquasecurity/defsec/pkg/providers/aws/workspaces"
	"github.com/khulnasoft-lab/vul-iac/pkg/scanners/cloudformation/parser"
)

// Adapt ...
func Adapt(cfFile parser.FileContext) workspaces.WorkSpaces {
	return workspaces.WorkSpaces{
		WorkSpaces: getWorkSpaces(cfFile),
	}
}
