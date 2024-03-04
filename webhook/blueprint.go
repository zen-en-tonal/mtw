package webhook

import (
	"github.com/google/uuid"
)

type Blueprint struct {
	ID          string
	Endpoint    string
	Method      string
	Auth        string
	Schema      string
	ContentType string
}

func (b Blueprint) options(defaults ...Option) (*[]Option, error) {
	var options []Option

	options = append(options, defaults...)

	if b.ID != "" {
		id, err := uuid.Parse(b.ID)
		if err != nil {
			return nil, err
		}
		options = append(options, WithID(id))
	}

	if b.Schema != "" {
		opt, err := WithSchema(b.Schema, b.ContentType)
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}

	if b.Auth != "" {
		options = append(options, WithAuth(b.Auth))
	}

	options = append(options, WithMethod(b.Method))

	return &options, nil
}
