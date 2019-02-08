package main

import (
	"go_v1/common"
	"hash/crc32"
	"sort"
	"sync"
)

var (
	//crc32q = crc32.MakeTable(0xD5828281)
	crc32q = crc32.MakeTable(0xD5868281)
)

// func SliceUppendAndSort(s []common.Int24, val common.Int24) []common.Int24 {
// 	if s == nil {
// 		s = make([]common.Int24, 0, 4)
// 	}
// 	var r []common.Int24
// 	if len(s) == cap(s) {
// 		r = make([]common.Int24, 0, cap(s)+4)
// 		r = append(r, s...)
// 		r = append(r, val)
// 	} else {
// 		r = append(s, val)
// 	}
// 	sort.Slice(r, func(i, j int) bool {
// 		return r[i].More(r[j])
// 	})

// 	return r
// }

func SliceUppendAndSort(s []common.Int24, val common.Int24) []common.Int24 {
	if s == nil {
		s = make([]common.Int24, 0, 4)
	}

	for _, v := range s {
		if v.Equal(val) {
			return s
		}
	}

	var r []common.Int24
	if len(s) == cap(s) {
		r = make([]common.Int24, 0, cap(s)+4)
		r = append(r, s...)
		r = append(r, val)
	} else {
		r = append(s, val)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i].More(r[j])
	})

	return r
}

func SliceRemove(s []common.Int24, val common.Int24) []common.Int24 {
	r := -1
	for i, v := range s {
		if v.Equal(val) {
			r = i
			break
		}
	}

	if r >= 0 {
		s = append(s[:r], s[r+1:]...)
	}

	return s
}

// type LikeIndex struct {
// 	mux  sync.RWMutex
// 	data map[common.Int24][]common.Int24
// }

// func LikeIndexCreate() *LikeIndex {
// 	r := LikeIndex{}
// 	r.data = make(map[common.Int24][]common.Int24)
// 	return &r
// }

// func (li *LikeIndex) AppendOne(from, to common.Int24) {
// 	li.mux.Lock()
// 	s, _ := li.data[to]

// 	li.data[to] = SliceUppendAndSort(s, from)

// 	li.mux.Unlock()
// }
// func (li *LikeIndex) Append(a *Account) {
// 	// li.mux.Lock()

// 	// if a.Likes != nil {
// 	// 	for _, v := range a.Likes {
// 	// 		s, _ := li.data[v.ID]
// 	// 		li.data[v.ID] = SliceUppendAndSort(s, a.ID)
// 	// 	}
// 	// }
// 	// li.mux.Unlock()
// }
// func (li *LikeIndex) RemoveOne(from, to common.Int24) {
// 	li.mux.Lock()
// 	s, ok := li.data[to]
// 	if ok {
// 		li.data[to] = SliceRemove(s, from)
// 	}
// 	li.mux.Unlock()
// }
// func (li *LikeIndex) Remove(a *Account) {
// 	// li.mux.Lock()

// 	// if a.Likes != nil {
// 	// 	for _, v := range a.Likes {
// 	// 		s, ok := li.data[v.ID]
// 	// 		if ok {
// 	// 			li.data[v.ID] = SliceRemove(s, a.ID)
// 	// 		}
// 	// 	}
// 	// }
// 	// li.mux.Unlock()
// }

func IdpGet(id common.Int24) uint16 {
	return uint16(id.Int() / 256)
}

func Idp2Get(id common.Int24) uint8 {
	return uint8(id.Int() / 8192)
}

type CommonSortKey struct {
	Idp  uint16
	Sspn uint8
}

type InterestSortKey struct {
	Idp        uint8
	Sspn       uint8
	InterestID uint8
}

type CitySortKey struct {
	Idp    uint8
	Sspn   uint8
	CityID uint16
}
type FnameSortKey struct {
	Idp     uint8
	Sspn    uint8
	FNameID uint16
}
type CountrySortKey struct {
	Idp       uint8
	Sspn      uint8
	CountryID uint8
}
type DomainSortKey struct {
	Idp      uint8
	Sspn     uint8
	DomainID uint8
}
type PhoneCodeSortKey struct {
	Idp         uint8
	Sspn        uint8
	PhoneCodeID uint8
}
type BirthYearSortKey struct {
	Idp       uint8
	Sspn      uint8
	BirthYear uint16
}

