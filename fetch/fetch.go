package fetch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
)

const (
	PATH     = "/media/sf_Z_DRIVE/serien"
	OUT_PATH = "/home/tim/Workspace/src/localflix-server-/fetch"
	OUT      = "DATA_DUMP"
	VIDEO    = ".*[.](avi|web|mkv|mp4)$"
	SUBTITLE = ".*[.](srt|idx|sub|sfv)$"
	IGNORE   = ".*[.](json)$"
)

type FileType interface {
	ReadDir(path string) ([]os.FileInfo, error)
	WriteJSON(filename string, data []byte) error
	PathRead() string
	PathWrite() string
}

type File struct{}

func (File) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}
func (File) WriteJSON(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0666)
}
func (File) PathRead() string {
	return PATH
}
func (File) PathWrite() string {
	return OUT_PATH
}

func Fetch() {
	fetch(File{})
}

func fetch(file FileType) []byte {
	series := []Serie{}

	path := file.PathRead()
	files, _ := file.ReadDir(path)

	for _, val := range files {
		if ignore(val.Name()) {
			continue
		}
		serie := Serie{Name: val.Name()}
		serie.fetchSeasons(path, file)
		series = append(series, serie)
	}

	js, _ := json.Marshal(series)
	file.WriteJSON(file.PathWrite()+"/"+OUT+".json", []byte(js))
	return js
}

type Serie struct {
	Name    string
	Seasons []*Season
}

type Season struct {
	Name      string
	SerieName string
	Episodes  []*Episode
}

type Episode struct {
	Name      string
	Extension string
	Subtitles []Subtitle
	Src       string
}

type Subtitle struct {
	Name string
}

func (season *Season) findEpisode(name string) (*Episode, bool) {
	for _, val := range season.Episodes {
		if val.Name == name {
			return val, true
		}
	}
	return &Episode{}, false
}

func (serie *Serie) fetchSeasons(path string, file FileType) {
	newPath := path + "/" + serie.Name
	files, _ := file.ReadDir(newPath)

	for i, val := range files {
		serie.Seasons = append(serie.Seasons, &Season{Name: val.Name(), SerieName: serie.Name})
		serie.Seasons[i].fetchEpisodes(newPath, file)
	}
}

func (season *Season) fetchEpisodes(path string, file FileType) {
	newPath := path + "/" + season.Name
	files, _ := file.ReadDir(newPath)

	for _, val := range files {
		episode := &Episode{
			Name:      filename(val.Name()),
			Extension: "",
			Subtitles: []Subtitle{},
		}

		if isValid(VIDEO, val) {
			ep, exist := season.findEpisode(filename(val.Name()))
			if !exist {
				season.Episodes = append(season.Episodes, episode)
				ep = episode
			}
			if len(ep.Extension) == 0 {
				ep.Extension = extension(val.Name())
			}
			ep.Src = newPath
		}

		if isValid(SUBTITLE, val) {
			name := filename(val.Name())

			ep, exist := season.findEpisode(name)
			if !exist {
				season.Episodes = append(season.Episodes, episode)
				ep = episode
			}
			ep.Subtitles = append(ep.Subtitles, Subtitle{Name: val.Name()})
		}
	}
}

func isValid(TYPE string, file os.FileInfo) bool {
	matched, _ := regexp.MatchString(TYPE, file.Name())
	return matched
}

func filename(name string) string {
	regex := regexp.MustCompile("^(.*)[.].*")
	return regex.ReplaceAllString(name, "$1")
}

func extension(name string) string {
	regex := regexp.MustCompile(".*[.](.*)$")
	return regex.ReplaceAllString(name, "$1")
}

func ignore(name string) bool {
	ok, err := regexp.Match(IGNORE, []byte(name))
	if err != nil {
		println(err)
	}
	return ok
}
