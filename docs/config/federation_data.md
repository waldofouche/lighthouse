---
icon: material/openid
title: Federation Data
---
<span class="badge badge-red" title="If this option is required or optional">required</span>

Under the `federation_data` option configuration related to OpenID Federation 
is set.

## `entity_id`
<span class="badge badge-purple" title="Value Type">uri</span>
<span class="badge badge-red" title="If this option is required or optional">required</span>

The `entity_id` option is used to set the Federation Entity ID.

??? file "config.yaml"

    ```yaml
    federation_data:
        entity_id: https://lighthouse.example.com
    ```


## `trust_anchors`
<span class="badge badge-purple" title="Value Type">list</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `trust_anchors` option is used to specify the Trust Anchors that should
be used.

??? file "config.yaml"

    ```yaml
    federation_data:
        trust_anchors:
            - entity_id: https://ta.example.com
            - entity_id: https://other-ta.example.org
              jwks: {...}
    ```

For each list element the following options are defined:

### `entity_id`
<span class="badge badge-purple" title="Value Type">uri</span>
<span class="badge badge-red" title="If this option is required or optional">required</span>

The `entity_id` of the Trust Anchor.

### `jwks`
<span class="badge badge-purple" title="Value Type">jwks</span>
<span class="badge badge-green" title="If this option is required or optional">recommended</span>

The `jwks` of the Trust Anchor that was obtained out-of-band. If omitted, it
will be obtained from the Trust Anchor's Entity Configuration and implicitly
trusted. In that case you are trusting TLS.

!!! tip

    We recommend to provide the `jwks` as `json`. `json` is valid `yaml` and 
    can just be included. This way you can pass the whole `jwks` in a single 
    line.

## `authority_hints`
<span class="badge badge-purple" title="Value Type">list of uris</span>
<span class="badge badge-green" title="If this option is required or optional">required, unless there are no superiors</span>

The `authority_hints` option is used to specify the Entity IDs of Federation
Entities that are direct superior to LightHouse and that issue a statement about LightHouse.

??? file "config.yaml"

    ```yaml
    federation_data:
        authority_hints:
            - https://ta.example.com
    ```

## `federation_entity_metadata`
<span class="badge badge-purple" title="Value Type">object / mapping</span>
<span class="badge badge-orange" title="If this option is required or optional">recommended</span>

The `federation_entity_metadata` option is used to set data that should be included in `metadata.federation_entity` 
inside the Entity's Entity Configuration.

The following options are available:

### `display_name`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-orange" title="If this option is required or optional">recommended</span>

The `display_name` option sets the Display Name of this Entity to be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            display_name: Example Trust Anchor
    ```

### `description`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `description` option sets the Description of this Entity to be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            description: "This is the Trust Anchor for the Example Federation."
    ```

### `keywords`
<span class="badge badge-purple" title="Value Type">list of string</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `keywords` option sets Keywords for this Entity that should be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            keywords:
                - TA
                - foo
                - bar
    ```

### `contacts`
<span class="badge badge-purple" title="Value Type">list of string</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `contacts` option sets the Contacts of this Entity to be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            contacts:
                - contact@example.com
    ```

### `logo_uri`
<span class="badge badge-purple" title="Value Type">uri</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `logo_uri` option sets the Logo URI of this Entity to be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            logo_uri: https://static.example.com/ta/logo.png
    ```

### `policy_uri`
<span class="badge badge-purple" title="Value Type">uri</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `policy_uri` option sets the Policy URI for this Entity to be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            policy_uri: https://ta.example.com/policy
    ```

### `information_uri`
<span class="badge badge-purple" title="Value Type">uri</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `information_uri` option sets the Information URI for this Entity to be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            information_uri: https://ta.example.com/about
    ```

### `organization_name`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `organization_name` option sets the Organization Name for this Entity to be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            organization_name: Example Organization
    ```

