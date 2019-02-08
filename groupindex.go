package main

import "sync"

func SexStatusGroupCreate(sex string, status string) uint8 {

	if sex == "" && status == "" {
		return 0
	}
	if sex == "f" && status == "" {
		return 1
	}
	if sex == "m" && status == "" {
		return 2
	}

	if sex == "" && status == "свободны" {
		return 3
	}
	if sex == "f" && status == "свободны" {
		return 4
	}
	if sex == "m" && status == "свободны" {
		return 5
	}

	if sex == "" && status == "заняты" {
		return 6
	}
	if sex == "f" && status == "заняты" {
		return 7
	}
	if sex == "m" && status == "заняты" {
		return 8
	}

	if sex == "" && status == "всё сложно" {
		return 9
	}
	if sex == "f" && status == "всё сложно" {
		return 10
	}
	if sex == "m" && status == "всё сложно" {
		return 11
	}

	return 0
}

func GSexGet(i uint8) string {
	t := i % 3
	if t == 0 {
		return ""
	}
	if t == 1 {
		return "f"
	}
	if t == 2 {
		return "m"
	}
	return ""
}
func GStatusGet(i uint8) string {

	t := i / 3
	if t == 0 {
		return ""
	}
	if t == 1 {
		return "свободны"
	}
	if t == 2 {
		return "заняты"
	}
	if t == 3 {
		return "всё сложно"
	}
	return ""
}

type GroupKey struct {
	SexStatus uint8
	Interest  uint8
	Country   uint8
	City      uint16
}

func AllWInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
			Country:   a.CountryID,
			City:      a.CityID,
		})
	}
	return res
}
func BirthCityInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
			City:      a.CityID,
		})
	}
	return res
}
func BirthCountryInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
			Country:   a.CountryID,
		})
	}
	return res
}
func JoinedCityInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
			City:      a.CityID,
		})
	}
	return res
}
func JoinedCountryInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
			Country:   a.CountryID,
		})
	}
	return res
}
func CityInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
			City:      a.CityID,
		})
	}
	return res
}
func CountryInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
			Country:   a.CountryID,
		})
	}
	return res
}
func BirthInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
		})
	}
	return res
}
func JoinedInterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
		})
	}
	return res
}

func AllWOInterestsKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
		Country:   a.CountryID,
		City:      a.CityID,
	}
}
func JoinedCityKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
		City:      a.CityID,
	}
}
func JoinedCountryKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
		Country:   a.CountryID,
	}
}
func BirthCityKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
		City:      a.CityID,
	}
}
func BirthCountryKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
		Country:   a.CountryID,
	}
}

func CountryKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
		Country:   a.CountryID,
	}
}
func CityKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
		City:      a.CityID,
	}
}
func BirthKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
	}
}
func JoinedKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
	}
}
func InterestsKey(a *Account) []GroupKey {
	res := make([]GroupKey, 0)
	for _, inter := range a.InterestsIDs.GetArr() {
		res = append(res, GroupKey{
			SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
			Interest:  inter,
		})
	}
	return res
}

func SexSatusOnlyKey(a *Account) GroupKey {
	return GroupKey{
		SexStatus: SexStatusGroupCreate(a.SexGet(), a.StatusGet()),
	}
}

type GroupIndex struct {
	mux sync.RWMutex

	//AllWInterests map[GroupKey]uint16

	CityInterests    map[uint16]map[GroupKey]uint16
	InterestsCity    map[uint8]map[GroupKey]uint16
	CountryInterests map[uint8]map[GroupKey]uint16
	InterestsCountry map[uint8]map[GroupKey]uint16

	BirthCityInterests    map[uint16]map[uint16]map[GroupKey]uint16
	BirthInterestsCity    map[uint16]map[uint8]map[GroupKey]uint16
	BirthCountryInterests map[uint16]map[uint8]map[GroupKey]uint16
	BirthInterestsCountry map[uint16]map[uint8]map[GroupKey]uint16

	JoinedCityInterests    map[uint16]map[uint16]map[GroupKey]uint16
	JoinedInterestsCity    map[uint16]map[uint8]map[GroupKey]uint16
	JoinedCountryInterests map[uint16]map[uint8]map[GroupKey]uint16
	JoinedInterestsCountry map[uint16]map[uint8]map[GroupKey]uint16

	BirthInterests  map[uint16]map[GroupKey]uint32
	JoinedInterests map[uint16]map[GroupKey]uint32

	//AllWOInterests map[GroupKey]uint16

	JoinedCity    map[uint16]map[GroupKey]uint16
	JoinedCountry map[uint16]map[GroupKey]uint16
	BirthCity     map[uint16]map[GroupKey]uint16
	BirthCountry  map[uint16]map[GroupKey]uint16

	City      map[GroupKey]uint32
	Country   map[GroupKey]uint32
	Birth     map[uint16]map[GroupKey]uint32
	Joined    map[uint16]map[GroupKey]uint32
	Interests map[GroupKey]uint32

	SexSatusOnly map[GroupKey]uint32
}

