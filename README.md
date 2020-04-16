# claack-ws
the go port of my uhh..unreleased nodejs typing server

# features

- highly performant websocket support
- stateless user sessions with [JWT](https://github.com/dgrijalva/jwt-go)
- postgresql backend
- service oriented
- recaptcha support
- "edge server" load balancing for websockets with [redis](https://redis.io/)
- packet rate limiting