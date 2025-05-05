package startup

const config = `
################################## NETWORK #####################################

# Listen on IPv6 loopback interface
bind ::1

# Only accept localhost connections if the default user has no password
protected-mode yes

# Accept TLS connections on the specified/default port
port {{ .Port }}
tls-port {{ .TLSPort }}

# The X.509 certificate and private key used to authenticate the server to connected clients
tls-cert-file {{ .TLSCertFile }}
tls-key-file {{ .TLSKeyFile }}

# The CA certificate or bundle used to authenticate TLS clients/peers
# This or tls-ca-cert-dir /etc/ssl/certs must be specified
tls-ca-cert-file {{ .TLSCACertFile }}

# Enable TLS for replication and cluster communication
tls-replication yes
tls-cluster yes

################################ SNAPSHOTTING  ################################

# Save a snapshot every 300 seconds (5 minutes) if 10 changes have been performed
# Snapshots are saved to the ./dump.rdb file
save 300 10

################################# REPLICATION #################################

{{ if .Replica }}
replicaof {{ .MasterIP }} {{ .MasterPort}}
{{ end }}

# The password that should be used to authenticate with the master node
# This value may be derived from the secret store (Infisical)
masterauth {{ .MasterAuth }}

# A special user that is capable of running PSYNC/other replication commands
# This is needed when using ACLs as the default user is not capable of running above commands
# Common permissions for the masteruser can be: ACL SETUSER replica-user ON >strong-password +psync +replconf +ping
# The value should be derived from the secret store (Infisical)
masteruser {{ .MasterUser }}

################################## SECURITY ###################################

# Practically, a Redis password is a shared secret between the client and server
# It should also be long, unguessable, and not memorized by any human to prevent brute force attacks

# Disable the default user
{{ if .EnableDefaultUser }}
user default on
{{ else }}
user default off
{{ end }}

# The user replicauser is allowed to execute commands required for replication purposes
{{ range $index, $user := .Users }}
user {{ $user.Username }} on >{{ $user.Password }} {{ range $cmd := $user.EnabledCommands }}{{ $cmd }} {{ end }}
{{ end }}

# The user userservice is allowed to reference keys starting with "userservice:"
# Its password is randomly generated and supplied at Redis Cluster creation
user userservice on ><password> -@all %R~userservice %W~userservice -@admin

################################ REDIS CLUSTER  ###############################

# Enable Redis Cluster as normal Redis nodes can't be part of a cluster
cluster-enabled yes

############################## MEMORY MANAGEMENT ################################

# When the maximum memory is reached, evict keys with an expire set using approximate LRU
maxmemory-policy volatile-lru


############################## APPEND ONLY MODE ###############################

# Use the Append Only File for better durability in the event of a disaster
appendonly yes

################################ SHUTDOWN #####################################

# Force DB saving operations even if no save points are configured
shutdown-on-sigint save
shutdown-on-sigterm save
`