type InterestCitySortKey struct {
	Idp        uint8
	Sspn       uint8
	InterestID uint8
	CityID     uint16
}
type CountryDomainSortKey struct {
	Idp       uint8
	Sspn      uint8
	DomainID  uint8
	CountryID uint8
}

func CommonSortKeyCreate(a *Account) CommonSortKey {
	return CommonSortKey{Idp: IdpGet(a.ID), Sspn: a.Sspn}
}

func InterestSortKeyCreate(a *Account) []InterestSortKey {
	res := make([]InterestSortKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, InterestSortKey{Idp: Idp2Get(a.ID), Sspn: a.Sspn, InterestID: inter})
	}

	return res
}
func CitySortKeyCreate(a *Account) CitySortKey {
	return CitySortKey{CityID: a.CityID, Idp: Idp2Get(a.ID), Sspn: a.Sspn}
}
func FnameSortKeyCreate(a *Account) FnameSortKey {
	return FnameSortKey{FNameID: a.FNameID, Idp: Idp2Get(a.ID), Sspn: a.Sspn}
}
func CountrySortKeyCreate(a *Account) CountrySortKey {
	return CountrySortKey{CountryID: a.CountryID, Idp: Idp2Get(a.ID), Sspn: a.Sspn}
}
func DomainSortKeyCreate(a *Account) DomainSortKey {
	return DomainSortKey{DomainID: a.DomainID, Idp: Idp2Get(a.ID), Sspn: a.Sspn}
}
func PhoneCodeSortKeyCreate(a *Account) PhoneCodeSortKey {
	return PhoneCodeSortKey{PhoneCodeID: a.PhoneCodeID, Idp: Idp2Get(a.ID), Sspn: a.Sspn}
}
func BirthYearSortKeyCreate(a *Account) BirthYearSortKey {
	return BirthYearSortKey{BirthYear: a.BirthYear, Idp: Idp2Get(a.ID), Sspn: a.Sspn}
}

func InterestCitySortKeyCreate(a *Account) []InterestCitySortKey {
	res := make([]InterestCitySortKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, InterestCitySortKey{Idp: Idp2Get(a.ID), Sspn: a.Sspn, InterestID: inter, CityID: a.CityID})
	}

	return res
}
func CountryDomainSortKeyCreate(a *Account) CountryDomainSortKey {
	return CountryDomainSortKey{DomainID: a.DomainID, CountryID: a.CountryID, Idp: Idp2Get(a.ID), Sspn: a.Sspn}
}

type GlobalIndex struct {
	mux     sync.RWMutex
	MaxIdp  uint16
	MaxIdp2 uint8

	CommonSortKeyData    map[CommonSortKey][]common.Int24
	InterestSortKeyData  map[InterestSortKey][]common.Int24
	CitySortKeyData      map[CitySortKey][]common.Int24
	FnameSortKeyData     map[FnameSortKey][]common.Int24
	CountrySortKeyData   map[CountrySortKey][]common.Int24
	DomainSortKeyData    map[DomainSortKey][]common.Int24
	PhoneCodeSortKeyData map[PhoneCodeSortKey][]common.Int24
	BirthYearSortKeyData map[BirthYearSortKey][]common.Int24

	InterestCitySortKeyData  map[InterestCitySortKey][]common.Int24
	CountryDomainSortKeyData map[CountryDomainSortKey][]common.Int24

	emails map[string]common.Int24
	phones map[string]common.Int24
}

