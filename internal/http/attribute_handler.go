package http

import (
	"context"
	"errors"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-attribute-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-attribute-service/internal/application/command"
	"github.com/Sokol111/ecommerce-attribute-service/internal/application/query"
	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/attribute"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
)

type attributeHandler struct {
	createHandler  command.CreateAttributeCommandHandler
	updateHandler  command.UpdateAttributeCommandHandler
	getByIDHandler query.GetAttributeByIDQueryHandler
	getListHandler query.GetAttributeListQueryHandler
}

func newAttributeHandler(
	createHandler command.CreateAttributeCommandHandler,
	updateHandler command.UpdateAttributeCommandHandler,
	getByIDHandler query.GetAttributeByIDQueryHandler,
	getListHandler query.GetAttributeListQueryHandler,
) *attributeHandler {
	return &attributeHandler{
		createHandler:  createHandler,
		updateHandler:  updateHandler,
		getByIDHandler: getByIDHandler,
		getListHandler: getListHandler,
	}
}

func toAttributeResponse(a *attribute.Attribute) httpapi.AttributeResponse {
	var options *[]httpapi.AttributeOption
	if len(a.Options) > 0 {
		opts := lo.Map(a.Options, func(opt attribute.Option, _ int) httpapi.AttributeOption {
			return httpapi.AttributeOption{
				Value:     opt.Value,
				Slug:      opt.Slug,
				ColorCode: opt.ColorCode,
				SortOrder: opt.SortOrder,
				Enabled:   opt.Enabled,
			}
		})
		options = &opts
	}

	return httpapi.AttributeResponse{
		Id:                a.ID,
		Version:           a.Version,
		Name:              a.Name,
		Slug:              a.Slug,
		Type:              httpapi.AttributeResponseType(a.Type),
		Unit:              a.Unit,
		DefaultFilterable: a.DefaultFilterable,
		DefaultSearchable: a.DefaultSearchable,
		SortOrder:         a.SortOrder,
		Enabled:           a.Enabled,
		Options:           options,
		CreatedAt:         a.CreatedAt,
		ModifiedAt:        a.ModifiedAt,
	}
}

func (h *attributeHandler) CreateAttribute(c context.Context, request httpapi.CreateAttributeRequestObject) (httpapi.CreateAttributeResponseObject, error) {
	options := lo.Map(lo.FromPtr(request.Body.Options), func(opt httpapi.AttributeOptionInput, _ int) command.OptionInput {
		return command.OptionInput{
			Value:     opt.Value,
			Slug:      opt.Slug,
			ColorCode: opt.ColorCode,
			SortOrder: lo.FromPtrOr(opt.SortOrder, 0),
			Enabled:   opt.Enabled,
		}
	})

	cmd := command.CreateAttributeCommand{
		ID:                uuidPtrToStringPtr(request.Body.Id),
		Name:              request.Body.Name,
		Slug:              request.Body.Slug,
		Type:              string(request.Body.Type),
		Unit:              request.Body.Unit,
		DefaultFilterable: request.Body.DefaultFilterable,
		DefaultSearchable: request.Body.DefaultSearchable,
		SortOrder:         lo.FromPtrOr(request.Body.SortOrder, 0),
		Enabled:           request.Body.Enabled,
		Options:           options,
	}

	created, err := h.createHandler.Handle(c, cmd)
	if err != nil {
		if errors.Is(err, attribute.ErrSlugAlreadyExists) {
			return httpapi.CreateAttribute409ApplicationProblemPlusJSONResponse{
				Status: 409,
				Type:   "about:blank",
				Title:  "Attribute with this slug already exists",
			}, nil
		}
		return nil, err
	}

	return httpapi.CreateAttribute200JSONResponse(toAttributeResponse(created)), nil
}