func GroupIndexCreate() *GroupIndex {
	gi := GroupIndex{

		//AllWInterests: make(map[GroupKey]uint16),

		CityInterests:    make(map[uint16]map[GroupKey]uint16),
		InterestsCity:    make(map[uint8]map[GroupKey]uint16),
		CountryInterests: make(map[uint8]map[GroupKey]uint16),
		InterestsCountry: make(map[uint8]map[GroupKey]uint16),

		BirthCityInterests:    make(map[uint16]map[uint16]map[GroupKey]uint16),
		BirthInterestsCity:    make(map[uint16]map[uint8]map[GroupKey]uint16),
		BirthCountryInterests: make(map[uint16]map[uint8]map[GroupKey]uint16),
		BirthInterestsCountry: make(map[uint16]map[uint8]map[GroupKey]uint16),

		JoinedCityInterests:    make(map[uint16]map[uint16]map[GroupKey]uint16),
		JoinedInterestsCity:    make(map[uint16]map[uint8]map[GroupKey]uint16),
		JoinedCountryInterests: make(map[uint16]map[uint8]map[GroupKey]uint16),
		JoinedInterestsCountry: make(map[uint16]map[uint8]map[GroupKey]uint16),

		BirthInterests:  make(map[uint16]map[GroupKey]uint32),
		JoinedInterests: make(map[uint16]map[GroupKey]uint32),

		//AllWOInterests: make(map[GroupKey]uint16),

		JoinedCity:    make(map[uint16]map[GroupKey]uint16),
		JoinedCountry: make(map[uint16]map[GroupKey]uint16),
		BirthCity:     make(map[uint16]map[GroupKey]uint16),
		BirthCountry:  make(map[uint16]map[GroupKey]uint16),

		City:      make(map[GroupKey]uint32),
		Country:   make(map[GroupKey]uint32),
		Birth:     make(map[uint16]map[GroupKey]uint32),
		Joined:    make(map[uint16]map[GroupKey]uint32),
		Interests: make(map[GroupKey]uint32),

		SexSatusOnly: make(map[GroupKey]uint32),
	}

	return &gi
}

