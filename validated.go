package govalidations

type Error struct {
	Name    string
	Message string
}

type Errors []*Error

type Validated struct {
	Object interface{}
	Errors Errors
}

func (vd *Validated) HasError() bool {
	return len(vd.Errors) > 0
}

func (es Errors) Has(name string) bool {
	for _, e := range es {
		if e.Name == name {
			return true
		}
	}
	return false
}

func (es Errors) IfHasThen(name string, result string) (r string) {
	if es.Has(name) {
		return result
	}
	return
}

func (es Errors) On(name string) (r string) {
	for _, e := range es {
		if e.Name == name {
			r = e.Message
			return
		}
	}
	return
}