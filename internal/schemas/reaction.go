package schemas

import "github.com/go-playground/validator/v10"

type SchemaReaction struct {
	Action   string `json:"action" binding:"required" validate:"required,oneof=add remove"`
	Reaction string `json:"reaction" binding:"required" validate:"required,oneof=like clap dislike heart"`
}

func (s SchemaReaction) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
