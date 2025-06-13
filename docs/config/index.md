# Config
LightHouse is configured through a single configuration file named `config.yaml`.

## Config File Location

LightHouse will search for this file at startup at different locations, the 
first file that is found will be used. Supported locations are:

- `config.yaml`

## Example Config File
The following is an example `config.yaml` file:

??? file "config.yaml"

    ```yaml
    server_port: 8765
    entity_id: "https://go-ia.fedservice.lh"
    authority_hints:
      - "https://trust-anchor.fedservice.lh/"
    signing_key_file: "/data/signing.key"
    organization_name: "GO oidc-fed Intermediate"
    data_location: "/data/data"
    human_readable_storage: true
    metadata_policy_file: "/data/metadata-policy.json"
    endpoints:
      fetch:
        path: "/fetch"
        url: "https://go-ia.fedservice.lh/fetch"
      list:
        path: "/list"
        url: "https://go-ia.fedservice.lh/list"
      resolve:
        path: "/resolve"
        url: "https://go-ia.fedservice.lh/resolve"
      trust_mark:
        path: "/trustmark"
        url: "https://go-ia.fedservice.lh/trustmark"
      trust_mark_status:
        path: "/trustmark/status"
        url: "https://go-ia.fedservice.lh/trustmark/status"
      trust_mark_list:
        path: "/trustmark/list"
        url: "https://go-ia.fedservice.lh/trustmark/list"
      enroll:
        path: "/enroll"
        url: "https://go-ia.fedservice.lh/enroll"
        checker:
            type: trust_mark
            config:
              trust_mark_type: https://go-ia.federservice.lh/tm/federation-member
              trust_anchors:
                - entity_id: https://go-ia.fedservice.lh
    trust_mark_specs:
      - trust_mark_type: "https://go-ia.federservice.lh/tm/federation-member"
        lifetime: 86400
        extra_claim: "example"
        checker:
          type: none
    trust_mark_issuers:
      "https://go-ia.federservice.lh/tm/federation-member":
        - "https://go-ia.fedservice.lh"
    trust_marks:
      - id: "https://go-ia.federservice.lh/tm/federation-member"
        trust_mark: "eyJhbGciOiJFUzUxMiIsImtpZCI6IlpsSFBmQXJTRnFGdjNHRlh3ZUptbmFkZDI4YTM4X3plcEJybEZkWHdIaTQiLCJ0eXAiOiJ0cnVzdC1tYXJrK2p3dCJ9.eyJleHAiOj..."
      - id: "https://trust-anchor.federservice.lh/tm/federation-member"
        trust_mark: "eyJhbGciOiJFUzUxMiIsImtpZCI6InpFLTlhVlhJanJZOUcxVU0tYURQVkxVR1RkWmFuOTk0NlJJUWhraWFjUVkiLCJ0eXAiOiJ0cnVzdC1tYXJrK2p3dCJ9.eyJleHAiO..."
    ```

## Configuration Sections

<div class="grid cards" markdown>

- :fontawesome-solid-person-digging: WIP

[//]: # (- [:material-server-network: Server]&#40;server.md&#41;)
[//]: # (- [:material-script-text: Logging]&#40;logging.md&#41;)
[//]: # (- [:simple-openid: Federation]&#40;federation.md&#41;)
[//]: # (- [:material-security: Auth]&#40;auth.md&#41;)
[//]: # (- [:material-cookie: Sessions]&#40;sessions.md&#41;)
[//]: # (- [:fontawesome-solid-person-digging: `debug_auth`]&#40;debug_auth.md&#41;)

</div>
