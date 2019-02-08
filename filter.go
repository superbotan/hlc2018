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

type FilterFQ struct {
	stime time.Time
	wtime time.Time

	limit int

	sex_eq string

	email_domain string
	domain_id    uint8
	email_lt     string
	email_gt     string

	status_eq  string
	status_neq string

	fname_eq    string
	fname_eq_id uint16
	fname_any   string
	fname_any_m map[uint16]bool
	fname_null  string

	sname_eq       string
	sname_eq_id    uint16
	sname_starts   string
	sname_starts_m map[uint16]bool
	sname_null     string

	phone_code  string
	phone_eq_id uint8
	phone_null  string

	country_eq    string
	country_eq_id uint8
	country_null  string

	city_eq    string
	city_eq_id uint16
	city_any   string
	city_any_m map[uint16]bool
	city_null  string

	birth_lt   int32
	birth_gt   int32
	birth_year uint16

	interests_contains   string
	interests_contains_m common.Int96
	interests_any        string
	interests_any_m      common.Int96

	likes_contains   string
	likes_contains_m map[common.Int24]bool

	premium_now  string
	premium_null string

	out_sex     bool
	out_status  bool
	out_fname   bool
	out_sname   bool
	out_phone   bool
	out_country bool
	out_city    bool
	out_birth   bool
	out_premium bool

	ex_domain   bool
	ex_fname    bool
	ex_sname    bool
	ex_phone    bool
	ex_country  bool
	ex_city     bool
	ex_birth    bool
	ex_interest bool

	p_one bool

	raw_query string
}

