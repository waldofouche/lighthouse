---
icon: material/police-badge
---

# Entity Checks
With the **Entity Checks** mechanism checks on an entity can be defined. 
One can define their own Entity Checks by implementing the `EntityChecker` 
interface and registering it through the `RegisterEntityChecker` function before
loading the config file. 

The following Entity Checks are already implemented and supported out of the 
box by LightHouse:
- [`none`](#none): Always forbids access
- [`trust_mark`](#trust-mark): Checks if the entity advertises a Trust Mark and verifies that it is valid
- [`trust_path`](#trust-path): Checks if there is a valid trust path from the entity to the defined Trust Anchor
- [`authority_hints`](#authority-hints): Checks if the entity's published
  `authority_hints` contain the defined Entity ID
- [`entity_id`](#entity-ids): Checks if the entity's `entity_id` is one of the 
  defined ones
- [`multiple_and`](#multiple): Used to combine multiple `EntityChecker` 
  using AND
- [`multiple_or`](#multiple): Used to combine multiple `EntityChecker` 
  using OR

In the following we describe in more details how to configure the different
Entity Checkers:

## None
No additional configuration applicable.

!!! file "Example"

    ```yaml
    checker:
      type: none
    ```

## Trust Mark
For a Trust Mark Entity Checker one must configure the Trust Mark Type of the
Trust Mark that should be checked. Additionally, one must provide either Trust
Anchors or the Trust Mark Issuer's jwks and in the case of delegation
information about the Trust Mark Owner.

### Config Parameters

| Claim                    | Necessity                                                       | Description                                                  |
|--------------------------|-----------------------------------------------------------------|--------------------------------------------------------------|
| `trust_mark_type`        | REQUIRED                                                        | The Trust Mark Type of the Trust Mark to check               |
| `trust_anchors`          | REQUIRED unless `trust_mark_issuer_jwks` is given               | A list of Trust Anchors used to verify the Trust Mark issuer |
| `trust_mark_issuer_jwks` | REQUIRED if `trust_anchors` is not given                        | The jwks of the Trust Mark Issuer                            |
| `trust_mark_owner`       | REQUIRED if `trust_anchors` is not given and delegation is used | Information about the Trust Mark Owner                       |

The `trust_anchors` claim is a list where each element can have the following
parameters:

| Claim       | Necessity | Description                                                                      |
|-------------|-----------|----------------------------------------------------------------------------------|
| `entity_id` | REQUIRED  | The Entity ID of the Trust Anchor                                                |
| `jwks`      | OPTIONAL  | The Trust Anchor's jwks; if omitted it is obtained from its Entity Configuration |

The `trust_mark_owner` claim has the following parameters:

| Claim       | Necessity | Description                           |
|-------------|-----------|---------------------------------------|
| `entity_id` | REQUIRED  | The Entity ID of the Trust Mark Owner |
| `jwks`      | REQUIRED  | The Trust Mark Owner's jwks           |


### Examples


=== ":material-file-code: Using Trust Anchor"

    ```yaml
    checker:
      type: trust_mark
      config:
        trust_mark_type: https://tm.example.org
        trust_anchors:
          - entity_id: https://ta.example.org
    ```

=== ":material-file-code: Using Trust Mark Issuer JWKS"

    ```yaml
    checker:
      type: trust_mark
      config:
        trust_mark_type: https://tm.example.org
        trust_mark_issuer_jwks: {"keys":[{"alg":"ES512","crv":"P-521","kid":"E6XirVKtuO2_76Ly8Lw1cS_W4FUfw_lx5M_z33aMO-I","kty":"EC","use":"sig","x":"AbZpRmHJVpqqJ2q4bFMPto5jVhReNe0toBHWm0y-AhdpqYIqLA-J3ICr_I42BgmC4pG9lQE4qU8mJjkX1I__PDK8","y":"AFl9aVDzsUJPbyxDe96FuLWJNYNOo68WcljWEXJ0QzsFaTDUtykNe1lf3UoOXQWnvNQ1eD2iyWTef1gRR9A6HOSI"}]}
    ```

=== ":material-file-code: Using Trust Mark Issuer JWKS and delegation"

    ```yaml
    checker:
      type: trust_mark
      config:
        trust_mark_type: https://tm.example.org
        trust_mark_issuer_jwks: {"keys":[{"alg":"ES512","crv":"P-521","kid":"E6XirVKtuO2_76Ly8Lw1cS_W4FUfw_lx5M_z33aMO-I","kty":"EC","use":"sig","x":"AbZpRmHJVpqqJ2q4bFMPto5jVhReNe0toBHWm0y-AhdpqYIqLA-J3ICr_I42BgmC4pG9lQE4qU8mJjkX1I__PDK8","y":"AFl9aVDzsUJPbyxDe96FuLWJNYNOo68WcljWEXJ0QzsFaTDUtykNe1lf3UoOXQWnvNQ1eD2iyWTef1gRR9A6HOSI"}]}
        trust_mark_owner:
          entity_id: https://ta.example.org
          jwks: {"keys":[{"alg":"ES512","crv":"P-521","kid":"gChx94HqIDTscqMzxDps6degt2j_Z7OrDsx0Fc24rKA","kty":"EC","use":"sig","x":"AAyVRMA84JsAtJ9z3qKVzgBN1DL8lDIrHRRYtnYiSkfe-i0V7W21QJ_VBBRF3kWFEYadRL9z4yJC7gYvsojF6p8C","y":"AYx1JCtCfrvNR8x8KibI2mQJKAsszjslfd8WlTha8lxtvncpg5c-UxjJgpCYRo3jwdvxUCa6LKHu0TzbUhKfFK8f"}]}
    ```

## Trust Path

For a trust path Entity Checker one must configure the Trust Anchors that should
be used to verify that there is an existing trust path to one of these Trust
Anchors.

### Config Parameters

| Claim           | Necessity | Description                                           |
|-----------------|-----------|-------------------------------------------------------|
| `trust_anchors` | REQUIRED  | A list of Trust Anchors used to verify the trust path |

The `trust_anchors` claim is a list where each element can have the following
parameters:

| Claim       | Necessity | Description                                                                     |
|-------------|-----------|---------------------------------------------------------------------------------|
| `entity_id` | REQUIRED  | The Entity ID of the Trust Anchor                                               |
| `jwks`      | OPTIONAL  | The Trust Anchors jwks; if omitted it is obtained from its Entity Configuration |


### Example

!!! file "Example"

    ```yaml
    checker:
      type: trust_path
      config:
        trust_anchors:
          - entity_id: https://ta.example.org
    ```

## Authority Hints

For an Authority Hints Entity Checker one must configure the Entity ID that
should be present in the authority hints.

### Config Parameters

| Claim       | Necessity | Description                                                          |
|-------------|-----------|----------------------------------------------------------------------|
| `entity_id` | REQUIRED  | The Entity ID that should be present in the entity's authority hints |

### Example

!!! file "Example"

    ```yaml
    checker:
      type: authority_hints
      config:
        entity_id: https://ia.example.org
    ```

## Entity IDs

For an Entity ID Entity Checker one must configure the Entity ID(s) that
are allowed.

### Config Parameters

| Claim        | Necessity | Description                  |
|--------------|-----------|------------------------------|
| `entity_ids` | REQUIRED  | A list of allowed Entity IDs |

### Example

!!! file "Example"

    ```yaml
    checker:
      type: entity_id
      config:
        entity_ids: 
          - https://op1.example.org
          - https://op2.example.org
    ```

## Multiple
To combine multiple Entity Checkers (either with AND or OR) one must provide all
Entity Checkers:

!!! file "Nested Example"

    ```yaml
    checker:
      type: multiple_and
      config:
        - type: trust_path
          config:
            trust_anchors:
              - entity_id: https://ta.example.org
        - type: multiple_or
          config:
            - type: trust_mark
              config: 
                trust_mark_type: https://tm.example.com
                trust_anchors:
                  - entity_id: https://ta.example.com
            - type: trust_mark
              config: 
                trust_mark_type: https://tm.example.org
                trust_anchors:
                  - entity_id: https://ta.example.org
    ```

