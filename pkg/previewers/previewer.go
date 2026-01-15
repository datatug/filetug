package previewers

type Previewer interface {
	GetMeta(path string) (meta *Meta)
}

type Meta struct {
	Groups []*MetaGroup
}
type MetaGroup struct {
	ID      string        `json:"id"`
	Title   string        `json:"title"`
	Records []*MetaRecord `json:"records"`
}

type MetaRecord struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Value string `json:"value"`
	//TitleAlign Align
	ValueAlign Align
}

type Align int

const (
	AlignLeft Align = iota
	AlignRight
)
