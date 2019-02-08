package main

import (
	"errors"
	"go_v1/common"
	"sort"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

type GroupFQ struct {
	limit int
	order int

	group_by_sex       bool
	group_by_status    bool
	group_by_interests bool
	group_by_country   bool
	group_by_city      bool

	sort []string

	sex         string
	status      string
	country     string
	city        string
	interests   string
	interests_m common.Int96
	country_id  uint8
	city_id     uint16
	birth       uint16
	likes       int32
	joined      uint16

	ex_interests bool
	ex_country   bool
	ex_city      bool
	ex_birth     bool
	ex_joined    bool
}

func groupFilterGet(ctx *fasthttp.RequestCtx) (fa *GroupFQ, err error) {

	fa = &GroupFQ{}

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
			return
		}
		if k == "order" && len(v) > 0 {
			l, err := strconv.Atoi(v)
			if err != nil {
				fail = true
				return
			}
			fa.order = l
			return
		}

		if k == "keys" && len(v) > 0 {
			fa.sort = make([]string, 0)
			for _, v1 := range strings.Split(v, ",") {
				fa.sort = append(fa.sort, v1)
				if v1 == "sex" {
					fa.group_by_sex = true
					continue
				}
				if v1 == "status" {
					fa.group_by_status = true
					continue
				}
				if v1 == "interests" {
					fa.group_by_interests = true
					fa.ex_interests = true
					continue
				}
				if v1 == "country" {
					fa.group_by_country = true
					fa.ex_country = true
					continue
				}
				if v1 == "city" {
					fa.group_by_city = true
					fa.ex_city = true
					continue
				}
				fail = true
			}

			if len(v) == 0 {
				fail = true
			}
			return
		}

		if k == "sex" && len(v) > 0 {
			fa.sex = v
			return
		}
		if k == "status" && len(v) > 0 {
			fa.status = v
			return
		}

		if k == "interests" && len(v) > 0 {
			fa.interests = v
			fa.interests_m = fa.interests_m.SetArr(interests.GetListArray(v))
			fa.ex_interests = true
			return
		}

		if k == "country" && len(v) > 0 {
			fa.country = v
			fa.country_id = countries.GetByName(v)
			fa.ex_country = true
			return
		}

		if k == "city" && len(v) > 0 {
			fa.city = v
			fa.city_id = cities.GetByName(v)
			fa.ex_city = true
			return
		}
		if k == "birth" && len(v) > 0 {
			l, err := strconv.Atoi(v)
			if err != nil {
				fail = true
				return
			}
			fa.birth = uint16(l)
			fa.ex_birth = true
			return
		}
		if k == "likes" && len(v) > 0 {
			l, err := strconv.Atoi(v)
			if err != nil {
				fail = true
				return
			}
			fa.likes = int32(l)
			return
		}
		if k == "joined" && len(v) > 0 {
			l, err := strconv.Atoi(v)
			if err != nil {
				fail = true
				return
			}
			fa.joined = uint16(l)
			fa.ex_joined = true
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

func groupLoadData(fa *GroupFQ) map[GroupKey]uint32 {
	res := make(map[GroupKey]uint32)

	index := make(map[GroupKey]uint16)
	index32 := make(map[GroupKey]uint32)

	if fa.likes != 0 {
		index_l := make(map[GroupKeyGlobal]uint32)
		if fa.interests != "" || fa.group_by_interests {
			for _, lks := range accounts.Get(common.Int24Create(fa.likes)).LikesBack {
				a := accounts.Get(lks)
				for _, key := range AllWInterestsGKey(a) {
					v, ok := index_l[key]
					if ok {
						index_l[key] = v + 1
					} else {
						index_l[key] = 1
					}
				}
			}
		} else {

			for _, lks := range accounts.Get(common.Int24Create(fa.likes)).LikesBack {

				a := accounts.Get(lks)
				key := AllWOInterestsGKey(a)
				v, ok := index_l[key]
				if ok {
					index_l[key] = v + 1
				} else {
					index_l[key] = 1
				}
			}

		}

		for k, v := range index_l {
			if v > 0 && fa.is_ok_gk(k) {
				key := fa.groupindex_gk(k)
				bv, ok := res[key]
				if ok {
					res[key] = bv + v
				} else {
					res[key] = v
				}

			}
		}

		return res
	} else {
		if !(fa.ex_interests || fa.ex_country || fa.ex_city || fa.ex_birth || fa.ex_joined) {
			index32 = groupindex.SexSatusOnly
		} else if fa.interests != "" && (len(fa.interests_m.GetArr()) == 0 || fa.interests_m.GetArr()[0] == 0) {
			// nothing
		} else if fa.city != "" && (fa.city_id == 0) {
			// nothing
		} else if fa.country != "" && (fa.country_id == 0) {
			// nothing
		} else if !(fa.ex_country || fa.ex_city || fa.ex_birth || fa.ex_joined) {
			if fa.interests == "" {
				index32 = groupindex.Interests
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.Interests[key]
					if ok {
						index32[key] = v
					}
				}
			}
		} else if !(fa.ex_interests || fa.ex_city || fa.ex_birth || fa.ex_joined) {
			if fa.country == "" {
				index32 = groupindex.Country
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, Country: fa.country_id}

					v, ok := groupindex.Country[key]
					if ok {
						index32[key] = v
					}
				}
			}
		} else if !(fa.ex_interests || fa.ex_country || fa.ex_birth || fa.ex_joined) {
			if fa.city == "" {
				index32 = groupindex.City
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, City: fa.city_id}

					v, ok := groupindex.City[key]
					if ok {
						index32[key] = v
					}
				}
			}
		} else if !(fa.ex_interests || fa.ex_country || fa.ex_city || fa.ex_joined) {
			{
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss}

					v, ok := groupindex.Birth[fa.birth][key]
					if ok {
						index32[key] = v
					}
				}
			}
		} else if !(fa.ex_interests || fa.ex_country || fa.ex_city || fa.ex_birth) {
			{
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss}

					v, ok := groupindex.Joined[fa.joined][key]
					if ok {
						index32[key] = v
					}
				}
			}
		} else if !(fa.ex_interests || fa.ex_city || fa.ex_joined) {
			if fa.country == "" {
				index = groupindex.BirthCountry[fa.birth]
			} else {
				for ss := uint8(0); ss < 12; ss++ {

					key := GroupKey{SexStatus: ss, Country: fa.country_id}

					v, ok := groupindex.BirthCountry[fa.birth][key]
					if ok {
						index[key] = v

					}
				}
			}
		} else if !(fa.ex_interests || fa.ex_country || fa.ex_joined) {
			if fa.city == "" {
				// ~ 1-2 ms
				index = groupindex.BirthCity[fa.birth]
			} else {
				for ss := uint8(0); ss < 12; ss++ {

					key := GroupKey{SexStatus: ss, City: fa.city_id}

					v, ok := groupindex.BirthCity[fa.birth][key]
					if ok {
						index[key] = v

					}
				}
			}
		} else if !(fa.ex_interests || fa.ex_city || fa.ex_birth) {
			if fa.country == "" {
				index = groupindex.JoinedCountry[fa.joined]
			} else {
				for ss := uint8(0); ss < 12; ss++ {

					key := GroupKey{SexStatus: ss, Country: fa.country_id}

					v, ok := groupindex.JoinedCountry[fa.joined][key]
					if ok {
						index[key] = v

					}
				}
			}
		} else if !(fa.ex_interests || fa.ex_country || fa.ex_birth) {
			if fa.city == "" {
				// ~ 1-2 ms
				index = groupindex.JoinedCity[fa.joined]
			} else {
				for ss := uint8(0); ss < 12; ss++ {

					key := GroupKey{SexStatus: ss, City: fa.city_id}

					v, ok := groupindex.JoinedCity[fa.joined][key]
					if ok {
						index[key] = v

					}
				}
			}
		} else if !(fa.ex_country || fa.ex_city || fa.ex_birth) {
			if fa.interests == "" {
				index32 = groupindex.JoinedInterests[fa.joined]
			} else {
				for ss := uint8(0); ss < 12; ss++ {

					key := GroupKey{SexStatus: ss, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.JoinedInterests[fa.joined][key]
					if ok {
						index32[key] = v

					}
				}
			}
		} else if !(fa.ex_country || fa.ex_city || fa.ex_joined) {
			if fa.interests == "" {
				index32 = groupindex.BirthInterests[fa.birth]
			} else {
				for ss := uint8(0); ss < 12; ss++ {

					key := GroupKey{SexStatus: ss, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.BirthInterests[fa.birth][key]
					if ok {
						index32[key] = v

					}
				}
			}
		} else if !(fa.ex_city || fa.ex_birth || fa.ex_joined) {
			if fa.country == "" && fa.interests == "" {
				// !++++!+++!+!!!
			} else if fa.country == "" && fa.interests != "" {
				index = groupindex.InterestsCountry[fa.interests_m.GetArr()[0]]
			} else if fa.country != "" && fa.interests == "" {
				index = groupindex.CountryInterests[fa.country_id]
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, Country: fa.country_id, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.InterestsCountry[fa.interests_m.GetArr()[0]][key]
					if ok {
						index[key] = v
					}
				}
			}
		} else if !(fa.ex_country || fa.ex_birth || fa.ex_joined) {
			if fa.city == "" && fa.interests == "" {
				// !++++!+++!+!!!
			} else if fa.city == "" && fa.interests != "" {
				index = groupindex.InterestsCity[fa.interests_m.GetArr()[0]]
			} else if fa.city != "" && fa.interests == "" {
				index = groupindex.CityInterests[fa.city_id]
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, City: fa.city_id, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.InterestsCity[fa.interests_m.GetArr()[0]][key]
					if ok {
						index[key] = v
					}
				}
			}
		} else if !(fa.ex_city || fa.ex_joined) {

			if fa.country == "" && fa.interests == "" {
				// !++++!+++!+!!!
			} else if fa.country == "" && fa.interests != "" {
				index = groupindex.BirthInterestsCountry[fa.birth][fa.interests_m.GetArr()[0]]
			} else if fa.country != "" && fa.interests == "" {
				index = groupindex.BirthCountryInterests[fa.birth][fa.country_id]
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, Country: fa.country_id, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.BirthInterestsCountry[fa.birth][fa.interests_m.GetArr()[0]][key]
					if ok {
						index[key] = v
					}
				}
			}
		} else if !(fa.ex_country || fa.ex_joined) {
			if fa.city == "" && fa.interests == "" {
				// !++++!+++!+!!!
			} else if fa.city == "" && fa.interests != "" {
				index = groupindex.BirthInterestsCity[fa.birth][fa.interests_m.GetArr()[0]]
			} else if fa.city != "" && fa.interests == "" {
				index = groupindex.BirthCityInterests[fa.birth][fa.city_id]
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, City: fa.city_id, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.BirthInterestsCity[fa.birth][fa.interests_m.GetArr()[0]][key]
					if ok {
						index[key] = v
					}
				}
			}
		} else if !(fa.ex_city || fa.ex_birth) {

			if fa.country == "" && fa.interests == "" {
				// !++++!+++!+!!!
			} else if fa.country == "" && fa.interests != "" {
				index = groupindex.JoinedInterestsCountry[fa.joined][fa.interests_m.GetArr()[0]]
			} else if fa.country != "" && fa.interests == "" {
				index = groupindex.JoinedCountryInterests[fa.joined][fa.country_id]
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, Country: fa.country_id, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.JoinedInterestsCountry[fa.joined][fa.interests_m.GetArr()[0]][key]
					if ok {
						index[key] = v
					}
				}
			}
		} else if !(fa.ex_country || fa.ex_birth) {
			if fa.city == "" && fa.interests == "" {
				// !++++!+++!+!!!
			} else if fa.city == "" && fa.interests != "" {
				index = groupindex.JoinedInterestsCity[fa.joined][fa.interests_m.GetArr()[0]]
			} else if fa.city != "" && fa.interests == "" {
				index = groupindex.JoinedCityInterests[fa.joined][fa.city_id]
			} else {
				for ss := uint8(0); ss < 12; ss++ {
					key := GroupKey{SexStatus: ss, City: fa.city_id, Interest: fa.interests_m.GetArr()[0]}

					v, ok := groupindex.JoinedInterestsCity[fa.joined][fa.interests_m.GetArr()[0]][key]
					if ok {
						index[key] = v
					}
				}
			}
		}
	}

	//fmt.Println(fa)

	for k, v := range index {
		if v > 0 && fa.is_ok(k) {
			key := fa.groupindex(k)
			bv, ok := res[key]
			if ok {
				res[key] = bv + uint32(v)
			} else {
				res[key] = uint32(v)
			}

		}
	}

	for k, v := range index32 {
		if v > 0 && fa.is_ok(k) {
			key := fa.groupindex(k)
			bv, ok := res[key]
			if ok {
				res[key] = bv + v
			} else {
				res[key] = v
			}

		}
	}

	return res
}

