package fetch

import (
	"encoding/json"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testFile struct{}

func (testFile) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func (testFile) WriteJSON(filename string, data []byte) error {
	return nil
}

func (testFile) PathRead() string {
	return "./test_data"
}

func (testFile) PathWrite() string {
	return "./test_data"
}

var _ = Describe("Fetch", func() {
	Describe("Fetch", func() {
		var series = []Serie{}

		BeforeEach(func() {
			js := fetch(testFile{})
			series = []Serie{}
			json.Unmarshal(js, &series)
		})

		Describe("serie", func() {
			It("should not be nil", func() {
				Expect(series).To(BeAssignableToTypeOf([]Serie{}))
				Expect(series[0]).ToNot(BeNil())
				Expect(series[1]).ToNot(BeNil())
			})

			It("should have 2 series", func() {
				Expect(len(series)).To(Equal(2))
				Expect(series[0]).To(BeAssignableToTypeOf(Serie{}))
				Expect(series[1]).To(BeAssignableToTypeOf(Serie{}))
			})
		})

		Describe("season", func() {
			It("should have 4 seasons in total", func() {
				for _, serie := range series {
					Expect(serie.Seasons).To(BeAssignableToTypeOf([]*Season{}))
					Expect(len(serie.Seasons)).To(Equal(2))
				}
			})
		})

		Describe("episode", func() {
			It("each season should have 1 episode", func() {
				for _, serie := range series {
					for _, season := range serie.Seasons {
						Expect(season.Episodes).To(BeAssignableToTypeOf([]*Episode{}))
						Expect(len(season.Episodes)).To(Equal(1))
					}
				}
			})
		})

		Describe("subtitles", func() {
			It("each episode should have 1 subtitle", func() {
				for _, serie := range series {
					for _, season := range serie.Seasons {
						for _, episode := range season.Episodes {
							Expect(episode.Subtitles).To(BeAssignableToTypeOf([]Subtitle{}))
							Expect(len(episode.Subtitles)).To(Equal(1))
						}
					}
				}
			})
		})
	})

	Describe("filename", func() {
		It("should work as intended", func() {
			Expect(filename("01.avi")).To(Equal("01"))
			Expect(filename(".01.avi")).To(Equal(".01"))
			Expect(filename(".01.avi.")).To(Equal(".01.avi"))
			Expect(filename("01.avi.")).To(Equal("01.avi"))
			Expect(filename(".avi.")).To(Equal(".avi"))
			Expect(filename(".avi..")).To(Equal(".avi."))
			Expect(filename("..avi..")).To(Equal("..avi."))
		})
		It("should return everything before the last dot", func() {
			name := "Game.of.Thrones.S06E07.HDTV.x264-KILLERS.mkv"
			Expect(filename(name)).To(Equal("Game.of.Thrones.S06E07.HDTV.x264-KILLERS"))
		})
	})
})
