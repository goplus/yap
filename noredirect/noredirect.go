package noredirect

import (
	"io"
	"net/http"
	"strings"
	"time"
	_ "unsafe"
)

// FileServer returns a handler that serves HTTP requests
// with the contents of the file system rooted at root.
//
// As a special case, the returned file server redirects any request
// ending in "/index.html" to the same path, without the final
// "index.html".
//
// To use the operating system's file system implementation,
// use http.Dir:
//
//	http.Handle("/", http.FileServer(http.Dir("/tmp")))
//
// To use an fs.FS implementation, use http.FS to convert it:
//
//	http.Handle("/", http.FileServer(http.FS(fsys)))
func FileServer(root http.FileSystem) http.Handler {
	return &fileHandler{root}
}

type fileHandler struct {
	root http.FileSystem
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	serveFile(w, r, f.root, upath)
}

//go:linkname serveContent net/http.serveContent
func serveContent(w http.ResponseWriter, r *http.Request, name string, modtime time.Time, sizeFunc func() (int64, error), content io.ReadSeeker)

//go:linkname toHTTPError net/http.toHTTPError
func toHTTPError(err error) (msg string, httpStatus int)

// name is '/'-separated, not filepath.Separator.
func serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string) {
	f, err := fs.Open(name)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	// serveContent will check modification time
	sizeFunc := func() (int64, error) {
		if size := d.Size(); size >= 0 {
			return size, nil
		}
		return fileSize(f)
	}
	serveContent(w, r, d.Name(), d.ModTime(), sizeFunc, f)
}

func fileSize(content http.File) (size int64, err error) {
	size, err = content.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	_, err = content.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	return size, nil
}
