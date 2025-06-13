# LightHouse - A Configurable OIDFed Trust Anchor


  ![Logo](./assets/logo_dm.svg#only-dark){ width="300" loading=lazy align=left }
  ![Logo](./assets/logo_lm.svg#only-light){ width="300" loading=lazy align=left } 


LightHouse helps you to navigate the wild and complex sea of OpenID 
Federation. Based on the
[go-oidfed implementation](https://github.com/go-oidfed/lib), LightHouse 
provides an easy to use, flexible, and configurable Trust Anchor, 
Intermediate Authority, Resolver, and / or Trust Mark Issuer.
By deploying LightHouse to your federation, entities will know there now is 
an Entity that will guide them and which they can put their trust in so they 
can safely drop anchor.

The [LightHouse source code](https://github.com/go-oidfed/lighthouse) can 
also be used 
as a starting point to implement your own Trust Anchor based on the 
[go-oidfed library](https://github.com/go-oidfed/lib).

However, the primary goal of lighthouse is to have an easy to set up 
Federation Authority that can be configured according to your needs and 
requirements.


## Read more

<div class="grid cards" markdown>

- [:material-file-cog: Configuration](config)
- [:material-map-check: Features](features)
- [:material-routes: Endpoints](features/endpoints.md)
- [:material-police-badge: Automatic Entity Checks](features/entity_checks.md)
- [:material-marker-check: Trust Marks](features/trustmarks.md)

</div>
