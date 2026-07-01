# 🗺️ Quick Navigation & Dev Cheat Sheet

## 1. Where Things Live (The Map)
* **`cmd/api/main.go`** ➡️ The starting point. Run this to boot the app.
* **`internal/domain/`** ➡️ **DB Models.** What your data looks like in MongoDB (`bson` tags).
* **`internal/dto/`** ➡️ **JSON Shapes.** What the frontend sends/receives (`json` tags).
* **`internal/repository/`** ➡️ **DB Queries.** Where you write MongoDB CRUD operations.
* **`internal/service/`** ➡️ **Business Logic.** The brain. Validation, rules, and logic go here.
* **`internal/handler/`** ➡️ **API Endpoints.** Takes the HTTP request, calls the Service, returns JSON.
* **`internal/routes/`** ➡️ **URLs.** Maps an endpoint (e.g., `POST /products`) to a Handler.
* **`internal/app/app.go`** ➡️ **The Wiring.** Where you connect Repos, Services, and Handlers together.

---

## 2. The Development Flow (How to Build a Feature)
When adding a new feature, always build from the database layer up to the router. Create or edit your files in this exact order:

**1. The Database Model** 
*   **File:** `internal/domain/product.go`
*   **Action:** Define the exact Go `struct` and `bson` tags for how the data looks in MongoDB.

**2. The API Payloads** 
*   **File:** `internal/dto/product.go`
*   **Action:** Define the `json` shapes for what the frontend will send (Requests) and receive (Responses).

**3. The Database Queries**
*   **File:** `internal/repository/product_repository.go`
*   **Action:** Write the MongoDB operations (`InsertOne`, `Find`, `Update`). 

**4. The Business Logic** 
*   **File:** `internal/service/product_service.go`
*   **Action:** Map the DTOs to the Domain model, enforce rules, generate slugs, and call the Repository.

**5. The HTTP Endpoint** 
*   **File:** `internal/handler/product_handler.go`
*   **Action:** Bind the incoming JSON to your DTO, call the Service, and return an HTTP status code (200, 400).

**6. The URLs & Wiring** 
*   **Files:** `internal/routes/public_routes.go` & `internal/app/app.go`
*   **Action:** Attach a URL path to your Handler, and inject the new Repo and Service into your app's startup process.

---

## 3. Adding a New Feature (5-Step Checklist)
*Example: Adding a `Review` model.*

- [ ] **1. DB Shape:** Create `domain/review.go` (Use `bson` tags).
- [ ] **2. API Shape:** Create `dto/review.go` (Use `json` tags for incoming payloads).
- [ ] **3. Queries:** Create `repository/review_repository.go` (Write `InsertOne`, `Find`, etc.).
- [ ] **4. Logic:** Create `service/review_service.go` (Map DTO -> Domain, enforce rules, call Repo).
- [ ] **5. Endpoint:** Create `handler/review_handler.go` (Parse HTTP JSON, call Service).
- [ ] **6. Expose & Wire:** Add the URL in `routes/` and link the new Repo/Service/Handler in `app/app.go`.
