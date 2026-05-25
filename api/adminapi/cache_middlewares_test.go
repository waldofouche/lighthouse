package adminapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-oidfed/lib/cache"
	"github.com/gofiber/fiber/v2"

	"github.com/go-oidfed/lighthouse/internal"
)

func setCacheEntry(t *testing.T, key string, value []byte) {
	t.Helper()
	_ = cache.Delete(key)
	if err := cache.Set(key, value, time.Minute); err != nil {
		t.Fatalf("failed to seed cache for %q: %v", key, err)
	}
	t.Cleanup(func() {
		_ = cache.Delete(key)
	})
}

func requireCacheEntry(t *testing.T, key string, wantSet bool, wantValue []byte) {
	t.Helper()
	var got []byte
	set, err := cache.Get(key, &got)
	if err != nil {
		t.Fatalf("failed to read cache key %q: %v", key, err)
	}
	if set != wantSet {
		t.Fatalf("expected cache present=%v for %q, got %v", wantSet, key, set)
	}
	if wantSet && !bytes.Equal(got, wantValue) {
		t.Fatalf("expected cached value %q for %q, got %q", string(wantValue), key, string(got))
	}
	if !wantSet && len(got) != 0 {
		t.Fatalf("expected empty cached value for cleared key %q, got %q", key, string(got))
	}
}

// TestEntityConfigurationCacheInvalidationMiddleware must NOT use t.Parallel().
// It operates on the global process-wide cache (cache.Set/Get/Delete), which is
// shared mutable state. Parallelizing these subtests would cause race conditions
// on the entity configuration cache key.
func TestEntityConfigurationCacheInvalidationMiddleware(t *testing.T) {
	cacheValue := []byte("entity-config-jwt")

	t.Run("SuccessDeletesCache", func(t *testing.T) {
		setEntityConfigurationCache(t, cacheValue)
		app := fiber.New()
		app.Post("/entity-config", entityConfigurationCacheInvalidationMiddleware, func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodPost, "/entity-config", http.NoBody)
		resp, bodyBytes := doRequest(t, app, req)
		requireStatus(t, resp, bodyBytes, http.StatusNoContent)
		requireEntityConfigurationCache(t, false, nil)
	})

	t.Run("FailureKeepsCache", func(t *testing.T) {
		setEntityConfigurationCache(t, cacheValue)
		app := fiber.New()
		app.Post("/entity-config", entityConfigurationCacheInvalidationMiddleware, func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusBadRequest)
		})

		req := httptest.NewRequest(http.MethodPost, "/entity-config", http.NoBody)
		resp, bodyBytes := doRequest(t, app, req)
		requireStatus(t, resp, bodyBytes, http.StatusBadRequest)
		requireEntityConfigurationCache(t, true, cacheValue)
	})

	// TODO: Currently there are no redirects. If we add any in the future, we should verify that they do not clear the cache.
	// t.Run("RedirectDoesNotClearCache", func(t *testing.T) {
	// 	setEntityConfigurationCache(t, cacheValue)

	// 	app := fiber.New()
	// 	app.Post("/entity-config", entityConfigurationCacheInvalidationMiddleware, func(c *fiber.Ctx) error {
	// 		return c.Redirect("/other", fiber.StatusMovedPermanently)
	// 	})

	// 	req := httptest.NewRequest(http.MethodPost, "/entity-config", http.NoBody)
	// 	resp, bodyBytes := doRequest(t, app, req)

	// 	requireStatus(t, resp, bodyBytes, fiber.StatusMovedPermanently)
	// 	requireEntityConfigurationCache(t, true, cacheValue)
	// })
}

// TestSubordinateStatementsCacheInvalidationMiddleware must NOT use t.Parallel().
// It operates on the global process-wide cache (cache.Set/Get/Delete), which is
// shared mutable state. Parallelizing these subtests would cause race conditions
// on the subordinate statement cache keys.
func TestSubordinateStatementsCacheInvalidationMiddleware(t *testing.T) {
	key123 := cache.Key(internal.CacheKeySubordinateStatement, "123")
	key456 := cache.Key(internal.CacheKeySubordinateStatement, "456")
	value123 := []byte("statement-123")
	value456 := []byte("statement-456")

	t.Run("SpecificSubordinateDeletesOnlyTarget", func(t *testing.T) {
		setCacheEntry(t, key123, value123)
		setCacheEntry(t, key456, value456)
		app := fiber.New()
		app.Delete("/subordinates/:subordinateID", subordinateStatementsCacheInvalidationMiddleware, func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodDelete, "/subordinates/123", http.NoBody)
		resp, bodyBytes := doRequest(t, app, req)
		requireStatus(t, resp, bodyBytes, http.StatusNoContent)
		requireCacheEntry(t, key123, false, nil)
		requireCacheEntry(t, key456, true, value456)
	})

	t.Run("CollectionSuccessClearsAll", func(t *testing.T) {
		setCacheEntry(t, key123, value123)
		setCacheEntry(t, key456, value456)
		app := fiber.New()
		app.Post("/subordinates", subordinateStatementsCacheInvalidationMiddleware, func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusCreated)
		})

		req := httptest.NewRequest(http.MethodPost, "/subordinates", http.NoBody)
		resp, bodyBytes := doRequest(t, app, req)
		requireStatus(t, resp, bodyBytes, http.StatusCreated)
		requireCacheEntry(t, key123, false, nil)
		requireCacheEntry(t, key456, false, nil)
	})

	t.Run("FailureKeepsAll", func(t *testing.T) {
		setCacheEntry(t, key123, value123)
		setCacheEntry(t, key456, value456)
		app := fiber.New()
		app.Delete("/subordinates/:subordinateID", subordinateStatementsCacheInvalidationMiddleware, func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusInternalServerError)
		})

		req := httptest.NewRequest(http.MethodDelete, "/subordinates/123", http.NoBody)
		resp, bodyBytes := doRequest(t, app, req)
		requireStatus(t, resp, bodyBytes, http.StatusInternalServerError)
		requireCacheEntry(t, key123, true, value123)
		requireCacheEntry(t, key456, true, value456)
	})
}
