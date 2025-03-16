// Code generated by gop (Go+); DO NOT EDIT.

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/goplus/yap/test"
	"github.com/goplus/yap/ydb"
	_ "github.com/goplus/yap/ydb/mysql"
	_ "github.com/goplus/yap/ydb/sqlite3"
	"github.com/qiniu/x/stringutil"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

const _ = true

type ArticleEntry struct {
	Id     string `CHAR(32) UNIQUE`
	Author string `CHAR(24) INDEX`
	Title  string
	Ctime  time.Time `DATETIME INDEX`
}
type Article struct {
	ArticleEntry
	Body []byte `LONGBLOB`
}
type Tag struct {
	Name string `CHAR(24) UNIQUE(doc)`
	Doc  string `CHAR(32)`
}
type User struct {
	Id       string `id CHAR(32) UNIQUE`
	Spwd     string
	Salt     string
	Nickname string
	Email    string    `CHAR(64) INDEX`
	Tel      string    `CHAR(16) INDEX`
	Born     time.Time `INDEX`
	Ctime    time.Time `DATETIME(6) INDEX`
}
type articles struct {
	ydb.Class
	*AppGen
}
type users struct {
	ydb.Class
	*AppGen
}
type AppGen struct {
	ydb.AppGen
}
//line ydb/demo/foo/articles_ydb.gox:22:1
func (this *articles) API_Add(doc Article) {
//line ydb/demo/foo/articles_ydb.gox:23:1
	this.Insert(doc)
}
//line ydb/demo/foo/articles_ydb.gox:26:1
func (this *articles) API_Get(docId string) (doc Article, err error) {
//line ydb/demo/foo/articles_ydb.gox:27:1
	this.Query("id=?", docId)
//line ydb/demo/foo/articles_ydb.gox:28:1
	this.Ret(&doc)
//line ydb/demo/foo/articles_ydb.gox:29:1
	err = this.LastErr()
//line ydb/demo/foo/articles_ydb.gox:30:1
	return
}
//line ydb/demo/foo/articles_ydb.gox:33:1
func (this *articles) API_SetTags(docId string, tags ...string) {
//line ydb/demo/foo/articles_ydb.gox:34:1
	var oldtags []string
//line ydb/demo/foo/articles_ydb.gox:35:1
	this.Query("tag.doc=?", docId)
//line ydb/demo/foo/articles_ydb.gox:36:1
	this.Ret("tag.name", &oldtags)
//line ydb/demo/foo/articles_ydb.gox:38:1
	tagsAdd, tagsDel := Diff(tags, oldtags)
//line ydb/demo/foo/articles_ydb.gox:39:1
	Info("oldtags:", oldtags, "tags:", tags, "add:", tagsAdd, "del:", tagsDel)
//line ydb/demo/foo/articles_ydb.gox:41:1
	this.Delete("tag.name=?", tagsDel)
//line ydb/demo/foo/articles_ydb.gox:42:1
	this.Insert("tag.doc", docId, "tag.name", tagsAdd)
}
//line ydb/demo/foo/articles_ydb.gox:45:1
func (this *articles) API_Tags(docId string) (tags []string) {
//line ydb/demo/foo/articles_ydb.gox:46:1
	this.Query("tag.doc=?", docId)
//line ydb/demo/foo/articles_ydb.gox:47:1
	this.Ret("tag.name", &tags)
//line ydb/demo/foo/articles_ydb.gox:48:1
	return
}
//line ydb/demo/foo/articles_ydb.gox:51:1
func (this *articles) API_ListByTag(tag string) (result []ArticleEntry) {
//line ydb/demo/foo/articles_ydb.gox:52:1
	var ids []string
//line ydb/demo/foo/articles_ydb.gox:53:1
	this.Query("tag.name=?", tag)
//line ydb/demo/foo/articles_ydb.gox:54:1
	this.Ret("tag.doc", &ids)
//line ydb/demo/foo/articles_ydb.gox:56:1
	this.Query("id=?", ids)
//line ydb/demo/foo/articles_ydb.gox:57:1
	this.Ret(&result)
//line ydb/demo/foo/articles_ydb.gox:58:1
	return
}
//line ydb/demo/foo/articles_ydb.gox:61:1
func (this *articles) API_ListByAuthor(author string) (result []ArticleEntry) {
//line ydb/demo/foo/articles_ydb.gox:62:1
	this.Query("author=?", author)
//line ydb/demo/foo/articles_ydb.gox:63:1
	this.Ret(&result)
//line ydb/demo/foo/articles_ydb.gox:64:1
	return
}
//line ydb/demo/foo/articles_ydb.gox:67
func (this *articles) Main() {
//line ydb/demo/foo/articles_ydb.gox:67:1
	this.Engine__0("mysql")
//line ydb/demo/foo/articles_ydb.gox:69:1
	ydb.Gopt_Sql_Gopx_Table[Article](this, "v0.1.0")
//line ydb/demo/foo/articles_ydb.gox:70:1
	ydb.Gopt_Sql_Gopx_Table[Tag](this, "v0.1.0")
//line ydb/demo/foo/articles_ydb.gox:72:1
	this.Use("article")
//line ydb/demo/foo/articles_ydb.gox:74:1
	doc1 := Article{}
//line ydb/demo/foo/articles_ydb.gox:75:1
	doc1.Id, doc1.Author, doc1.Title = "123", "abc", "title1"
//line ydb/demo/foo/articles_ydb.gox:76:1
	this.Gop_Exec("add", doc1)
//line ydb/demo/foo/articles_ydb.gox:77:1
	this.Ret()
//line ydb/demo/foo/articles_ydb.gox:79:1
	this.Gop_Exec("add", doc1)
//line ydb/demo/foo/articles_ydb.gox:80:1
	this.Ret(ydb.ErrDuplicated)
//line ydb/demo/foo/articles_ydb.gox:82:1
	doc2 := Article{}
//line ydb/demo/foo/articles_ydb.gox:83:1
	doc2.Id, doc2.Author, doc2.Title = "124", "efg", "title2"
//line ydb/demo/foo/articles_ydb.gox:84:1
	this.Gop_Exec("add", doc2)
//line ydb/demo/foo/articles_ydb.gox:85:1
	this.Ret()
//line ydb/demo/foo/articles_ydb.gox:87:1
	doc3 := Article{}
//line ydb/demo/foo/articles_ydb.gox:88:1
	doc3.Id, doc3.Author, doc3.Title = "125", "efg", "title3"
//line ydb/demo/foo/articles_ydb.gox:89:1
	this.Gop_Exec("add", doc3)
//line ydb/demo/foo/articles_ydb.gox:90:1
	this.Ret()
//line ydb/demo/foo/articles_ydb.gox:92:1
	doc4 := Article{}
//line ydb/demo/foo/articles_ydb.gox:93:1
	doc4.Id, doc4.Author, doc4.Title = "225", "abc", "title4"
//line ydb/demo/foo/articles_ydb.gox:94:1
	this.Gop_Exec("add", doc4)
//line ydb/demo/foo/articles_ydb.gox:95:1
	this.Ret()
//line ydb/demo/foo/articles_ydb.gox:97:1
	doc5 := Article{}
//line ydb/demo/foo/articles_ydb.gox:98:1
	doc5.Id, doc5.Author, doc5.Title = "555", "abc", "title5"
//line ydb/demo/foo/articles_ydb.gox:99:1
	this.Gop_Exec("add", doc5)
//line ydb/demo/foo/articles_ydb.gox:100:1
	this.Ret()
//line ydb/demo/foo/articles_ydb.gox:102:1
	this.Gop_Exec("get", doc1.Id)
//line ydb/demo/foo/articles_ydb.gox:103:1
	this.Ret(doc1)
//line ydb/demo/foo/articles_ydb.gox:105:1
	this.Gop_Exec("get", doc2.Id)
//line ydb/demo/foo/articles_ydb.gox:106:1
	this.Ret(doc2)
//line ydb/demo/foo/articles_ydb.gox:108:1
	this.Gop_Exec("get", "unknown")
//line ydb/demo/foo/articles_ydb.gox:109:1
	test.Gopt_Case_MatchAny(this, this.Out(1), ydb.ErrNoRows)
//line ydb/demo/foo/articles_ydb.gox:111:1
	this.Gop_Exec("setTags", doc1.Id, "tag1", "tag2")
//line ydb/demo/foo/articles_ydb.gox:112:1
	this.Ret()
//line ydb/demo/foo/articles_ydb.gox:114:1
	this.Gop_Exec("tags", doc1.Id)
//line ydb/demo/foo/articles_ydb.gox:115:1
	this.Ret(test.Set__0("tag2", "tag1"))
//line ydb/demo/foo/articles_ydb.gox:117:1
	this.Gop_Exec("setTags", doc1.Id, "tag1", "tag3")
//line ydb/demo/foo/articles_ydb.gox:118:1
	this.Ret()
//line ydb/demo/foo/articles_ydb.gox:120:1
	this.Gop_Exec("tags", doc1.Id)
//line ydb/demo/foo/articles_ydb.gox:121:1
	this.Ret(test.Set__0("tag1", "tag3"))
//line ydb/demo/foo/articles_ydb.gox:123:1
	this.Gop_Exec("setTags", doc2.Id, "tag1", "tag5")
//line ydb/demo/foo/articles_ydb.gox:124:1
	this.Gop_Exec("setTags", doc3.Id, "tag1", "tag3")
//line ydb/demo/foo/articles_ydb.gox:125:1
	this.Gop_Exec("setTags", doc4.Id, "tag2", "tag3")
//line ydb/demo/foo/articles_ydb.gox:126:1
	this.Gop_Exec("setTags", doc5.Id, "tag5", "tag3")
//line ydb/demo/foo/articles_ydb.gox:128:1
	this.Gop_Exec("listByTag", "tag1")
//line ydb/demo/foo/articles_ydb.gox:129:1
	this.Ret(test.Set__2(doc3.ArticleEntry, doc1.ArticleEntry, doc2.ArticleEntry))
//line ydb/demo/foo/articles_ydb.gox:131:1
	this.Gop_Exec("listByTag", "tag3")
//line ydb/demo/foo/articles_ydb.gox:132:1
	this.Ret(test.Set__2(doc3.ArticleEntry, doc4.ArticleEntry, doc1.ArticleEntry, doc5.ArticleEntry))
//line ydb/demo/foo/articles_ydb.gox:134:1
	this.Gop_Exec("listByAuthor", "eft")
//line ydb/demo/foo/articles_ydb.gox:135:1
	this.Ret(test.Set__2(doc2.ArticleEntry, doc3.ArticleEntry))
}

