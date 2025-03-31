package schemas

import "github.com/go-playground/validator/v10"

type Attachment struct {
	FileName string `json:"name" validate:"required,min=3,max=100"`
	FileType string `json:"type" validate:"required,min=3,max=100"`
	FileSize int64  `json:"size" validate:"required"`
	FileUrl  string `json:"url" validate:"required,url"`
}

func (s *Attachment) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type Attachments []Attachment

type SelectAttachment struct {
	ID       uint   `json:"id"`
	FileUrl  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}
