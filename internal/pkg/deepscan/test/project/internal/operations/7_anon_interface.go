package operations

func AnonAnyInterface7(x interface{}) {
	_ = x
}

func AnonMethodInterface7(x interface{ Hello() }) {
	x.Hello()
}
