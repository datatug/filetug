package ftfav

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"

	"github.com/filetug/filetug/pkg/filetug/ftsettings"
	"gopkg.in/yaml.v3"
)

type Favorite struct {
	Store       url.URL `json:"store,omitempty" yaml:"store,omitempty"`
	Path        string  `json:"path" yaml:"path"`
	Shortcut    rune    `json:"shortcut,omitempty" yaml:"shortcut,omitempty"`
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
}

func (f Favorite) Key() string {
	key := f.Store
	key.Path = filepath.Join(key.Path, f.Path)
	return key.String()
}

const favoritesFileName = "datatug-favorites.yaml"

var favoritesFilePath string

var GetDatatugUserDir = ftsettings.GetDatatugUserDir

func init() {
	datatugUserDir, err := GetDatatugUserDir()
	if err == nil {
		favoritesFilePath = filepath.Join(datatugUserDir, favoritesFileName)
	}
}

var errUserHomeDirIsUnknown = errors.New("user home directory is unknown")

func GetFavorites() (favorites []Favorite, err error) {
	if favoritesFilePath == "" {
		return nil, errUserHomeDirIsUnknown
	}
	data, err := os.ReadFile(favoritesFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Favorite{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Favorite{}, nil
	}
	err = yaml.Unmarshal(data, &favorites)
	if err != nil {
		return nil, err
	}
	return favorites, nil
}

func AddFavorite(f Favorite) (err error) {
	if favoritesFilePath == "" {
		return errUserHomeDirIsUnknown
	}
	favorites, err := GetFavorites()
	if err != nil {
		return err
	}
	favorites = append(favorites, f)
	return writeFavorites(favorites)
}

func DeleteFavorite(f Favorite) (err error) {
	if favoritesFilePath == "" {
		return errUserHomeDirIsUnknown
	}
	favorites, err := GetFavorites()
	if err != nil {
		return err
	}
	deleteKey := f.Key()
	updated := make([]Favorite, 0, len(favorites))
	for _, item := range favorites {
		if item.Key() == deleteKey {
			continue
		}
		updated = append(updated, item)
	}
	return writeFavorites(updated)
}

func writeFavorites(favorites []Favorite) error {
	data, err := yaml.Marshal(favorites)
	if err != nil {
		return err
	}
	dir := filepath.Dir(favoritesFilePath)
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}
	return os.WriteFile(favoritesFilePath, data, 0o644)
}
