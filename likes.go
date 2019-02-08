package main

// import (
// 	"go_v1/common"
// )

// // type LikeElement struct {
// // 	Ts int32
// // 	Id common.Int24
// // }

// type LikesindexV2 struct {
// 	Likes map[common.Int24]*AccountLikeDirect
// }

// func LikesindexV2Create() *LikesindexV2 {
// 	return &LikesindexV2{Likes: make(map[common.Int24]*AccountLikeDirect)}
// }

// func (li *LikesindexV2) Append(idFrom common.Int24, ts int32, id common.Int24) {

// 	ald, ok := li.Likes[idFrom]
// 	if !ok {
// 		ald = AccountLikeDirectCreate()
// 		li.Likes[idFrom] = ald
// 	}

// 	ald.Append(ts, id)

// }

// type AccountLikeDirect struct {
// 	Likes   *[]LikeElement
// 	Correct map[common.Int24]uint16
// }

// func AccountLikeDirectCreate() *AccountLikeDirect {
// 	return &AccountLikeDirect{Likes: &[]LikeElement{}, Correct: make(map[common.Int24]uint16)}
// }

// func (ald *AccountLikeDirect) Append(ts int32, id common.Int24) {

// 	corr, ok := ald.Correct[id]
// 	pos := -1
// 	step := -1

// 	for i, v := range *ald.Likes { // Двоичный поиск бы сюда
// 		if v.Id.More(id) {
// 			step = i
// 		}
// 		if v.Id == id {
// 			pos = i
// 			break
// 		}
// 	}

// 	if pos == -1 {
// 		s := make([]LikeElement, 0, len(*ald.Likes)+1)
// 		if step >= 0 {
// 			s = append(s, (*ald.Likes)[0:step+1]...)
// 			s = append(s, LikeElement{Id: id, Ts: ts})
// 			s = append(s, (*ald.Likes)[step+1:]...)
// 		} else {
// 			s = append(s, LikeElement{Id: id, Ts: ts})
// 			s = append(s, (*ald.Likes)[0:]...)
// 		}
// 		ald.Likes = &s
// 	} else if ok {
// 		tsn := (int64((*ald.Likes)[pos].Ts)*int64(corr) + int64(ts)) / int64(corr+1)
// 		(*ald.Likes)[pos] = LikeElement{Id: id, Ts: int32(tsn)}
// 		//ald.Correct[id] = corr + 1
// 		i_li_map++
// 	} else {
// 		tsn := (int64((*ald.Likes)[pos].Ts) + int64(ts)) / 2
// 		(*ald.Likes)[pos] = LikeElement{Id: id, Ts: int32(tsn)}
// 		//ald.Correct[id] = 2
// 		i_li_map++
// 	}
// }
