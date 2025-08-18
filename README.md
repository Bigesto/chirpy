## Hello there 👋

This is a tutored project aiming at crafting a **Twitter copycat** (in a very, *very*, simple way).

The database is a local **PostgreSQL** — you'll need to create one to have your own **Chirpy**.  
Requested packages are listed in the `go.mod` file.

---

## API Endpoints

### Health & Metrics
- **GET `/api/healthz`**  
  Returns a confirmation if the server is online.

- **GET `/admin/metrics`**  
  Returns the number of requests the site has received so far. *(No admin authentication required.)*

- **POST `/admin/reset`**  
  Resets the entire database **and** the hits counter.  
  🔒 Requires admin authentication.
  ⚠️ **Be careful** — there are no safeguards.

---

### Users

- **POST `/api/users`**  
  Creates a new user and returns credentials in the JSON response body.  
  **Request body:**
  ```json
  {
    "email": "user@example.com",
    "password": "123passwordtype!"
  }
  ```

- **POST `/api/login`**  
  Logs in a user and returns credentials in the JSON response body.  
  **Request body:** *(same as above)*

- **PUT `/api/users`**  
  Updates a user and returns credentials.  
  **Request body:** *(same as above)*

---

### Chirps

- **POST `/api/chirps`**  
  Posts a new chirp. Requires authentication.  
  **Headers:**
  ```
  Authorization: Bearer <accesstoken>
  ```
  **Request body:**
  ```json
  {
    "body": "Text of your chirp",
    "user_id": "your ID"
  }
  ```

- **GET `/api/chirps`**  
  Retrieves **all** chirps as an array.

- **GET `/api/chirps?author_id={author_id}`**  
  Retrieves all chirps from a specific user.  
  *(Same structure as above, but filtered.)*

- **GET `/api/chirps/{chirpID}`**  
  Retrieves a single chirp by its ID.

- **DELETE `/api/chirps/{chirpID}`**  
  Deletes a chirp by its ID. Requires authentication as the chirp creator.  
  **Headers:**
  ```
  Authorization: Bearer <accesstoken>
  ```

---

### Authentication & Tokens

- **POST `/api/refresh`**  
  Refreshes an access token using a refresh token (from user credentials).  
  Returns the new token in the response body.  
  **Headers:**
  ```
  Authorization: Bearer <refreshtoken>
  ```

- **POST `/api/revoke`**  
  Revokes a specific token.  
  **Headers:**
  ```
  Authorization: Bearer <refreshtokentorevoke>
  ```

---

### Webhooks

- **POST `/api/polka/webhooks`**  
  Upgrades a user to *Red Membership* via our partner **Polka**.  
  **Request body:**
  ```json
  {
    "event": "user.upgraded",
    "data": {
      "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
    }
  }
  ```

| Method | Endpoint                          | Description                                                   | Auth Required |
|--------|------------------------------------|---------------------------------------------------------------|--------------|
| GET    | `/api/healthz`                     | Check if server is online                                     | No           |
| GET    | `/admin/metrics`                   | Get total request count                                       | No           |
| POST   | `/admin/reset`                     | Reset database and hits counter                               | Yes (Admin)  |
| POST   | `/api/users`                       | Create a new user                                             | No           |
| POST   | `/api/login`                       | Log in a user                                                 | No           |
| PUT    | `/api/users`                       | Update a user's credentials                                   | Yes          |
| POST   | `/api/chirps`                      | Create a new chirp                                            | Yes          |
| GET    | `/api/chirps`                      | Get all chirps                                                | No           |
| GET    | `/api/chirps?author_id={author_id}`| Get all chirps by a specific user                             | No           |
| GET    | `/api/chirps/{chirpID}`            | Get a chirp by ID                                             | No           |
| DELETE | `/api/chirps/{chirpID}`            | Delete a chirp by ID (must be the author)                     | Yes          |
| POST   | `/api/refresh`                     | Refresh access token using refresh token                      | Yes          |
| POST   | `/api/revoke`                      | Revoke a specific token                                       | Yes          |
| POST   | `/api/polka/webhooks`              | Upgrade user to *Red Membership* via Polka webhook            | No           |
