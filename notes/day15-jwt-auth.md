# Day 15 – JWT Authentication

## What we built

| File | Change |
|------|--------|
| `auth.go` | `jwtSecret` constant + `authMiddleware` |
| `handlers.go` | `loginHandler` — verifies credentials, signs and returns a JWT |
| `routes.go` | `POST /login` public; write routes protected via `app.Group` |
| `handlers_test.go` | `getAuthToken` helper, `TestLoginHandler`, updated write-route tests |

---

## JWT fundamentals

A JWT is three base64-encoded parts joined by dots:

```
Header.Payload.Signature
```

- **Header** — algorithm + token type (`{"alg":"HS256","typ":"JWT"}`)
- **Payload** — claims: who, when it expires, custom fields (`username`, `exp`)
- **Signature** — HMAC(Header + Payload, secret) — proves authenticity

The server never stores the token. It just re-verifies the signature on every request.

---

## Creating a token (`loginHandler`)

```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "username": req.Username,
    "exp":      time.Now().Add(time.Hour * 24).Unix(),
})
tokenString, err := token.SignedString([]byte(jwtSecret))
```

- `NewWithClaims` builds the token in memory (unsigned)
- `SignedString` serializes + signs it → produces the final `xxx.yyy.zzz` string
- `exp` must be a Unix timestamp (integer seconds) — use `.Unix()`
- Secret must be `[]byte` because HMAC operates on raw bytes

---

## Validating a token (`authMiddleware`)

```go
token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method")
    }
    return []byte(jwtSecret), nil
})

if err != nil || !token.Valid {
    return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
}
```

**Why a callback for the key?**
The library calls it after parsing the header but before verifying the signature, so you can inspect the token and choose the right key. Also supports multi-key setups.

**Why check `token.Method`?**
Prevents the `alg:none` attack — an attacker could forge a token with no signature if you don't explicitly assert the algorithm.

---

## Go syntax refresher

### Type assertion (comma-ok form)
```go
value, ok := someInterface.(ConcreteType)
// ok = true  → holds that type
// ok = false → does not (never panics)
```

### If with init statement
```go
if statement; condition { }
// statement runs first, variable scoped to the if block
```

Combined:
```go
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
    // token is NOT signed with HMAC → reject
}
```

---

## Route protection patterns

### Per-route (fine-grained)
```go
app.Post("/tasks", authMiddleware, createTaskHandler)
```

### Group (cleaner for multiple routes sharing middleware)
```go
tasks := app.Group("/tasks", authMiddleware)
tasks.Post("/", createTaskHandler)
tasks.Put("/:id", updateTaskHandler)
tasks.Delete("/:id", deleteTaskHandler)
```

GET routes registered separately remain public.

---

## Where middleware lives

| Type | Location | Registered via |
|------|----------|----------------|
| Global (logger, CORS, request ID) | `main.go` | `app.Use()` |
| Feature-specific (auth) | own file (`auth.go`) | route-level or group |

---

## HTTP status codes — 400 vs 401

| Code | Meaning | Example |
|------|---------|---------|
| `400 Bad Request` | Request is structurally broken | Unparseable JSON |
| `401 Unauthorized` | Request is valid but credentials are wrong/missing | Empty or bad password |

Sending `{}` to `/login` → valid JSON, wrong credentials → **401**, not 400.

---

## Testing authenticated routes

Add a helper to avoid repeating login logic:

```go
func getAuthToken(t *testing.T, app *fiber.App) string {
    t.Helper()
    body := `{"username":"admin","password":"secret"}`
    req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    var result map[string]string
    readBody(t, resp, &result)
    return result["token"]
}
```

Use it in any test that hits a protected route:
```go
token := getAuthToken(t, app)
req.Header.Set("Authorization", "Bearer "+token)
```