func (gi *GroupIndex) Append(a *Account) {
	{
		// for _, key := range AllWInterestsKey(a) {
		// 	v, ok := gi.AllWInterests[key]
		// 	if ok {
		// 		gi.AllWInterests[key] = v + 1
		// 	} else {
		// 		gi.AllWInterests[key] = 1
		// 	}
		// }
	}
	{ // Birth
		{
			{
				_, okm := gi.BirthInterestsCity[a.BirthYear]
				if !okm {
					gi.BirthCityInterests[a.BirthYear] = make(map[uint16]map[GroupKey]uint16)
					gi.BirthInterestsCity[a.BirthYear] = make(map[uint8]map[GroupKey]uint16)
					gi.BirthCountryInterests[a.BirthYear] = make(map[uint8]map[GroupKey]uint16)
					gi.BirthInterestsCountry[a.BirthYear] = make(map[uint8]map[GroupKey]uint16)
				}
			}
			for _, key := range BirthCityInterestsKey(a) {
				{
					_, okb := gi.BirthInterestsCity[a.BirthYear][key.Interest]
					if !okb {
						gi.BirthInterestsCity[a.BirthYear][key.Interest] = make(map[GroupKey]uint16)
					}
					v, ok := gi.BirthInterestsCity[a.BirthYear][key.Interest][key]
					if ok {
						gi.BirthInterestsCity[a.BirthYear][key.Interest][key] = v + 1
					} else {
						gi.BirthInterestsCity[a.BirthYear][key.Interest][key] = 1
					}
				}
				{
					_, okb := gi.BirthCityInterests[a.BirthYear][key.City]
					if !okb {
						gi.BirthCityInterests[a.BirthYear][key.City] = make(map[GroupKey]uint16)
					}

					v, ok := gi.BirthCityInterests[a.BirthYear][key.City][key]
					if ok {
						gi.BirthCityInterests[a.BirthYear][key.City][key] = v + 1
					} else {
						gi.BirthCityInterests[a.BirthYear][key.City][key] = 1
					}
				}
			}
		}
		{
			for _, key := range BirthCountryInterestsKey(a) {
				{
					_, okb := gi.BirthInterestsCountry[a.BirthYear][key.Interest]
					if !okb {
						gi.BirthInterestsCountry[a.BirthYear][key.Interest] = make(map[GroupKey]uint16)
					}
					v, ok := gi.BirthInterestsCountry[a.BirthYear][key.Interest][key]
					if ok {
						gi.BirthInterestsCountry[a.BirthYear][key.Interest][key] = v + 1
					} else {
						gi.BirthInterestsCountry[a.BirthYear][key.Interest][key] = 1
					}
				}
				{
					_, okb := gi.BirthCountryInterests[a.BirthYear][key.Country]
					if !okb {
						gi.BirthCountryInterests[a.BirthYear][key.Country] = make(map[GroupKey]uint16)
					}

					v, ok := gi.BirthCountryInterests[a.BirthYear][key.Country][key]
					if ok {
						gi.BirthCountryInterests[a.BirthYear][key.Country][key] = v + 1
					} else {
						gi.BirthCountryInterests[a.BirthYear][key.Country][key] = 1
					}
				}
			}
		}
	}
	{ // Joined
		{
			{
				_, okm := gi.JoinedInterestsCity[a.JoinedYear]
				if !okm {
					gi.JoinedCityInterests[a.JoinedYear] = make(map[uint16]map[GroupKey]uint16)
					gi.JoinedInterestsCity[a.JoinedYear] = make(map[uint8]map[GroupKey]uint16)
					gi.JoinedCountryInterests[a.JoinedYear] = make(map[uint8]map[GroupKey]uint16)
					gi.JoinedInterestsCountry[a.JoinedYear] = make(map[uint8]map[GroupKey]uint16)
				}
			}
			for _, key := range JoinedCityInterestsKey(a) {
				{
					_, okb := gi.JoinedInterestsCity[a.JoinedYear][key.Interest]
					if !okb {
						gi.JoinedInterestsCity[a.JoinedYear][key.Interest] = make(map[GroupKey]uint16)
					}
					v, ok := gi.JoinedInterestsCity[a.JoinedYear][key.Interest][key]
					if ok {
						gi.JoinedInterestsCity[a.JoinedYear][key.Interest][key] = v + 1
					} else {
						gi.JoinedInterestsCity[a.JoinedYear][key.Interest][key] = 1
					}
				}
				{
					_, okb := gi.JoinedCityInterests[a.JoinedYear][key.City]
					if !okb {
						gi.JoinedCityInterests[a.JoinedYear][key.City] = make(map[GroupKey]uint16)
					}

					v, ok := gi.JoinedCityInterests[a.JoinedYear][key.City][key]
					if ok {
						gi.JoinedCityInterests[a.JoinedYear][key.City][key] = v + 1
					} else {
						gi.JoinedCityInterests[a.JoinedYear][key.City][key] = 1
					}
				}
			}
		}
		{
			for _, key := range JoinedCountryInterestsKey(a) {
				{
					_, okb := gi.JoinedInterestsCountry[a.JoinedYear][key.Interest]
					if !okb {
						gi.JoinedInterestsCountry[a.JoinedYear][key.Interest] = make(map[GroupKey]uint16)
					}
					v, ok := gi.JoinedInterestsCountry[a.JoinedYear][key.Interest][key]
					if ok {
						gi.JoinedInterestsCountry[a.JoinedYear][key.Interest][key] = v + 1
					} else {
						gi.JoinedInterestsCountry[a.JoinedYear][key.Interest][key] = 1
					}
				}
				{
					_, okb := gi.JoinedCountryInterests[a.JoinedYear][key.Country]
					if !okb {
						gi.JoinedCountryInterests[a.JoinedYear][key.Country] = make(map[GroupKey]uint16)
					}

					v, ok := gi.JoinedCountryInterests[a.JoinedYear][key.Country][key]
					if ok {
						gi.JoinedCountryInterests[a.JoinedYear][key.Country][key] = v + 1
					} else {
						gi.JoinedCountryInterests[a.JoinedYear][key.Country][key] = 1
					}
				}
			}
		}
	}
	{
		{
			for _, key := range CityInterestsKey(a) {
				{
					_, okb := gi.InterestsCity[key.Interest]
					if !okb {
						gi.InterestsCity[key.Interest] = make(map[GroupKey]uint16)
					}
					v, ok := gi.InterestsCity[key.Interest][key]
					if ok {
						gi.InterestsCity[key.Interest][key] = v + 1
					} else {
						gi.InterestsCity[key.Interest][key] = 1
					}
				}
				{
					_, okb := gi.CityInterests[key.City]
					if !okb {
						gi.CityInterests[key.City] = make(map[GroupKey]uint16)
					}

					v, ok := gi.CityInterests[key.City][key]
					if ok {
						gi.CityInterests[key.City][key] = v + 1
					} else {
						gi.CityInterests[key.City][key] = 1
					}
				}
			}
		}
		{
			for _, key := range CountryInterestsKey(a) {
				{
					_, okb := gi.InterestsCountry[key.Interest]
					if !okb {
						gi.InterestsCountry[key.Interest] = make(map[GroupKey]uint16)
					}
					v, ok := gi.InterestsCountry[key.Interest][key]
					if ok {
						gi.InterestsCountry[key.Interest][key] = v + 1
					} else {
						gi.InterestsCountry[key.Interest][key] = 1
					}
				}
				{
					_, okb := gi.CountryInterests[key.Country]
					if !okb {
						gi.CountryInterests[key.Country] = make(map[GroupKey]uint16)
					}

					v, ok := gi.CountryInterests[key.Country][key]
					if ok {
						gi.CountryInterests[key.Country][key] = v + 1
					} else {
						gi.CountryInterests[key.Country][key] = 1
					}
				}
			}
		}
	}
	{
		_, okb := gi.BirthInterests[a.BirthYear]
		if !okb {
			gi.BirthInterests[a.BirthYear] = make(map[GroupKey]uint32)
		}

		for _, key := range BirthInterestsKey(a) {
			v, ok := gi.BirthInterests[a.BirthYear][key]
			if ok {
				gi.BirthInterests[a.BirthYear][key] = v + 1
			} else {
				gi.BirthInterests[a.BirthYear][key] = 1
			}
		}
	}
	{
		_, okb := gi.JoinedInterests[a.JoinedYear]
		if !okb {
			gi.JoinedInterests[a.JoinedYear] = make(map[GroupKey]uint32)
		}

		for _, key := range JoinedInterestsKey(a) {
			v, ok := gi.JoinedInterests[a.JoinedYear][key]
			if ok {
				gi.JoinedInterests[a.JoinedYear][key] = v + 1
			} else {
				gi.JoinedInterests[a.JoinedYear][key] = 1
			}
		}
	}

	{
		// key := AllWOInterestsKey(a)
		// v, ok := gi.AllWOInterests[key]
		// if ok {
		// 	gi.AllWOInterests[key] = v + 1
		// } else {
		// 	gi.AllWOInterests[key] = 1
		// }
	}
	{
		key := JoinedCityKey(a)

		_, okb := gi.JoinedCity[a.JoinedYear]
		if !okb {
			gi.JoinedCity[a.JoinedYear] = make(map[GroupKey]uint16)
		}

		v, ok := gi.JoinedCity[a.JoinedYear][key]
		if ok {
			gi.JoinedCity[a.JoinedYear][key] = v + 1
		} else {
			gi.JoinedCity[a.JoinedYear][key] = 1
		}
	}
	{
		key := JoinedCountryKey(a)

		_, okb := gi.JoinedCountry[a.JoinedYear]
		if !okb {
			gi.JoinedCountry[a.JoinedYear] = make(map[GroupKey]uint16)
		}

		v, ok := gi.JoinedCountry[a.JoinedYear][key]
		if ok {
			gi.JoinedCountry[a.JoinedYear][key] = v + 1
		} else {
			gi.JoinedCountry[a.JoinedYear][key] = 1
		}
	}
	{
		key := BirthCityKey(a)

		_, okb := gi.BirthCity[a.BirthYear]
		if !okb {
			gi.BirthCity[a.BirthYear] = make(map[GroupKey]uint16)
		}

		v, ok := gi.BirthCity[a.BirthYear][key]
		if ok {
			gi.BirthCity[a.BirthYear][key] = v + 1
		} else {
			gi.BirthCity[a.BirthYear][key] = 1
		}
	}
	{
		key := BirthCountryKey(a)

		_, okb := gi.BirthCountry[a.BirthYear]
		if !okb {
			gi.BirthCountry[a.BirthYear] = make(map[GroupKey]uint16)
		}

		v, ok := gi.BirthCountry[a.BirthYear][key]
		if ok {
			gi.BirthCountry[a.BirthYear][key] = v + 1
		} else {
			gi.BirthCountry[a.BirthYear][key] = 1
		}
	}

	{
		key := CityKey(a)
		v, ok := gi.City[key]
		if ok {
			gi.City[key] = v + 1
		} else {
			gi.City[key] = 1
		}
	}
	{
		key := CountryKey(a)
		v, ok := gi.Country[key]
		if ok {
			gi.Country[key] = v + 1
		} else {
			gi.Country[key] = 1
		}
	}
	{
		_, okb := gi.Birth[a.BirthYear]
		if !okb {
			gi.Birth[a.BirthYear] = make(map[GroupKey]uint32)
		}

		key := BirthKey(a)
		v, ok := gi.Birth[a.BirthYear][key]
		if ok {
			gi.Birth[a.BirthYear][key] = v + 1
		} else {
			gi.Birth[a.BirthYear][key] = 1
		}
	}
	{
		_, okb := gi.Joined[a.JoinedYear]
		if !okb {
			gi.Joined[a.JoinedYear] = make(map[GroupKey]uint32)
		}

		key := JoinedKey(a)
		v, ok := gi.Joined[a.JoinedYear][key]
		if ok {
			gi.Joined[a.JoinedYear][key] = v + 1
		} else {
			gi.Joined[a.JoinedYear][key] = 1
		}
	}
	{
		for _, key := range JoinedInterestsKey(a) {
			v, ok := gi.Interests[key]
			if ok {
				gi.Interests[key] = v + 1
			} else {
				gi.Interests[key] = 1
			}
		}
	}

	{
		key := SexSatusOnlyKey(a)
		v, ok := gi.SexSatusOnly[key]
		if ok {
			gi.SexSatusOnly[key] = v + 1
		} else {
			gi.SexSatusOnly[key] = 1
		}
	}
}

