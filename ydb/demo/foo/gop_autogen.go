package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/goplus/yap/ydb"
	_ "github.com/goplus/yap/ydb/mysql"
	_ "github.com/goplus/yap/ydb/sqlite3"
	"math/rand"
	"strconv"
	"time"
)

const _ = true

type User struct {
	Id       string `id CHAR(32) UNIQUE`
	Spwd     string
	Salt     string
	Nickname string
	Email    string    `INDEX`
	Tel      string    `INDEX`
	Born     time.Time `INDEX`
	Ctime    time.Time `DATETIME(6) INDEX`
}
type ArticleEntry struct {
	Id     string `UNIQUE`
	Author string `INDEX`
	Title  string
	Ctime  time.Time `DATATIME INDEX`
}
type Article struct {
	ArticleEntry
	Body []byte `LONGBLOB`
}
type Tag struct {
	Name    string `UNIQUE(article)`
	Article string
}
type foo struct {
	ydb.Sql
}

func main() {
//line ydb/demo/foo/foo_ydb.gox:82:1
	ydb.Gopt_AppGen_Main(new(ydb.AppGen), new(foo))
}

var ErrNoEmailAndTel = errors.New("no email and telephone")
var rnd = rand.New(rand.NewSource(time.Now().UnixMicro()))
//line ydb/demo/foo/foo.gop:19:1
func Rand() string {
//line ydb/demo/foo/foo.gop:20:1
	return strconv.FormatInt(rnd.Int63(), 36)
}
//line ydb/demo/foo/foo.gop:23:1
func Hs256(pwd string, salt string) string {
//line ydb/demo/foo/foo.gop:24:1
	b := hmac.New(sha256.New, []byte(salt)).Sum([]byte(pwd))
//line ydb/demo/foo/foo.gop:25:1
	return base64.RawURLEncoding.EncodeToString(b)
}
//line ydb/demo/foo/foo_ydb.gox:36
func (this *foo) Main() {
//line ydb/demo/foo/foo_ydb.gox:36:1
	this.Engine__0("mysql")
//line ydb/demo/foo/foo_ydb.gox:38:1
	ydb.Gopt_Sql_Gopx_Table[User](this, "user v0.1.0")
//line ydb/demo/foo/foo_ydb.gox:40:1
	ydb.Gopt_Sql_Gopx_Table[Article](this, "v0.1.0")
//line ydb/demo/foo/foo_ydb.gox:41:1
	this.From("oldart v0.9.1", func() {
	})
//line ydb/demo/foo/foo_ydb.gox:46:1
	ydb.Gopt_Sql_Gopx_Table[Tag](this, "v0.1.0")
//line ydb/demo/foo/foo_ydb.gox:48:1
	this.Class("Users", func() {
//line ydb/demo/foo/foo_ydb.gox:49:1
		this.Use("user")
//line ydb/demo/foo/foo_ydb.gox:51:1
		this.Api("register", func(id string, pwd string, nickname string, email string, tel string, ctime time.Time) error {
//line ydb/demo/foo/foo_ydb.gox:52:1
			if email == "" && tel == "" {
//line ydb/demo/foo/foo_ydb.gox:53:1
				return ErrNoEmailAndTel
			}
//line ydb/demo/foo/foo_ydb.gox:55:1
			this.Limit__2(3, "email=?", email)
//line ydb/demo/foo/foo_ydb.gox:56:1
			this.Limit__2(3, "tel=?", tel)
//line ydb/demo/foo/foo_ydb.gox:58:1
			salt := Rand()
//line ydb/demo/foo/foo_ydb.gox:59:1
			spwd := Hs256(pwd, salt)
//line ydb/demo/foo/foo_ydb.gox:60:1
			this.Insert__1(&User{Id: id, Spwd: spwd, Salt: salt, Nickname: nickname, Email: email, Tel: tel, Ctime: ctime})
//line ydb/demo/foo/foo_ydb.gox:61:1
			return nil
		})
//line ydb/demo/foo/foo_ydb.gox:63:1
		this.Call__1("user", "pwd", "nickname", "", "", time.Now())
//line ydb/demo/foo/foo_ydb.gox:64:1
		this.Ret__1(ErrNoEmailAndTel)
//line ydb/demo/foo/foo_ydb.gox:65:1
		this.Call__1("user", "pwd", "nickname", "user@foo.com", "", time.Now())
//line ydb/demo/foo/foo_ydb.gox:66:1
		this.Ret__0(nil)
//line ydb/demo/foo/foo_ydb.gox:67:1
		this.Call__1("user", "pwd", "nickname", "user@foo.com", "13500000000", time.Now())
//line ydb/demo/foo/foo_ydb.gox:68:1
		this.Ret__1(ydb.ErrDuplicated)
//line ydb/demo/foo/foo_ydb.gox:70:1
		this.Api("login", func(id string, pwd string) bool {
//line ydb/demo/foo/foo_ydb.gox:71:1
			var spwd, salt string
//line ydb/demo/foo/foo_ydb.gox:72:1
			this.Query__1("id=?", id)
//line ydb/demo/foo/foo_ydb.gox:73:1
			this.Ret__1("salt", &salt, "spwd", &spwd)
//line ydb/demo/foo/foo_ydb.gox:74:1
			return Hs256(pwd, salt) == spwd
		})
//line ydb/demo/foo/foo_ydb.gox:76:1
		this.Call__1("", "")
//line ydb/demo/foo/foo_ydb.gox:77:1
		this.Ret__1(false)
//line ydb/demo/foo/foo_ydb.gox:78:1
		this.Call__1("user", "pwd")
//line ydb/demo/foo/foo_ydb.gox:79:1
		this.Ret__1(true)
	})
//line ydb/demo/foo/foo_ydb.gox:82:1
	this.Class("Articles", func() {
//line ydb/demo/foo/foo_ydb.gox:83:1
		this.Use("article")
//line ydb/demo/foo/foo_ydb.gox:85:1
		this.Api("listByTag", func(tag string) (result []ArticleEntry) {
//line ydb/demo/foo/foo_ydb.gox:86:1
			var ids []string
//line ydb/demo/foo/foo_ydb.gox:87:1
			this.Query__1("tag.name=?", tag)
//line ydb/demo/foo/foo_ydb.gox:88:1
			this.Ret__1("tag.article", &ids)
//line ydb/demo/foo/foo_ydb.gox:90:1
			this.Query__1("id=?", ids)
//line ydb/demo/foo/foo_ydb.gox:91:1
			this.Ret__1(&result)
//line ydb/demo/foo/foo_ydb.gox:92:1
			return
		})
//line ydb/demo/foo/foo_ydb.gox:95:1
		this.Api("listByAuthor", func(author string) (result []ArticleEntry) {
//line ydb/demo/foo/foo_ydb.gox:96:1
			this.Query__1("author=?", author)
//line ydb/demo/foo/foo_ydb.gox:97:1
			this.Ret__1(&result)
//line ydb/demo/foo/foo_ydb.gox:98:1
			return
		})
	})
}
