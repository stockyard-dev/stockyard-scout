package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Scan struct{ID string `json:"id"`;Project string `json:"project"`;Ecosystem string `json:"ecosystem,omitempty"`;Status string `json:"status"`;TotalDeps int `json:"total_deps"`;Outdated int `json:"outdated"`;Vulnerable int `json:"vulnerable"`;Report string `json:"report,omitempty"`;CreatedAt string `json:"created_at"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"scout.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS scans(id TEXT PRIMARY KEY,project TEXT NOT NULL,ecosystem TEXT DEFAULT '',status TEXT DEFAULT 'pending',total_deps INTEGER DEFAULT 0,outdated INTEGER DEFAULT 0,vulnerable INTEGER DEFAULT 0,report TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(s *Scan)error{s.ID=genID();s.CreatedAt=now();if s.Status==""{s.Status="pending"};_,err:=d.db.Exec(`INSERT INTO scans VALUES(?,?,?,?,?,?,?,?,?)`,s.ID,s.Project,s.Ecosystem,s.Status,s.TotalDeps,s.Outdated,s.Vulnerable,s.Report,s.CreatedAt);return err}
func(d *DB)Get(id string)*Scan{var s Scan;if d.db.QueryRow(`SELECT * FROM scans WHERE id=?`,id).Scan(&s.ID,&s.Project,&s.Ecosystem,&s.Status,&s.TotalDeps,&s.Outdated,&s.Vulnerable,&s.Report,&s.CreatedAt)!=nil{return nil};return &s}
func(d *DB)List()[]Scan{rows,_:=d.db.Query(`SELECT * FROM scans ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Scan;for rows.Next(){var s Scan;rows.Scan(&s.ID,&s.Project,&s.Ecosystem,&s.Status,&s.TotalDeps,&s.Outdated,&s.Vulnerable,&s.Report,&s.CreatedAt);o=append(o,s)};return o}
func(d *DB)Update(id string,s *Scan)error{_,err:=d.db.Exec(`UPDATE scans SET status=?,total_deps=?,outdated=?,vulnerable=?,report=? WHERE id=?`,s.Status,s.TotalDeps,s.Outdated,s.Vulnerable,s.Report,id);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM scans WHERE id=?`,id);return err}
type Stats struct{Scans int `json:"scans"`;Projects int `json:"projects"`}
func(d *DB)Stats()Stats{var s Stats;d.db.QueryRow(`SELECT COUNT(*) FROM scans`).Scan(&s.Scans);d.db.QueryRow(`SELECT COUNT(DISTINCT project) FROM scans`).Scan(&s.Projects);return s}
