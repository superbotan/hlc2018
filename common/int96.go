package common

type Int96 [12]byte

func Int96Create(arr []uint8) (r Int96) {
	i := Int96{}
	i = i.SetArr(arr)
	return i
}

func (i Int96) Set(n uint8) (r Int96) {

	i[n/8] |= (1 << (n % 8))

	return i
}

func (i Int96) SetArr(arr []uint8) (r Int96) {

	for _, n := range arr {
		i = i.Set(n)
	}

	return i
}

func (i Int96) GetArr() (arr []uint8) {
	arr = make([]uint8, 0)

	var n uint8
	for n = 0; n < 96; n++ {
		if i[n/8]&(1<<(n%8)) > 0 {
			arr = append(arr, n)
		}
	}

	return arr
}

func (i Int96) Get(n uint8) bool {
	return i[n/8]&(1<<(n%8)) > 0
}

func (i Int96) Contains(il Int96) bool {
	for n := 0; n < 12; n++ {
		if i[n]&il[n] > 0 {
			return true
		}
	}
	return false
}

func (i Int96) AllIn(il Int96) bool {
	for _, n := range i.GetArr() {
		if !il.Get(n) {
			return false
		}
	}
	return true
}

func (i Int96) ComplCount(il Int96) (r uint8) {

	for _, n := range i.GetArr() {
		if il.Get(n) {
			r = r + 1
		}
	}
	return r
}

func (i Int96) ComplCountArr(arr []uint8) (r uint8) {

	for _, n := range arr {
		if i.Get(n) {
			r = r + 1
		}
	}
	return r
}
