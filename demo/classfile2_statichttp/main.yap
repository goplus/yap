import "github.com/qiniu/x/http/fs"

static "/", fs.http("https://goplus.org"), false
run ":8080"
