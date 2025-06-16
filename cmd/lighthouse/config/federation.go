package config

import (
	"encoding/json"

	oidfed "github.com/go-oidfed/lib"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/go-oidfed/lighthouse/internal/utils/fileutils"
)

type federationConf struct {
	EntityID           string              `yaml:"entity_id"`
	ClientName         string              `yaml:"client_name"`
	LogoURI            string              `yaml:"logo_uri"`
	Scopes             []string            `yaml:"scopes"`
	TrustAnchors       oidfed.TrustAnchors `yaml:"trust_anchors"`
	AuthorityHints     []string            `yaml:"authority_hints"`
	OrganizationName   string              `yaml:"organization_name"`
	KeyStorage         string              `yaml:"key_storage"`
	UseResolveEndpoint bool                `yaml:"use_resolve_endpoint"`

	MetadataPolicyFile    string                                       `yaml:"metadata_policy_file"`
	MetadataPolicy        *oidfed.MetadataPolicies                     `yaml:"-"`
	TrustMarks            []*oidfed.EntityConfigurationTrustMarkConfig `yaml:"trust_marks"`
	TrustMarkIssuers      oidfed.AllowedTrustMarkIssuers               `yaml:"trust_mark_issuers"`
	TrustMarkOwners       oidfed.TrustMarkOwners                       `yaml:"trust_mark_owners"`
	ConfigurationLifetime int64                                        `yaml:"configuration_lifetime"`
}

var defaultFederationConf = federationConf{
	ConfigurationLifetime: 24 * 3600,
}

func (c *federationConf) validate() error {
	if c.EntityID == "" {
		return errors.New("error in federation conf: entity_id must be specified")
	}
	if c.MetadataPolicyFile == "" {
		log.Warn("federation conf: metadata_policy_file not set")
	} else {
		policyContent, err := fileutils.ReadFile(c.MetadataPolicyFile)
		if err != nil {
			return errors.Wrap(err, "error reading metadata_policy file")
		}
		if err = json.Unmarshal(policyContent, &c.MetadataPolicy); err != nil {
			return errors.Wrap(err, "error unmarshalling metadata_policy")
		}
	}
	return nil
}
