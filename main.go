package main

import (
	"encoding/json"
	"fmt"
	"os"
	"bytes"
	"os/signal"
	"strings"
	"syscall"

	"github.com/prometheus/common/log"
	"github.com/valyala/fasthttp"
)

var (
	isFirstPhaseFinish = false
	
	filterAccess = sync.RWMutex{}
	filtersCache = map[string][]byte{}
	
	groupAccess = sync.RWMutex{}
	groupsCache = map[string][]byte{}
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {

	if string(ctx.URI().Path()) == "/accounts/filter/" {
		accountsFilter(ctx)
		return
	}
	if string(ctx.URI().Path()) == "/accounts/group/" {
		accountsGroup(ctx)
		return
	}

	if string(ctx.URI().Path()) == "/accounts/new/" {
		accountsNew(ctx)
		return
	}
	if string(ctx.URI().Path()) == "/accounts/likes/" {
		accountsLikes(ctx)
		return
	}
	if strings.Index(string(ctx.URI().Path()), "/accounts/") == 0 && strings.Index(string(ctx.URI().Path()), "/suggest/") > 0 {
		accountsSuggest(ctx)
		return
	}
	if strings.Index(string(ctx.URI().Path()), "/accounts/") == 0 && strings.Index(string(ctx.URI().Path()), "/recommend/") > 0 {
		accountsRecommend(ctx)
		return
	}

	if strings.Index(string(ctx.URI().Path()), "/accounts/") == 0 {
		accountsupd(ctx)
		return
	}

	ctx.Response.SetStatusCode(404)
}
func main() {

	api := &fasthttp.Server{
		Handler: fastHTTPHandler,
	}

	go readData()

	serverErrors := make(chan error, 1)
	go func() {
		log.Infof("Listen and serve :80")
		//serverErrors <- api.ListenAndServe(":80")
		serverErrors <- api.ListenAndServe(":8080")
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Can`t start server; %v", err)

	case <-osSignals:
		log.Infof("Start shutdown...")
		go func() {
			if err := api.Shutdown(); err != nil {
				log.Infof("Graceful shutdown did not complete in 5s : %v", err)
			}
		}()
	}

	fmt.Println("Good Bye")
}

func accountsFilter(ctx *fasthttp.RequestCtx) {
	
	cacheKey := getCacheKey(ctx.QueryArgs())
	if isFirstPhaseFinish {

		filterAccess.RLock()
		res, ok := filtersCache[cacheKey]
		filterAccess.RUnlock()

		if ok {
			ctx.SetContentType("application/json")
			ctx.Write(res)
			return
		}
	}

	fa, err := accountsFilterGet(ctx)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		return
	}

	ids := preLoadData(fa)

	r := filterMakeResult(ids, fa)

	b, err := json.Marshal(r)

	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	ctx.Response.Header.Add("content-type", "application/json")

	fmt.Fprint(ctx, string(b))
	
	if isFirstPhaseFinish {
		filterAccess.Lock()
		filtersCache[cacheKey] = b
		filterAccess.Unlock()
	}

}

func accountsGroup(ctx *fasthttp.RequestCtx) {
	
	cacheKey := getCacheKey(ctx.QueryArgs())
	if isFirstPhaseFinish {

		groupAccess.RLock()
		res, ok := groupsCache[cacheKey]
		groupAccess.RUnlock()

		if ok {
			ctx.SetContentType("application/json")
			ctx.Write(res)
			return
		}
	}

	fa, err := groupFilterGet(ctx)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		return
	}

	res_data := groupLoadData(fa)

	r := fa.get_answer(res_data)

	b, err := json.Marshal(r)

	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	ctx.Response.Header.Add("content-type", "application/json")

	fmt.Fprint(ctx, string(b))
	
	if isFirstPhaseFinish {
		groupAccess.Lock()
		groupsCache[cacheKey] = b
		groupAccess.Unlock()
	}
}

func accountsNew(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != "POST" {
		ctx.Response.SetStatusCode(404)
		return
	}

	accounts_add(ctx)
	isFirstPhaseFinish = true
}
func accountsLikes(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != "POST" {
		ctx.Response.SetStatusCode(404)
		return
	}

	accounts_add_likes(ctx)
}

func accountsupd(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != "POST" {
		ctx.Response.SetStatusCode(404)
		return
	}

	accounts_upd(ctx)
}

func accountsSuggest(ctx *fasthttp.RequestCtx) {

	if string(ctx.Method()) != "GET" {
		ctx.Response.SetStatusCode(404)
		return
	}

	fa, err := suggestFilterGet(ctx)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		return
	}
	if fa.fail_found {
		ctx.Response.SetStatusCode(404)
		return
	}

	res_data := suggestLoadData_old(fa)

	r := fa.get_answer(res_data)

	b, err := json.Marshal(r)

	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	// ctx.Response.SetStatusCode(404)
	// return

	ctx.Response.Header.Add("content-type", "application/json")

	fmt.Fprint(ctx, string(b))
}

func accountsRecommend(ctx *fasthttp.RequestCtx) {

	if string(ctx.Method()) != "GET" {
		ctx.Response.SetStatusCode(404)
		return
	}

	fa, err := recommendFilterGet(ctx)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		return
	}
	if fa.fail_found {
		ctx.Response.SetStatusCode(404)
		return
	}

	res_data := recommendLoadData(fa)

	r := fa.get_answer(res_data)

	b, err := json.Marshal(r)

	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	// ctx.Response.SetStatusCode(404)
	// return

	ctx.Response.Header.Add("content-type", "application/json")

	fmt.Fprint(ctx, string(b))
}

func getCacheKey(args *fasthttp.Args) string {
	cacheKey := bytes.Buffer{}
	args.VisitAll(func(k, v []byte) {
		if string(k) == "query_id" {
			return
		}
		cacheKey.Write(k)
		cacheKey.Write(v)
	})
	return cacheKey.String()
}
