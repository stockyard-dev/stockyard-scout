package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-scout/internal/server";"github.com/stockyard-dev/stockyard-scout/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8630"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./scout-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("scout: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Scout — Self-hosted dependency scanner\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n\n",port,port,dataDir)
log.Printf("scout: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
