package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"path/filepath"
	"strings"
)

const (
	LanguageCtxKey = "language"
)

//go:embed lang
var lang embed.FS
var defaultLanguage = language.English

type ITranslationService interface {
	Translate(ctx context.Context, key string, param ...any) string
}

type service struct {
	translations map[language.Tag]map[string]string
}

func (s *service) Translate(ctx context.Context, key string, param ...any) string {
	acceptLanguage := ctx.Value(LanguageCtxKey)
	lang := defaultLanguage
	if acceptLanguage != nil {
		parsedLang, _, _ := language.ParseAcceptLanguage(acceptLanguage.(string))
		if len(parsedLang) > 0 {
			for _, tag := range parsedLang {
				_, ok := s.translations[tag]
				if !ok {
					continue
				}
				lang = tag
				break
			}
		} else {
			log.Warn().Msg("Failed to parse language")
		}
	} else {
		log.Warn().Msg("No language on context")
	}

	value, ok := s.translations[lang][key]
	if !ok {
		return key
	}
	if len(param) > 0 {
		switch v := param[0].(type) {
		case interface{}:
			if v == nil {
				return value
			}
		case []interface{}:
			if len(v) == 0 {
				return value
			}
		default:
			return fmt.Sprintf(value, param...)
		}
	}
	return value
}

func New() ITranslationService {
	s := &service{translations: map[language.Tag]map[string]string{}}
	dir, err := lang.ReadDir("lang")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load translation files")
	}
	for _, entry := range dir {
		if entry.IsDir() {
			log.Warn().Msgf("Ignoring directory lang/%s", entry.Name())
			continue
		}
		file, err := lang.Open(filepath.Join("lang", entry.Name()))
		if err != nil {
			log.Warn().Err(err).Msg("Failed to load translation file. Skipping")
			continue
		}
		name := strings.ReplaceAll(entry.Name(), filepath.Ext(entry.Name()), "")
		translations := map[string]string{}
		err = json.NewDecoder(file).Decode(&translations)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse translation file. Skipping")
			continue
		}
		lang, err := language.Parse(name)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse translation file name. Skipping")
			continue
		}
		s.translations[lang] = translations
	}
	return s
}
