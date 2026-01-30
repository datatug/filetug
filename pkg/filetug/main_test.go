package filetug

import (
	"testing"

	"github.com/filetug/filetug/pkg/filetug/navigator"
	"go.uber.org/mock/gomock"
)

func TestSetupApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := navigator.NewMockApp(ctrl)
	app.EXPECT().EnableMouse(gomock.Any()).AnyTimes()
	app.EXPECT().SetRoot(gomock.Any(), gomock.Any()).AnyTimes()
	SetupApp(app)
}
