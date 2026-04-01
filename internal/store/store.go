package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type Site struct{ID int64 `json:"id"`;Name string `json:"name"`;BaseURL string `json:"base_url"`;Schedule string `json:"schedule"`;LastCrawlAt *time.Time `json:"last_crawl_at"`;CreatedAt time.Time `json:"created_at"`}
type CrawlRun struct{ID int64 `json:"id"`;SiteID int64 `json:"site_id"`;StartedAt time.Time `json:"started_at"`;FinishedAt *time.Time `json:"finished_at"`;PagesChecked int `json:"pages_checked"`;BrokenLinks int `json:"broken_links"`;SSLIssues int `json:"ssl_issues"`;Status string `json:"status"`}
type Issue struct{ID int64 `json:"id"`;RunID int64 `json:"run_id"`;URL string `json:"url"`;IssueType string `json:"issue_type"`;StatusCode int `json:"status_code"`;Detail string `json:"detail"`;FoundAt time.Time `json:"found_at"`}
func Open(dataDir string)(*DB,error){
    if err:=os.MkdirAll(dataDir,0755);err!=nil{return nil,fmt.Errorf("mkdir: %w",err)}
    dsn:=filepath.Join(dataDir,"scout.db")+"?_journal_mode=WAL&_busy_timeout=5000"
    db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)}
    db.SetMaxOpenConns(1);if err:=migrate(db);err!=nil{return nil,fmt.Errorf("migrate: %w",err)}
    return &DB{db},nil}
func migrate(db *sql.DB)error{
    _,err:=db.Exec(`CREATE TABLE IF NOT EXISTS sites(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL,base_url TEXT NOT NULL,schedule TEXT DEFAULT 'daily',last_crawl_at DATETIME,created_at DATETIME DEFAULT CURRENT_TIMESTAMP);
    CREATE TABLE IF NOT EXISTS crawl_runs(id INTEGER PRIMARY KEY AUTOINCREMENT,site_id INTEGER NOT NULL,started_at DATETIME DEFAULT CURRENT_TIMESTAMP,finished_at DATETIME,pages_checked INTEGER DEFAULT 0,broken_links INTEGER DEFAULT 0,ssl_issues INTEGER DEFAULT 0,status TEXT DEFAULT 'running');
    CREATE TABLE IF NOT EXISTS issues(id INTEGER PRIMARY KEY AUTOINCREMENT,run_id INTEGER NOT NULL,url TEXT NOT NULL,issue_type TEXT NOT NULL,status_code INTEGER DEFAULT 0,detail TEXT DEFAULT '',found_at DATETIME DEFAULT CURRENT_TIMESTAMP);`);return err}
func(db *DB)ListSites()([]Site,error){rows,err:=db.Query(`SELECT id,name,base_url,schedule,last_crawl_at,created_at FROM sites ORDER BY created_at DESC`);if err!=nil{return nil,err};defer rows.Close();var out[]Site;for rows.Next(){var s Site;rows.Scan(&s.ID,&s.Name,&s.BaseURL,&s.Schedule,&s.LastCrawlAt,&s.CreatedAt);out=append(out,s)};return out,nil}
func(db *DB)CreateSite(s *Site)error{res,err:=db.Exec(`INSERT INTO sites(name,base_url,schedule)VALUES(?,?,?)`,s.Name,s.BaseURL,s.Schedule);if err!=nil{return err};s.ID,_=res.LastInsertId();return nil}
func(db *DB)GetSite(id int64)(*Site,error){s:=&Site{};err:=db.QueryRow(`SELECT id,name,base_url,schedule,last_crawl_at,created_at FROM sites WHERE id=?`,id).Scan(&s.ID,&s.Name,&s.BaseURL,&s.Schedule,&s.LastCrawlAt,&s.CreatedAt);if err==sql.ErrNoRows{return nil,nil};return s,err}
func(db *DB)DeleteSite(id int64)error{_,err:=db.Exec(`DELETE FROM sites WHERE id=?`,id);return err}
func(db *DB)CreateRun(r *CrawlRun)error{res,err:=db.Exec(`INSERT INTO crawl_runs(site_id,status)VALUES(?,'running')`,r.SiteID);if err!=nil{return err};r.ID,_=res.LastInsertId();return nil}
func(db *DB)FinishRun(r *CrawlRun)error{_,err:=db.Exec(`UPDATE crawl_runs SET finished_at=CURRENT_TIMESTAMP,pages_checked=?,broken_links=?,ssl_issues=?,status='done' WHERE id=?`,r.PagesChecked,r.BrokenLinks,r.SSLIssues,r.ID);return err}
func(db *DB)LogIssue(i *Issue)error{_,err:=db.Exec(`INSERT INTO issues(run_id,url,issue_type,status_code,detail)VALUES(?,?,?,?,?)`,i.RunID,i.URL,i.IssueType,i.StatusCode,i.Detail);return err}
func(db *DB)ListRuns(siteID int64)([]CrawlRun,error){rows,err:=db.Query(`SELECT id,site_id,started_at,finished_at,pages_checked,broken_links,ssl_issues,status FROM crawl_runs WHERE site_id=? ORDER BY started_at DESC LIMIT 20`,siteID);if err!=nil{return nil,err};defer rows.Close();var out[]CrawlRun;for rows.Next(){var r CrawlRun;rows.Scan(&r.ID,&r.SiteID,&r.StartedAt,&r.FinishedAt,&r.PagesChecked,&r.BrokenLinks,&r.SSLIssues,&r.Status);out=append(out,r)};return out,nil}
func(db *DB)ListIssues(runID int64)([]Issue,error){rows,err:=db.Query(`SELECT id,run_id,url,issue_type,status_code,detail,found_at FROM issues WHERE run_id=? ORDER BY found_at`,runID);if err!=nil{return nil,err};defer rows.Close();var out[]Issue;for rows.Next(){var i Issue;rows.Scan(&i.ID,&i.RunID,&i.URL,&i.IssueType,&i.StatusCode,&i.Detail,&i.FoundAt);out=append(out,i)};return out,nil}
func(db *DB)UpdateLastCrawl(id int64)error{_,err:=db.Exec(`UPDATE sites SET last_crawl_at=CURRENT_TIMESTAMP WHERE id=?`,id);return err}
func(db *DB)CountSites()(int,error){var n int;db.QueryRow(`SELECT COUNT(*) FROM sites`).Scan(&n);return n,nil}
func(db *DB)CountIssues()(int,error){var n int;db.QueryRow(`SELECT COUNT(*) FROM issues`).Scan(&n);return n,nil}
