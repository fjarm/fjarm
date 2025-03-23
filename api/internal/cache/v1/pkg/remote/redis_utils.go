package remote

// Redis processes Lua scripts in a single-threaded manner, blocking all other operations until the script is finished
// executing. This means that we can use a Lua script atomically and no other client can modify the lock state.
const safeDeleteScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`