func GlobalIndexCreate() *GlobalIndex {
	r := GlobalIndex{
		CommonSortKeyData:    make(map[CommonSortKey][]common.Int24),
		InterestSortKeyData:  make(map[InterestSortKey][]common.Int24),
		CitySortKeyData:      make(map[CitySortKey][]common.Int24),
		FnameSortKeyData:     make(map[FnameSortKey][]common.Int24),
		CountrySortKeyData:   make(map[CountrySortKey][]common.Int24),
		DomainSortKeyData:    make(map[DomainSortKey][]common.Int24),
		PhoneCodeSortKeyData: make(map[PhoneCodeSortKey][]common.Int24),
		BirthYearSortKeyData: make(map[BirthYearSortKey][]common.Int24),

		InterestCitySortKeyData:  make(map[InterestCitySortKey][]common.Int24),
		CountryDomainSortKeyData: make(map[CountryDomainSortKey][]common.Int24),

		emails: make(map[string]common.Int24),
		phones: make(map[string]common.Int24),
	}
	return &r
}

func (asi *GlobalIndex) exists_email(email string) bool {
	//hashe := crc32.Checksum([]byte(email+"salt"), crc32q)
	hashe := EmailStructCreate(email).HashStringGet()
	//hashe := crc32.Checksum([]byte(email+"salt"), crc32q)
	asi.mux.RLock()
	_, ok := asi.emails[hashe]
	asi.mux.RUnlock()
	return ok
}
func (asi *GlobalIndex) email_ext_get(email string) common.Int24 {
	hashe := EmailStructCreate(email).HashStringGet()
	//hashe := crc32.Checksum([]byte(email+"salt"), crc32q)
	asi.mux.RLock()
	v, ok := asi.emails[hashe]
	asi.mux.RUnlock()
	if !ok {
		return common.Int24Create(0)
	}
	return v
}

func (asi *GlobalIndex) exists_phone(phone string) bool {
	//hashp := crc32.Checksum([]byte(phone+"salt"), crc32q)
	hashp := phone
	asi.mux.RLock()
	_, ok := asi.phones[hashp]
	asi.mux.RUnlock()
	return ok
}
func (asi *GlobalIndex) phone_ext_get(phone string) common.Int24 {
	//hashp := crc32.Checksum([]byte(phone+"salt"), crc32q)
	hashp := phone
	asi.mux.RLock()
	v, ok := asi.phones[hashp]
	asi.mux.RUnlock()
	if !ok {
		return common.Int24Create(0)
	}
	return v
}

func (gi *GlobalIndex) Append(a *Account) {
	gi.mux.Lock()

	//hashe := crc32.Checksum([]byte(a.StringEmail()+"salt"), crc32q)
	hashe := a.EmailStructCreate().HashStringGet()
	gi.emails[hashe] = a.ID
	if a.Phone != "" {
		//hashp := crc32.Checksum([]byte(a.Phone+"salt"), crc32q)
		hashp := a.Phone
		gi.phones[hashp] = a.ID
	}

	{
		key := CommonSortKeyCreate(a)
		if gi.MaxIdp < key.Idp {
			gi.MaxIdp = key.Idp
		}
		s, _ := gi.CommonSortKeyData[key]
		gi.CommonSortKeyData[key] = SliceUppendAndSort(s, a.ID)
	}

	if Idp2Get(a.ID) > gi.MaxIdp2 {
		gi.MaxIdp2 = Idp2Get(a.ID)
	}

	{
		keys := InterestSortKeyCreate(a)
		for _, key := range keys {
			s, _ := gi.InterestSortKeyData[key]
			gi.InterestSortKeyData[key] = SliceUppendAndSort(s, a.ID)
		}
	}

	{
		key := CitySortKeyCreate(a)
		s, _ := gi.CitySortKeyData[key]
		gi.CitySortKeyData[key] = SliceUppendAndSort(s, a.ID)
	}
	{
		key := FnameSortKeyCreate(a)
		s, _ := gi.FnameSortKeyData[key]
		gi.FnameSortKeyData[key] = SliceUppendAndSort(s, a.ID)
	}

	{
		key := CountrySortKeyCreate(a)
		s, _ := gi.CountrySortKeyData[key]
		gi.CountrySortKeyData[key] = SliceUppendAndSort(s, a.ID)
	}

	{
		key := DomainSortKeyCreate(a)
		s, _ := gi.DomainSortKeyData[key]
		gi.DomainSortKeyData[key] = SliceUppendAndSort(s, a.ID)
	}
	{
		key := PhoneCodeSortKeyCreate(a)
		s, _ := gi.PhoneCodeSortKeyData[key]
		gi.PhoneCodeSortKeyData[key] = SliceUppendAndSort(s, a.ID)
	}

	{
		key := BirthYearSortKeyCreate(a)
		s, _ := gi.BirthYearSortKeyData[key]
		gi.BirthYearSortKeyData[key] = SliceUppendAndSort(s, a.ID)
	}

	{
		keys := InterestCitySortKeyCreate(a)
		for _, key := range keys {
			s, _ := gi.InterestCitySortKeyData[key]
			gi.InterestCitySortKeyData[key] = SliceUppendAndSort(s, a.ID)
		}
	}
	{
		key := CountryDomainSortKeyCreate(a)
		s, _ := gi.CountryDomainSortKeyData[key]
		gi.CountryDomainSortKeyData[key] = SliceUppendAndSort(s, a.ID)
	}

	gi.mux.Unlock()
}

