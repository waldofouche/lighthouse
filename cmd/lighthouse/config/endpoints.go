package config

import (
	"github.com/fatih/structs"
	oidfed "github.com/go-oidfed/lib"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/go-oidfed/lighthouse"
	"github.com/go-oidfed/lighthouse/internal/utils"
)

// Endpoints holds configuration for the different possible endpoints
type Endpoints struct {
	FetchEndpoint                      lighthouse.EndpointConf `yaml:"fetch"`
	ListEndpoint                       lighthouse.EndpointConf `yaml:"list"`
	ResolveEndpoint                    lighthouse.EndpointConf `yaml:"resolve"`
	TrustMarkStatusEndpoint            lighthouse.EndpointConf `yaml:"trust_mark_status"`
	TrustMarkedEntitiesListingEndpoint lighthouse.EndpointConf `yaml:"trust_mark_list"`
	TrustMarkEndpoint                  trustMarkEndpointConf   `yaml:"trust_mark"`
	HistoricalKeysEndpoint             lighthouse.EndpointConf `yaml:"historical_keys"`

	EnrollmentEndpoint        extendedEndpointConfig  `yaml:"enroll"`
	EnrollmentRequestEndpoint lighthouse.EndpointConf `yaml:"enroll_request"`
	TrustMarkRequestEndpoint  lighthouse.EndpointConf `yaml:"trust_mark_request"`
	EntityCollectionEndpoint  lighthouse.EndpointConf `yaml:"entity_collection"`
}

type extendedEndpointConfig struct {
	lighthouse.EndpointConf `yaml:",inline"`
	CheckerConfig           lighthouse.EntityCheckerConfig `yaml:"checker"`
}

type trustMarkEndpointConf struct {
	lighthouse.EndpointConf `yaml:",inline"`
	TrustMarkSpecs          []extendedTrustMarkSpec `yaml:"trust_mark_specs"`
}

type extendedTrustMarkSpec struct {
	CheckerConfig        lighthouse.EntityCheckerConfig `yaml:"checker"`
	oidfed.TrustMarkSpec `yaml:",inline"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (e *extendedTrustMarkSpec) UnmarshalYAML(node *yaml.Node) error {
	type forChecker struct {
		CheckerConfig lighthouse.EntityCheckerConfig `yaml:"checker"`
	}
	mm := e.TrustMarkSpec
	var fc forChecker

	if err := node.Decode(&fc); err != nil {
		return errors.WithStack(err)
	}
	if err := node.Decode(&mm); err != nil {
		return errors.WithStack(err)
	}
	extra := make(map[string]interface{})
	if err := node.Decode(&extra); err != nil {
		return errors.WithStack(err)
	}
	s1 := structs.New(fc)
	s2 := structs.New(mm)
	for _, tag := range utils.FieldTagNames(s1.Fields(), "yaml") {
		delete(extra, tag)
	}
	for _, tag := range utils.FieldTagNames(s2.Fields(), "yaml") {
		delete(extra, tag)
	}
	if len(extra) == 0 {
		extra = nil
	}

	mm.Extra = extra
	e.TrustMarkSpec = mm
	e.CheckerConfig = fc.CheckerConfig
	e.IncludeExtraClaimsInInfo = true
	return nil
}
