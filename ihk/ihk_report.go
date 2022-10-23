package ihk

import (
	"bytes"
	"mime/multipart"
	"time"
)

type Report struct {
	Token           string
	StartDate       time.Time
	EndDate         time.Time
	InstructorEmail string
	ContentActivity string
	ContentSubjects string
	ContentTraining string
	Department      string
}

func (r *Report) ToMultipartFormData(action string) (*bytes.Buffer, string, error) {
	buffer := &bytes.Buffer{}
	w := multipart.NewWriter(buffer)
	if err := w.WriteField("token", r.Token); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("lfdnr", "0"); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("edtvon", r.StartDate.Format("02.01.2006")); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("edtbis", r.EndDate.Format("02.01.2006")); err != nil {
		return nil, "", err
	}
	// Duplicated fields here are intentional
	if err := w.WriteField("ausbabschnitt", r.Department); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("ausbabschnitt", r.Department); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("ausbMail", r.InstructorEmail); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("ausbMail2", r.InstructorEmail); err != nil {
		return nil, "", err
	}

	if err := w.WriteField("ausbinhalt1", r.ContentActivity); err != nil {
		return nil, "", err
	}

	// Not generally relevant but need to send them
	if err := w.WriteField("stdMo", "0"); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("stdDi", "0"); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("stdMi", "0"); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("stdDo", "0"); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("stdFr", "0"); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("stdSa", "0"); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("stdSo", "0"); err != nil {
		return nil, "", err
	}

	// Real content
	if err := w.WriteField("ausbinhalt2", r.ContentActivity); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("ausbinhalt12", "null"); err != nil {
		return nil, "", err
	}

	if err := w.WriteField("ausbinhalt3", r.ContentActivity); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("ausbinhalt13", "null"); err != nil {
		return nil, "", err
	}

	if w, err := w.CreateFormFile("file", ""); err != nil {
		_, _ = w.Write(nil)
	}

	if err := w.WriteField(action, ""); err != nil {
		return nil, "", err
	}

	_ = w.Close()

	return buffer, w.Boundary(), nil
}
