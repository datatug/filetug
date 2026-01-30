package filetug

import (
	"testing"

	"github.com/filetug/filetug/pkg/filetug/navigator"
	"go.uber.org/mock/gomock"
)

func TestSetupApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := navigator.NewMockApp(ctrl)
	app.EXPECT().QueueUpdateDraw(gomock.Any()).MinTimes(1)
	app.EXPECT().EnableMouse(true)
	app.EXPECT().SetRoot(gomock.Any(), true).Times(1)
	SetupApp(app)
}
