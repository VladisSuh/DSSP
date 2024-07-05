package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	pools := NewPoolManager()
	cr := &ChainOfResponsibility{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("login.html")
		if err != nil {
			http.Error(w, "Could not read HTML file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "text/html")
		if _, err := io.Copy(w, file); err != nil {
			http.Error(w, "Failed to send HTML file", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")
		if authenticate(username, password) {
			http.SetCookie(w, &http.Cookie{
				Name:    "authenticated",
				Value:   "true",
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			})
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": true}`)
		} else {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false}`)
		}
	})

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cookie, err := r.Cookie("authenticated"); err != nil || cookie.Value != "true" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	http.Handle("/commands", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("commands.html")
		if err != nil {
			http.Error(w, "Could not read HTML file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "text/html")
		if _, err := io.Copy(w, file); err != nil {
			http.Error(w, "Failed to send HTML file", http.StatusInternalServerError)
		}
	})))

	http.HandleFunc("/run-command", func(w http.ResponseWriter, r *http.Request) {
		command := r.URL.Query().Get("command")
		if command == "" {
			http.Error(w, `{"error": "Missing command parameter"}`, http.StatusBadRequest)
			return
		}
		if err := runCommand(pools, command, cr); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error executing command: %s"}`, err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Command executed successfully"}`)
	})

	http.HandleFunc("/get-info", func(w http.ResponseWriter, r *http.Request) {
		type Info struct {
			Pools map[string][]string `json:"pools"`
		}
		info := Info{Pools: make(map[string][]string)}
		for poolName, pool := range pools.Pools {
			for schemaName := range pool.Schemas {
				info.Pools[poolName] = append(info.Pools[poolName], schemaName)
			}
		}
		data, err := json.Marshal(info)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Error getting info: %s"}`, err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			file, err := os.Open("registration.html")
			if err != nil {
				http.Error(w, "Ошибка при чтении HTML файла", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if _, err := io.Copy(w, file); err != nil {
				http.Error(w, "Ошибка при отправке HTML файла", http.StatusInternalServerError)
				return
			}
			return
		}

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")
			users[username] = password
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
