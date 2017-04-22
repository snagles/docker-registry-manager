package manager

type Repository struct {
	Name string
	Tags map[string]*Tag
}
