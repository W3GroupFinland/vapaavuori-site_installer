package app_base

import (
	"net/http"
)

type FileServer struct {
	Handler http.Handler
	Acl     []HttpAclHandlerFunc
}

func NewFileServer(dir string, prefix string) *FileServer {
	fs := http.Dir(dir)
	fileHandler := http.FileServer(fs)

	return &FileServer{Handler: http.StripPrefix(prefix, fileHandler)}
}

func (fs *FileServer) AddAcl(fn HttpAclHandlerFunc) *FileServer {
	fs.Acl = append(fs.Acl, fn)

	return fs
}

func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, fn := range fs.Acl {
		accepted := fn(w, r)
		if !accepted {
			http.Error(w, "403 Forbidden.", 403)
			return
		}
	}

	fs.Handler.ServeHTTP(w, r)
}
