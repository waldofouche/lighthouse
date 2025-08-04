package lighthouse

import (
	"fmt"

	"github.com/go-oidfed/lib"
	"github.com/go-oidfed/lib/apimodel"
	"github.com/gofiber/fiber/v2"
	"tideland.dev/go/slices"
)

// TODO allow limiting the collection endpoint to certain trust anchors

const defaultPagingLimit = 100

// AddEntityCollectionEndpoint adds an entity collection endpoint
func (fed *LightHouse) AddEntityCollectionEndpoint(endpoint EndpointConf) {
	if fed.Metadata.FederationEntity.Extra == nil {
		fed.Metadata.FederationEntity.Extra = make(map[string]interface{})
	}
	fed.Metadata.FederationEntity.Extra["federation_collection_endpoint"] = endpoint.ValidateURL(fed.FederationEntity.EntityID)
	if endpoint.Path == "" {
		return
	}
	fed.server.Get(
		endpoint.Path, func(ctx *fiber.Ctx) error {
			var req apimodel.EntityCollectionRequest
			if err := ctx.QueryParser(&req); err != nil {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("could not parse request parameters: " + err.Error()))
			}
			if req.FromEntityID != "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorUnsupportedParameter("parameter 'from_entity_id' is not yet supported"))
			}
			if req.TrustAnchor == "" {
				req.TrustAnchor = fed.FederationEntity.EntityID
			}
			if req.Limit != 0 {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorUnsupportedParameter("parameter 'limit' is not yet supported"))
			}
			if wantedButNotSupported := slices.Subtract(
				req.EntityClaims, []string{
					"entity_id",
					"entity_types",
					"ui_infos",
					"trust_marks",
				},
			); len(wantedButNotSupported) > 0 {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(
					oidfed.ErrorUnsupportedParameter(
						fmt.Sprintf(
							"parameter 'entity_claims' contains the following unsupported values: %+v",
							wantedButNotSupported,
						),
					),
				)
			}
			if wantedButNotSupported := slices.Subtract(
				req.UIClaims, []string{
					"display_name",
					"description",
					"keywords",
					"logo_uri",
					"policy_uri",
					"information_uri",
				},
			); len(wantedButNotSupported) > 0 {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(
					oidfed.ErrorUnsupportedParameter(
						fmt.Sprintf(
							"parameter 'ui_claims' contains the following unsupported values: %+v",
							wantedButNotSupported,
						),
					),
				)
			}
			collector := oidfed.SimpleEntityCollector{}
			entities := collector.CollectEntities(req)

			res := oidfed.EntityCollectionResponse{
				FederationEntities: entities,
			}
			return ctx.JSON(res)
		},
	)
}
