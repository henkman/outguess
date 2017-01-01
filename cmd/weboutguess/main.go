package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/henkman/outguess"
	"github.com/valyala/fasthttp"
)

func index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	ctx.WriteString(`<html>
<head>
<title>outguess</title>
</head>
<body>
<div>
	<h2>get embedded data from jpg</h2>
	<form action="/get" method="post" enctype="multipart/form-data">
		<label>
		jpg
		<input type="file" name="file" />
		</label><br/>
		<label>
		key (leave empty for default)
		<input type="text" name="key" />
		</label><br/>
		<input type="submit" value="get"/>
	</form>
</div>
</body>
</html>
`)
}

func get(ctx *fasthttp.RequestCtx) {
	filehead, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := filehead.Open()
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	var key []byte
	_key := ctx.FormValue("key")
	if len(_key) > 0 {
		key = _key
	}
	if err := outguess.Get(file, ctx.Response.BodyWriter(), key); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	_host    string
	_noindex bool
)

func init() {
	flag.StringVar(&_host, "host", "0.0.0.0:8080", "listen address")
	flag.BoolVar(&_noindex, "noindex", false, "set if you want no index page")
	flag.Parse()
}

func main() {
	router := fasthttprouter.New()
	if !_noindex {
		router.GET("/", index)
	}
	router.POST("/get", get)
	log.Fatal(fasthttp.ListenAndServe(_host, router.Handler))
}
