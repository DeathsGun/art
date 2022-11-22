package report

import (
	"context"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/i18n"
)

type Category int

const (
	Activity Category = iota
	Training
	Subjects
)

func (c Category) Text(ctx context.Context) string {
	service := di.Instance[i18n.ITranslationService]("i18n")
	switch c {
	case Activity:
		return service.Translate(ctx, "CATEGORY_ACTIVITY")
	case Subjects:
		return service.Translate(ctx, "CATEGORY_SUBJECTS")
	case Training:
		return service.Translate(ctx, "CATEGORY_TRAINING")
	}
	return ""
}
