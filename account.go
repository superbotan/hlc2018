package main

import (
	"strings"
	"sync"
	"time"
)

import (
	"go_v1/common"
	"strconv"
)

var i_li_map int
var i_li int

type AccountStorage struct {
	mx   sync.RWMutex
	data [1700000]*Account
}

func AccountStorageCreate() *AccountStorage {
	as := AccountStorage{}

	return &as
}

func (as *AccountStorage) Append(a *Account) {
	if a.ID.More(common.Int24Create(1700000 - 1)) {
		return
	}
	as.mx.Lock()
	as.data[a.ID.Int()] = a
	as.mx.Unlock()
}

func (as *AccountStorage) Get(id common.Int24) *Account {
	if id.More(common.Int24Create(1700000 - 1)) {
		return nil
	}
	a := as.data[id.Int()]

	return a
}

func (as *AccountStorage) GetLock(id common.Int24) *Account {
	if id.More(common.Int24Create(1700000 - 1)) {
		return nil
	}
	as.mx.RLock()
	as.mx.RUnlock()
	a := as.data[id.Int()]

	return a
}

func SexStatusPremIntCreate(sex string, status string, is_premium_now bool) uint8 {
	r := uint8(0)

	if sex == "f" {
		r = r + 1
	}
	if is_premium_now {
		r = r + 2
	}
	if status == "заняты" {
		r = r + 4
	}
	if status == "всё сложно" {
		r = r + 8
	}

	return r
}

func SexGet(i uint8) string {
	if i&uint8(1) > 0 {
		return "f"
	} else {
		return "m"
	}
}
func PremiumNowGet(i uint8) bool {
	return i&uint8(2) > 0
}
func StatusGet(i uint8) string {
	if i&uint8(4) > 0 {
		return "заняты"
	} else if i&uint8(8) > 0 {
		return "всё сложно"
	} else {
		return "свободны"
	}
}

func (a *Account) SexGet() string {
	return SexGet(a.Sspn)
}
func (a *Account) PremiumNowGet() bool {
	return PremiumNowGet(a.Sspn)
}
func (a *Account) StatusGet() string {
	return StatusGet(a.Sspn)
}

func (a *Account) HasPremium() bool {
	return a.Premium.Start > 0
}

type LikeElement struct {
	Ts int32
	Id common.Int24
}

func (a *Account) MailString() string {
	return a.FirstEmail + "@" + domains.GetByID(a.DomainID)
}

type Account struct {
	ID common.Int24

	//Sspn Sex Status PN
	Sspn     uint8
	DomainID uint8

	CountryID   uint8
	PhoneCodeID uint8

	InterestsIDs common.Int96

	CityID uint16

	FNameID uint16
	SNameID uint16

	BirthYear  uint16
	JoinedYear uint16
	Birth      int32
	Joined     int32

	Premium struct {
		Start  int32 // нижней границей 01.01.2018
		Finish int32
	}

	FirstEmail string
	Phone      string

	Likes        []LikeElement
	LikesCorrect map[common.Int24]uint16
	LikesBack    []common.Int24
}

func (a *Account) AppendLike(ts int32, id common.Int24) {
	i_li++

	corr, ok := a.LikesCorrect[id]
	pos := -1
	step := -1

	for i, v := range a.Likes { // Двоичный поиск бы сюда
		if v.Id.More(id) {
			step = i
		}
		if v.Id == id {
			pos = i
			break
		}
	}

	if pos == -1 {
		s := make([]LikeElement, 0, len(a.Likes)+1)
		if step >= 0 {
			s = append(s, (a.Likes)[0:step+1]...)
			s = append(s, LikeElement{Id: id, Ts: ts})
			s = append(s, (a.Likes)[step+1:]...)
		} else {
			s = append(s, LikeElement{Id: id, Ts: ts})
			s = append(s, (a.Likes)[0:]...)
		}
		a.Likes = s
	} else if ok {
		tsn := (int64((a.Likes)[pos].Ts)*int64(corr) + int64(ts)) / int64(corr+1)
		(a.Likes)[pos] = LikeElement{Id: id, Ts: int32(tsn)}
		a.LikesCorrect[id] = corr + 1
		i_li_map++
	} else {
		tsn := (int64((a.Likes)[pos].Ts) + int64(ts)) / 2
		(a.Likes)[pos] = LikeElement{Id: id, Ts: int32(tsn)}
		a.LikesCorrect[id] = 2
		i_li_map++
	}
}