func (fa *GroupFQ) groupindex(k GroupKey) GroupKey {
	res := GroupKey{}

	sex := ""
	status := ""

	if fa.group_by_sex {
		sex = GSexGet(k.SexStatus)
	}
	if fa.group_by_status {
		status = GStatusGet(k.SexStatus)
	}

	res.SexStatus = SexStatusGroupCreate(sex, status)

	if fa.group_by_interests && fa.interests == "" {
		res.Interest = k.Interest
	}

	if fa.group_by_city && fa.city == "" {
		res.City = k.City
	}

	if fa.group_by_country && fa.country == "" {
		res.Country = k.Country
	}

	return res
}

type GroupKeyGlobal struct {
	SexStatus uint8
	Interest  uint8
	Country   uint8
	City      uint16
	Birth     uint16
	Joined    uint16
}

func AllWInterestsGKey(a *Account) []GroupKeyGlobal {
	res := make([]GroupKeyGlobal, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKeyGlobal{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
			Country:   a.CountryID,
			City:      a.CityID,
			Birth:     a.BirthYear,
			Joined:    a.JoinedYear,
		})
	}
	return res
}
func AllWOInterestsGKey(a *Account) GroupKeyGlobal {
	return GroupKeyGlobal{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
		Country:   a.CountryID,
		City:      a.CityID,
		Birth:     a.BirthYear,
		Joined:    a.JoinedYear,
	}
}

