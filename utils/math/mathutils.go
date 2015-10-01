package math



/**
I am surprised that I have to define these. . . Maybe i just didn't look hard enough for a lib.
 */

//http://blog.plover.com/math/choose.html
func NChoseK(n, k uint) uint64 {
	uN := uint64(n)
	uK := uint64(k)
	if uK > uN {
		return 0
	} else if uK == 0 {
		return 1
	}

	var  r uint64 = 1

	for d := uint64(1) ; d <= uK; d++ {
		r *= uN
		r /= d
		uN--
	}

	return r
}