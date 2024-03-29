### Domain Driven Design

DDD is an approach to designing systems, usually software, that emphasises creating a common
language between domain experts and system builders.

![](./docs/ddd/ddd.jpeg)

> It's a way of thinking about software development that focuses on the business
> problem being solved and the domain of the business problem.
>
> It's not a technology or a methodology. It's a way of thinking about how you build and
> design your software.
>
> - Eric Evans

**Domain** - an area of interest or an area over which a person has control.
> A domain is an arbitrary boundary around some subset of concepts in the universe.
>
> They are subjective and aren't mutually exclusive.

![](./docs/ddd/domains.webp)

When modelling systems we have to choose the most appropriate domain boundaries with which to align our software and
organisational boundaries.

**Subdomain** - a domain within a larger domain.

**Ubitiquitous Language** - a language structured around the domain model and used by
all team members to connect all the activities of the team with the software.

![](./docs/ddd/ubiquitous-language.png)

**Bounded Context** - distinct area within a business domain where concepts, language and
rules are clearly defined and consistent.

![](./docs/ddd/bc.jpeg)

> Ideal way: when single subdomain matches a single bounded context.
>
> But still can have multiple bounded contexts within a single subdomain.

**Context Map** - a high-level map of the bounded contexts and their relationships.

**Entity** - an object that have an identity, which is independent of its attributes
throughout its lifetime.

> E.g: person, bank account, car.

**Value objects** - something that is immutable, has no identifier and is replaceable
by another value object with the same value type.

> We can think about VO as a complex attribute of an entity.
>
> E.g: address (street, city, country), money (amount, currency).
>
> Again, don't have their own ID, represented in the same table with the entity.

| Property       | Entity | Value Object |
|----------------|--------|--------------|
| Identity       | +      | -            |
| Mutability     | +      | -            |
| Lifespan       | +      | -            |
| Replaceability | -      | +            |

**Aggregate** - encapsulates a group of entities and value objects that are treated as a unit.

> Every request of changing the state of the entity must go through the aggregate.
>
> E.g: shopping cart - consists of items, counter of total price (for performance reasons).
>
> Now when we want to add an item to the cart, we need to update the counter of the total price.
>
> Aggregate is often called a consistency/transactional boundary, because he groups together
> multiple entities which require atomic and consistent update.

**Invariant** - a condition that must always be true for the aggregate to be valid.

> If invariant is invalid, our data is inconsistent.
>
> To make aggregate valid always we need to use transactions.

**Eventual consistency** - a consistency model used in distributed computing, that guarantees
that when an update is made to a data store, it will eventually be reflected in all nodes
that store the data, but does not guarantee that all nodes will be updated at the same time.

> E.g: when we add an item to the cart, we don't need to update the counter of the total price
> immediately, we can do it by passing value to the event bus and then asynchronously update
> the counter.
