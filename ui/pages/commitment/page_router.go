package pages

import (
	"net/http"

	"github.com/dskart/waterfall-engine/ui/pkg/router"
	"github.com/gorilla/mux"
)

func init() {
	router.PageHandleFunc("/commitment/{id}", http.MethodGet, getPage)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "no id", http.StatusBadRequest)
		return
	}

	if err := Page(id).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
