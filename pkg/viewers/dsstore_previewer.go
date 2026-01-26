package viewers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/filetug/filetug/pkg/files"
	"github.com/filetug/filetug/pkg/fsutils"
	"github.com/strongo/dsstore"
)

var _ Previewer = (*DsstorePreviewer)(nil)

type DsstorePreviewer struct {
	TextPreviewer
}

func NewDsstorePreviewer() *DsstorePreviewer {
	previewer := NewTextPreviewer()
	return &DsstorePreviewer{
		TextPreviewer: *previewer,
	}
}

func (p DsstorePreviewer) Preview(entry files.EntryWithDirPath, data []byte, queueUpdateDraw func(func())) {

	fullName := entry.FullName()
	data, err := fsutils.ReadFileData(fullName, 0)
	if err != nil {
		return
	}
	bufferRead := bytes.NewBuffer(data)
	var s dsstore.Store
	err = s.Read(bufferRead)
	if err != nil {
		p.showError(fmt.Sprintf("Failed to read %s: %s", entry.Name(), err.Error()))
		return
	}
	var sb strings.Builder
	for _, r := range s.Records {
		sb.WriteString(fmt.Sprintf("%s: %s\n", r.FileName, r.Type))
	}
	data = []byte(sb.String())
	p.Preview(entry, data, queueUpdateDraw)
}
