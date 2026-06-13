package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

// cleanup biar gak memory leak
func cleanupClients() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 5*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	c, exists := clients[ip]

	if !exists {
		limiter := rate.NewLimiter(1, 5) // 1 req/sec burst 5
		clients[ip] = &client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	c.lastSeen = time.Now()
	return c.limiter
}

func RateLimit(next http.Handler) http.Handler {
	go cleanupClients()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr

		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w,
				"Terlalu banyak request!",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
