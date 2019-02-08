package common

type Int24 struct {
	c uint8
	b uint8
	a uint8
}

var (
	Int24Zero = Int24Create(0)
)

func Int24Create(i int32) Int24 {
	r := Int24{}

	r.a = uint8(i % 256)
	r.b = uint8((i / 256) % 256)
	r.c = uint8(((i / 256) / 256) % 256)

	return r
}
func (i Int24) Int() int32 {
	r := int32(i.a) + (int32(i.b)+int32(i.c)*256)*256

	return r
}

func (i Int24) Equal(j Int24) bool {
	return i.a == j.a && i.b == j.b && i.c == j.c
}

// Less i < j
func (i Int24) Less(j Int24) bool {
	if i.c != j.c {
		return i.c < j.c
	}
	if i.b != j.b {
		return i.b < j.b
	}

	return i.a < j.a
}

// More i > j
func (i Int24) More(j Int24) bool {
	if i.c != j.c {
		return i.c > j.c
	}
	if i.b != j.b {
		return i.b > j.b
	}

	return i.a > j.a
}

func (i Int24) Plus(j Int24) Int24 {
	return Int24Create(i.Int() + j.Int())
}
func (i Int24) Minus(j Int24) Int24 {
	return Int24Create(i.Int() - j.Int())
}

func (i Int24) C() uint8 {
	return i.c
}

func (i Int24) BA() uint16 {
	return uint16(i.a) + (uint16(i.b) * 256)
}
