package components

import (
	"net/http"

	uiErrors "github.com/dskart/waterfall-engine/ui/pkg/errors"
	"github.com/dskart/waterfall-engine/ui/pkg/middleware"
	"github.com/dskart/waterfall-engine/ui/pkg/router"
)

func init() {
	router.ComponentHandleFunc("/commitments_table", http.MethodGet, GetCommitmentsTable)
}

func GetCommitmentsTable(w http.ResponseWriter, r *http.Request) {
	sess := middleware.CtxSession(r.Context())
	commitments, sanitizedErr := sess.GetCommitments()
	if sanitizedErr != nil {
		http.Error(w, sanitizedErr.Error(), uiErrors.ErrorHTTPStatus(sanitizedErr))
		return
	}

	if err := CommitmentsTable(commitments).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
