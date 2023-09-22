package helm

import (
	"github.com/khulnasoft-lab/defsec/pkg/scanners/options"
	"github.com/khulnasoft-lab/vul-iac/pkg/scanners/helm/parser"
)

type ConfigurableHelmScanner interface {
	options.ConfigurableScanner
	AddParserOptions(options ...options.ParserOption)
}

func ScannerWithValuesFile(paths ...string) options.ScannerOption {
	return func(s options.ConfigurableScanner) {
		if helmScanner, ok := s.(ConfigurableHelmScanner); ok {
			helmScanner.AddParserOptions(parser.OptionWithValuesFile(paths...))
		}
	}
}

func ScannerWithValues(values ...string) options.ScannerOption {
	return func(s options.ConfigurableScanner) {
		if helmScanner, ok := s.(ConfigurableHelmScanner); ok {
			helmScanner.AddParserOptions(parser.OptionWithValues(values...))
		}
	}
}

func ScannerWithFileValues(values ...string) options.ScannerOption {
	return func(s options.ConfigurableScanner) {
		if helmScanner, ok := s.(ConfigurableHelmScanner); ok {
			helmScanner.AddParserOptions(parser.OptionWithFileValues(values...))
		}
	}
}

func ScannerWithStringValues(values ...string) options.ScannerOption {
	return func(s options.ConfigurableScanner) {
		if helmScanner, ok := s.(ConfigurableHelmScanner); ok {
			helmScanner.AddParserOptions(parser.OptionWithStringValues(values...))
		}
	}
}

func ScannerWithAPIVersions(values ...string) options.ScannerOption {
	return func(s options.ConfigurableScanner) {
		if helmScanner, ok := s.(ConfigurableHelmScanner); ok {
			helmScanner.AddParserOptions(parser.OptionWithAPIVersions(values...))
		}
	}
}