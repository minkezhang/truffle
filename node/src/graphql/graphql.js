"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ApiDataPartsFragmentDoc = exports.ProviderType = exports.CorpusType = exports.ApiType = void 0;
var ApiType;
(function (ApiType) {
    ApiType["ApiKitsu"] = "API_KITSU";
    ApiType["ApiMal"] = "API_MAL";
    ApiType["ApiNone"] = "API_NONE";
    ApiType["ApiSpotify"] = "API_SPOTIFY";
    ApiType["ApiSteam"] = "API_STEAM";
    ApiType["ApiTruffle"] = "API_TRUFFLE";
})(ApiType || (exports.ApiType = ApiType = {}));
var CorpusType;
(function (CorpusType) {
    CorpusType["CorpusAlbum"] = "CORPUS_ALBUM";
    CorpusType["CorpusAnime"] = "CORPUS_ANIME";
    CorpusType["CorpusAnimeFilm"] = "CORPUS_ANIME_FILM";
    CorpusType["CorpusBook"] = "CORPUS_BOOK";
    CorpusType["CorpusFilm"] = "CORPUS_FILM";
    CorpusType["CorpusGame"] = "CORPUS_GAME";
    CorpusType["CorpusManga"] = "CORPUS_MANGA";
    CorpusType["CorpusNone"] = "CORPUS_NONE";
    CorpusType["CorpusShortStory"] = "CORPUS_SHORT_STORY";
    CorpusType["CorpusTv"] = "CORPUS_TV";
})(CorpusType || (exports.CorpusType = CorpusType = {}));
var ProviderType;
(function (ProviderType) {
    ProviderType["ProviderAmazonStreaming"] = "PROVIDER_AMAZON_STREAMING";
    ProviderType["ProviderAppleTv"] = "PROVIDER_APPLE_TV";
    ProviderType["ProviderCrunchyroll"] = "PROVIDER_CRUNCHYROLL";
    ProviderType["ProviderDisneyPlus"] = "PROVIDER_DISNEY_PLUS";
    ProviderType["ProviderGooglePlay"] = "PROVIDER_GOOGLE_PLAY";
    ProviderType["ProviderHboMax"] = "PROVIDER_HBO_MAX";
    ProviderType["ProviderHulu"] = "PROVIDER_HULU";
    ProviderType["ProviderNetflix"] = "PROVIDER_NETFLIX";
    ProviderType["ProviderNone"] = "PROVIDER_NONE";
    ProviderType["ProviderOther"] = "PROVIDER_OTHER";
    ProviderType["ProviderSteam"] = "PROVIDER_STEAM";
})(ProviderType || (exports.ProviderType = ProviderType = {}));
exports.ApiDataPartsFragmentDoc = { "kind": "Document", "definitions": [{ "kind": "FragmentDefinition", "name": { "kind": "Name", "value": "APIDataParts" }, "typeCondition": { "kind": "NamedType", "name": { "kind": "Name", "value": "APIData" } }, "selectionSet": { "kind": "SelectionSet", "selections": [{ "kind": "Field", "name": { "kind": "Name", "value": "api" } }, { "kind": "Field", "name": { "kind": "Name", "value": "id" } }, { "kind": "Field", "name": { "kind": "Name", "value": "corpus" } }, { "kind": "Field", "name": { "kind": "Name", "value": "titles" }, "selectionSet": { "kind": "SelectionSet", "selections": [{ "kind": "Field", "name": { "kind": "Name", "value": "title" } }] } }, { "kind": "Field", "name": { "kind": "Name", "value": "providers" } }, { "kind": "Field", "name": { "kind": "Name", "value": "tags" } }, { "kind": "Field", "name": { "kind": "Name", "value": "aux" }, "selectionSet": { "kind": "SelectionSet", "selections": [{ "kind": "InlineFragment", "typeCondition": { "kind": "NamedType", "name": { "kind": "Name", "value": "AuxAnime" } }, "selectionSet": { "kind": "SelectionSet", "selections": [{ "kind": "Field", "name": { "kind": "Name", "value": "studios" } }] } }, { "kind": "InlineFragment", "typeCondition": { "kind": "NamedType", "name": { "kind": "Name", "value": "AuxManga" } }, "selectionSet": { "kind": "SelectionSet", "selections": [{ "kind": "Field", "name": { "kind": "Name", "value": "authors" } }, { "kind": "Field", "name": { "kind": "Name", "value": "artists" } }] } }] } }, { "kind": "Field", "name": { "kind": "Name", "value": "tracker" }, "selectionSet": { "kind": "SelectionSet", "selections": [{ "kind": "InlineFragment", "typeCondition": { "kind": "NamedType", "name": { "kind": "Name", "value": "TrackerAnime" } }, "selectionSet": { "kind": "SelectionSet", "selections": [{ "kind": "Field", "name": { "kind": "Name", "value": "season" } }, { "kind": "Field", "name": { "kind": "Name", "value": "episode" } }, { "kind": "Field", "name": { "kind": "Name", "value": "last_updated" } }] } }, { "kind": "InlineFragment", "typeCondition": { "kind": "NamedType", "name": { "kind": "Name", "value": "TrackerManga" } }, "selectionSet": { "kind": "SelectionSet", "selections": [{ "kind": "Field", "name": { "kind": "Name", "value": "volume" } }, { "kind": "Field", "name": { "kind": "Name", "value": "chapter" } }, { "kind": "Field", "name": { "kind": "Name", "value": "last_updated" } }] } }] } }] } }] };