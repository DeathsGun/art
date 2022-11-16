package text

import (
	"bytes"
	"context"
	"github.com/deathsgun/art/config/model"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/report"
	"time"
)

type impl struct {
}

func (i *impl) ContentType() string {
	return "text/plain; charset=utf-8"
}

func (i *impl) GetStartDate(_ context.Context) (time.Time, error) {
	return time.Now(), nil
}

func (i *impl) Id() string {
	return "PROVIDER_TEXT"
}

func (i *impl) Logo() string {
	return "file-text.svg"
}

func (i *impl) Capabilities() []provider.Capability {
	return []provider.Capability{
		provider.TypeExport,
	}
}

func (i *impl) ValidateConfig(_ context.Context, _ *model.ProviderConfig) error {
	return nil
}

func (i *impl) Export(ctx context.Context, report *report.Report) ([]byte, error) {
	buffer := &bytes.Buffer{}
	buffers, err := report.Format(ctx, true)
	if err != nil {
		return nil, err
	}
	for category, text := range buffers {
		_, err = buffer.WriteString(category.Text(ctx) + ":\n")
		if err != nil {
			return nil, err
		}
		if len(text) == 0 {
			continue
		}
		_, err = buffer.Write(text)
		if err != nil {
			return nil, err
		}
		_, err = buffer.WriteString("\n")
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func NewProvider() provider.ExportProvider {
	return &impl{}
}
