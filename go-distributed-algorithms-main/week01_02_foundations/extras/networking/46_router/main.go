package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	router.HandleFunc("POST /hello/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("user_id")

		w.Write([]byte(fmt.Sprintf("Hello, %s!", userID)))
	})

	fmt.Println("Server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