func (a *Account) AppendLikeBack(id common.Int24) {
	pos := -1
	step := -1

	for i, v := range a.LikesBack { // Двоичный поиск бы сюда
		if v.More(id) {
			step = i
		}
		if v == id {
			pos = i
			break
		}
	}

	if pos == -1 {
		s := make([]common.Int24, 0, len(a.LikesBack)+1)
		if step >= 0 {
			s = append(s, (a.LikesBack)[0:step+1]...)
			s = append(s, id)
			s = append(s, (a.LikesBack)[step+1:]...)
		} else {
			s = append(s, id)
			s = append(s, (a.LikesBack)[0:]...)
		}
		a.LikesBack = s
	}
}

func (a *Account) RemoveLikeBack(id common.Int24) {
	pos := -1

	for i, v := range a.LikesBack { // Двоичный поиск бы сюда
		if v == id {
			pos = i
			break
		}
	}

	if pos == -1 {
		a.LikesBack = append(a.LikesBack[:pos], a.LikesBack[pos:]...)
	}
}

func (a *Account) FillOtherLikeBack() {
	for _, v := range a.Likes {
		ains := accounts.Get(v.Id)

		ains.AppendLikeBack(a.ID)
	}
}
func (a *Account) RemoveOtherLikeBack() {
	for _, v := range a.Likes {
		ains := accounts.Get(v.Id)

		ains.RemoveLikeBack(a.ID)
	}
}

type EmailStruct struct {
	FirstEmail string
	DomainID   uint8
}

func EmailStructCreate(email string) EmailStruct {
	res := EmailStruct{}
	domainstart := strings.Index(email, "@") + 1
	r0 := []rune(email)
	domain := string(r0[domainstart:])
	res.DomainID = domains.Append(domain)
	res.FirstEmail = string(r0[:domainstart-1])

	return res
}
func (es EmailStruct) HashStringGet() string {
	return es.FirstEmail + "@" + strconv.Itoa(int(es.DomainID))

}
func (a *Account) EmailStructCreate() EmailStruct {
	res := EmailStruct{}

	res.DomainID = a.DomainID
	res.FirstEmail = a.FirstEmail

	return res
}

func (a *Account) StringEmail() string {

	return a.FirstEmail + "@" + domains.GetByID(a.DomainID)
}

func AccountCreate(aa AccountAdd) Account {

	a := Account{}

	a.Likes = make([]LikeElement, 0)
	a.LikesBack = make([]common.Int24, 0)
	a.LikesCorrect = make(map[common.Int24]uint16)

	a.ID = common.Int24Create(aa.ID)

	a.Premium.Finish = aa.Premium.Finish
	a.Premium.Start = aa.Premium.Start

	isPremiumNow := a.Premium.Start <= GetTimeNow() && a.Premium.Finish >= GetTimeNow()

	a.Sspn = SexStatusPremIntCreate(aa.Sex, aa.Status, isPremiumNow)

	a.Birth = aa.Birth
	bdate := time.Unix(int64(a.Birth), 0)
	a.BirthYear = uint16(bdate.Year())

	a.Joined = aa.Joined
	jdate := time.Unix(int64(a.Joined), 0)
	a.JoinedYear = uint16(jdate.Year())

	if aa.Phone != "" {
		runes := []rune(aa.Phone)
		phonecodestart := strings.Index(aa.Phone, "(") + 1
		phonecode := string(runes[phonecodestart : phonecodestart+3])
		a.Phone = aa.Phone
		a.PhoneCodeID = phonecodes.Append(phonecode)
	}

	domainstart := strings.Index(aa.Email, "@") + 1
	r0 := []rune(aa.Email)
	domain := string(r0[domainstart:])
	a.DomainID = domains.Append(domain)
	a.FirstEmail = string(r0[:domainstart-1])

	a.CountryID = countries.Append(aa.Country)
	a.CityID = cities.Append(aa.City)

	a.FNameID = fnames.Append(aa.Fname)
	a.SNameID = snames.Append(aa.Sname)

	if aa.Interests != nil {
		for _, inter := range aa.Interests {
			a.InterestsIDs = a.InterestsIDs.Set(interests.Append(inter))
		}
	}

	if aa.Likes != nil {
		//a.Likes = make([]LikeItem, 0, len(aa.Likes)+4)
		for _, like := range aa.Likes {
			(&a).AppendLike(like.Dt, common.Int24Create(like.ID))
			// 	//a.Likes = append(a.Likes, LikeItem{Dt: like.Dt, ID: common.Int24Create(like.ID)})
			// 	//likeindex2.Append(a.ID, like.Dt, common.Int24Create(like.ID))
		}
	}

	return a
}

