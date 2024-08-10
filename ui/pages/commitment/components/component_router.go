package components

import (
	"net/http"
	"strconv"

	uiErrors "github.com/dskart/waterfall-engine/ui/pkg/errors"
	"github.com/dskart/waterfall-engine/ui/pkg/middleware"
	"github.com/dskart/waterfall-engine/ui/pkg/router"
)

func init() {
	router.ComponentHandleFunc("/commitment/distributions_table", http.MethodGet, GetDistributionsTable)
	router.ComponentHandleFunc("/commitment/breadcrumbs", http.MethodGet, GetBreadcrumbs)
	router.ComponentHandleFunc("/commitment/stats", http.MethodGet, GetStats)
}

func GetDistributionsTable(w http.ResponseWriter, r *http.Request) {
	rawCommitmentId := r.URL.Query().Get("commitmentId")
	if rawCommitmentId == "" {
		http.Error(w, "missing commitmentId", http.StatusBadRequest)
		return
	}

	commitmentId, err := strconv.Atoi(rawCommitmentId)
	if err != nil {
		http.Error(w, "commitmentId is not a int", http.StatusBadRequest)
		return
	}

	sess := middleware.CtxSession(r.Context())
	distributions, sanitizedErr := sess.GetDistributionsByCommitmentId(commitmentId)
	if sanitizedErr != nil {
		http.Error(w, sanitizedErr.Error(), uiErrors.ErrorHTTPStatus(sanitizedErr))
		return
	}

	if err := DistributionsTable(distributions).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetBreadcrumbs(w http.ResponseWriter, r *http.Request) {
	rawCommitmentId := r.URL.Query().Get("commitmentId")
	if rawCommitmentId == "" {
		http.Error(w, "missing commitmentId", http.StatusBadRequest)
		return
	}

	commitmentId, err := strconv.Atoi(rawCommitmentId)
	if err != nil {
		http.Error(w, "commitmentId is not a int", http.StatusBadRequest)
		return
	}

	sess := middleware.CtxSession(r.Context())
	commitment, sanitizedErr := sess.GetCommitmentById(commitmentId)
	if sanitizedErr != nil {
		http.Error(w, sanitizedErr.Error(), uiErrors.ErrorHTTPStatus(sanitizedErr))
		return
	}

	if err := BreadCrumbs(commitment.EntityName).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	rawCommitmentId := r.URL.Query().Get("commitmentId")
	if rawCommitmentId == "" {
		http.Error(w, "missing commitmentId", http.StatusBadRequest)
		return
	}

	commitmentId, err := strconv.Atoi(rawCommitmentId)
	if err != nil {
		http.Error(w, "commitmentId is not a int", http.StatusBadRequest)
		return
	}

	sess := middleware.CtxSession(r.Context())
	stats, sanitizedErr := sess.GetStatsByCommitmentId(commitmentId)
	if sanitizedErr != nil {
		http.Error(w, sanitizedErr.Error(), uiErrors.ErrorHTTPStatus(sanitizedErr))
		return
	}

	if err := Stats(stats).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
