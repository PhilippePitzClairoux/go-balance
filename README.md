# Go Load Balancer
This program is a simple http load balancer.
It supports subdomains and paths.

## Common used terms in the project
### Trigger 
A trigger is either a path or a subdomain.
Paths have precedence over subdomains. \
Here's an example :

```yaml
"/gsearch"          : "https://google.com/search"
"google.localhost" : "https://google.com/search"
```

when a HTTP call is made, if the path `/gsearch` is found,
the request will be redirected to https://google.com/search

when a HTTP call is made, if the path does not match `/gsearch`
and the Host header contains `google.localhost`, the request will
be forwarded to `https://google.com/search`

### Target
Target is a terme that refers to the target server that will be called
once the trigger matched.

### Pool
A pool is a structure that defines many targets. Can be usefull
for load-balancing purposes or fail-overs. When a trigger has a
pool as a target, depending on the `distribution_type`, every call
will return a different target.

## Pool Distribution Types
There are a few supported distribution types :

### FAIL_OVER (safe)
When a request is made, the load-balancer pings the server using `test_connection` configuration.
If the server is up, we call it. If it isn't up, we recursively test the next server untill we find one.
The max amount of retries is currently hardcoded at 3, but eventually this value will be
configurable.

### ROUND_ROBIN (unsafe)
When a request is made, the load-balancer returns a server. The counter will be incremented by
one at every request (regardless of the outcome). If for some reason the server is down, the request will
just fail. It's the clients duty to retry the connection.
Ignores `test_connection`

### ONE_CONNECTION_SERVER_POOL (unsafe)
You can use a ONE_CONNECTION_SERVER_POOL, but it's not recommended.
It's mostly used internally for targets that have a single host. If you still want to use it,
by all means, go for it.
Ignores `test_connection`