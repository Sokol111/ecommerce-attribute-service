package command

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/attribute"
)

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
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attribute: %w", err)
	}

	if err := h.repo.Insert(ctx, a); err != nil {
		return nil, fmt.Errorf("failed to insert attribute: %w", err)
	}

	return a, nil
}
