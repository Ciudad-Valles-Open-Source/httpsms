package repositories

import (
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	gocache "github.com/patrickmn/go-cache"
)

// UserStatusCacheTTL is how long the status of a user is cached in memory. A
// status change made through the API updates the cache of the instance which
// handled it immediately; other API instances (when running several replicas)
// can serve the previous status for at most this long.
const UserStatusCacheTTL = 5 * time.Second

// userStatusDeleted marks a deleted user in the cache
const userStatusDeleted = entities.UserStatus("")

// UserStatusCache caches the entities.UserStatus of users by ID. It is safe for concurrent use.
type UserStatusCache struct {
	cache *gocache.Cache
	ttl   time.Duration
}

// NewUserStatusCache creates a new UserStatusCache
func NewUserStatusCache(ttl time.Duration) *UserStatusCache {
	return &UserStatusCache{
		cache: gocache.New(ttl, 2*ttl),
		ttl:   ttl,
	}
}

// get returns the cached status, found is false when nothing is cached
func (c *UserStatusCache) get(userID entities.UserID) (status entities.UserStatus, found bool) {
	if c == nil {
		return "", false
	}
	value, found := c.cache.Get(userID.String())
	if !found {
		return "", false
	}
	status, ok := value.(entities.UserStatus)
	return status, ok
}

// add caches a status loaded from the database. It never replaces a value set
// by set() in the meantime, so a slow read can't overwrite a newer status.
func (c *UserStatusCache) add(userID entities.UserID, status entities.UserStatus) {
	if c == nil {
		return
	}
	_ = c.cache.Add(userID.String(), status, c.ttl)
}

// set stores the status written to the database, replacing any cached value.
func (c *UserStatusCache) set(userID entities.UserID, status entities.UserStatus) {
	if c == nil {
		return
	}
	c.cache.Set(userID.String(), status, c.ttl)
}
