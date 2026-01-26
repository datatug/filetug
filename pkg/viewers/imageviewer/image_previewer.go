package imageviewer

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"

	"github.com/filetug/filetug/pkg/files"
	"github.com/filetug/filetug/pkg/viewers"
	"github.com/rivo/tview"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
)

var _ viewers.Previewer = (*ImagePreviewer)(nil)

type ImagePreviewer struct {
	metaTable *viewers.MetaTable
}

func NewImagePreviewer() *ImagePreviewer {
	return &ImagePreviewer{
		metaTable: viewers.NewMetaTable(),
	}
}

func (p ImagePreviewer) Preview(entry files.EntryWithDirPath, _ []byte, queueUpdateDraw func(func())) {
	fullName := entry.FullName()
	meta := p.GetMeta(fullName)
	if meta != nil {
		queueUpdateDraw(func() {
			metaTable := viewers.NewMetaTable()
			metaTable.SetMeta(meta)
		})
	}
}

func (p ImagePreviewer) Meta() tview.Primitive {
	return p.metaTable
}

func (p ImagePreviewer) Main() tview.Primitive {
	return nil
}

func (p ImagePreviewer) GetMeta(path string) (meta *viewers.Meta) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return
	}
	main := viewers.MetaGroup{
		ID:    "main",
		Title: "Format: " + strings.ToUpper(format),
	}
	main.Records = append(main.Records,
		&viewers.MetaRecord{
			ID:         "width",
			Title:      "Width",
			Value:      strconv.Itoa(cfg.Width),
			ValueAlign: viewers.AlignRight,
		},
		&viewers.MetaRecord{
			ID:         "height",
			Title:      "Height",
			Value:      strconv.Itoa(cfg.Height),
			ValueAlign: viewers.AlignRight,
		},
	)
	return &viewers.Meta{
		Groups: []*viewers.MetaGroup{
			&main,
		},
	}
}
