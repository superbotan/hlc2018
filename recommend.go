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

type RecommendFQ struct {
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

func recommendFilterGet(ctx *fasthttp.RequestCtx) (fa *RecommendFQ, err error) {

	fa = &RecommendFQ{}

	fa.stime = time.Now()
	fa.wtime = fa.stime

	path := string(ctx.URI().Path())
	path = strings.Replace(path, "/accounts/", "", 1)
	path = strings.Replace(path, "/recommend/", "", 1)

	id, err := strconv.Atoi(path)
	if err != nil {
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

func (fa *RecommendFQ) check_wait() bool {
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

type compat struct {
	is_premium bool
	status     uint8
	int_count  uint8
	delta_age  int32
	id         common.Int24
}

func (a *Account) CreateCompatFromID(id common.Int24, status uint8, iarr []uint8, is_premium bool) compat {
	res := compat{id: id, status: status, is_premium: is_premium}

	ac := accounts.Get(id)

	res.delta_age = common.Abs(a.Birth - ac.Birth)

	res.int_count = ac.InterestsIDs.ComplCountArr(iarr)

	return res
}

func (fa RecommendFQ) CreateCompatFromID(id common.Int24, status uint8, iarr []uint8, is_premium bool) (c compat, ok bool) {
	ac := accounts.Get(id)

	if fa.country != "" && ac.CountryID != fa.country_id {
		return c, false
	}

	res := compat{id: id, status: status, is_premium: is_premium}

	res.delta_age = common.Abs(fa.a.Birth - ac.Birth)

	res.int_count = ac.InterestsIDs.ComplCountArr(iarr)

	return res, true
}

func (c compat) More(c2 compat) bool {
	if c.id == common.Int24Zero {
		return false
	}

	if c2.id == common.Int24Zero {
		return true
	}

	if c.is_premium && !c2.is_premium {
		return true
	}

	if c.is_premium == c2.is_premium {
		if c.status < c2.status {
			return true
		}
		if c.status == c2.status {
			if c.int_count > c2.int_count {
				return true
			}
			if c.int_count == c2.int_count {
				if c.delta_age < c2.delta_age {
					return true
				}
				if c.delta_age == c2.delta_age {
					return c.id.Less(c2.id)
				}
			}
		}
	}

	return false
}

func SexBackGet(i uint8) string {
	if i&uint8(1) > 0 {
		return "m"
	} else {
		return "f"
	}
}
func (a *Account) SexBackGet() string {
	return SexBackGet(a.Sspn)
}

func recommendLoadData100City(fa *RecommendFQ) []compat {
	res_data := make([]compat, 0)

	iarr := fa.a.InterestsIDs.GetArr()
	liarr := uint8(len(iarr))

	if liarr == 0 {
		return res_data
	}

	sex := fa.a.SexBackGet()

	sspn := SexStatusPremIntCreate(sex, "свободны", true)

	for idp := index.MaxIdp2; idp >= 0; idp-- {

		key := InterestCitySortKey{Idp: idp, Sspn: sspn, InterestID: iarr[0], CityID: fa.city_id}
		s, ok := index.InterestCitySortKeyData[key]
		if ok {
			for _, id := range s {
				c := fa.a.CreateCompatFromID(id, 0, iarr, true)
				if c.int_count == liarr {
					res_data = append(res_data, c)
				}
			}
		}

		if idp == 0 {
			break
		}
	}

	return res_data
}

func recommendLoadData100(fa *RecommendFQ) []compat {
	res_data := make([]compat, 0)

	iarr := fa.a.InterestsIDs.GetArr()
	liarr := uint8(len(iarr))

	if liarr == 0 {
		return res_data
	}

	sex := fa.a.SexBackGet()

	sspn := SexStatusPremIntCreate(sex, "свободны", true)

	for idp := index.MaxIdp2; idp >= 0; idp-- {
		key := InterestSortKey{Idp: idp, Sspn: sspn, InterestID: iarr[0]}
		s, ok := index.InterestSortKeyData[key]
		if ok {
			for _, id := range s {
				c := fa.a.CreateCompatFromID(id, 0, iarr, true)
				if c.int_count == liarr {
					res_data = append(res_data, c)
				}
			}
		}

		if idp == 0 {
			break
		}
	}

	return res_data
}

func recommendLoadData100Country(fa *RecommendFQ) []compat {
	res_data := make([]compat, 0)

	iarr := fa.a.InterestsIDs.GetArr()
	liarr := uint8(len(iarr))

	if liarr == 0 {
		return res_data
	}

	sex := fa.a.SexBackGet()

	sspn := SexStatusPremIntCreate(sex, "свободны", true)

	for idp := index.MaxIdp2; idp >= 0; idp-- {
		key := InterestSortKey{Idp: idp, Sspn: sspn, InterestID: iarr[0]}
		s, ok := index.InterestSortKeyData[key]
		if ok {
			for _, id := range s {
				c, okc := fa.CreateCompatFromID(id, 0, iarr, true)
				if okc {
					if c.int_count == liarr {
						res_data = append(res_data, c)
					}
				}
			}
		}

		if idp == 0 {
			break
		}
	}

	return res_data
}

func StatusWeight(status string) uint8 {
	if status == "свободны" {
		return 0
	}
	if status == "всё сложно" {
		return 1
	}
	//"заняты"
	return 2

}

func recommendLoadDataCity(res_data []compat, fa *RecommendFQ, status string, premium_now bool) []compat {

	sw := StatusWeight(status)

	iarr := fa.a.InterestsIDs.GetArr()
	liarr := uint8(len(iarr))

	if liarr == 0 {
		return res_data
	}

	sex := fa.a.SexBackGet()

	sspn := SexStatusPremIntCreate(sex, status, premium_now)

	for idp := index.MaxIdp2; idp >= 0; idp-- {
		fa.check_wait()
		for _, ire := range iarr {
			key := InterestCitySortKey{Idp: idp, Sspn: sspn, InterestID: ire, CityID: fa.city_id}
			s, ok := index.InterestCitySortKeyData[key]
			if ok {
				for _, id := range s {
					// _, ex := dubl[id]
					// if !ex {
					//	dubl[id] = struct{}{}
					c := fa.a.CreateCompatFromID(id, sw, iarr, premium_now)
					res_data = RecombineResult(res_data, c)
					//res_data = append(res_data, c)
					//}
				}
			}
		}

		if idp == 0 {
			break
		}
	}

	return res_data
}

func recommendLoadDataOt(res_data []compat, fa *RecommendFQ, status string, premium_now bool) []compat {

	sw := StatusWeight(status)

	iarr := fa.a.InterestsIDs.GetArr()
	liarr := uint8(len(iarr))

	if liarr == 0 {
		return res_data
	}

	sex := fa.a.SexBackGet()

	sspn := SexStatusPremIntCreate(sex, status, premium_now)

	for idp := index.MaxIdp2; idp >= 0; idp-- {
		fa.check_wait()
		for _, ire := range iarr {
			key := InterestSortKey{Idp: idp, Sspn: sspn, InterestID: ire}
			s, ok := index.InterestSortKeyData[key]
			if ok {
				for _, id := range s {

					c := fa.a.CreateCompatFromID(id, sw, iarr, premium_now)

					res_data = RecombineResult(res_data, c)
				}
			}
		}

		if idp == 0 {
			break
		}
	}

	return res_data
}

func recommendLoadDataCountry(res_data []compat, fa *RecommendFQ, status string, premium_now bool) []compat {

	sw := StatusWeight(status)

	iarr := fa.a.InterestsIDs.GetArr()
	liarr := uint8(len(iarr))

	if liarr == 0 {
		return res_data
	}

	sex := fa.a.SexBackGet()

	sspn := SexStatusPremIntCreate(sex, status, premium_now)

	for idp := index.MaxIdp2; idp >= 0; idp-- {
		fa.check_wait()
		for _, ire := range iarr {
			key := InterestSortKey{Idp: idp, Sspn: sspn, InterestID: ire}
			s, ok := index.InterestSortKeyData[key]
			if ok {
				for _, id := range s {
					c, okc := fa.CreateCompatFromID(id, sw, iarr, premium_now)
					if okc {

						res_data = RecombineResult(res_data, c)

					}
				}
			}
		}

		if idp == 0 {
			break
		}
	}

	return res_data
}

func CreateMinForm(res_data []compat) []compat {

	sort.Slice(res_data, func(i int, j int) bool {

		return res_data[i].More(res_data[j])
	})

	res := make([]compat, 22, 22)

	for i, v := range res_data {
		if i < 22 {
			res[i] = v
		}
	}

	return res
}

func RecombineResult(res_data []compat, add compat) []compat {

	if res_data[21].More(add) {
		return res_data
	}

	for _, v := range res_data {
		if v.id == add.id {
			return res_data
		}
	}

	res_data[21] = add
	sort.Slice(res_data, func(i int, j int) bool {

		return res_data[i].More(res_data[j])
	})
	// pos := -1
	// for i, c := range res_data {
	// 	if c.id == add.id {
	// 		return res_data
	// 	}
	// 	if c.More(add) {
	// 		pos = i
	// 	} else {
	// 		break
	// 	}
	// }

	// if pos < 22 {
	// 	for i := 20; i <= -1; i-- {
	// 		if i > pos {
	// 			res_data[i+1] = res_data[i]
	// 		} else if i == pos {
	// 			res_data[i+1] = add
	// 		} else {
	// 			break
	// 		}

	// 	}
	// }
	return res_data
}

func RealElementsCount(res_data []compat) int {
	r := 0

	for _, c := range res_data {
		if c.id != common.Int24Zero {
			r++
		}
	}

	return r
}

func recommendLoadData(fa *RecommendFQ) []compat {
	res_data := make([]compat, 0)

	if fa.city != "" && fa.city_id == 0 {
		return res_data
	}
	if fa.country != "" && fa.country_id == 0 {
		return res_data
	}

	// 1) Выбираем по наименьшему интересу полу премиум налу и городу/стране/просто
	// у кого int96 выдаёт 100 % совпадения по интересам их пихаем в список
	// и если в нём меньше чем нужно, то
	// 1 Б) берем всех, сортируем по интересам (кол-во совпадений) исключая дубляжи далее меняем статусы, потом убираем премиум

	if fa.city != "" {
		res_data = recommendLoadData100City(fa)
		if fa.limit <= len(res_data) {
			return res_data
		}
		res_data = CreateMinForm(res_data)
		res_data = recommendLoadDataCity(res_data, fa, "свободны", true)
		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCity(res_data, fa, "всё сложно", true)
		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCity(res_data, fa, "заняты", true)

		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCity(res_data, fa, "свободны", false)
		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCity(res_data, fa, "всё сложно", false)
		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCity(res_data, fa, "заняты", false)

		return res_data
	}

	if fa.country != "" {
		res_data = recommendLoadData100Country(fa)

		if fa.limit <= len(res_data) {
			return res_data
		}
		res_data = CreateMinForm(res_data)
		res_data = recommendLoadDataCountry(res_data, fa, "свободны", true)
		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCountry(res_data, fa, "всё сложно", true)
		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCountry(res_data, fa, "заняты", true)

		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCountry(res_data, fa, "свободны", false)
		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCountry(res_data, fa, "всё сложно", false)
		if fa.limit <= RealElementsCount(res_data) {
			return res_data
		}
		res_data = recommendLoadDataCountry(res_data, fa, "заняты", false)

		return res_data
	}

	// Без определённого ограничения
	res_data = recommendLoadData100(fa)

	if fa.limit <= len(res_data) {
		return res_data
	}
	res_data = CreateMinForm(res_data)
	res_data = recommendLoadDataOt(res_data, fa, "свободны", true)
	if fa.limit <= RealElementsCount(res_data) {
		return res_data
	}
	res_data = recommendLoadDataOt(res_data, fa, "всё сложно", true)
	if fa.limit <= RealElementsCount(res_data) {
		return res_data
	}
	res_data = recommendLoadDataOt(res_data, fa, "заняты", true)

	if fa.limit <= RealElementsCount(res_data) {
		return res_data
	}
	res_data = recommendLoadDataOt(res_data, fa, "свободны", false)
	if fa.limit <= RealElementsCount(res_data) {
		return res_data
	}
	res_data = recommendLoadDataOt(res_data, fa, "всё сложно", false)
	if fa.limit <= RealElementsCount(res_data) {
		return res_data
	}
	res_data = recommendLoadDataOt(res_data, fa, "заняты", false)

	return res_data
}

func (fa *RecommendFQ) get_answer(res_data []compat) map[string][]map[string]interface{} {

	sort.Slice(res_data, func(i int, j int) bool {

		return res_data[i].More(res_data[j])
	})

	// for _, v := range res_data {
	// 	fmt.Println(v)
	// }

	res := make(map[string][]map[string]interface{})
	d := make([]map[string]interface{}, 0)
	res["accounts"] = d

	for i := range res_data {
		if i >= fa.limit {
			break
		}
		if res_data[i].id == common.Int24Zero {
			continue
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

		elem["birth"] = a.Birth

		if a.Premium.Start != 0 {
			prem := make(map[string]interface{})
			prem["start"] = a.Premium.Start
			prem["finish"] = a.Premium.Finish
			elem["premium"] = prem
		}

		res["accounts"] = append(res["accounts"], elem)
	}

	return res
}
