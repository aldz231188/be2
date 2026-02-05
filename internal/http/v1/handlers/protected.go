package handlers

// import (
// 	"be2/internal/grpcutil"
// 	"encoding/json"
// 	"net/http"
// )

// type Protected struct{}

// func NewProtected() *Protected { return &Protected{} }

// func (h *Protected) GetProfile(w http.ResponseWriter, r *http.Request) {
// 	uid, _ := r.Context().Value(grpcutil.CtxUserID).(int64)
// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "user_id": uid})
// }
