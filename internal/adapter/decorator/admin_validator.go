package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgarSucre/chat/internal/adapter"
	"github.com/edgarSucre/chat/internal/domain"
	"github.com/go-playground/validator/v10"
)

type (
	AdminValidator struct {
		uc        adapter.AdminUseCase
		validator *validator.Validate
	}

	validation map[string]string
)

var (
	tags validation = map[string]string{
		"required": "is required",
		"gt":       "must be greater than",
		"gte":      "must be greater or equal to",
	}
)

func NewAdminValidator(uc adapter.AdminUseCase) adapter.AdminUseCase {
	return &AdminValidator{
		uc:        uc,
		validator: validator.New(),
	}
}

func (av *AdminValidator) CreateUser(ctx context.Context, params domain.UserParam) (domain.UserResponse, error) {

	err := av.validator.Struct(params)

	if err != nil {
		return domain.UserResponse{}, getValidationError(err)
	}

	return av.uc.CreateUser(ctx, params)
}

func getValidationError(err error) *domain.Err {

	if orig, ok := err.(validator.ValidationErrors); ok {
		msg := getMessage(orig)

		return domain.WrapErrorf(
			err,
			domain.ErrorCodeInvalidParams,
			msg,
		)
	}

	return domain.WrapErrorf(
		err,
		domain.ErrorCodeInvalidParams,
		"invalid aarams",
	)
}

func getMessage(errors []validator.FieldError) (msg string) {

	for _, er := range errors {
		tag := tags[er.Tag()]
		field := er.StructField()
		param := er.Param()

		if len(msg) > 0 {
			msg = strings.TrimSpace(fmt.Sprintf("%s, %s %s %s", msg, field, tag, param))
			continue
		}
		msg = strings.TrimSpace(fmt.Sprintf("%s %s %s", field, tag, param))
	}

	return msg
}