### `organization_uri`
<span class="badge badge-purple" title="Value Type">uri</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `organization_uri` option sets the Organization URI for this Entity to be included in the Federation Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            organization_uri: https://example.com
    ```

### `extra`
<span class="badge badge-purple" title="Value Type">object / mapping</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `extra` option is used to set additional data that should be included the Federation 
Entity Metadata.

??? file "config.yaml"

    ```yaml
    federation_data:
        federation_entity_metadata:
            extra:
                foo: bar
                level: 2
    ```

## `metadata_policy_file`
<span class="badge badge-purple" title="Value Type">file path</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `metadata_policy_file` option is used to set a metadata policy that is applicable to all subordinates. The 
passed file must contain the Metadata Policy as json per OpenID Federation Specification. 
It is optional to provide this option, but if provided the file must exist and contain valid Metadata Policy.

??? file "config.yaml"

    ```yaml
    federation_data:
        metadata_policy_file: /path/to/metadata-policy.json
    ```


## `constraints`
<span class="badge badge-purple" title="Value Type">object / mapping</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `constraints` option is used to set the Constraints object that should be included in the Entity Configuration. 
The configuration of this option is in line / analogous to how Constraints are defined in the OpenID Federation 
Specification.

??? file "config.yaml"

    ```yaml
    federation_data:
        constraints:
            max_path_len: 2
            naming_constraints:
                permitted:
                    - .example.com
                excluded:
                    - east.example.com
            allowed_entity_types:
                - openid_provider
                - openid_relying_party
    ```

## `crit`
<span class="badge badge-purple" title="Value Type">list of strings</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `crit` option is used to set which additional claims are critical as per OpenID Federation Specification.


??? file "config.yaml"

    ```yaml
    federation_data:
        crit:
            - foobar
    ```

## `metadata_policy_crit`
<span class="badge badge-purple" title="Value Type">list of strings</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `crit` option is used to set which additional metadata policy operators are critical as per OpenID Federation 
Specification.


??? file "config.yaml"

    ```yaml
    federation_data:
        metadata_policy_crit:
            - remove
    ```

## `trust_marks`
<span class="badge badge-purple" title="Value Type">list of trust mark configs</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `trust_marks` option is used to set Trust Marks (about LightHouse) that should be published
in the Entity Configuration.

??? file "config.yaml"

    ```yaml
    federation_data:
        trust_marks:
            - trust_mark_type: https://example.com/tm
              trust_mark_issuer: https://example.com/tmi
              refresh: true
              min_lifetime: 300
              refresh_grace_period: 7200
    ```

Each Trust Mark Config has the following options defined:

### `trust_mark_type`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-red" title="If this option is required or optional">required</span>

The `trust_mark_type` option sets the Identifier for the type of this Trust
Mark.

### `trust_mark_issuer`
<span class="badge badge-purple" title="Value Type">uri</span>
<span class="badge badge-red" title="If this option is required or optional">required if `trust_mark_jwt` not given</span>

The `trust_mark_issuer` option is used to set the Entity ID of the Trust
Mark Issuer of this Trust Mark.

Either a Trust Mark JWT (`trust_mark_jwt`) must be given or the Trust Mark
Issuer (`trust_mark_issuer`).

If this option is given, [`refresh`](#refresh) will be set to `true` and LightHouse
will obtain Trust Mark JWTs for this Trust Mark Type dynamically.

### `trust_mark_jwt`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-red" title="If this option is required or optional">required if `trust_mark_issuer` not given</span>

The `trust_mark_jwt` option is used to set a Trust Mark JWT string. This
will be published in the Entity Configuration.
If the set Trust Mark JWT expires, it either must be manually updated before
expiration, or automatic refreshing must be enabled through the [`refresh`](#refresh)
option.

### `refresh`
<span class="badge badge-purple" title="Value Type">boolean</span>
<span class="badge badge-blue" title="Default Value">`false`</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `refresh` option indicates if this Trust Mark should automatically be
refreshed. If set to `true`, LightHouse will fetch a new Trust Mark JWT from
the Trust Mark Issuer when the
old one expires, assuring that always a valid Trust Mark JWT is published in
the Entity Configuration.

### `min_lifetime`
<span class="badge badge-purple" title="Value Type">[duration](index.md#time-duration-configuration-options)</span>
<span class="badge badge-blue" title="Default Value">10 seconds</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `min_lifetime` option is used to set a minimum lifetime in seconds on
this Trust Mark. If [`refresh`](#refresh) is set to `true` LightHouse will assure
that the Trust Mark JWT published in the Entity Configuration will not
expire before this lifetime whenever an Entity Configuration is requested.

### `refresh_grace_period`
<span class="badge badge-purple" title="Value Type">[duration](index.md#time-duration-configuration-options)</span>
<span class="badge badge-blue" title="Default Value">1 hour</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `refresh_grace_period` option is used to set a grace period given in
seconds. If [`refresh`](#refresh) is
set to `true`, LightHouse checks if the Trust Mark expires within the defined grace
period, whenever its Entity Configuration is requested. If the Trust Mark
expires within the grace period the old (but still valid) Trust Mark JWT
will still be included in the Entity Configuration, but in parallel LightHouse
will refresh it by requesting a new Trust Mark JWT from the Trust Mark Issuer.

This allows LightHouse to proactively request Trust Mark JWTs that are expiring
soon in the background.

## `trust_mark_issuers`
<span class="badge badge-purple" title="Value Type">object / mapping</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `trust_mark_issuers` option is used to set the allowed trust mark issuers within this federation. The 
configuration of this option is in line with the format in the OpenID Federation Specification.


??? file "config.yaml"

    ```yaml
    federation_data:
        trust_mark_issuers:
           "https://openid.net/certification/op": []
            "https://refeds.org/sirtfi":
                - https://example.org
    ```


## `trust_mark_owners`
<span class="badge badge-purple" title="Value Type">object / mapping</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `trust_mark_owners` option is used to set the trust mark owners recognized within this federation. The
configuration of this option is in line with the format in the OpenID Federation Specification.


??? file "config.yaml"

    ```yaml
    federation_data:
        trust_mark_owners:
            "https://refeds.org/sirtfi":
                sub: https://refeds.org
                jwks: {"keys":[{"alg":"RS256","e":"AQAB","kid":"key1","kty":"RSA","n":"pnXBOusEANuug6ewezb9J_...","use":"sig"}]}
    ```

## `extra_entity_configuration_data`
<span class="badge badge-purple" title="Value Type">object / mapping</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `extra_entity_configuration_data` option is used to set additional data that should be included in the 
Entity Configuration.

??? file "config.yaml"

    ```yaml
    federation_data:
        extra_entity_configuration_data:
            foo: bar
            level: 2
    ```

## `configuration_lifetime`
<span class="badge badge-purple" title="Value Type">[duration](index.md#time-duration-configuration-options)</span>
<span class="badge badge-blue" title="Default Value">1 day</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `configuration_lifetime` option sets the lifetime of Entity Configurations, i.e. this options defines for how long 
the Entity Configuration JWTs are valid.


??? file "config.yaml"

    ```yaml
    federation_data:
        configuration_lifetime: 1w
    ```
`
