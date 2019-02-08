package main

import (
	json "encoding/json"
	fmt "fmt"
	"go_v1/common"
	"strconv"
	"strings"
	"sync"

	"github.com/myfantasy/mfe"
	"github.com/valyala/fasthttp"
)

var (
	muw sync.Mutex
)

func accounts_add_likes(ctx *fasthttp.RequestCtx) {
	// muw.Lock()
	// defer muw.Unlock()

	var d LikesAdd

	if err := json.Unmarshal(ctx.PostBody(), &d); err != nil {
		ctx.Response.SetStatusCode(400)
		return
	}

	for _, v := range d.Likes {
		if accounts.GetLock(common.Int24Create(v.Liker)) == nil {
			ctx.Response.SetStatusCode(400)
			return
		}
		if accounts.GetLock(common.Int24Create(v.Likee)) == nil {
			ctx.Response.SetStatusCode(400)
			return
		}
		if v.Ts >= GetTimeNow() {
			ctx.Response.SetStatusCode(400)
			return
		}
	}

	func() {
		muw.Lock()
		defer muw.Unlock()

		var a *Account
		var ai *Account

		for _, v := range d.Likes {
			a = accounts.GetLock(common.Int24Create(v.Liker))
			a.AppendLike(v.Ts, common.Int24Create(v.Likee))

			ai = accounts.GetLock(common.Int24Create(v.Likee))

			ai.AppendLikeBack(a.ID)
		}
	}()
	ctx.Response.SetStatusCode(202)
	fmt.Fprint(ctx, "{}")
}

func accounts_add(ctx *fasthttp.RequestCtx) {
	// muw.Lock()
	// defer muw.Unlock()

	var d AccountAdd

	if err := json.Unmarshal(ctx.PostBody(), &d); err != nil {
		ctx.Response.SetStatusCode(400)
		//fmt.Fprint(ctx, "1")
		return
	}

	if accounts.GetLock(common.Int24Create(d.ID)) != nil {
		ctx.Response.SetStatusCode(400)
		// fmt.Fprintln(ctx, d.Id)
		// fmt.Fprint(ctx, users.GetLock(d.Id))
		return
	}

	if d.Email == "" || !(strings.Index(d.Email, "@") > 0 && strings.Index(d.Email, ".") > 0) {
		ctx.Response.SetStatusCode(400)
		// fmt.Fprint(ctx, "6")
		return
	}

	if d.Phone != "" && index.exists_phone(d.Phone) {
		ctx.Response.SetStatusCode(400)
		// fmt.Fprint(ctx, "3")
		return
	}

	if d.Sex == "" || !mfe.InS(d.Sex, "m", "f") {
		ctx.Response.SetStatusCode(400)
		// fmt.Fprint(ctx, "4")
		return
	}
	if d.Status == "" || !mfe.InS(d.Status, "свободны", "заняты", "всё сложно") {
		ctx.Response.SetStatusCode(400)
		// fmt.Fprint(ctx, "5")
		return
	}

	if index.exists_email(d.Email) {
		ctx.Response.SetStatusCode(400)
		// fmt.Fprint(ctx, "3")
		return
	}

	if d.Likes != nil {
		for _, v := range d.Likes {
			if accounts.GetLock(common.Int24Create(v.ID)) == nil {
				ctx.Response.SetStatusCode(400)
				return
			}
		}
	}

	go func() {
		muw.Lock()
		defer muw.Unlock()

		a := AccountCreate(d)

		accounts.Append(&a)
		index.Append(&a)
		groupindex.Append(&a)
		(&a).FillOtherLikeBack()
	}()

	ctx.Response.SetStatusCode(201)
	fmt.Fprint(ctx, "{}")
}

func accounts_upd(ctx *fasthttp.RequestCtx) {
	// muw.Lock()
	// defer muw.Unlock()

	path := string(ctx.URI().Path())
	if len(path) < 12 {
		ctx.Response.SetStatusCode(404)
		return
	}
	r0 := []rune(path)
	id_str := string(r0[10 : len(r0)-1])

	id, err := strconv.Atoi(id_str)

	if err != nil {
		ctx.Response.SetStatusCode(404)
		return
	}
	ac := accounts.GetLock(common.Int24Create(int32(id)))
	if ac == nil {
		ctx.Response.SetStatusCode(404)
		return
	}

	var d AccountUpd

	if err := json.Unmarshal(ctx.PostBody(), &d); err != nil {
		ctx.Response.SetStatusCode(400)
		return
	}

	if d.Phone != "" && index.exists_phone(d.Phone) && ac.ID != index.phone_ext_get(d.Phone) {
		ctx.Response.SetStatusCode(400)
		return
	}

	if d.Sex != "" && !mfe.InS(d.Sex, "m", "f") {
		ctx.Response.SetStatusCode(400)
		return
	}
	if d.Status != "" && !mfe.InS(d.Status, "свободны", "заняты", "всё сложно") {
		ctx.Response.SetStatusCode(400)
		return
	}
	if d.Email != "" && !(strings.Index(d.Email, "@") > 0 && strings.Index(d.Email, ".") > 0) {
		ctx.Response.SetStatusCode(400)
		return
	}

	if d.Email != "" && index.exists_email(d.Email) && ac.ID != index.email_ext_get(d.Email) {
		ctx.Response.SetStatusCode(400)
		return
	}

	if d.Likes != nil {
		for _, v := range d.Likes {
			if accounts.GetLock(common.Int24Create(v.ID)) == nil {
				ctx.Response.SetStatusCode(400)
				return
			}
		}
	}

	go func() {
		muw.Lock()
		defer muw.Unlock()

		if d.Likes != nil {

			ac.RemoveOtherLikeBack()

		}

		groupindex.Remove(ac)
		index.Remove(ac)

		an := ac.AccountUpdate(&d)

		accounts.Append(&an)

		index.Append(&an)
		groupindex.Append(&an)

		if d.Likes != nil {

			(&an).FillOtherLikeBack()

		}
	}()
	ctx.Response.SetStatusCode(202)
	fmt.Fprint(ctx, "{}")
}
