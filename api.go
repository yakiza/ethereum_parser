package ethereum_parser

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpHandlers struct {
	service service
}

func (h *HttpHandlers) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block, err := h.service.GetCurrentBlock(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte((err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(block); err != nil {
		http.Error(w, fmt.Sprintf("error building the responsse, %v", err), http.StatusInternalServerError)
	}

}

func (h *HttpHandlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	// Can add a check that responds accordingly if address missing
	hasSubscribed, err := h.service.Subscribe(r.Context(), r.URL.Query().Get(addressParam))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fmt.Sprintf("subscribed %v", hasSubscribed)); err != nil {
		http.Error(w, fmt.Sprintf("error building the responsse, %v", err), http.StatusInternalServerError)
	}
}

func (h *HttpHandlers) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.service.GetTransactions(r.Context(), r.URL.Query().Get(addressParam))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, fmt.Sprintf("error building the responsse, %v", err), http.StatusInternalServerError)
	}

}

func NewHandlers(service service) HttpHandlers {
	return HttpHandlers{
		service: service,
	}
}

const addressParam = "address"
