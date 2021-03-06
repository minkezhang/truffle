package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"

	dpb "github.com/minkezhang/truffle/api/go/database"
)

type TrackerT int
type AuxDataT int

const (
	AuxDataNone AuxDataT = iota
	AuxDataVideo
	AuxDataBook
	AuxDataGame
	AuxDataAudio

	TrackerNone TrackerT = iota
	TrackerVideo
	TrackerBook
)

var (
	AuxDataL = map[dpb.Corpus]AuxDataT{
		dpb.Corpus_CORPUS_TV:         AuxDataVideo,
		dpb.Corpus_CORPUS_ANIME:      AuxDataVideo,
		dpb.Corpus_CORPUS_FILM:       AuxDataVideo,
		dpb.Corpus_CORPUS_ANIME_FILM: AuxDataVideo,

		dpb.Corpus_CORPUS_GAME: AuxDataGame,

		dpb.Corpus_CORPUS_MANGA:       AuxDataBook,
		dpb.Corpus_CORPUS_BOOK:        AuxDataBook,
		dpb.Corpus_CORPUS_SHORT_STORY: AuxDataBook,

		dpb.Corpus_CORPUS_ALBUM: AuxDataAudio,
	}

	TrackerL = map[dpb.Corpus]TrackerT{
		dpb.Corpus_CORPUS_ANIME: TrackerVideo,
		dpb.Corpus_CORPUS_TV:    TrackerVideo,

		dpb.Corpus_CORPUS_MANGA: TrackerBook,
		dpb.Corpus_CORPUS_BOOK:  TrackerBook,
	}
)

func ETag(epb *dpb.Entry) ([]byte, error) {
	epb = proto.Clone(epb).(*dpb.Entry)

	epb.Id = nil
	epb.Etag = nil

	data, err := prototext.Marshal(epb)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot marshal input proto: %v", err)
	}

	s := md5.New()
	io.WriteString(s, string(data))

	// b64 string of the etag is easier to read.
	return []byte(
		base64.URLEncoding.EncodeToString(
			s.Sum(nil),
		),
	), nil
}

func ID(id *dpb.LinkedID) string {
	var parts []string
	_, api, _ := strings.Cut(dpb.API_name[int32(id.GetApi())], "_")
	if !map[dpb.API]bool{
		dpb.API_API_UNKNOWN: true,
		dpb.API_API_TRUFFLE: true,
	}[id.GetApi()] && api != "" {
		parts = append(parts, strings.ToLower(api))
	}
	parts = append(parts, id.GetId())
	return strings.Join(parts, ":")
}

func ToEnum(prefix string, suffix string) string {
	return fmt.Sprintf("%v_%v", prefix, strings.ReplaceAll(strings.ToUpper(suffix), " ", "_"))
}

func Clean(epb *dpb.Entry) *dpb.Entry {
	epb = proto.Clone(epb).(*dpb.Entry)

	lids := map[string]*dpb.LinkedID{}
	for _, lid := range epb.GetLinkedIds() {
		lids[ID(lid)] = lid
	}
	epb.LinkedIds = nil
	for _, lid := range lids {
		epb.LinkedIds = append(epb.GetLinkedIds(), lid)
	}

	providers := map[dpb.Provider]bool{}
	for _, p := range epb.GetProviders() {
		providers[p] = true
	}
	epb.Providers = nil
	for p, _ := range providers {
		epb.Providers = append(epb.GetProviders(), p)
	}

	f := func(candidates []string) []string {
		lookup := map[string]bool{}
		for _, c := range candidates {
			lookup[c] = true
		}
		var unique []string
		for c, _ := range lookup {
			unique = append(unique, c)
		}
		return unique
	}

	epb.Titles = f(epb.GetTitles())
	switch AuxDataL[epb.GetCorpus()] {
	case AuxDataVideo:
		epb.GetAuxDataVideo().Studios = f(epb.GetAuxDataVideo().GetStudios())
		epb.GetAuxDataVideo().Writers = f(epb.GetAuxDataVideo().GetWriters())
		epb.GetAuxDataVideo().Directors = f(epb.GetAuxDataVideo().GetDirectors())
	case AuxDataAudio:
		epb.GetAuxDataAudio().Composers = f(epb.GetAuxDataAudio().GetComposers())
	case AuxDataBook:
		epb.GetAuxDataBook().Authors = f(epb.GetAuxDataBook().GetAuthors())
	case AuxDataGame:
		epb.GetAuxDataGame().Studios = f(epb.GetAuxDataGame().GetStudios())
		epb.GetAuxDataGame().Studios = f(epb.GetAuxDataGame().GetStudios())
		epb.GetAuxDataGame().Writers = f(epb.GetAuxDataGame().GetWriters())
	}

	return epb
}