func (fa *GroupFQ) groupindex_gk(k GroupKeyGlobal) GroupKey {
	res := GroupKey{}

	sex := ""
	status := ""

	if fa.group_by_sex {
		sex = GSexGet(k.SexStatus)
	}
	if fa.group_by_status {
		status = GStatusGet(k.SexStatus)
	}

	res.SexStatus = SexStatusGroupCreate(sex, status)

	if fa.group_by_interests && fa.interests == "" {
		res.Interest = k.Interest
	}

	if fa.group_by_city && fa.city == "" {
		res.City = k.City
	}

	if fa.group_by_country && fa.country == "" {
		res.Country = k.Country
	}

	return res
}

func (fa *GroupFQ) is_ok_gk(k GroupKeyGlobal) bool {
	if fa.birth != 0 && k.Birth != fa.birth {
		//fmt.Println(a.birth)
		return false
	}
	if fa.sex != "" && (fa.sex != GSexGet(k.SexStatus)) {
		//fmt.Println(a.sex)
		return false
	}
	if fa.joined != 0 && k.Joined != fa.joined {
		//fmt.Println(a.joined)
		return false
	}
	if fa.status != "" && (fa.status != GStatusGet(k.SexStatus)) {
		//fmt.Println(a.status)
		return false
	}

	if fa.country != "" && fa.country_id != k.Country {
		//fmt.Println(a.country)
		return false
	}
	if fa.city != "" && fa.city_id != k.City {
		//fmt.Println(a.city)
		return false
	}

	if fa.interests != "" && fa.interests_m.GetArr()[0] != k.Interest {
		//fmt.Println(a.city)
		return false
	}

	return true
}

