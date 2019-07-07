package session

import (
	"time"
	"sync"
	"fmt"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/utils"
)

//用一个map来存取多条session
var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano()/1000000
}

//删除过期session
func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

//加载DB中的session
func LoadSessionsFromDB() *sync.Map{
	r, err:=dbops.RetrieveAllSessions()
	if err!=nil {
		return nil
	}
	r.Range(func(k, v interface{}) bool {
		ss:=v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
	return sessionMap
}

//创建一个新的session
func GenerateNewSessionId(un string) string {
	id, _:=utils.NewUUID()
	ct:=nowInMilli()
	ttl:=ct + 30*60*1000//severside session valid time：30 min
	ss:=&defs.SimpleSession{Username: un, TTL:ttl}
	sessionMap.Store(id, ss)
	err:=dbops.InsertSession(id, ttl, un)
	if err!=nil{
		return fmt.Sprintf("Error of GenerateNewSessionId: %s", err)
	}
	return id
}


//处理过期session
func IsSessionExpired(sid string) (string, bool) {
	ss, ok:=sessionMap.Load(sid)
	if ok{
		ct:=nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}
		return ss.(*defs.SimpleSession).Username, false
	}
	return "", true
}