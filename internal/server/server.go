package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-scout/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux}
func New(db *store.DB)*Server{s:=&Server{db:db,mux:http.NewServeMux()}
s.mux.HandleFunc("GET /api/scans",s.list);s.mux.HandleFunc("POST /api/scans",s.create);s.mux.HandleFunc("GET /api/scans/{id}",s.get);s.mux.HandleFunc("PUT /api/scans/{id}",s.update);s.mux.HandleFunc("DELETE /api/scans/{id}",s.del)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"scans":oe(s.db.List())})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var sc store.Scan;json.NewDecoder(r.Body).Decode(&sc);if sc.Project==""{we(w,400,"project required");return};s.db.Create(&sc);wj(w,201,s.db.Get(sc.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){sc:=s.db.Get(r.PathValue("id"));if sc==nil{we(w,404,"not found");return};wj(w,200,sc)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");ex:=s.db.Get(id);if ex==nil{we(w,404,"not found");return};var sc store.Scan;json.NewDecoder(r.Body).Decode(&sc);s.db.Update(id,&sc);wj(w,200,s.db.Get(id))}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"scout","scans":st.Scans})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
