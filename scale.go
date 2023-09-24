//go:build tinygo || js || wasm

package scale

import (
	"embed"
	"fmt"
	signature "github.com/loopholelabs/scale-signature-http"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
)

//go:embed out
//go:embed out/_next
//go:embed out/_next/static/**/*.js
//go:embed out/_next/static/chunks/**/*.js
var nextjs embed.FS

func Scale(ctx *signature.Context) (*signature.Context, error) {
	uri, err := url.Parse(ctx.Request().URI())
	if err != nil {
		return ctx, fmt.Errorf("error parsing uri: %w", err)
	}

	if uri.Path == "" || uri.Path == "/" {
		uri.Path = "/index.html"
	}

	f, err := nextjs.Open("out" + uri.Path)
	if err != nil {
		ctx.Response().SetStatusCode(http.StatusNotFound)
		ctx.Response().SetBody(fmt.Sprintf("404: %s not found", uri.Path))
		return ctx, nil
	}
	body, err := io.ReadAll(f)
	if err != nil {
		return ctx, fmt.Errorf("error reading %s: %w", uri.Path, err)
	}

	_ = mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
	_ = mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mimetype := mime.TypeByExtension(filepath.Ext(uri.Path))
	if mimetype == "" {
		mimetype = http.DetectContentType(body[:512])
	}
	ctx.Response().Headers().Set("content-type", []string{mimetype})

	ctx.Response().SetBodyBytes(body)
	return ctx.Next()
}