func accountsFilterGet(ctx *fasthttp.RequestCtx) (fa *FilterFQ, err error) {

	fa = &FilterFQ{}

	fa.stime = time.Now()
	fa.wtime = fa.stime

	fa.raw_query = ctx.URI().String()

	fa.birth_gt = -900000000
	fa.birth_lt = -900000000

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

		if k == "sex_eq" && len(v) > 0 {
			fa.sex_eq = v
			fa.out_sex = true
			return
		}

		if k == "email_domain" && len(v) > 0 {
			fa.email_domain = v
			fa.domain_id = domains.GetByName(fa.email_domain)
			fa.ex_domain = true
			return
		}
		if k == "email_lt" && len(v) > 0 {
			fa.email_lt = v
			return
		}
		if k == "email_gt" && len(v) > 0 {
			fa.email_gt = v
			return
		}

		if k == "status_eq" && len(v) > 0 {
			fa.status_eq = v
			fa.out_status = true
			return
		}
		if k == "status_neq" && len(v) > 0 {
			fa.status_neq = v
			fa.out_status = true
			return
		}

		if k == "fname_eq" && len(v) > 0 {
			fa.fname_eq = v
			fa.fname_eq_id = fnames.GetByName(fa.fname_eq)
			fa.out_fname = true
			fa.ex_fname = true
			return
		}
		if k == "fname_any" && len(v) > 0 {
			fa.fname_any = v
			fa.fname_any_m = fnames.GetList(v)
			fa.out_fname = true
			fa.ex_fname = true
			return
		}
		if k == "fname_null" && len(v) > 0 {
			fa.fname_null = v
			fa.out_fname = true
			fa.ex_fname = fa.fname_null == "1"
			return
		}

		if k == "sname_eq" && len(v) > 0 {
			fa.sname_eq = v
			fa.sname_eq_id = snames.GetByName(fa.sname_eq)
			fa.out_sname = true
			fa.ex_sname = true
			return
		}
		if k == "sname_starts" && len(v) > 0 {
			fa.sname_starts = v
			fa.sname_starts_m = snames.GetLikeValue(v)
			fa.out_sname = true
			return
		}
		if k == "sname_null" && len(v) > 0 {
			fa.sname_null = v
			fa.out_sname = true
			fa.ex_sname = fa.sname_null == "1"
			return
		}

		if k == "phone_code" && len(v) > 0 {
			fa.phone_code = v
			fa.phone_eq_id = phonecodes.GetByName(fa.phone_code)
			fa.out_phone = true
			fa.ex_phone = true
			return
		}
		if k == "phone_null" && len(v) > 0 {
			fa.phone_null = v
			fa.out_phone = true
			fa.ex_phone = fa.phone_null == "1"
			return
		}

		if k == "country_eq" && len(v) > 0 {
			fa.country_eq = v
			fa.country_eq_id = countries.GetByName(fa.country_eq)
			fa.out_country = true
			fa.ex_country = true
			return
		}
		if k == "country_null" && len(v) > 0 {
			fa.country_null = v
			fa.out_country = true
			fa.ex_country = fa.country_null == "1"
			return
		}

		if k == "city_eq" && len(v) > 0 {
			fa.city_eq = v
			fa.city_eq_id = cities.GetByName(fa.city_eq)
			fa.out_city = true
			fa.ex_city = true
			return
		}
		if k == "city_any" && len(v) > 0 {
			fa.city_any = v
			fa.city_any_m = cities.GetList(v)
			fa.out_city = true
			fa.ex_city = true
			return
		}
		if k == "city_null" && len(v) > 0 {
			fa.city_null = v
			fa.out_city = true
			fa.ex_city = fa.city_null == "1"
			return
		}

		if k == "birth_lt" && len(v) > 0 {
			l, err := strconv.Atoi(v)
			if err != nil {
				fail = true
				return
			}
			fa.birth_lt = int32(l)
			fa.out_birth = true
			return
		}
		if k == "birth_gt" && len(v) > 0 {
			l, err := strconv.Atoi(v)
			if err != nil {
				fail = true
				return
			}
			fa.birth_gt = int32(l)
			fa.out_birth = true
			return
		}
		if k == "birth_year" && len(v) > 0 {
			l, err := strconv.Atoi(v)
			if err != nil {
				fail = true
				return
			}
			fa.birth_year = uint16(l)
			fa.out_birth = true
			fa.ex_birth = true
			return
		}

		if k == "interests_contains" && len(v) > 0 {
			fa.interests_contains = v
			fa.interests_contains_m = fa.interests_contains_m.SetArr(interests.GetListArray(v))
			fa.ex_interest = true
			return
		}
		if k == "interests_any" && len(v) > 0 {
			fa.interests_any = v
			fa.interests_any_m = fa.interests_any_m.SetArr(interests.GetListArray(v))
			fa.ex_interest = true
			return
		}

		if k == "likes_contains" && len(v) > 0 {
			fa.likes_contains = v
			fa.likes_contains_m = make(map[common.Int24]bool)
			for _, v := range strings.Split(v, ",") {
				l, err := strconv.Atoi(v)
				if err != nil {
					fail = true
					return
				}
				fa.likes_contains_m[common.Int24Create(int32(l))] = true
			}
			return
		}

		if k == "premium_now" && len(v) > 0 {
			fa.premium_now = v
			fa.out_premium = true
			return
		}
		if k == "premium_null" && len(v) > 0 {
			fa.premium_null = v
			fa.out_premium = true
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
func (fa *FilterFQ) check_wait() bool {
	if time.Since(fa.wtime) > 500*time.Microsecond {
		time.Sleep(10 * time.Nanosecond)
		fa.wtime = time.Now()
	}
	if !fa.p_one && time.Since(fa.stime) > 10*time.Millisecond {
		//fmt.Println("flt", fa.raw_query)
		fa.p_one = true
	}
	return false
}

func (fa *FilterFQ) checkSspn(sspn uint8) bool {
	return (fa.sex_eq == "" || fa.sex_eq == SexGet(sspn)) &&
		(fa.status_eq == "" || fa.status_eq == StatusGet(sspn)) &&
		(fa.status_neq == "" || fa.status_neq != StatusGet(sspn)) &&
		(fa.premium_now == "" || PremiumNowGet(sspn))
}

func preLoadData(fa *FilterFQ) []common.Int24 {

	if fa.limit <= 0 {
		return make([]common.Int24, 0)
	}

	res := make([]common.Int24, 0, fa.limit)

	if fa.likes_contains != "" {
		for k, _ := range fa.likes_contains_m {
			fls := make([][]common.Int24, 0, 1)
			fls = append(fls, accounts.Get(k).LikesBack)

			res = fa.filterCheck(fls, res)
			return res
		}
	}

	// if fa.email_lt != "" || fa.email_gt != "" {
	// 	return res
	// }
	// if fa.birth_lt != -900000000 || fa.birth_gt != -900000000 {
	// 	return res
	// }
	// if fa.sname_starts != "" {
	// 	return res
	// }

	if fa.interests_any != "" && (fa.city_eq != "" || fa.city_null == "1") {
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for _, inter := range fa.interests_any_m.GetArr() {
				for i := uint8(0); i <= 12; i++ {
					if fa.checkSspn(i) {
						key := InterestCitySortKey{Idp: idp, Sspn: i, InterestID: inter, CityID: fa.city_eq_id}
						s, ok := index.InterestCitySortKeyData[key]
						if ok {
							fls = append(fls, s)
						}
					}
				}
			}
			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}
	if fa.interests_contains != "" && (fa.city_eq != "" || fa.city_null == "1") {
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			inter := fa.interests_contains_m.GetArr()[0]
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := InterestCitySortKey{Idp: idp, Sspn: i, InterestID: inter, CityID: fa.city_eq_id}
					s, ok := index.InterestCitySortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)

			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.interests_any != "" && (fa.city_any != "") {
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for _, inter := range fa.interests_any_m.GetArr() {
				for city := range fa.city_any_m {
					for i := uint8(0); i <= 12; i++ {
						if fa.checkSspn(i) {
							key := InterestCitySortKey{Idp: idp, Sspn: i, InterestID: inter, CityID: city}
							s, ok := index.InterestCitySortKeyData[key]
							if ok {
								fls = append(fls, s)
							}
						}
					}
				}
			}
			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}
	if fa.interests_contains != "" && (fa.city_any != "") {
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			inter := fa.interests_contains_m.GetArr()[0]
			for city := range fa.city_any_m {
				for i := uint8(0); i <= 12; i++ {
					if fa.checkSspn(i) {
						key := InterestCitySortKey{Idp: idp, Sspn: i, InterestID: inter, CityID: city}
						s, ok := index.InterestCitySortKeyData[key]
						if ok {
							fls = append(fls, s)
						}
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.interests_any != "" {
		// fa.ex_country ||
		// !(fa.ex_city || fa.ex_birth || fa.ex_fname || fa.ex_sname || fa.ex_domain || fa.ex_phone) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for _, inter := range fa.interests_any_m.GetArr() {
				for i := uint8(0); i <= 12; i++ {
					if fa.checkSspn(i) {
						key := InterestSortKey{Idp: idp, Sspn: i, InterestID: inter}
						s, ok := index.InterestSortKeyData[key]
						if ok {
							fls = append(fls, s)
						}
					}
				}
			}
			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}
	if fa.interests_contains != "" {
		// !(fa.ex_city || fa.ex_birth || fa.ex_country || fa.ex_fname || fa.ex_sname || fa.ex_domain || fa.ex_phone) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			inter := fa.interests_contains_m.GetArr()[0]
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := InterestSortKey{Idp: idp, Sspn: i, InterestID: inter}
					s, ok := index.InterestSortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.city_eq != "" || fa.city_null == "1" {
		// !(fa.ex_birth || fa.ex_country || fa.ex_fname || fa.ex_sname || fa.ex_domain || fa.ex_phone || fa.ex_interest) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := CitySortKey{Idp: idp, Sspn: i, CityID: fa.city_eq_id}
					s, ok := index.CitySortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.city_any != "" {
		//!(fa.ex_birth || fa.ex_country || fa.ex_fname || fa.ex_sname || fa.ex_domain || fa.ex_phone || fa.ex_interest) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for city := range fa.city_any_m {
				for i := uint8(0); i <= 12; i++ {
					if fa.checkSspn(i) {
						key := CitySortKey{Idp: idp, Sspn: i, CityID: city}
						s, ok := index.CitySortKeyData[key]
						if ok {
							fls = append(fls, s)
						}
					}
				}
			}
			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.fname_eq != "" || fa.fname_null == "1" {
		//!(fa.ex_birth || fa.ex_country || fa.ex_city || fa.ex_sname || fa.ex_domain || fa.ex_phone || fa.ex_interest) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := FnameSortKey{Idp: idp, Sspn: i, FNameID: fa.fname_eq_id}
					s, ok := index.FnameSortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.fname_any != "" {
		//!(fa.ex_birth || fa.ex_country || fa.ex_city || fa.ex_sname || fa.ex_domain || fa.ex_phone || fa.ex_interest) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for fname := range fa.fname_any_m {
				for i := uint8(0); i <= 12; i++ {
					if fa.checkSspn(i) {
						key := FnameSortKey{Idp: idp, Sspn: i, FNameID: fname}
						s, ok := index.FnameSortKeyData[key]
						if ok {
							fls = append(fls, s)
						}
					}
				}
			}
			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.domain_id != 0 && (fa.country_eq != "" || fa.country_null == "1") {
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := CountryDomainSortKey{Idp: idp, Sspn: i, CountryID: fa.country_eq_id, DomainID: fa.domain_id}
					s, ok := index.CountryDomainSortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.country_eq != "" || fa.country_null == "1" {
		//!(fa.ex_city || fa.ex_birth || fa.ex_fname || fa.ex_sname || fa.ex_phone || fa.ex_domain || fa.ex_interest) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := CountrySortKey{Idp: idp, Sspn: i, CountryID: fa.country_eq_id}
					s, ok := index.CountrySortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.birth_year != 0 {
		// Возможные тормоза
		// fa.ex_country ||
		// fa.ex_phone ||
		// fa.ex_domain ||
		// !(fa.ex_city || fa.ex_fname || fa.ex_sname || fa.ex_interest) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := BirthYearSortKey{Idp: idp, Sspn: i, BirthYear: fa.birth_year}
					s, ok := index.BirthYearSortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.phone_code != "" || fa.phone_null == "1" {
		// Возможные тормоза
		// fa.ex_birth ||
		// !(fa.ex_city || fa.ex_fname || fa.ex_sname || fa.ex_country || fa.ex_domain || fa.ex_interest) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := PhoneCodeSortKey{Idp: idp, Sspn: i, PhoneCodeID: fa.phone_eq_id}
					s, ok := index.PhoneCodeSortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if fa.email_domain != "" {
		// !(fa.ex_city || fa.ex_birth || fa.ex_country || fa.ex_fname || fa.ex_sname || fa.ex_phone || fa.ex_interest) &&
		for idp := index.MaxIdp2; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			fls := make([][]common.Int24, 0, 0)
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) {
					key := DomainSortKey{Idp: idp, Sspn: i, DomainID: fa.domain_id}
					s, ok := index.DomainSortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}

			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	if true {
		// final mode
		//if !(fa.ex_city || fa.ex_birth || fa.ex_country || fa.ex_fname || fa.ex_sname || fa.ex_domain || fa.ex_phone || fa.ex_interest) {
		for idp := index.MaxIdp; idp >= 0 && fa.limit > len(res); idp-- {
			if fa.check_wait() {
				return res
			}
			//fmt.Println("fc", idp)
			//fmt.Println(len(index.CommonSortKeyData))

			//fmt.Println(fa)

			fls := make([][]common.Int24, 0, 0)
			for i := uint8(0); i <= 12; i++ {
				if fa.checkSspn(i) || true {
					key := CommonSortKey{Idp: idp, Sspn: i}
					s, ok := index.CommonSortKeyData[key]
					if ok {
						fls = append(fls, s)
					}
				}
			}
			//res = append(res, fa.filterCheck(fls, len(res))...)
			res = fa.filterCheck(fls, res)
			if idp == 0 {
				break
			}
		}

		return res
	}

	return res
}

func (fa *FilterFQ) filterCheck(fls [][]common.Int24, inpq []common.Int24) []common.Int24 {

	if fa.limit-len(inpq) <= 0 {
		return inpq
	}

	l := len(fls)

	if l == 0 {
		return inpq
	}
	if l == 1 {
		for i := 0; i < len((fls[0])); i++ {
			if filterOne(fls[0][i], fa) {
				inpq = append(inpq, fls[0][i])
			}
			if len(inpq) >= (fa.limit) {
				break
			}
		}
		return inpq
	}

	ixs := make([]int, l)

	current_j := -1
	c_val := common.Int24Create(0)

	for {
		for j := 0; j < l; j++ {
			if ixs[j] != -1 && ixs[j] < len(fls[j]) {
				if current_j == -1 {
					current_j = j
					c_val = fls[j][ixs[j]]
				} else if c_val.Less(fls[j][ixs[j]]) {
					current_j = j
					c_val = fls[j][ixs[j]]
				}

			} else {
				ixs[j] = -1
			}
		}

		if current_j == -1 {
			break
		}

		if filterOne(c_val, fa) {
			if len(inpq) == 0 || inpq[len(inpq)-1] != c_val {
				inpq = append(inpq, c_val)

				if len(inpq) >= fa.limit {
					break
				}
			}
		}

		ixs[current_j] = ixs[current_j] + 1
		current_j = -1
		c_val = common.Int24Create(0)
	}

	return inpq
}

func filterMakeResult(ids []common.Int24, fa *FilterFQ) map[string][]map[string]interface{} {

	res := make(map[string][]map[string]interface{})

	res["accounts"] = make([]map[string]interface{}, 0)

	sort.Slice(ids, func(i, j int) bool {
		return ids[i].More(ids[j])
	})

	for i := range ids {
		if i >= fa.limit {
			break
		}

		a := accounts.Get(ids[i])

		elem := make(map[string]interface{})

		elem["id"] = a.ID.Int()
		elem["email"] = a.StringEmail()

		if fa.out_sex {
			elem["sex"] = a.SexGet()
		}
		if fa.out_status {
			elem["status"] = a.StatusGet()
		}
		if fa.out_fname {
			if a.FNameID != 0 {
				elem["fname"] = fnames.GetByID(a.FNameID)
			}
		}
		if fa.out_sname {
			if a.SNameID != 0 {
				elem["sname"] = snames.GetByID(a.SNameID)
			}
		}
		if fa.out_phone {
			if a.Phone != "" {
				elem["phone"] = a.Phone
			}
		}
		if fa.out_country {
			if a.CountryID != 0 {
				elem["country"] = countries.GetByID(a.CountryID)
			}
		}
		if fa.out_city {
			if a.CityID != 0 {
				elem["city"] = cities.GetByID(a.CityID)
			}
		}
		if fa.out_birth {
			elem["birth"] = a.Birth
		}
		if fa.out_premium {
			if a.Premium.Start != 0 {
				prem := make(map[string]interface{})
				prem["start"] = a.Premium.Start
				prem["finish"] = a.Premium.Finish
				elem["premium"] = prem
			}
		}

		res["accounts"] = append(res["accounts"], elem)
	}

	return res
}

func filterOne(i common.Int24, fa *FilterFQ) bool {

	a := accounts.Get(i)

	// if i.Int() == 29997 {
	// 	fmt.Println(a)
	// }

	if fa.sex_eq != "" && a.SexGet() != fa.sex_eq {
		return false
	}

	if fa.email_domain != "" && a.DomainID != fa.domain_id {
		return false
	}
	if fa.email_lt != "" && a.FirstEmail >= fa.email_lt {
		return false
	}
	if fa.email_gt != "" && a.FirstEmail < fa.email_gt {
		return false
	}

	if fa.status_eq != "" && a.StatusGet() != fa.status_eq {
		return false
	}
	if fa.status_neq != "" && a.StatusGet() == fa.status_neq {
		return false
	}

	if fa.fname_eq != "" && (a.FNameID == 0 || a.FNameID != fa.fname_eq_id) {
		return false
	}
	if fa.fname_any != "" { // any
		_, ok := fa.fname_any_m[a.FNameID]
		if !ok {
			return false
		}
	}
	if fa.fname_null == "1" && a.FNameID != 0 || fa.fname_null == "0" && a.FNameID == 0 {
		return false
	}

	if fa.sname_eq != "" && (a.SNameID == 0 || a.SNameID != fa.sname_eq_id) {
		return false
	}
	if fa.sname_starts != "" { // starts
		_, ok := fa.sname_starts_m[a.SNameID]
		if !ok {
			return false
		}
	}
	if fa.sname_null == "1" && a.SNameID != 0 || fa.sname_null == "0" && a.SNameID == 0 {
		return false
	}

	if fa.phone_code != "" && (a.PhoneCodeID == 0 || a.PhoneCodeID != fa.phone_eq_id) {
		return false
	}
	if fa.phone_null == "1" && a.PhoneCodeID != 0 || fa.phone_null == "0" && a.PhoneCodeID == 0 {
		return false
	}

	if fa.country_eq != "" && (a.CountryID == 0 || a.CountryID != fa.country_eq_id) {
		return false
	}
	if fa.country_null == "1" && a.CountryID != 0 || fa.country_null == "0" && a.CountryID == 0 {
		return false
	}

	if fa.city_eq != "" && (a.CityID == 0 || a.CityID != fa.city_eq_id) {
		return false
	}
	if fa.city_any != "" { // any
		_, ok := fa.city_any_m[a.CityID]
		if !ok {
			return false
		}
	}
	if fa.city_null == "1" && a.CityID != 0 || fa.city_null == "0" && a.CityID == 0 {
		return false
	}

	if fa.birth_lt != -900000000 && a.Birth > fa.birth_lt {
		return false
	}
	if fa.birth_gt != -900000000 && a.Birth < fa.birth_gt {
		return false
	}
	if fa.birth_year != 0 && a.BirthYear != fa.birth_year {
		return false
	}

	if fa.interests_contains != "" { // interests_contains
		if !fa.interests_contains_m.AllIn(a.InterestsIDs) {

			return false
		}
		// fmt.Print(fa.interests_contains_m, "\t\t")
		// fmt.Println(a.InterestsIDs)
	}
	if fa.interests_any != "" { // interests_any
		if !fa.interests_any_m.Contains(a.InterestsIDs) {
			// fmt.Print(fa.interests_any_m, "\t\t")
			// fmt.Println(a.InterestsIDs)
			return false
		}
	}

	if fa.likes_contains != "" {

		m := make(map[common.Int24]bool)

		for _, u := range a.Likes {
			m[u.Id] = true
		}
		for u := range fa.likes_contains_m {
			_, ok := m[u]
			if !ok {
				return false
			}
		}
	}

	if fa.premium_now != "" && (!a.PremiumNowGet()) {
		return false
	}

	if fa.premium_null == "1" && a.HasPremium() || fa.premium_null == "0" && !a.HasPremium() {
		return false
	}

	return true
}
