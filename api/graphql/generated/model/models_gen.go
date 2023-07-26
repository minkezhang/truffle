// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Aux interface {
	IsAux()
}

type AuxAlbum struct {
	Studio    []string `json:"studio,omitempty"`
	Composers []string `json:"composers,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Producers []string `json:"producers,omitempty"`
}

func (AuxAlbum) IsAux() {}

type AuxAnime struct {
	Composors []string `json:"composors,omitempty"`
	Directors []string `json:"directors,omitempty"`
	Studios   []string `json:"studios,omitempty"`
	Writers   []string `json:"writers,omitempty"`
}

func (AuxAnime) IsAux() {}

type AuxAnimeFilm struct {
	Composors []string `json:"composors,omitempty"`
	Directors []string `json:"directors,omitempty"`
	Studios   []string `json:"studios,omitempty"`
	Writers   []string `json:"writers,omitempty"`
}

func (AuxAnimeFilm) IsAux() {}

type AuxBook struct {
	Authors []string `json:"authors,omitempty"`
}

func (AuxBook) IsAux() {}

type AuxFilm struct {
	Cinematographers []string `json:"cinematographers,omitempty"`
	Composors        []string `json:"composors,omitempty"`
	Directors        []string `json:"directors,omitempty"`
	Editors          []string `json:"editors,omitempty"`
	Writers          []string `json:"writers,omitempty"`
}

func (AuxFilm) IsAux() {}

type AuxGame struct {
	Artists     []string `json:"artists,omitempty"`
	Composers   []string `json:"composers,omitempty"`
	Designers   []string `json:"designers,omitempty"`
	Developers  []string `json:"developers,omitempty"`
	Directors   []string `json:"directors,omitempty"`
	Programmers []string `json:"programmers,omitempty"`
	Publishers  []string `json:"publishers,omitempty"`
	Writers     []string `json:"writers,omitempty"`
}

func (AuxGame) IsAux() {}

type AuxManga struct {
	Illustrators []string `json:"illustrators,omitempty"`
	Writers      []string `json:"writers,omitempty"`
}

func (AuxManga) IsAux() {}

type AuxShortStory struct {
	Authors []string `json:"authors,omitempty"`
}

func (AuxShortStory) IsAux() {}

type AuxTv struct {
	Cinematographers []string `json:"cinematographers,omitempty"`
	Creators         []string `json:"creators,omitempty"`
	Composers        []string `json:"composers,omitempty"`
	Editors          []string `json:"editors,omitempty"`
	Writers          []string `json:"writers,omitempty"`
}

func (AuxTv) IsAux() {}

type Entry struct {
	ID       string      `json:"id"`
	Metadata []*Metadata `json:"metadata"`
	Corpus   CorpusType  `json:"corpus"`
}

type Metadata struct {
	API       APIType        `json:"api"`
	ID        string         `json:"id"`
	Titles    []*Title       `json:"titles"`
	Score     *float64       `json:"score,omitempty"`
	Providers []ProviderType `json:"providers,omitempty"`
	Aux       Aux            `json:"aux"`
	Tags      []string       `json:"tags,omitempty"`
}

type NewTodo struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type Title struct {
	Language string `json:"language"`
	Title    string `json:"title"`
}

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
	User *User  `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type APIType string

const (
	APITypeAPINone     APIType = "API_NONE"
	APITypeAPITruffle  APIType = "API_TRUFFLE"
	APITypeAPIMal      APIType = "API_MAL"
	APITypeAPIKitsuAPI APIType = "API_KITSU_API"
	APITypeAPISteam    APIType = "API_STEAM"
	APITypeAPISpotify  APIType = "API_SPOTIFY"
)

var AllAPIType = []APIType{
	APITypeAPINone,
	APITypeAPITruffle,
	APITypeAPIMal,
	APITypeAPIKitsuAPI,
	APITypeAPISteam,
	APITypeAPISpotify,
}

func (e APIType) IsValid() bool {
	switch e {
	case APITypeAPINone, APITypeAPITruffle, APITypeAPIMal, APITypeAPIKitsuAPI, APITypeAPISteam, APITypeAPISpotify:
		return true
	}
	return false
}

func (e APIType) String() string {
	return string(e)
}

func (e *APIType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = APIType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid APIType", str)
	}
	return nil
}

