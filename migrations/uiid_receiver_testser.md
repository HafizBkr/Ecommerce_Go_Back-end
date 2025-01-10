 // r.Get("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
   //     id := chi.URLParam(r, "id")
    //     fmt.Printf("URL complète: %s\n", r.URL.Path)
    //     fmt.Printf("ID extrait: %s\n", id)
        
    //     if id == "" {
    //         w.WriteHeader(http.StatusBadRequest)
    //         w.Write([]byte("ID non trouvé dans l'URL"))
    //         return
    //     }
        
    //     w.Write([]byte(fmt.Sprintf("ID trouvé: %s", id)))
    // })
	// r.Route("/categories", func(r chi.Router) {
	// 	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Println("Route GET /categories/{id} appelée")
	// 		categoryHandler.HandleGetCategoryByID(w, r)
	// 	})
	// })