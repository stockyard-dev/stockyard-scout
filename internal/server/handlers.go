package server
import("encoding/json";"net/http";"net/url";"strconv";"strings";"sync";"time";"github.com/stockyard-dev/stockyard-scout/internal/store")
func(s *Server)handleListSites(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListSites();if list==nil{list=[]store.Site{}};writeJSON(w,200,list)}
func(s *Server)handleCreateSite(w http.ResponseWriter,r *http.Request){
    if !s.limits.IsPro(){n,_:=s.db.CountSites();if n>=2{writeError(w,403,"free tier: 2 sites max");return}}
    var site store.Site;json.NewDecoder(r.Body).Decode(&site)
    if site.BaseURL==""{writeError(w,400,"base_url required");return}
    if site.Name==""{site.Name=site.BaseURL};if site.Schedule==""{site.Schedule="daily"}
    if err:=s.db.CreateSite(&site);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,site)}
func(s *Server)handleDeleteSite(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteSite(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListRuns(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListRuns(id);if list==nil{list=[]store.CrawlRun{}};writeJSON(w,200,list)}
func(s *Server)handleListIssues(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListIssues(id);if list==nil{list=[]store.Issue{}};writeJSON(w,200,list)}
func(s *Server)handleCrawl(w http.ResponseWriter,r *http.Request){
    id,_:=strconv.ParseInt(r.PathValue("id"),10,64)
    site,_:=s.db.GetSite(id);if site==nil{writeError(w,404,"site not found");return}
    run:=&store.CrawlRun{SiteID:id};s.db.CreateRun(run)
    go crawlSite(s.db,run,site)
    writeJSON(w,202,map[string]interface{}{"run_id":run.ID,"status":"started","message":"Crawling "+site.BaseURL})}
func crawlSite(db *store.DB,run *store.CrawlRun,site *store.Site){
    client:=&http.Client{Timeout:10*time.Second}
    visited:=map[string]bool{}
    toVisit:=[]string{site.BaseURL}
    base,err:=url.Parse(site.BaseURL);if err!=nil{db.FinishRun(run);return}
    var mu sync.Mutex
    for len(toVisit)>0&&run.PagesChecked<50{
        current:=toVisit[0];toVisit=toVisit[1:]
        mu.Lock();if visited[current]{mu.Unlock();continue};visited[current]=true;mu.Unlock()
        run.PagesChecked++
        resp,err:=client.Get(current)
        if err!=nil{db.LogIssue(&store.Issue{RunID:run.ID,URL:current,IssueType:"unreachable",Detail:err.Error()});run.BrokenLinks++;continue}
        defer resp.Body.Close()
        if resp.StatusCode>=400{db.LogIssue(&store.Issue{RunID:run.ID,URL:current,IssueType:"broken",StatusCode:resp.StatusCode});run.BrokenLinks++;continue}
        // Check SSL
        if strings.HasPrefix(current,"https")&&resp.TLS==nil{db.LogIssue(&store.Issue{RunID:run.ID,URL:current,IssueType:"ssl_issue",Detail:"TLS not established"});run.SSLIssues++}
        // Extract links from body (simple)
        if strings.Contains(resp.Header.Get("Content-Type"),"html"){
            buf:=make([]byte,1<<17)
            n,_:=resp.Body.Read(buf)
            body:=string(buf[:n])
            for _,link:=range extractLinks(body,base){
                mu.Lock();if!visited[link]{toVisit=append(toVisit,link)};mu.Unlock()
            }
        }
    }
    db.UpdateLastCrawl(site.ID);db.FinishRun(run)}
func extractLinks(html string,base *url.URL)[]string{
    var links[]string
    for i:=0;i<len(html)-6;i++{
        if html[i:i+6]=="href=\""||html[i:i+6]=="href='"||html[i:i+5]=="href="{
            start:=i+5;if html[i+4]=='"'||html[i+4]=='\''{start=i+6}
            quote:='"';if start<len(html)&&html[start-1]=='\''{quote='\''}
            end:=strings.IndexByte(html[start:],byte(quote));if end<0{continue}
            href:=html[start:start+end]
            if strings.HasPrefix(href,"#")||strings.HasPrefix(href,"mailto:")||strings.HasPrefix(href,"javascript:"){continue}
            u,err:=base.Parse(href);if err!=nil||u.Host!=base.Host{continue}
            u.Fragment=""
            links=append(links,u.String())}}
    return links}
func(s *Server)handleStats(w http.ResponseWriter,r *http.Request){si,_:=s.db.CountSites();is,_:=s.db.CountIssues();writeJSON(w,200,map[string]interface{}{"sites":si,"issues":is})}
