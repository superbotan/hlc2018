package main

import (
	"errors"
	"go_v1/common"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type SuggestFQ struct {
	stime time.Time
	wtime time.Time

	limit int

	country string
	city    string
	sex     string

	country_id uint8
	city_id    uint16

	id common.Int24

	fail_found bool
	p_one      bool

	a *Account
}

func suggestFilterGet(ctx *fasthttp.RequestCtx) (fa *SuggestFQ, err error) {

	fa = &SuggestFQ{}

	fa.stime = time.Now()
	fa.wtime = fa.stime

	path := string(ctx.URI().Path())
	path = strings.Replace(path, "/accounts/", "", 1)
	path = strings.Replace(path, "/suggest/", "", 1)

	id, err := strconv.Atoi(path)
	if err != nil {
		fa.fail_found = true
		return fa, nil
	}

	if id <= 0 || id >= 1700000 {
		fa.fail_found = true
		return fa, nil
	}

	fa.id = common.Int24Create(int32(id))

	a := accounts.Get(fa.id)
	if a == nil {
		fa.fail_found = true
		return fa, nil
	}

	fa.a = a
	fa.sex = a.SexGet()

	fail := false

	ctx.URI().QueryArgs().VisitAll(func(key []byte, value []byte) {

		k := string(key)
		v := string(value)

		if k == "query_id" {
			return
		}

		if k == "limit" && len(v) > 0 {
			l, err := strconv.Atoi(v)
			if err != nil {
				fail = true
				return
			}
			fa.limit = l

			if l <= 0 {
				fail = true
			}
			return
		}

		if k == "country" && len(v) > 0 {
			fa.country = v
			fa.country_id = countries.GetByName(v)
			return
		}

		if k == "city" && len(v) > 0 {
			fa.city = v
			fa.city_id = cities.GetByName(v)
			return
		}

		fail = true
		return

	})

	if fail {
		return fa, errors.New("")
	}

	return fa, nil
}

func (a1 *Account) Similary(a2 *Account) (similary float64) {

	i := 0
	j := 0

	l1 := len(a1.Likes)
	l2 := len(a2.Likes)

	for i < l1 && j < l2 {
		if a1.Likes[i].Id == a2.Likes[j].Id {
			if a1.Likes[i].Ts == a2.Likes[j].Ts {
				similary = similary + 1
			} else if a2.Likes[j].Ts > a1.Likes[i].Ts {
				similary = similary + float64(1)/float64(a2.Likes[j].Ts-a1.Likes[i].Ts)
			} else {
				similary = similary + float64(1)/float64(a1.Likes[i].Ts-a2.Likes[j].Ts)
			}

			i++
			j++
		} else if a1.Likes[i].Id.More(a2.Likes[j].Id) {
			i++
		} else {
			j++
		}
	}

	return similary
}

func (a1 *Account) SimilaryNotLiked(a2 *Account) []common.Int24 {
	res := make([]common.Int24, 0)

	i := 0
	j := 0

	l1 := len(a1.Likes)
	l2 := len(a2.Likes)

	for i < l1 && j < l2 {
		if a1.Likes[i].Id == a2.Likes[j].Id {
			i++
			j++
		} else if a1.Likes[i].Id.More(a2.Likes[j].Id) {
			i++
		} else {
			res = append(res, a2.Likes[j].Id)
			j++
		}
	}

	for j < l2 {
		res = append(res, a2.Likes[j].Id)
		j++
	}

	return res
}

func (a1 *Account) SimilaryNotLikedAppend(a2 *Account, prom []sim_res_prom, simil float64, pos int, cnt int, slen int) {

	i := 0
	j := 0

	l1 := len(a1.Likes)
	l2 := len(a2.Likes)

	ex := false
	k := 0

	for i < l1 && j < l2 {
		if cnt == 0 {
			return
		}
		if a1.Likes[i].Id == a2.Likes[j].Id {
			i++
			j++
		} else if a1.Likes[i].Id.More(a2.Likes[j].Id) {
			i++
		} else {
			k = 0
			ex = false
			for k < slen {
				if prom[k].id == a2.Likes[j].Id {
					ex = true
					if prom[k].similary < simil {
						prom[k].similary = simil
						prom[k].a = a2
						prom[k].id = a2.Likes[j].Id
						cnt--
					}
					break
				}
				k++
			}
			if !ex {
				prom[pos].similary = simil
				prom[pos].a = a2
				prom[pos].id = a2.Likes[j].Id
				pos++
				cnt--
			}
			j++

		}
	}

	for j < l2 {
		if cnt == 0 {
			return
		}
		k = 0
		ex = false
		for k < slen {
			if prom[k].id == a2.Likes[j].Id {
				ex = true
				if prom[k].similary < simil {
					prom[k].similary = simil
					prom[k].a = a2
					prom[k].id = a2.Likes[j].Id
					cnt--
				}
				break
			}
			k++
		}
		if !ex {
			prom[pos].similary = simil
			prom[pos].a = a2
			prom[pos].id = a2.Likes[j].Id
			pos++
			cnt--
		}
		j++
	}

}

type sim_res_group struct {
	a        *Account
	similary float64
}

type sim_res_out struct {
	similary float64
	id       common.Int24
}

type sim_res_prom struct {
	similary float64
	a        *Account
	id       common.Int24
}

func (fa *SuggestFQ) check_wait() bool {
	if time.Since(fa.wtime) > 200*time.Microsecond {
		time.Sleep(10 * time.Nanosecond)
		fa.wtime = time.Now()
	}
	if !fa.p_one && time.Since(fa.stime) > 10*time.Millisecond {
		//fmt.Println("flt", fa.raw_query)
		fa.p_one = true
	}
	return false
}

func (fa *SuggestFQ) AppendResProm(prom []sim_res_prom, simil float64, asim *Account) {
	if prom[len(prom)-1].similary > simil {
		return
	}

	for _, v := range prom {
		if v.id == asim.ID {
			return
		}
	}

	fa.a.SimilaryNotLikedAppend(asim, prom, simil, fa.limit, fa.limit, fa.limit)

	sort.Slice(prom, func(i int, j int) bool {
		return prom[i].similary > prom[j].similary || prom[i].similary == prom[j].similary && prom[i].id.More(prom[j].id)
	})
}

func suggestLoadData(fa *SuggestFQ) []sim_res_out {

	if fa.city != "" && fa.city_id == 0 {
		return make([]sim_res_out, 0)
	}
	if fa.country != "" && fa.country_id == 0 {
		return make([]sim_res_out, 0)
	}

	plist := make([]sim_res_prom, fa.limit*2, fa.limit*2)

	var asim *Account

	for _, id_lk := range fa.a.Likes {
		for _, id_s := range accounts.Get(id_lk.Id).LikesBack {
			if id_s != fa.id {
				asim = accounts.Get(id_s)
				if fa.is_ok(asim) {

					simil := fa.a.Similary(asim)
					fa.AppendResProm(plist, simil, asim)

					//plist = append(plist, sim_res_group{similary: simil, a: asim})

				}
			}
		}
		fa.check_wait()
	}

	res_out := make([]sim_res_out, 0)

	for _, lp := range plist {

		if lp.id == common.Int24Zero {
			break
		}

		res_out = append(res_out, sim_res_out{similary: lp.similary, id: lp.id})

		if len(res_out) >= fa.limit {
			break
		}
	}

	return res_out
}

func suggestLoadData_old(fa *SuggestFQ) []sim_res_out {

	if fa.city != "" && fa.city_id == 0 {
		return make([]sim_res_out, 0)
	}
	if fa.country != "" && fa.country_id == 0 {
		return make([]sim_res_out, 0)
	}

	plist := make([]sim_res_group, 0)

	var asim *Account

	prep_sim := make(map[common.Int24]struct{})
	for _, id_lk := range fa.a.Likes {
		for _, id_s := range accounts.Get(id_lk.Id).LikesBack {
			if id_s != fa.id {
				asim = accounts.Get(id_s)
				if fa.is_ok(asim) {
					_, ok := prep_sim[id_s]
					if !ok {
						prep_sim[id_s] = struct{}{}

						simil := fa.a.Similary(asim)
						plist = append(plist, sim_res_group{similary: simil, a: asim})
					}
				}
				fa.check_wait()
			}
		}
		//fa.check_wait()
	}

	sort.Slice(plist, func(i int, j int) bool {
		return plist[i].similary > plist[j].similary
	})

	res_out := make([]sim_res_out, 0)
	res_used := make(map[common.Int24]struct{})

	for _, lp := range plist {
		for _, id := range fa.a.SimilaryNotLiked(lp.a) {
			_, ok := res_used[id]
			if !ok {
				res_used[id] = struct{}{}
				res_out = append(res_out, sim_res_out{similary: lp.similary, id: id})
			}
		}

		if len(res_out) >= fa.limit {
			break
		}
	}

	return res_out
}

func (fa *SuggestFQ) is_ok(a *Account) bool {
	if fa.city != "" && a.CityID != fa.city_id {
		return false
	}
	if fa.country != "" && a.CountryID != fa.country_id {
		return false
	}
	if fa.sex == a.SexGet() {
		return true
	}
	return true
}

func (fa *SuggestFQ) get_answer(res_data []sim_res_out) map[string][]map[string]interface{} {

	sort.Slice(res_data, func(i int, j int) bool {
		//return fa.order == -1 && l[i].c > l[j].c || fa.order == 1 && l[i].c < l[j].c
		if res_data[i].similary > res_data[j].similary {
			return true
		} else if res_data[i].similary == res_data[j].similary {
			return res_data[i].id.More(res_data[j].id)
		}
		return false
	})

	res := make(map[string][]map[string]interface{})
	d := make([]map[string]interface{}, 0)
	res["accounts"] = d

	for i := range res_data {
		if i >= fa.limit {
			break
		}

		a := accounts.Get(res_data[i].id)

		elem := make(map[string]interface{})

		elem["id"] = a.ID.Int()
		elem["email"] = a.StringEmail()
		elem["status"] = a.StatusGet()
		if a.SNameID != 0 {
			elem["sname"] = snames.GetByID(a.SNameID)
		}
		if a.FNameID != 0 {
			elem["fname"] = fnames.GetByID(a.FNameID)
		}

		res["accounts"] = append(res["accounts"], elem)
	}

	return res
}
