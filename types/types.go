package types

import (
	"strconv"
	"strings"
)

type MovieInfo struct {
	Id    int
	Title string `json:"original_title"`
}

type MovieMeta struct {
	Page       int
	TotalPages int         `json:"total_pages"`
	Results    []MovieInfo `json:"results"`
}

type ProdComp struct {
	Name          string
	Id            int
	OriginCountry string `json:"origin_country"`
}

type ProdCountry struct {
	Name string
}

type Language struct {
	Name string
}

type Genre struct {
	Id   int
	Name string
}

type Movie struct {
	MovieInfo
	Adult               bool
	Budget              int
	Homepage            string
	Genres              []Genre
	OriginalLanguage    string `json:"original_language"`
	Overview            string
	Popularity          float64
	ProductionCompanies []ProdComp    `json:"production_companies"`
	ProductionCountries []ProdCountry `json:"production_countries"`
	ReleaseDate         string        `json:"release_date"`
	Revenue             int
	Runtime             int
	Languages           []Language `json:"spoken_languages"`
	Status              string
	Tagline             string
	Title               string
}

func Headers() []string {
	return []string{
		"id", "title", "release_date", "budget", "revenue", "status", "tagline",
		"runtime", "popularity", "original_language", "homepage", "adult", "overview",
		"languages", "production_countries", "production_companies", "genres",
	}
}

func (m *Movie) String() []string {
	fields := []string{
		strconv.Itoa(m.Id), m.Title, m.ReleaseDate, strconv.Itoa(m.Budget), strconv.Itoa(m.Revenue),
		m.Status, m.Tagline, strconv.Itoa(m.Runtime), strconv.FormatFloat(m.Popularity, 'E', -1, 64),
		m.OriginalLanguage, m.Homepage, strconv.FormatBool(m.Adult), m.Overview,
	}

	languages := []string{}
	for _, language := range m.Languages {
		languages = append(languages, language.Name)
	}
	fields = append(fields, strings.Join(languages, ","))

	prodCountries := []string{}
	for _, country := range m.ProductionCountries {
		prodCountries = append(prodCountries, country.Name)
	}
	fields = append(fields, strings.Join(prodCountries, ","))

	prodCompanies := []string{}
	for _, prod := range m.ProductionCompanies {
		prodCompanies = append(prodCompanies, prod.Name)
	}
	fields = append(fields, strings.Join(prodCompanies, ","))

	genres := []string{}
	for _, genre := range m.Genres {
		genres = append(genres, genre.Name)
	}
	fields = append(fields, strings.Join(genres, ","))

	return fields
}