var ErrNoEmailAndTel = errors.New("no email and telephone")
var rnd = rand.New(rand.NewSource(time.Now().UnixMicro()))
//line ydb/demo/foo/foo.gop:23:1
func Rand() string {
//line ydb/demo/foo/foo.gop:24:1
	return strconv.FormatInt(rnd.Int63(), 36)
}
//line ydb/demo/foo/foo.gop:27:1
func Hs256(pwd string, salt string) string {
//line ydb/demo/foo/foo.gop:28:1
	b := hmac.New(sha256.New, []byte(salt)).Sum([]byte(pwd))
//line ydb/demo/foo/foo.gop:29:1
	return base64.RawURLEncoding.EncodeToString(b)
}
//line ydb/demo/foo/foo.gop:32:1
func Diff(new []string, old []string) (add []string, del []string) {
//line ydb/demo/foo/foo.gop:33:1
	sort.Strings(new)
//line ydb/demo/foo/foo.gop:34:1
	sort.Strings(old)
//line ydb/demo/foo/foo.gop:35:1
	return stringutil.Diff(new, old)
}
//line ydb/demo/foo/foo.gop:38:1
// Info calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Info(args ...interface{}) {
//line ydb/demo/foo/foo.gop:41:1
	log.Println(args...)
}
//line ydb/demo/foo/users_ydb.gox:19:1
func (this *users) API_Register(id string, pwd string, nickname string, email string, tel string, ctime time.Time) error {
//line ydb/demo/foo/users_ydb.gox:20:1
	if email == "" && tel == "" {
//line ydb/demo/foo/users_ydb.gox:21:1
		return ErrNoEmailAndTel
	}
//line ydb/demo/foo/users_ydb.gox:23:1
	this.Limit__1(3, "email=?", email)
//line ydb/demo/foo/users_ydb.gox:24:1
	this.Limit__1(3, "tel=?", tel)
//line ydb/demo/foo/users_ydb.gox:26:1
	salt := Rand()
//line ydb/demo/foo/users_ydb.gox:27:1
	spwd := Hs256(pwd, salt)
//line ydb/demo/foo/users_ydb.gox:28:1
	this.Insert(&User{Id: id, Spwd: spwd, Salt: salt, Nickname: nickname, Email: email, Tel: tel, Ctime: ctime})
//line ydb/demo/foo/users_ydb.gox:29:1
	return nil
}
//line ydb/demo/foo/users_ydb.gox:32:1
func (this *users) API_Login(id string, pwd string) bool {
//line ydb/demo/foo/users_ydb.gox:33:1
	var spwd, salt string
//line ydb/demo/foo/users_ydb.gox:34:1
	this.Query("id=?", id)
//line ydb/demo/foo/users_ydb.gox:35:1
	this.Ret("salt", &salt, "spwd", &spwd)
//line ydb/demo/foo/users_ydb.gox:36:1
	if this.NoRows() {
//line ydb/demo/foo/users_ydb.gox:37:1
		return false
	}
//line ydb/demo/foo/users_ydb.gox:39:1
	return Hs256(pwd, salt) == spwd
}
//line ydb/demo/foo/users_ydb.gox:42
func (this *users) Main() {
//line ydb/demo/foo/users_ydb.gox:42:1
	this.Engine__0("sqlite3")
//line ydb/demo/foo/users_ydb.gox:44:1
	ydb.Gopt_Sql_Gopx_Table[User](this, "user v0.1.0")
//line ydb/demo/foo/users_ydb.gox:46:1
	this.Use("user")
//line ydb/demo/foo/users_ydb.gox:48:1
	this.Gop_Exec("register", "user", "pwd", "nickname", "", "", time.Now())
//line ydb/demo/foo/users_ydb.gox:49:1
	this.Ret(ErrNoEmailAndTel)
//line ydb/demo/foo/users_ydb.gox:50:1
	this.Gop_Exec("register", "user", "pwd", "nickname", "user@foo.com", "", time.Now())
//line ydb/demo/foo/users_ydb.gox:51:1
	this.Ret(nil)
//line ydb/demo/foo/users_ydb.gox:52:1
	this.Gop_Exec("register", "user", "pwd", "nickname", "user@foo.com", "13500000000", time.Now())
//line ydb/demo/foo/users_ydb.gox:53:1
	this.Ret(ydb.ErrDuplicated)
//line ydb/demo/foo/users_ydb.gox:55:1
	this.Gop_Exec("login", "", "")
//line ydb/demo/foo/users_ydb.gox:56:1
	this.Ret(false)
//line ydb/demo/foo/users_ydb.gox:57:1
	this.Gop_Exec("login", "user", "pwd")
//line ydb/demo/foo/users_ydb.gox:58:1
	this.Ret(true)
}
func (this *AppGen) Main() {
	ydb.Gopt_AppGen_Main(this, new(articles), new(users))
}
func main() {
	new(AppGen).Main()
}
