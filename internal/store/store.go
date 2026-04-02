package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Target struct {
	ID string `json:"id"`
	Name string `json:"name"`
	URL string `json:"url"`
	Type string `json:"type"`
	Interval int `json:"interval_seconds"`
	Status string `json:"status"`
	LastCheck string `json:"last_check"`
	LastResult string `json:"last_result"`
	FailCount int `json:"fail_count"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"scout.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS targets(id TEXT PRIMARY KEY,name TEXT NOT NULL,url TEXT DEFAULT '',type TEXT DEFAULT 'http',interval_seconds INTEGER DEFAULT 60,status TEXT DEFAULT 'active',last_check TEXT DEFAULT '',last_result TEXT DEFAULT 'pending',fail_count INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Target)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO targets(id,name,url,type,interval_seconds,status,last_check,last_result,fail_count,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.URL,e.Type,e.Interval,e.Status,e.LastCheck,e.LastResult,e.FailCount,e.CreatedAt);return err}
func(d *DB)Get(id string)*Target{var e Target;if d.db.QueryRow(`SELECT id,name,url,type,interval_seconds,status,last_check,last_result,fail_count,created_at FROM targets WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.URL,&e.Type,&e.Interval,&e.Status,&e.LastCheck,&e.LastResult,&e.FailCount,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Target{rows,_:=d.db.Query(`SELECT id,name,url,type,interval_seconds,status,last_check,last_result,fail_count,created_at FROM targets ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Target;for rows.Next(){var e Target;rows.Scan(&e.ID,&e.Name,&e.URL,&e.Type,&e.Interval,&e.Status,&e.LastCheck,&e.LastResult,&e.FailCount,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Target)error{_,err:=d.db.Exec(`UPDATE targets SET name=?,url=?,type=?,interval_seconds=?,status=?,last_check=?,last_result=?,fail_count=? WHERE id=?`,e.Name,e.URL,e.Type,e.Interval,e.Status,e.LastCheck,e.LastResult,e.FailCount,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM targets WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM targets`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Target{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["type"];ok&&v!=""{where+=" AND type=?";args=append(args,v)}
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,url,type,interval_seconds,status,last_check,last_result,fail_count,created_at FROM targets WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Target;for rows.Next(){var e Target;rows.Scan(&e.ID,&e.Name,&e.URL,&e.Type,&e.Interval,&e.Status,&e.LastCheck,&e.LastResult,&e.FailCount,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM targets GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
