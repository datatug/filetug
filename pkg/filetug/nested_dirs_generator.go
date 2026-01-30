package filetug

import (
	"context"
	"fmt"
	"path"
	"sync"

	"github.com/filetug/filetug/pkg/files"
	"github.com/filetug/filetug/pkg/sneatv"
	"github.com/rivo/tview"
)

type nestedDirsGeneratorPanel struct {
	*sneatv.Boxed
	flex *tview.Flex
	form *tview.Form
}

func newNestedDirsGeneratorPanel(nav *Navigator, active tview.Primitive) tview.Primitive {

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	form := tview.NewForm()
	flex.AddItem(form, 0, 1, true)

	form.AddInputField("Depth", "20", 0, nil, nil)
	form.AddInputField("SubDirs", "10", 0, nil, nil)
	form.AddInputField("FilesPerDir", "10", 0, nil, nil)
	form.AddInputField("File Size (bytes)", "1024", 0, nil, nil)
	form.AddButton("Generate", func() {})
	form.AddButton("Cancel", func() {
		nav.right.SetContent(nav.previewer)
		if active != nil {
			nav.app.SetFocus(active)
		}
	})

	//var spacer tview.Primitive = nil
	//buttons := tview.NewFlex()
	//rows.AddItem(buttons, 3, 0, false)
	//
	//generateBtn := tview.NewButton("Generate")
	//buttons.AddItem(generateBtn, 0, len(generateBtn.GetLabel()), false)
	//
	//buttons.AddItem(spacer, 1, 0, false)
	//
	//cancelBtn := tview.NewButton("Cancel")
	//buttons.AddItem(cancelBtn, 0, len(cancelBtn.GetLabel()), false)
	//cancelBtn.SetSelectedFunc(func() {
	//})

	p := nestedDirsGeneratorPanel{
		Boxed: sneatv.NewBoxed(flex),
		flex:  flex,
		form:  form,
	}
	return &p
}

func GeneratedNestedDirs(ctx context.Context, store files.Store, dirPath, subDirNameFormat string, depth, subDirsCount int) (err error) {
	if err = store.CreateDir(ctx, dirPath); err != nil {
		return err
	}
	if subDirNameFormat == "" {
		subDirNameFormat = "Directory%d"
	}
	if depth == 0 {
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(subDirsCount)
	for i := 0; i < subDirsCount; i++ {
		i := i
		go func() {
			defer wg.Done()
			subDirName := fmt.Sprintf(subDirNameFormat, i)
			subDirPath := path.Join(dirPath, subDirName)
			subDirErr := GeneratedNestedDirs(ctx, store, subDirPath, subDirNameFormat, depth-1, subDirsCount)
			if subDirErr != nil {
				err = subDirErr
			}
		}()
	}
	wg.Wait()
	return nil
}
