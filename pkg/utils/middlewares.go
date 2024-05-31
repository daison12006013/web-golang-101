package utils

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"web-golang-101/pkg/db"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", GetEnvWithDefault("CORS_ORIGIN", "*"))
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		EnableCors(&w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get memory stats before processing the request
		var m1 runtime.MemStats
		runtime.ReadMemStats(&m1)

		next.ServeHTTP(w, r)

		// Get memory stats after processing the request
		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)

		duration := time.Since(start)

		// Calculate the memory used to process the request (in bytes)
		memUsed := float64(m2.Alloc - m1.Alloc)

		// Convert memory usage to MB
		memUsedMB := fmt.Sprintf("%.4f", memUsed/1024/1024) + "MB"

		Logger().Info().
			Str("memory_used", memUsedMB).
			Float64("duration", duration.Seconds()).
			Str("url", r.URL.String()).
			Str("method", r.Method).
			Msg("Incoming request")
	})
}

func RejectNonSpecificDomain(domain string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host != domain && !strings.HasSuffix(r.Host, "."+domain) {
				http.Error(w, "Unacceptable domain", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CaptureErrors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var err error

				switch t := r.(type) {
				case ErrorStatus:
					http.Error(w, t.Message, t.Status)
					err = t.Error
				case error:
					err = t
				default:
					err = fmt.Errorf("unknown error: %v", r)
				}

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				sentry.CaptureException(err)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

// KeyByRealIP copied from httprate.KeyByRealIP
// this is to fully support other platforms, especially CF-Connecting-IP
func KeyByRealIP(r *http.Request) (string, error) {
	var ip string

	if fcip := r.Header.Get(GetEnvWithDefault("IP_CLIENT_HEADER_KEY", "CF-Connecting-IP")); fcip != "" {
		ip = fcip
	} else if tcip := r.Header.Get("True-Client-IP"); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else {
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
	}

	return ip, nil
}

func RequireApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var apikey string
		combiKeys := []string{"apiKey", "apikey", "api-key", "api_key"}

		for _, combiKey := range combiKeys {
			apikey = r.URL.Query().Get(combiKey)
			if apikey != "" {
				break
			}
		}

		if apikey == "" {
			apikey = r.Header.Get("X-API-KEY")
		}

		if apikey == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}

		conn, err := db.NewConnection()
		if err != nil {
			http.Error(w, "Error connecting to database", http.StatusInternalServerError)
			return
		}
		defer conn.DB.Close()

		q := conn.NewQuery()

		exists, err := q.ApiKeyExists(context.Background(), apikey)
		if err != nil {
			http.Error(w, "Error checking API key", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetDefaultMiddlewares(r *chi.Mux) {
	rejectNonSpecificDomain := os.Getenv("REJECT_NON_SPECIFIC_DOMAIN") // Reject requests from non-specific domains
	if rejectNonSpecificDomain != "" {
		r.Use(RejectNonSpecificDomain(rejectNonSpecificDomain))
	}

	rateLimitStr := GetEnvWithDefault("RATE_LIMIT", "5") // Rate limit in requests per second
	if rateLimitStr != "" {
		rateLimit, err := strconv.Atoi(rateLimitStr)
		if err != nil {
			log.Fatal("Invalid rate limit value:", err)
		}
		r.Use(RateLimit(rateLimit))
	}

	if GetEnvWithDefault("ENABLE_CAPTURE_ERRORS", "true") == "true" {
		r.Use(CaptureErrors)
	}

	if GetEnvWithDefault("ENABLE_CORS_MIDDLEWARE", "true") == "true" {
		r.Use(CorsMiddleware)
	}

	if GetEnvWithDefault("ENABLE_LOG_REQUEST", "true") == "true" {
		r.Use(LogRequest)
	}
}

func RateLimit(rateLimit int) func(http.Handler) http.Handler {
	return httprate.Limit(rateLimit, 1*time.Second,
		httprate.WithKeyFuncs(KeyByRealIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			resp := NewResponse(w)
			resp.WriteErrorResponse("Too many request", http.StatusTooManyRequests)
		}),
	)
}

func WafMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := NewResponse(w)

		sqlInjectionPattern := regexp.MustCompile(`(?i)(?:'|\b)(?:--|select\b|update\b|drop\b|insert\b|delete\b|or\b|and\b|exec\b|execute\b|union\b|truncate\b|declare\b)`)

		if sqlInjectionPattern.MatchString(r.URL.RawQuery) {
			resp.WriteErrorResponse("Forbidden", http.StatusForbidden)
			return
		}

		xssPattern := regexp.MustCompile(`(?i)<\s*script[^>]*>|<\s*/\s*script\s*>|<\s*img[^>]*>|<\s*a[^>]*>|<\s*body[^>]*>|<\s*iframe[^>]*>|<\s*input[^>]*>|<\s*form[^>]*>|<\s*style[^>]*>|<\s*svg[^>]*>|<\s*link[^>]*>|<\s*object[^>]*>|<\s*embed[^>]*>|<\s*frame[^>]*>|<\s*frameset[^>]*>|<\s*applet[^>]*>|<\s*meta[^>]*>|<\s*layer[^>]*>|<\s*ilayer[^>]*>|<\s*bgsound[^>]*>|<\s*base[^>]*>|<\s*xml[^>]*>|<\s*import[^>]*>|<\s*isindex[^>]*>|<\s*textarea[^>]*>|<\s*div[^>]*>`)

		if xssPattern.MatchString(r.URL.RawQuery) {
			resp.WriteErrorResponse("Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