func (fa *GroupFQ) is_ok(k GroupKey) bool {
	// if fa.birth != 0 && k.Birth != fa.birth {
	// 	//fmt.Println(a.birth)
	// 	return false
	// }
	if fa.sex != "" && (fa.sex != GSexGet(k.SexStatus)) {
		//fmt.Println(a.sex)
		return false
	}
	// if fa.joined != 0 && k.Joined != fa.joined {
	// 	//fmt.Println(a.joined)
	// 	return false
	// }
	if fa.status != "" && (fa.status != GStatusGet(k.SexStatus)) {
		//fmt.Println(a.status)
		return false
	}

	if fa.country != "" && fa.country_id != k.Country {
		//fmt.Println(a.country)
		return false
	}
	if fa.city != "" && fa.city_id != k.City {
		//fmt.Println(a.city)
		return false
	}

	if fa.interests != "" && fa.interests_m.GetArr()[0] != k.Interest {
		//fmt.Println(a.city)
		return false
	}

	return true
}

type gr_item struct {
	k GroupKey
	c uint32
}

func (fa *GroupFQ) more(i gr_item, j gr_item) bool {
	if fa.order == -1 && i.c > j.c || fa.order == 1 && i.c < j.c {
		return true
	} else if i.c == j.c && fa.sort != nil {
		for _, s := range fa.sort {
			if s == "sex" {
				si := GSexGet(i.k.SexStatus)
				sj := GSexGet(j.k.SexStatus)
				if fa.order == -1 && si > sj || fa.order == 1 && si < sj {
					return true
				} else if si == sj {
					continue
				}

			}
			if s == "status" {
				si := GStatusGet(i.k.SexStatus)
				sj := GStatusGet(j.k.SexStatus)
				if fa.order == -1 && si > sj || fa.order == 1 && si < sj {
					return true
				} else if si == sj {
					continue
				}

			}
			if s == "interests" {
				si := interests.GetByID(i.k.Interest)
				sj := interests.GetByID(j.k.Interest)
				if fa.order == -1 && si > sj || fa.order == 1 && si < sj {
					return true
				} else if si == sj {
					continue
				}

			}
			if s == "country" {
				si := countries.GetByID(i.k.Country)
				sj := countries.GetByID(j.k.Country)
				if fa.order == -1 && si > sj || fa.order == 1 && si < sj {
					return true
				} else if si == sj {
					continue
				}

			}
			if s == "city" {
				si := cities.GetByID(i.k.City)
				sj := cities.GetByID(j.k.City)
				if fa.order == -1 && si > sj || fa.order == 1 && si < sj {
					return true
				} else if si == sj {
					continue
				}

			}
			return false
		}
	}

	return false
}

