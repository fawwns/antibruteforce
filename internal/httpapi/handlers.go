package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/fawwns/antibruteforce/internal/logger"
	"github.com/fawwns/antibruteforce/internal/service"
)

type Handler struct {
	svc *service.AntiBruteForceService
}

type CheckResponse struct {
	OK bool `json:"ok"`
}

func NewHandler(svc *service.AntiBruteForceService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	password := r.URL.Query().Get("password")
	ip := r.URL.Query().Get("ip")

	logger.Info.Printf("Check request: login=%s ip=%s", login, ip)

	if ip == "" {
		http.Error(w, "ip required", http.StatusBadRequest)
		return
	}

	ok := h.svc.Check(login, password, ip)

	response := CheckResponse{OK: ok}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ResetBucket(w http.ResponseWriter, r *http.Request) {
	h.svc.ResetBuckets()
	w.Write([]byte(`{"result": "ok"}`))
}

func (h *Handler) AddWhitelist(w http.ResponseWriter, r *http.Request) {
	cidr := r.URL.Query().Get("cidr")
	if cidr == "" {
		http.Error(w, "cidr required", http.StatusBadRequest)
		return
	}

	if err := h.svc.AddToWhitelist(cidr); err != nil {
		logger.Warn.Printf("Whitelist add error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("Add to whitelist: %s", cidr)
	w.Write([]byte(`{"result": "ok"}`))
}

func (h *Handler) AddBlacklist(w http.ResponseWriter, r *http.Request) {
	cidr := r.URL.Query().Get("cidr")
	if cidr == "" {
		http.Error(w, "cidr required", http.StatusBadRequest)
		return
	}

	if err := h.svc.AddToBlacklist(cidr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("Add to Blacklist: %s", cidr)
	w.Write([]byte(`{"result": "ok"}`))
}

func (h *Handler) RemoveFromWhitelist(w http.ResponseWriter, r *http.Request) {
	cidr := r.URL.Query().Get("cidr")
	if cidr == "" {
		http.Error(w, "cidr required", http.StatusBadRequest)
		return
	}

	if err := h.svc.RemoveFromWhitelist(cidr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(`{"result":"ok"}`))
}
func (h *Handler) RemoveFromBlacklist(w http.ResponseWriter, r *http.Request) {
	cidr := r.URL.Query().Get("cidr")
	if cidr == "" {
		http.Error(w, "cidr required", http.StatusBadRequest)
		return
	}

	if err := h.svc.RemoveFromBlacklist(cidr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(`{"result":"ok"}`))
}