func (e APIType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type CorpusType string

const (
	CorpusTypeCorpusNone       CorpusType = "CORPUS_NONE"
	CorpusTypeCorpusAnime      CorpusType = "CORPUS_ANIME"
	CorpusTypeCorpusAnimeFilm  CorpusType = "CORPUS_ANIME_FILM"
	CorpusTypeCorpusManga      CorpusType = "CORPUS_MANGA"
	CorpusTypeCorpusBook       CorpusType = "CORPUS_BOOK"
	CorpusTypeCorpusShortStory CorpusType = "CORPUS_SHORT_STORY"
	CorpusTypeCorpusFilm       CorpusType = "CORPUS_FILM"
	CorpusTypeCorpusTv         CorpusType = "CORPUS_TV"
	CorpusTypeCorpusAlbum      CorpusType = "CORPUS_ALBUM"
	CorpusTypeCorpusGame       CorpusType = "CORPUS_GAME"
)

var AllCorpusType = []CorpusType{
	CorpusTypeCorpusNone,
	CorpusTypeCorpusAnime,
	CorpusTypeCorpusAnimeFilm,
	CorpusTypeCorpusManga,
	CorpusTypeCorpusBook,
	CorpusTypeCorpusShortStory,
	CorpusTypeCorpusFilm,
	CorpusTypeCorpusTv,
	CorpusTypeCorpusAlbum,
	CorpusTypeCorpusGame,
}

func (e CorpusType) IsValid() bool {
	switch e {
	case CorpusTypeCorpusNone, CorpusTypeCorpusAnime, CorpusTypeCorpusAnimeFilm, CorpusTypeCorpusManga, CorpusTypeCorpusBook, CorpusTypeCorpusShortStory, CorpusTypeCorpusFilm, CorpusTypeCorpusTv, CorpusTypeCorpusAlbum, CorpusTypeCorpusGame:
		return true
	}
	return false
}

func (e CorpusType) String() string {
	return string(e)
}

func (e *CorpusType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CorpusType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CorpusType", str)
	}
	return nil
}

func (e CorpusType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProviderType string

const (
	ProviderTypeProviderNone            ProviderType = "PROVIDER_NONE"
	ProviderTypeProviderOther           ProviderType = "PROVIDER_OTHER"
	ProviderTypeProviderCrunchyroll     ProviderType = "PROVIDER_CRUNCHYROLL"
	ProviderTypeProviderDisneyPlus      ProviderType = "PROVIDER_DISNEY_PLUS"
	ProviderTypeProviderHulu            ProviderType = "PROVIDER_HULU"
	ProviderTypeProviderHboMax          ProviderType = "PROVIDER_HBO_MAX"
	ProviderTypeProviderNetflix         ProviderType = "PROVIDER_NETFLIX"
	ProviderTypeProviderAmazonStreaming ProviderType = "PROVIDER_AMAZON_STREAMING"
	ProviderTypeProviderGooglePlay      ProviderType = "PROVIDER_GOOGLE_PLAY"
	ProviderTypeProviderAppleTv         ProviderType = "PROVIDER_APPLE_TV"
	ProviderTypeProviderSteam           ProviderType = "PROVIDER_STEAM"
)

var AllProviderType = []ProviderType{
	ProviderTypeProviderNone,
	ProviderTypeProviderOther,
	ProviderTypeProviderCrunchyroll,
	ProviderTypeProviderDisneyPlus,
	ProviderTypeProviderHulu,
	ProviderTypeProviderHboMax,
	ProviderTypeProviderNetflix,
	ProviderTypeProviderAmazonStreaming,
	ProviderTypeProviderGooglePlay,
	ProviderTypeProviderAppleTv,
	ProviderTypeProviderSteam,
}

func (e ProviderType) IsValid() bool {
	switch e {
	case ProviderTypeProviderNone, ProviderTypeProviderOther, ProviderTypeProviderCrunchyroll, ProviderTypeProviderDisneyPlus, ProviderTypeProviderHulu, ProviderTypeProviderHboMax, ProviderTypeProviderNetflix, ProviderTypeProviderAmazonStreaming, ProviderTypeProviderGooglePlay, ProviderTypeProviderAppleTv, ProviderTypeProviderSteam:
		return true
	}
	return false
}

func (e ProviderType) String() string {
	return string(e)
}

func (e *ProviderType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProviderType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProviderType", str)
	}
	return nil
}

func (e ProviderType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
