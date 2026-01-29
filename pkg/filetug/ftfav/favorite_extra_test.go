package ftfav

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFavorites_GetFavorites_InvalidYaml(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "favorites.yaml")
	oldPath := favoritesFilePath
	favoritesFilePath = tempPath
	defer func() {
		favoritesFilePath = oldPath
	}()

	err := os.WriteFile(tempPath, []byte("invalid: ["), 0o644)
	assert.NoError(t, err)

	favorites, err := GetFavorites()
	assert.Nil(t, favorites)
	assert.Error(t, err)
}

func TestFavorites_GetFavorites_InvalidStoreURL(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "favorites.yaml")
	oldPath := favoritesFilePath
	favoritesFilePath = tempPath
	defer func() {
		favoritesFilePath = oldPath
	}()

	data := []byte("- store: \"http://[::1\"\n  path: /tmp\n")
	err := os.WriteFile(tempPath, data, 0o644)
	assert.NoError(t, err)

	favorites, err := GetFavorites()
	assert.Nil(t, favorites)
	assert.Error(t, err)
}

func TestFavorites_GetFavorites_FileNotExists(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "missing.yaml")
	oldPath := favoritesFilePath
	favoritesFilePath = tempPath
	defer func() {
		favoritesFilePath = oldPath
	}()

	favorites, err := GetFavorites()
	assert.NoError(t, err)
	assert.Len(t, favorites, 0)
}

func TestFavorites_GetFavorites_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "favorites.yaml")
	oldPath := favoritesFilePath
	favoritesFilePath = tempPath
	defer func() {
		favoritesFilePath = oldPath
	}()

	err := os.WriteFile(tempPath, []byte(""), 0o644)
	assert.NoError(t, err)

	favorites, err := GetFavorites()
	assert.NoError(t, err)
	assert.Len(t, favorites, 0)
}

func TestFavorites_AddDelete_EmptyPath(t *testing.T) {
	oldPath := favoritesFilePath
	favoritesFilePath = ""
	defer func() {
		favoritesFilePath = oldPath
	}()

	addErr := AddFavorite(Favorite{Path: "/tmp"})
	assert.ErrorIs(t, addErr, errUserHomeDirIsUnknown)

	deleteErr := DeleteFavorite(Favorite{Path: "/tmp"})
	assert.ErrorIs(t, deleteErr, errUserHomeDirIsUnknown)
}

func TestFavorites_AddDelete_GetFavoritesError(t *testing.T) {
	tempDir := t.TempDir()
	oldPath := favoritesFilePath
	favoritesFilePath = tempDir
	defer func() {
		favoritesFilePath = oldPath
	}()

	addErr := AddFavorite(Favorite{Path: "/tmp"})
	assert.Error(t, addErr)

	deleteErr := DeleteFavorite(Favorite{Path: "/tmp"})
	assert.Error(t, deleteErr)
}

func TestFavorites_WriteFavorites_MkdirError(t *testing.T) {
	tempDir := t.TempDir()
	parentFile := filepath.Join(tempDir, "parent")
	oldPath := favoritesFilePath
	favoritesFilePath = filepath.Join(parentFile, "favorites.yaml")
	defer func() {
		favoritesFilePath = oldPath
	}()

	err := os.WriteFile(parentFile, []byte("x"), 0o644)
	assert.NoError(t, err)

	writeErr := writeFavorites([]Favorite{{Path: "/tmp"}})
	assert.Error(t, writeErr)
}

func TestFavorites_WriteFavorites_MarshalError(t *testing.T) {
	oldPath := favoritesFilePath
	oldMarshal := yamlMarshal
	favoritesFilePath = filepath.Join(t.TempDir(), "favorites.yaml")
	defer func() {
		favoritesFilePath = oldPath
		yamlMarshal = oldMarshal
	}()

	yamlMarshal = func(in any) ([]byte, error) {
		_ = in
		return nil, assert.AnError
	}

	writeErr := writeFavorites([]Favorite{{Path: "/tmp"}})
	assert.Error(t, writeErr)
}

func TestFavorites_DeleteFavorite_KeepsOtherItems(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "favorites.yaml")
	oldPath := favoritesFilePath
	favoritesFilePath = tempPath
	defer func() {
		favoritesFilePath = oldPath
	}()

	first := Favorite{Path: "/first"}
	second := Favorite{Path: "/second"}

	err := AddFavorite(first)
	assert.NoError(t, err)
	err = AddFavorite(second)
	assert.NoError(t, err)

	err = DeleteFavorite(first)
	assert.NoError(t, err)

	favorites, err := GetFavorites()
	assert.NoError(t, err)
	assert.Len(t, favorites, 1)
	assert.Equal(t, "/second", favorites[0].Path)
}
