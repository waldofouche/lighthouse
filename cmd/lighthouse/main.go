package main

import (
	"os"

	"github.com/go-oidfed/lib"
	"github.com/lestrrat-go/jwx/v3/jwa"
	log "github.com/sirupsen/logrus"

	"github.com/go-oidfed/lighthouse"
	"github.com/go-oidfed/lighthouse/cmd/lighthouse/config"
	"github.com/go-oidfed/lighthouse/internal/logger"
)

func main() {
	configFile := "config.yaml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	config.Load(configFile)
	logger.Init()
	log.Info("Loaded Config")
	c := config.Get()
	initKey()
	log.Println("Loaded signing key")
	for _, tmc := range c.Federation.TrustMarks {
		if err := tmc.Verify(
			c.Federation.EntityID, c.Endpoints.TrustMarkEndpoint.ValidateURL(c.Federation.EntityID),
			oidfed.NewTrustMarkSigner(signingKey, jwa.ES512()),
		); err != nil {
			log.Fatal(err)
		}
	}

	subordinateStorage, trustMarkedEntitiesStorage, err := config.LoadStorageBackends(c.Storage)
	if err != nil {
		log.Fatal(err)
	}

	lh, err := lighthouse.NewLightHouse(
		c.Federation.EntityID, c.Federation.AuthorityHints,
		&oidfed.Metadata{
			FederationEntity: &oidfed.FederationEntityMetadata{
				OrganizationName: c.Federation.OrganizationName,
				LogoURI:          c.Federation.LogoURI,
			},
		},
		signingKey, jwa.ES512(), c.Federation.ConfigurationLifetime, lighthouse.SubordinateStatementsConfig{
			MetadataPolicies:             nil,
			SubordinateStatementLifetime: 3600,
			// TODO read all of this from config or a storage backend
		}, nil,
	)
	if err != nil {
		panic(err)
	}

	lh.MetadataPolicies = c.Federation.MetadataPolicy
	// TODO other constraints etc.

	lh.TrustMarkIssuers = c.Federation.TrustMarkIssuers
	lh.TrustMarkOwners = c.Federation.TrustMarkOwners
	lh.TrustMarks = c.Federation.TrustMarks

	var trustMarkCheckerMap map[string]lighthouse.EntityChecker
	if specLen := len(c.Endpoints.TrustMarkEndpoint.TrustMarkSpecs); specLen > 0 {
		specs := make([]oidfed.TrustMarkSpec, specLen)
		for i, s := range c.Endpoints.TrustMarkEndpoint.TrustMarkSpecs {
			specs[i] = s.TrustMarkSpec
			if s.CheckerConfig.Type != "" {
				if trustMarkCheckerMap == nil {
					trustMarkCheckerMap = make(map[string]lighthouse.EntityChecker)
				}
				trustMarkCheckerMap[s.TrustMarkType], err = lighthouse.EntityCheckerFromEntityCheckerConfig(
					s.CheckerConfig,
				)
				if err != nil {
					panic(err)
				}
			}
		}
		lh.TrustMarkIssuer = oidfed.NewTrustMarkIssuer(
			c.Federation.EntityID, lh.GeneralJWTSigner.TrustMarkSigner(),
			specs,
		)
	}
	log.Println("Initialized Entity")

	if endpoint := c.Endpoints.FetchEndpoint; endpoint.IsSet() {
		lh.AddFetchEndpoint(endpoint, subordinateStorage)
	}
	if endpoint := c.Endpoints.ListEndpoint; endpoint.IsSet() {
		lh.AddSubordinateListingEndpoint(endpoint, subordinateStorage, trustMarkedEntitiesStorage)
	}
	if endpoint := c.Endpoints.ResolveEndpoint; endpoint.IsSet() {
		lh.AddResolveEndpoint(endpoint)
	}
	if endpoint := c.Endpoints.TrustMarkStatusEndpoint; endpoint.IsSet() {
		lh.AddTrustMarkStatusEndpoint(endpoint, trustMarkedEntitiesStorage)
	}
	if endpoint := c.Endpoints.TrustMarkedEntitiesListingEndpoint; endpoint.IsSet() {
		lh.AddTrustMarkedEntitiesListingEndpoint(endpoint, trustMarkedEntitiesStorage)
	}
	if endpoint := c.Endpoints.TrustMarkEndpoint; endpoint.IsSet() {
		lh.AddTrustMarkEndpoint(endpoint.EndpointConf, trustMarkedEntitiesStorage, trustMarkCheckerMap)
	}
	if endpoint := c.Endpoints.TrustMarkRequestEndpoint; endpoint.IsSet() {
		lh.AddTrustMarkRequestEndpoint(endpoint, trustMarkedEntitiesStorage)
	}
	if endpoint := c.Endpoints.EnrollmentEndpoint; endpoint.IsSet() {
		var checker lighthouse.EntityChecker
		if checkerConfig := endpoint.CheckerConfig; checkerConfig.Type != "" {
			checker, err = lighthouse.EntityCheckerFromEntityCheckerConfig(checkerConfig)
			if err != nil {
				panic(err)
			}
		}
		lh.AddEnrollEndpoint(endpoint.EndpointConf, subordinateStorage, checker)
	}
	if endpoint := c.Endpoints.EnrollmentRequestEndpoint; endpoint.IsSet() {
		lh.AddEnrollRequestEndpoint(endpoint, subordinateStorage)
	}
	if endpoint := c.Endpoints.EntityCollectionEndpoint; endpoint.IsSet() {
		lh.AddEntityCollectionEndpoint(endpoint)
	}
	log.Info("Added Endpoints")

	lh.Start(config.Get().Server)
}
