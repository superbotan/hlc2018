package main

type AccountAdd struct {
	ID    int32  `json:"id,required"`
	Email string `json:"email,required"` // uniq

	Fname string `json:"fname,omitempty"`
	Sname string `json:"sname,omitempty"`

	Phone string `json:"phone,omitempty"` // uniq
	Sex   string `json:"sex,required"`    // "m" означает мужской пол, а "f" - женский
	Birth int32  `json:"birth,required"`

	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`

	Joined int32  `json:"joined,required"` // 01.01.2011, сверху 01.01.2018
	Status string `json:"status,required"` // "свободны", "заняты", "всё сложно".

	Interests []string `json:"interests,omitempty"`

	Premium struct {
		Start  int32 `json:"start,required"` // нижней границей 01.01.2018
		Finish int32 `json:"finish,required"`
	} `json:"premium,omitempty"`

	Likes []struct {
		Dt int32 `json:"ts,required"`
		ID int32 `json:"id,required"`
	} `json:"likes,omitempty"`
}

type AccountFile struct {
	Accounts []AccountAdd `json:"accounts,required"`
}

type LikesAdd struct {
	Likes []struct {
		Liker int32 `json:"liker,required"`
		Likee int32 `json:"likee,required"`
		Ts    int32 `json:"ts,required"`
	} `json:"likes,required"`
}

type AccountUpd struct {
	Birth     *int32   `json:"birth,omitempty"`
	Email     string   `json:"email,omitempty"`
	Fname     string   `json:"fname,omitempty"`
	Sname     string   `json:"sname,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Sex       string   `json:"sex,omitempty"`
	Country   string   `json:"country,omitempty"`
	City      string   `json:"city,omitempty"`
	Joined    int32    `json:"joined,omitempty"`
	Status    string   `json:"status,omitempty"`
	Interests []string `json:"interests,omitempty"`

	Premium struct {
		Start  int32 `json:"start,required"`
		Finish int32 `json:"finish,required"`
	} `json:"premium,omitempty"`

	Likes []struct {
		Dt int32 `json:"ts,required"`
		ID int32 `json:"id,required"`
	} `json:"likes,omitempty"`
}
