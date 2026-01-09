package sneatv

func WithSeparator(separator string) func(bc *Breadcrumbs) {
	return func(bc *Breadcrumbs) {
		bc.separator = separator
	}
}
