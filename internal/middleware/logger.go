package middleware

import (
	"log"
	"net/http"
	"time"
)

// ANSI Color Codes untuk terminal
const (
	Reset       = "\033[0m"
	ColorInfo   = "\033[36m" // Cyan
	ColorSucc   = "\033[32m" // Green
	ColorWarn   = "\033[33m" // Yellow
	ColorErr    = "\033[31m" // Red
	ColorMethod = "\033[35m" // Magenta
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int // <-- Tambah field ini untuk hitung ukuran response
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Intercept fungsi Write buat ngitung sizenya
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n // Jumlahin bytes yang ditulis
	return n, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default ke 200 kalau ga panggil WriteHeader
		}

		// Request Masuk
		log.Printf(
			"%s[REQ]%s %s%-7s%s %s | ip=%s",
			ColorInfo, Reset,
			ColorMethod, r.Method, Reset,
			r.URL.Path,
			r.RemoteAddr,
		)

		next.ServeHTTP(rw, r)

		// Tentukan warna status code berdasarkan angkanya
		var statusColor string
		switch {
		case rw.statusCode >= 500:
			statusColor = ColorErr
		case rw.statusCode >= 400:
			statusColor = ColorWarn
		default:
			statusColor = ColorSucc
		}

		// Request Selesai (Sekarang nampilin bytesWritten juga)
		log.Printf(
			"%s[RES]%s %s%d%s %s%-7s%s %s | %v | %d bytes",
			ColorInfo, Reset,
			statusColor, rw.statusCode, Reset,
			ColorMethod, r.Method, Reset,
			r.URL.Path,
			time.Since(start),
			rw.bytesWritten, // <-- Muncul di paling kanan log
		)
	})
}
