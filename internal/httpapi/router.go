package httpapi

import "net/http"

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/check", h.Check)

	mux.HandleFunc("/whitelist/add", h.AddWhitelist)
	mux.HandleFunc("/blacklist/add", h.AddBlacklist)
	mux.HandleFunc("/whitelist/remove", h.RemoveFromWhitelist)
	mux.HandleFunc("/blacklist/remove", h.RemoveFromBlacklist)
	mux.HandleFunc("/bucket/reset", h.ResetBucket)
	return mux
}
