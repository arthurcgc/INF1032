package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/arthurcgc/go-scraper/types"
)

var (
	api_key = "84f39195ae013547986894bbd4652ffd"
)

func MoviesByYear(year int, page int) (types.MovieMeta, error) {
	req, err := http.NewRequest("GET", "http://api.themoviedb.org/3/discover/movie", nil)
	if err != nil {
		return types.MovieMeta{}, err
	}
	query := req.URL.Query()
	query.Add("api_key", api_key)
	query.Add("primary_release_year", strconv.Itoa(year))
	query.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return types.MovieMeta{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.MovieMeta{}, err
	}
	movies := types.MovieMeta{}
	err = json.Unmarshal(body, &movies)
	if err != nil {
		return types.MovieMeta{}, err
	}

	return movies, nil
}

func FullMovie(movieId int) (types.Movie, error) {
	req, err := http.NewRequest("GET", "http://api.themoviedb.org/3/movie/"+strconv.Itoa(movieId), nil)
	if err != nil {
		return types.Movie{}, err
	}
	query := req.URL.Query()
	query.Add("api_key", api_key)
	req.URL.RawQuery = query.Encode()

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return types.Movie{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.Movie{}, err
	}
	movie := types.Movie{}
	err = json.Unmarshal(body, &movie)
	if err != nil {
		return types.Movie{}, err
	}

	return movie, nil
}

func WriteMovies(metaMovies types.MovieMeta, writer *types.InternalWriter) error {
	for _, movie := range metaMovies.Results {
		fullMovie, err := FullMovie(movie.Id)
		if err != nil {
			continue
		}
		if fullMovie.Budget > 0 && fullMovie.Revenue > 0 {
			err = writer.WriteMovie(fullMovie)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ParsePages(wg *sync.WaitGroup, year int, csvWriter *types.InternalWriter) {
	for pageCount := 1; ; pageCount++ {
		movieIds, err := MoviesByYear(year, pageCount)
		fmt.Printf("year: %d\t page: %d\t total pages: %d\n", year, pageCount, movieIds.TotalPages)
		if err != nil {
			continue
		}
		err = WriteMovies(movieIds, csvWriter)
		if err != nil {
			panic(err)
		}
		if pageCount > movieIds.TotalPages {
			break
		}
	}
	wg.Done()

}

func main() {
	csvWriter, err := types.NewWriter()
	if err != nil {
		panic(err)
	}
	err = csvWriter.WriteHeader()
	if err != nil {
		panic(err)
	}
	defer csvWriter.CloseFile()

	var wg sync.WaitGroup
	for year := 1900; year < 2021; year++ {
		wg.Add(1)
		go ParsePages(&wg, year, csvWriter)
	}
	wg.Wait()
}