func (fa *GroupFQ) get_answer(res_data map[GroupKey]uint32) map[string][]map[string]interface{} {
	l := make([]gr_item, 0)

	for k, c := range res_data {
		l = append(l, gr_item{k: k, c: c})
	}

	sort.Slice(l, func(i int, j int) bool {
		//return fa.order == -1 && l[i].c > l[j].c || fa.order == 1 && l[i].c < l[j].c
		return fa.more(l[i], l[j])
	})

	res := make(map[string][]map[string]interface{})
	d := make([]map[string]interface{}, 0)
	res["groups"] = d

	for i := 0; i < len(l) && i < fa.limit; i++ {
		k := l[i].k
		c := l[i].c

		elem := make(map[string]interface{})

		elem["count"] = c
		if fa.group_by_status {
			elem["status"] = GStatusGet(k.SexStatus)
		}
		if fa.group_by_sex {
			elem["sex"] = GSexGet(k.SexStatus)
		}

		if fa.group_by_city && fa.city == "" && k.City != 0 {
			elem["city"] = cities.GetByID(k.City)
		}

		if fa.group_by_country && fa.country == "" && k.Country != 0 {
			elem["country"] = countries.GetByID(k.Country)
		}

		if fa.group_by_interests && fa.interests == "" && k.Interest != 0 {
			elem["interests"] = interests.GetByID(k.Interest)
		}

		res["groups"] = append(res["groups"], elem)
	}

	return res
}
