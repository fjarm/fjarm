package remote

import (
	"time"

	"github.com/redis/rueidis"
)

// NewRedisClient creates a new Redis client using rueidis.
func NewRedisClient(
	authCredentials rueidis.AuthCredentials,
	addrs []string,
) (rueidis.Client, error) {
	client, err := rueidis.NewClient(
		rueidis.ClientOption{
			// TODO(2025-03-09): Supply AuthCredentialsFn to provide username and password for ACL support.
			AuthCredentialsFn: func(authCtx rueidis.AuthCredentialsContext) (rueidis.AuthCredentials, error) {
				return authCredentials, nil
			},
			// When running Sentinel mode, all node addresses need to be individually supplied. In Cluster mode, only
			// the one address needs to be supplied like - "redis-cluster.railway.internal:6379".
			InitAddress: addrs,
			ClientTrackingOptions: []string{
				// This is the default value. Keys mentioned in read operations aren't cached. Caching must be
				// proactively turned on immediately before the actual command to enable client-side caching.
				"OPTIN",
			},
			// TODO(2025-03-09): Allow specifying CacheSizeEachConn when client-side caching is enabled.
			BlockingPoolCleanup: 30 * time.Second,
			MaxFlushDelay:       0,
			// TODO(2025-03-09): Set ShardsRefreshInterval to non-zero value after enabling Redis Cluster.
			//ClusterOption:         rueidis.ClusterOption{
			//	ShardsRefreshInterval: 0,
			//},
			DisableCache:          true, // Disable client-side caching.
			DisableAutoPipelining: true, // Manual pipelining can be enabled using client.DoMulti().
			// Toggled to true for read-only clients. But this should be accomplished using ACLs.
			ReplicaOnly: false,
		},
	)
	return client, err
}
