---
icon: material/marker-check
---


## Trust Mark Issuance
The issuance of Trust Marks boils down to "if you are on the list of entities
that can obtain this Trust Mark, we will issue the Trust Mark".
Therefore, our Trust Mark Issuer implementation manages for each Trust Mark a 
list of entities that can obtain this Trust Mark.

It is possible to use the [Entity Checks](entity_checks.md) mechanism to
dynamically add entities to that list. I.e. any `EntityChecker` can be used on
the Trust Mark endpoint, resulting in the following behavior of the Trust Mark Issuer:

```mermaid
graph TD
    A[Trust Mark Request] --> B{Subject<br>already<br>in the list?};
    B --> |Yes| C[Trust Mark Issued];
    B --> |No| D{Entity Checks<br>defined?};
    D --> |No| E[No Trust Mark Issued];
    D --> |Yes| F{Evaluate<br>Checks};
    F --> |Negative| E;
    F --> |Positive| G[Add Subject to list];
    G --> C;
```