# Ayomendaki Admin — Implementation Plan + Audit + Unit Tests

## Overview

Complete implementation plan covering:
1. Audit all existing features (smoke tests)
2. User Management (RBAC: superadmin / admin / user)
3. Comprehensive unit tests (117+ test cases)
4. Run, fix, commit, push

---

## Phase A: Audit — Smoke Test Semua Route

**File:** `internal/handler/smoke_test.go`

### Route Matrix

| Group | Routes | Method |
|-------|--------|--------|
| Auth | `/login`, `/logout` | GET, POST |
| Dashboard | `/` | GET |
| Trips | `/trips`, `/trips/{id}` | GET, POST, PUT, DELETE |
| Schedules | `/trips/{id}/schedules`, `/schedules/{id}` | GET, POST, PUT, DELETE |
| Bookings | `/bookings`, `/bookings/{id}/status` | GET, PATCH |
| Meeting Points | `/meeting-points` | GET, POST, PUT, DELETE |
| Packages | `/packages` | GET, POST, PUT, DELETE |
| Reports | `/reports/revenue` | GET |
| Profile | `/profile` | GET, PUT |
| 404 | `/nonexistent` | GET |

**~33 smoke tests.** Hit every endpoint → verify 200/303.

---

## Phase B: User Management (Fitur Baru)

### Schema

```sql
-- Add role to operators table
role TEXT NOT NULL DEFAULT 'admin'
-- Values: 'superadmin', 'admin', 'user'
```

### Seed Data

| Username | Password | Role | Description |
|----------|----------|------|-------------|
| superadmin | admin123 | superadmin | Full access |
| admin | admin123 | admin | Manage trips, bookings, etc |
| user | user123 | user | Reports & profile only |

### Role Access Matrix

| Feature | superadmin | admin | user |
|---------|-----------|-------|------|
| Dashboard | ✅ | ✅ | ❌ |
| Trips CRUD | ✅ | ✅ | ❌ |
| Schedules CRUD | ✅ | ✅ | ❌ |
| Bookings CRUD | ✅ | ✅ | ❌ |
| Packages CRUD | ✅ | ✅ | ❌ |
| Meeting Points CRUD | ✅ | ✅ | ❌ |
| Reports | ✅ | ✅ | ✅ |
| Profile | ✅ | ✅ | ✅ |
| Users CRUD | ✅ | ❌ | ❌ |

### Files

| File | Action |
|------|--------|
| `internal/database/db.go` | Add `role` column to `operators` table |
| `internal/database/seed.go` | Seed 3 users with different roles |
| `internal/model/operator.go` | Add `Role` field to `Operator` |1
| `internal/common/common.go` | Add `GetRole()` + role constants |
| `internal/server/middleware.go` | Add `roleMiddleware()` for RBAC |
| `internal/server/server.go` | Register user routes |
| `internal/handler/users.go` | NEW file — user CRUD handler |
| `web/templates/users/index.html` | NEW — user list page |
| `web/templates/users/form.html` | NEW — user create/edit form |
| `web/templates/layout.html` | Conditional sidebar menu based on role |

### Sidebar (conditional)

```
superadmin:   Dashboard | Trip | Booking | Laporan | Paket | Meeting Point | Users | Pengaturan
admin:        Dashboard | Trip | Booking | Laporan | Paket | Meeting Point | Pengaturan
user:         Laporan | Pengaturan
```

---

## Phase C: Unit Test Files

| # | File Path | Tests | Coverage |
|---|-----------|-------|----------|
| 1 | `internal/auth/auth_test.go` | 9 | Password + Session |
| 2 | `internal/handler/helpers_test.go` | — | Shared DB + mock setup |
| 3 | `internal/handler/auth_test.go` | 6 | Login/Logout |
| 4 | `internal/handler/trips_test.go` | 13 | Trip CRUD + packages + MPs |
| 5 | `internal/handler/schedules_test.go` | 8 | Schedule CRUD |
| 6 | `internal/handler/bookings_test.go` | 10 | Booking + status |
| 7 | `internal/handler/meeting_points_test.go` | 10 | MP CRUD |
| 8 | `internal/handler/packages_test.go` | 9 | Package CRUD |
| 9 | `internal/handler/reports_test.go` | 2 | Revenue report |
| 10 | `internal/handler/profile_test.go` | 5 | Profile settings |
| 11 | `internal/handler/dashboard_test.go` | 3 | Dashboard |
| 12 | `internal/handler/users_test.go` | 9 | User CRUD + RBAC |
| 13 | `internal/handler/smoke_test.go` | 33 | Route audit |
| **Total** | **13 files** | **~117** | |

### Test Infrastructure (`helpers_test.go`)

```go
type mockRenderer struct {
    lastTemplate string
    lastData     map[string]interface{}
    loginData    interface{}
}
func (r *mockRenderer) RenderTemplate(w, name, data) { r.lastTemplate = name; r.lastData = data }
func (r *mockRenderer) RenderLogin(w, name, data) { r.loginData = data }

func setupTestDB(t *testing.T) *sql.DB { ... }  // :memory: + migrate + seed
func setupHandler(t, db) *Handler { ... }        // Handler with mock renderer
func authRequest(t, method, path, body, cookie) *http.Request { ... }
func assertStatus(t, w, code) { ... }
func assertContains(t, body, sub) { ... }
func assertRedirect(t, w, target) { ... }
```

---

## Phase D: Execute & Iterate

```bash
# Step 1 — Run smoke tests (audit)
go test ./internal/handler/ -v -run Smoke -count=1

# Step 2 — Run all unit tests
go test ./internal/... -v -count=1

# Step 3 — Fix failures → repeat until green

# Step 4 — Coverage
go test ./internal/... -cover

# Step 5 — Commit & push
git add -A
git commit -m "feat: user management RBAC + comprehensive unit tests"
git push origin master
```

---

## Target

- **0 failing tests**
- **Coverage > 70%**
- **All ~33 API endpoints smoke-tested**
- **RBAC for 3 roles: superadmin, admin, user**
