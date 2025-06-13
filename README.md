![](logo_dm.svg)

# LightHouse - A Go-Based Trust Anchor / Intermediate Authority / Trust Mark Issuer

LightHouse helps you to navigate the wild and complex sea of OpenID Federation.

LightHouse is a flexible and configurable OpenID Federation Entity. It can be configured and deployed as a 
Trust Anchor / Intermediate Authority / Resolver / Trust Mark Issuer or everything at the same time.
LightHouse uses the [go-oidfed/lib](https://github.com/go-oidfed/lib) oidfed library.

LightHouse also can be used to build your own federation entity on top of the existing implementation.

## Documentation

For more information please refer to the Documentation at https://go-oidfed.github.io/lighthouse/

## Configuration

The configuration of LightHouse is explained in details at https://go-oidfed.github.io/lighthouse/config/.

## Docker Images

Docker images are available at [docker hub under `oidfed/lighthouse`](https://hub.docker.com/r/oidfed/lighthouse).

## Related Activities

- The go oidfed library at https://github.com/go-oidfed/lib contains:
    - The basic go-oidfed library with the core oidfed functionalities.
    - It can be used to build all kind of oidfed capable entities.
    - LightHouse uses this library
- The whoami-rp repository at https://github.com/go-oidfed/whoami-rp contains:
    - A simple - but not very useful - example RP.
- The OFFA repository at https://github.com/go-oidfed/offa:
    - OFFA stands for Openid Federation Forward Auth
    - OFFA can be deployed next to existing services to add oidfed
      authentication to services that do not natively support it.
    - OFFA can be used with Apache, Caddy, NGINX, and Traefik.