func (gi *GlobalIndex) Remove(a *Account) {
	gi.mux.Lock()

	//hashe := crc32.Checksum([]byte(a.StringEmail()+"salt"), crc32q)
	hashe := a.EmailStructCreate().HashStringGet()
	delete(gi.emails, hashe)
	if a.Phone != "" {
		//hashp := crc32.Checksum([]byte(a.Phone+"salt"), crc32q)
		hashp := a.Phone
		delete(gi.phones, hashp)
	}

	{
		key := CommonSortKeyCreate(a)
		if gi.MaxIdp < key.Idp {
			gi.MaxIdp = key.Idp
		}
		s, _ := gi.CommonSortKeyData[key]
		gi.CommonSortKeyData[key] = SliceRemove(s, a.ID)
	}

	if Idp2Get(a.ID) > gi.MaxIdp2 {
		gi.MaxIdp2 = Idp2Get(a.ID)
	}

	{
		keys := InterestSortKeyCreate(a)
		for _, key := range keys {
			s, _ := gi.InterestSortKeyData[key]
			gi.InterestSortKeyData[key] = SliceRemove(s, a.ID)
		}
	}

	{
		key := CitySortKeyCreate(a)
		s, _ := gi.CitySortKeyData[key]
		gi.CitySortKeyData[key] = SliceRemove(s, a.ID)
	}
	{
		key := FnameSortKeyCreate(a)
		s, _ := gi.FnameSortKeyData[key]
		gi.FnameSortKeyData[key] = SliceRemove(s, a.ID)
	}

	{
		key := CountrySortKeyCreate(a)
		s, _ := gi.CountrySortKeyData[key]
		gi.CountrySortKeyData[key] = SliceRemove(s, a.ID)
	}

	{
		key := DomainSortKeyCreate(a)
		s, _ := gi.DomainSortKeyData[key]
		gi.DomainSortKeyData[key] = SliceRemove(s, a.ID)
	}
	{
		key := PhoneCodeSortKeyCreate(a)
		s, _ := gi.PhoneCodeSortKeyData[key]
		gi.PhoneCodeSortKeyData[key] = SliceRemove(s, a.ID)
	}

	{
		key := BirthYearSortKeyCreate(a)
		s, _ := gi.BirthYearSortKeyData[key]
		gi.BirthYearSortKeyData[key] = SliceRemove(s, a.ID)
	}

	{
		keys := InterestCitySortKeyCreate(a)
		for _, key := range keys {
			s, _ := gi.InterestCitySortKeyData[key]
			gi.InterestCitySortKeyData[key] = SliceRemove(s, a.ID)
		}
	}
	{
		key := CountryDomainSortKeyCreate(a)
		s, _ := gi.CountryDomainSortKeyData[key]
		gi.CountryDomainSortKeyData[key] = SliceRemove(s, a.ID)
	}

	gi.mux.Unlock()
}
