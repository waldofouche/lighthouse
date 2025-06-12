package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fatih/structs"
	"github.com/go-oidfed/lib"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/go-oidfed/lighthouse"
	"github.com/go-oidfed/lighthouse/internal"
	"github.com/go-oidfed/lighthouse/storage"
)

// Config holds configuration for the entity
type Config struct {
	ServerPort            int                                          `yaml:"server_port"`
	EntityID              string                                       `yaml:"entity_id"`
	LogoURI               string                                       `yaml:"logo_uri"`
	AuthorityHints        []string                                     `yaml:"authority_hints"`
	MetadataPolicyFile    string                                       `yaml:"metadata_policy_file"`
	MetadataPolicy        *oidfed.MetadataPolicies                     `yaml:"-"`
	SigningKeyFile        string                                       `yaml:"signing_key_file"`
	ConfigurationLifetime int64                                        `yaml:"configuration_lifetime"`
	OrganizationName      string                                       `yaml:"organization_name"`
	DataLocation          string                                       `yaml:"data_location"`
	ReadableStorage       bool                                         `yaml:"human_readable_storage"`
	EnableDebugLog        bool                                         `yaml:"enable_debug_log"`
	Endpoints             Endpoints                                    `yaml:"endpoints"`
	TrustMarkSpecs        []extendedTrustMarkSpec                      `yaml:"trust_mark_specs"`
	TrustMarks            []*oidfed.EntityConfigurationTrustMarkConfig `yaml:"trust_marks"`
	TrustMarkIssuers      oidfed.AllowedTrustMarkIssuers               `yaml:"trust_mark_issuers"`
	TrustMarkOwners       oidfed.TrustMarkOwners                       `yaml:"trust_mark_owners"`
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
	for _, tag := range internal.FieldTagNames(s1.Fields(), "yaml") {
		delete(extra, tag)
	}
	for _, tag := range internal.FieldTagNames(s2.Fields(), "yaml") {
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

// Endpoints holds configuration for the different possible endpoints
type Endpoints struct {
	FetchEndpoint                      lighthouse.EndpointConf `yaml:"fetch"`
	ListEndpoint                       lighthouse.EndpointConf `yaml:"list"`
	ResolveEndpoint                    lighthouse.EndpointConf `yaml:"resolve"`
	TrustMarkStatusEndpoint            lighthouse.EndpointConf `yaml:"trust_mark_status"`
	TrustMarkedEntitiesListingEndpoint lighthouse.EndpointConf `yaml:"trust_mark_list"`
	TrustMarkEndpoint                  lighthouse.EndpointConf `yaml:"trust_mark"`
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

var c Config

// Get returns the Config
func Get() Config {
	return c
}

// Load loads the config from the given file
func Load(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(content, &c); err != nil {
		log.Fatal(err)
	}
	if c.EnableDebugLog {
		oidfed.EnableDebugLogging()
	}
	if c.EntityID == "" {
		log.Fatal("entity_id not set")
	}
	if c.SigningKeyFile == "" {
		log.Fatal("signing_key_file not set")
	}
	if c.ConfigurationLifetime == 0 {
		c.ConfigurationLifetime = 24 * 60 * 60
	}
	if c.DataLocation == "" {
		log.Fatal("data_location not set")
	}
	if c.MetadataPolicyFile == "" {
		log.Println("WARNING: metadata_policy_file not set")
	} else {
		policyContent, err := os.ReadFile(c.MetadataPolicyFile)
		if err != nil {
			log.Fatal(err)
		}
		if err = json.Unmarshal(policyContent, &c.MetadataPolicy); err != nil {
			log.Fatal(err)
		}
	}
}

// LoadStorageBackends loads and returns the storage backends for the passed Config
func LoadStorageBackends(c Config) (
	subordinateStorage storage.SubordinateStorageBackend,
	trustMarkedEntitiesStorage storage.TrustMarkedEntitiesStorageBackend, err error,
) {
	if c.ReadableStorage {
		warehouse := storage.NewFileStorage(c.DataLocation)
		subordinateStorage = warehouse.SubordinateStorage()
		trustMarkedEntitiesStorage = warehouse.TrustMarkedEntitiesStorage()
	} else {
		warehouse, err := storage.NewBadgerStorage(c.DataLocation)
		if err != nil {
			return nil, nil, err
		}
		subordinateStorage = warehouse.SubordinateStorage()
		trustMarkedEntitiesStorage = warehouse.TrustMarkedEntitiesStorage()
	}
	log.Println("Loaded storage backend")
	return
}
