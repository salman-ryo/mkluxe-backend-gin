A Repository is us writing a wrapper around the db client so in future if we change the db from Mongo to SQL we don't have to change all the business service logic to use those, instead we can update the repository.

## Without a repository

Imagine your service looks like this:

```go
func (s *CategoryService) CreateCategory(ctx context.Context, cat *domain.Category) error {
    _, err := s.db.Collection("categories").InsertOne(ctx, cat)
    return err
}
```

Now your service knows:

* MongoDB collections
* `InsertOne()`
* BSON
* MongoDB query syntax

Your business logic and database logic are mixed together.

If you later switched to PostgreSQL, you'd have to rewrite every service that talks to MongoDB.

---

## With a repository

Instead:

```go
func (s *CategoryService) CreateCategory(ctx context.Context, cat *domain.Category) error {
    return s.categoryRepo.Create(ctx, cat)
}
```

The service doesn't know **how** the category is stored.

It only knows:

> "Ask the repository to save it."

The repository handles the implementation details.

---

## Think of it like this

The service speaks in business language:

```text
CreateCategory()

DeleteCategory()

FindCategoryBySlug()
```

The repository translates that into database language.

For MongoDB:

```go
collection.InsertOne(...)
collection.FindOne(...)
collection.DeleteOne(...)
```

For PostgreSQL it might become:

```sql
INSERT INTO categories ...

SELECT * FROM categories ...

DELETE FROM categories ...
```

The service code doesn't change.

---

## It also centralizes database code

Suppose ten different services need to look up a category by slug.

Without a repository:

```text
Service A
    Mongo query

Service B
    Mongo query

Service C
    Mongo query

Service D
    Mongo query
```

Now the same query is duplicated everywhere.

With a repository:

```text
Service A
      │
Service B
      │
Service C
      │
Service D
      │
      ▼
repo.GetBySlug()
```

If you optimize the query later, you only change one place.

---

## Another benefit: testing

Repositories also make testing much easier.

Suppose your service is:

```go
func (s *CategoryService) PublishCategory(...) {
    cat, err := s.repo.GetByID(...)
    ...
}
```

During tests, you don't need a real MongoDB instance.

You can create a fake repository:

```go
type FakeCategoryRepository struct{}

func (f *FakeCategoryRepository) GetByID(...) (*domain.Category, error) {
    return &domain.Category{
        Name: "Shoes",
    }, nil
}
```

Now you're testing your business logic without involving the database.

---

## Is the goal always to swap databases?

Interestingly, **not usually**.

People often say:

> "It's so we can switch databases."

That's true, but in practice it's a relatively rare reason.

The more valuable benefits are:

* **Separation of concerns**: Services contain business rules, repositories contain persistence logic.
* **Cleaner code**: MongoDB-specific code lives in one place.
* **Reusability**: Common queries are written once and reused.
* **Testability**: Services can be tested by mocking or faking the repository.

So I like to think of it this way:

> The repository is an **adapter** between your application's business logic and your database.

Your service shouldn't care whether the data comes from:

* MongoDB,
* PostgreSQL,
* Redis,
* a REST API,
* or even a JSON file.

It just asks:

```go
category, err := repo.GetByID(ctx, id)
```

and lets the repository figure out how to fulfill that request. That's the core idea behind the repository pattern.
