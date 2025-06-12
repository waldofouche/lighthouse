package lighthouse

import (
	"crypto"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-oidfed/lib/oidfedconst"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lestrrat-go/jwx/v3/jwa"

	"github.com/go-oidfed/lib"
	"github.com/go-oidfed/lib/cache"
	"github.com/go-oidfed/lib/unixtime"

	"github.com/go-oidfed/lighthouse/internal"
	"github.com/go-oidfed/lighthouse/storage"
)

const entityConfigurationCachePeriod = 5 * time.Second

// EndpointConf is a type for configuring an endpoint with an internal and external path
type EndpointConf struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// IsSet returns a bool indicating if this endpoint was configured or not
func (c EndpointConf) IsSet() bool {
	return c.Path != "" || c.URL != ""
}

// ValidateURL validates that an external URL is set,
// and if not prefixes the internal path with the passed rootURL and sets it
// at the external url
func (c *EndpointConf) ValidateURL(rootURL string) string {
	if c.URL == "" {
		c.URL, _ = url.JoinPath(rootURL, c.Path)
	}
	return c.URL
}

// LightHouse is a type a that represents a federation entity that can have multiple purposes (TA/IA + TMI, etc.)
type LightHouse struct {
	*oidfed.FederationEntity
	*oidfed.TrustMarkIssuer
	*oidfed.GeneralJWTSigner
	SubordinateStatementsConfig
	server *fiber.App
}

// SubordinateStatementsConfig is a type for setting MetadataPolicies and additional attributes that should go into the
// SubordinateStatements issued by this LightHouse
type SubordinateStatementsConfig struct {
	MetadataPolicies             *oidfed.MetadataPolicies
	SubordinateStatementLifetime int64
	Constraints                  *oidfed.ConstraintSpecification
	CriticalExtensions           []string
	MetadataPolicyCrit           []oidfed.PolicyOperatorName
	Extra                        map[string]any
}

// NewLightHouse creates a new LightHouse
func NewLightHouse(
	entityID string, authorityHints []string, metadata *oidfed.Metadata,
	privateSigningKey crypto.Signer, signingAlg jwa.SignatureAlgorithm, configurationLifetime int64,
	stmtConfig SubordinateStatementsConfig,
) (
	*LightHouse,
	error,
) {
	generalSigner := oidfed.NewGeneralJWTSigner(privateSigningKey, signingAlg)
	fed, err := oidfed.NewFederationEntity(
		entityID, authorityHints, metadata, generalSigner.EntityStatementSigner(), configurationLifetime,
	)
	if err != nil {
		return nil, err
	}
	if fed.Metadata == nil {
		fed.Metadata = &oidfed.Metadata{}
	}
	if fed.Metadata.FederationEntity == nil {
		fed.Metadata.FederationEntity = &oidfed.FederationEntityMetadata{}
	}
	server := fiber.New()
	server.Use(recover.New())
	server.Use(compress.New())
	server.Use(logger.New())
	entity := &LightHouse{
		FederationEntity:            fed,
		TrustMarkIssuer:             oidfed.NewTrustMarkIssuer(entityID, generalSigner.TrustMarkSigner(), nil),
		GeneralJWTSigner:            generalSigner,
		SubordinateStatementsConfig: stmtConfig,
		server:                      server,
	}
	server.Get(
		"/.well-known/openid-federation", func(ctx *fiber.Ctx) error {
			cacheKey := cache.Key(cache.KeyEntityConfiguration, entityID)
			var cached []byte
			set, err := cache.Get(cacheKey, &cached)
			if err != nil {
				ctx.Status(fiber.StatusInternalServerError)
				return ctx.JSON(oidfed.ErrorServerError(err.Error()))
			}
			if set {
				ctx.Set(fiber.HeaderContentType, oidfedconst.ContentTypeEntityStatement)
				return ctx.Send(cached)
			}
			jwt, err := entity.EntityConfigurationJWT()
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(oidfed.ErrorServerError(err.Error()))
			}
			err = cache.Set(cacheKey, jwt, entityConfigurationCachePeriod)
			if err != nil {
				log.Println(err.Error())
			}
			ctx.Set(fiber.HeaderContentType, oidfedconst.ContentTypeEntityStatement)
			return ctx.Send(jwt)
		},
	)
	return entity, nil
}

// HttpHandlerFunc returns an http.HandlerFunc for serving all the necessary endpoints
func (fed LightHouse) HttpHandlerFunc() http.HandlerFunc {
	return adaptor.FiberApp(fed.server)
}

// Listen starts an http server at the specific address for serving all the
// necessary endpoints
func (fed LightHouse) Listen(addr string) error {
	return fed.server.Listen(addr)
}

// CreateSubordinateStatement returns a oidfed.EntityStatementPayload for the passed storage.SubordinateInfo
func (fed LightHouse) CreateSubordinateStatement(subordinate *storage.SubordinateInfo) oidfed.EntityStatementPayload {
	now := time.Now()
	return oidfed.EntityStatementPayload{
		Issuer:             fed.FederationEntity.EntityID,
		Subject:            subordinate.EntityID,
		IssuedAt:           unixtime.Unixtime{Time: now},
		ExpiresAt:          unixtime.Unixtime{Time: now.Add(time.Duration(fed.SubordinateStatementLifetime) * time.Second)},
		SourceEndpoint:     fed.Metadata.FederationEntity.FederationFetchEndpoint,
		JWKS:               subordinate.JWKS,
		Metadata:           subordinate.Metadata,
		MetadataPolicy:     fed.MetadataPolicies,
		Constraints:        fed.Constraints,
		CriticalExtensions: fed.CriticalExtensions,
		MetadataPolicyCrit: fed.MetadataPolicyCrit,
		TrustMarks:         subordinate.TrustMarks,
		Extra:              internal.MergeMaps(true, fed.Extra, subordinate.Extra),
	}
}
