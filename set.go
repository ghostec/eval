package eval

type Set map[interface{}]bool

func NewSet(elements ...interface{}) Set {
	set := map[interface{}]bool{}

	for _, el := range elements {
		set[el] = true
	}

	return set
}

func (set Set) Contains(el interface{}) bool {
	_, ok := set[el]
	return ok
}