func (gi *GroupIndex) Remove(a *Account) {
	{
		// for _, key := range AllWInterestsKey(a) {
		// 	v, ok := gi.AllWInterests[key]
		// 	if ok {
		// 		gi.AllWInterests[key] = v - 1
		// 	}
		// }
	}
	{ // Birth
		{
			for _, key := range BirthCityInterestsKey(a) {
				{
					v, ok := gi.BirthInterestsCity[a.BirthYear][key.Interest][key]
					if ok {
						gi.BirthInterestsCity[a.BirthYear][key.Interest][key] = v - 1
					}
				}
				{
					v, ok := gi.BirthCityInterests[a.BirthYear][key.City][key]
					if ok {
						gi.BirthCityInterests[a.BirthYear][key.City][key] = v - 1
					}
				}
			}
		}
		{
			for _, key := range BirthCountryInterestsKey(a) {
				{
					v, ok := gi.BirthInterestsCountry[a.BirthYear][key.Interest][key]
					if ok {
						gi.BirthInterestsCountry[a.BirthYear][key.Interest][key] = v - 1
					}
				}
				{
					v, ok := gi.BirthCountryInterests[a.BirthYear][key.Country][key]
					if ok {
						gi.BirthCountryInterests[a.BirthYear][key.Country][key] = v - 1
					}
				}
			}
		}
	}
	{ // Joined
		{
			for _, key := range JoinedCityInterestsKey(a) {
				{
					v, ok := gi.JoinedInterestsCity[a.JoinedYear][key.Interest][key]
					if ok {
						gi.JoinedInterestsCity[a.JoinedYear][key.Interest][key] = v - 1
					}
				}
				{
					v, ok := gi.JoinedCityInterests[a.JoinedYear][key.City][key]
					if ok {
						gi.JoinedCityInterests[a.JoinedYear][key.City][key] = v - 1
					}
				}
			}
		}
		{
			for _, key := range JoinedCountryInterestsKey(a) {
				{
					v, ok := gi.JoinedInterestsCountry[a.JoinedYear][key.Interest][key]
					if ok {
						gi.JoinedInterestsCountry[a.JoinedYear][key.Interest][key] = v - 1
					}
				}
				{
					v, ok := gi.JoinedCountryInterests[a.JoinedYear][key.Country][key]
					if ok {
						gi.JoinedCountryInterests[a.JoinedYear][key.Country][key] = v - 1
					}
				}
			}
		}
	}
	{
		{
			for _, key := range CityInterestsKey(a) {
				{
					v, ok := gi.InterestsCity[key.Interest][key]
					if ok {
						gi.InterestsCity[key.Interest][key] = v - 1
					}
				}
				{
					v, ok := gi.CityInterests[key.City][key]
					if ok {
						gi.CityInterests[key.City][key] = v - 1
					}
				}
			}
		}
		{
			for _, key := range CountryInterestsKey(a) {
				{

					v, ok := gi.InterestsCountry[key.Interest][key]
					if ok {
						gi.InterestsCountry[key.Interest][key] = v - 1
					}
				}
				{
					v, ok := gi.CountryInterests[key.Country][key]
					if ok {
						gi.CountryInterests[key.Country][key] = v - 1
					}
				}
			}
		}
	}
	{
		for _, key := range BirthInterestsKey(a) {
			v, ok := gi.BirthInterests[a.BirthYear][key]
			if ok {
				gi.BirthInterests[a.BirthYear][key] = v - 1
			}
		}
	}
	{
		for _, key := range JoinedInterestsKey(a) {
			v, ok := gi.JoinedInterests[a.JoinedYear][key]
			if ok {
				gi.JoinedInterests[a.JoinedYear][key] = v - 1
			}
		}
	}

	{
		// key := AllWOInterestsKey(a)
		// v, ok := gi.AllWOInterests[key]
		// if ok {
		// 	gi.AllWOInterests[key] = v - 1
		// }
	}
	{
		key := JoinedCityKey(a)
		v, ok := gi.JoinedCity[a.JoinedYear][key]
		if ok {
			gi.JoinedCity[a.JoinedYear][key] = v - 1
		}
	}
	{
		key := JoinedCountryKey(a)
		v, ok := gi.JoinedCountry[a.JoinedYear][key]
		if ok {
			gi.JoinedCountry[a.JoinedYear][key] = v - 1
		}
	}
	{
		key := BirthCityKey(a)
		v, ok := gi.BirthCity[a.BirthYear][key]
		if ok {
			gi.BirthCity[a.BirthYear][key] = v - 1
		}
	}
	{
		key := BirthCountryKey(a)
		v, ok := gi.BirthCountry[a.BirthYear][key]
		if ok {
			gi.BirthCountry[a.BirthYear][key] = v - 1
		}
	}

	{
		key := CityKey(a)
		v, ok := gi.City[key]
		if ok {
			gi.City[key] = v - 1
		}
	}
	{
		key := CountryKey(a)
		v, ok := gi.Country[key]
		if ok {
			gi.Country[key] = v - 1
		}
	}
	{
		key := BirthKey(a)
		v, ok := gi.Birth[a.BirthYear][key]
		if ok {
			gi.Birth[a.BirthYear][key] = v - 1
		}
	}
	{
		key := JoinedKey(a)
		v, ok := gi.Joined[a.JoinedYear][key]
		if ok {
			gi.Joined[a.JoinedYear][key] = v - 1
		}
	}

	{
		for _, key := range JoinedInterestsKey(a) {
			v, ok := gi.Interests[key]
			if ok {
				gi.Interests[key] = v - 1
			}
		}
	}

	{
		key := SexSatusOnlyKey(a)
		v, ok := gi.SexSatusOnly[key]
		if ok {
			gi.SexSatusOnly[key] = v - 1
		}
	}
}
