package command

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/attribute"
)

type OptionInput struct {
	Value     string
	Slug      string
	ColorCode *string
	SortOrder int
	Enabled   bool
}

type CreateAttributeCommand struct {
	ID                *string
	Name              string
	Slug              string
	Type              string
	Unit              *string
	DefaultFilterable bool
	DefaultSearchable bool
	SortOrder         int
	Enabled           bool
	Options           []OptionInput
}

type CreateAttributeCommandHandler interface {
	Handle(ctx context.Context, cmd CreateAttributeCommand) (*attribute.Attribute, error)
}

type createAttributeHandler struct {
	repo attribute.Repository
}

func NewCreateAttributeHandler(repo attribute.Repository) CreateAttributeCommandHandler {
	return &createAttributeHandler{
		repo: repo,
	}
}

func (h *createAttributeHandler) Handle(ctx context.Context, cmd CreateAttributeCommand) (*attribute.Attribute, error) {
	options := lo.Map(cmd.Options, func(opt OptionInput, _ int) attribute.Option {
		return attribute.Option{
			Value:     opt.Value,
			Slug:      opt.Slug,
			ColorCode: opt.ColorCode,
			SortOrder: opt.SortOrder,
			Enabled:   opt.Enabled,
		}
	})

	a, err := attribute.NewAttribute(
		lo.FromPtrOr(cmd.ID, ""),
		cmd.Name,
		cmd.Slug,
		attribute.AttributeType(cmd.Type),
		cmd.Unit,
		cmd.DefaultFilterable,
		cmd.DefaultSearchable,
		cmd.SortOrder,
		cmd.Enabled,
		options,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attribute: %w", err)
	}

	if err := h.repo.Insert(ctx, a); err != nil {
		return nil, fmt.Errorf("failed to insert attribute: %w", err)
	}

	return a, nil
}