func (h *attributeHandler) GetAttributeById(c context.Context, request httpapi.GetAttributeByIdRequestObject) (httpapi.GetAttributeByIdResponseObject, error) {
	q := query.GetAttributeByIDQuery{ID: request.Id}

	found, err := h.getByIDHandler.Handle(c, q)
	if errors.Is(err, persistence.ErrEntityNotFound) {
		return httpapi.GetAttributeById404ApplicationProblemPlusJSONResponse{
			Status: 404,
			Type:   "about:blank",
			Title:  "Attribute not found",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return httpapi.GetAttributeById200JSONResponse(toAttributeResponse(found)), nil
}

func (h *attributeHandler) GetAttributeList(c context.Context, request httpapi.GetAttributeListRequestObject) (httpapi.GetAttributeListResponseObject, error) {
	// Default sort and order
	sort := "sortOrder"
	order := "asc"

	// Override with request params if provided
	if request.Params.Sort != nil {
		sort = string(*request.Params.Sort)
	}

	if request.Params.Order != nil {
		order = string(*request.Params.Order)
	}

	var attrType *string
	if request.Params.Type != nil {
		t := string(*request.Params.Type)
		attrType = &t
	}

	q := query.GetAttributeListQuery{
		Page:    request.Params.Page,
		Size:    request.Params.Size,
		Enabled: request.Params.Enabled,
		Type:    attrType,
		Sort:    sort,
		Order:   order,
	}

	result, err := h.getListHandler.Handle(c, q)
	if err != nil {
		return nil, err
	}

	response := httpapi.GetAttributeList200JSONResponse{
		Items: make([]httpapi.AttributeResponse, 0, len(result.Items)),
		Page:  result.Page,
		Size:  result.Size,
		Total: int(result.Total),
	}

	for _, a := range result.Items {
		response.Items = append(response.Items, toAttributeResponse(a))
	}

	return response, nil
}

func (h *attributeHandler) UpdateAttribute(c context.Context, request httpapi.UpdateAttributeRequestObject) (httpapi.UpdateAttributeResponseObject, error) {
	options := lo.Map(lo.FromPtr(request.Body.Options), func(opt httpapi.AttributeOptionInput, _ int) command.OptionInput {
		return command.OptionInput{
			Value:     opt.Value,
			Slug:      opt.Slug,
			ColorCode: opt.ColorCode,
			SortOrder: lo.FromPtrOr(opt.SortOrder, 0),
			Enabled:   opt.Enabled,
		}
	})

	cmd := command.UpdateAttributeCommand{
		ID:                request.Body.Id.String(),
		Version:           request.Body.Version,
		Name:              request.Body.Name,
		Slug:              request.Body.Slug,
		Type:              string(request.Body.Type),
		Unit:              request.Body.Unit,
		DefaultFilterable: request.Body.DefaultFilterable,
		DefaultSearchable: request.Body.DefaultSearchable,
		SortOrder:         lo.FromPtrOr(request.Body.SortOrder, 0),
		Enabled:           request.Body.Enabled,
		Options:           options,
	}

	updated, err := h.updateHandler.Handle(c, cmd)
	if err != nil {
		if errors.Is(err, persistence.ErrEntityNotFound) {
			return httpapi.UpdateAttribute404ApplicationProblemPlusJSONResponse{
				Status: 404,
				Type:   "about:blank",
				Title:  "Attribute not found",
			}, nil
		}
		if errors.Is(err, persistence.ErrOptimisticLocking) {
			return httpapi.UpdateAttribute412ApplicationProblemPlusJSONResponse{
				Status: 412,
				Type:   "about:blank",
				Title:  "Version mismatch",
			}, nil
		}
		if errors.Is(err, attribute.ErrSlugAlreadyExists) {
			return httpapi.UpdateAttribute409ApplicationProblemPlusJSONResponse{
				Status: 409,
				Type:   "about:blank",
				Title:  "Attribute with this slug already exists",
			}, nil
		}
		return nil, err
	}

	return httpapi.UpdateAttribute200JSONResponse(toAttributeResponse(updated)), nil
}

// Stub implementation for delete (not implemented yet)

func (h *attributeHandler) DeleteAttribute(c context.Context, request httpapi.DeleteAttributeRequestObject) (httpapi.DeleteAttributeResponseObject, error) {
	return httpapi.DeleteAttribute500ApplicationProblemPlusJSONResponse{
		Status: 500,
		Type:   "about:blank",
		Title:  "Not implemented",
	}, nil
}