func (aold *Account) AccountUpdate(au *AccountUpd) Account {

	a := Account{}

	a.LikesBack = aold.LikesBack

	a.ID = aold.ID

	if au.Premium.Finish != 0 {
		a.Premium.Finish = au.Premium.Finish
		a.Premium.Start = au.Premium.Start
	} else {
		a.Premium.Finish = aold.Premium.Finish
		a.Premium.Start = aold.Premium.Start
	}

	isPremiumNow := a.Premium.Start <= GetTimeNow() && a.Premium.Finish >= GetTimeNow()

	sex := au.Sex
	if sex == "" {
		sex = aold.SexGet()
	}
	status := au.Status
	if status == "" {
		status = aold.StatusGet()
	}

	a.Sspn = SexStatusPremIntCreate(sex, status, isPremiumNow)

	if au.Birth != nil {
		a.Birth = *au.Birth
		bdate := time.Unix(int64(a.Birth), 0)
		a.BirthYear = uint16(bdate.Year())
	} else {
		a.Birth = aold.Birth
		a.BirthYear = aold.BirthYear
	}

	if au.Joined != 0 {
		a.Joined = au.Joined
		jdate := time.Unix(int64(a.Joined), 0)
		a.JoinedYear = uint16(jdate.Year())
	} else {
		a.Joined = aold.Joined
		a.JoinedYear = aold.JoinedYear
	}

	if au.Phone != "" {
		runes := []rune(au.Phone)
		phonecodestart := strings.Index(au.Phone, "(") + 1
		phonecode := string(runes[phonecodestart : phonecodestart+3])
		a.Phone = au.Phone
		a.PhoneCodeID = phonecodes.Append(phonecode)
	} else {
		a.Phone = aold.Phone
		a.PhoneCodeID = aold.PhoneCodeID
	}

	if au.Email != "" {
		domainstart := strings.Index(au.Email, "@") + 1
		r0 := []rune(au.Email)
		domain := string(r0[domainstart:])
		a.DomainID = domains.Append(domain)
		a.FirstEmail = string(r0[:domainstart-1])
	} else {
		a.FirstEmail = aold.FirstEmail
		a.DomainID = aold.DomainID
	}

	if au.Country != "" {
		a.CountryID = countries.Append(au.Country)
	} else {
		a.CountryID = aold.CountryID
	}
	if au.City != "" {
		a.CityID = cities.Append(au.City)
	} else {
		a.CityID = aold.CityID
	}

	if au.Fname != "" {
		a.FNameID = fnames.Append(au.Fname)
	} else {
		a.FNameID = aold.FNameID
	}
	if au.Sname != "" {
		a.SNameID = snames.Append(au.Sname)
	} else {
		a.SNameID = aold.SNameID
	}

	if au.Interests != nil {
		for _, inter := range au.Interests {
			a.InterestsIDs = a.InterestsIDs.Set(interests.Append(inter))
		}
	} else {
		a.InterestsIDs = aold.InterestsIDs
	}

	if au.Likes != nil {
		a.Likes = make([]LikeElement, 0)
		a.LikesCorrect = make(map[common.Int24]uint16)

		for _, like := range au.Likes {
			(&a).AppendLike(like.Dt, common.Int24Create(like.ID))
		}
	} else {
		a.Likes = aold.Likes
		a.LikesCorrect = aold.LikesCorrect
	}

	return a
}
