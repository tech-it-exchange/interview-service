package mappers

type Mappers struct {
	Http *HttpMapper
}

func NewMappers() *Mappers {
	httpMapper := NewHttpMapper()

	return &Mappers{
		Http: httpMapper,
	}
}
