package http

import (
	"context"
	"errors"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-attribute-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-attribute-service/internal/application/command"
	"github.com/Sokol111/ecommerce-attribute-service/internal/application/query"
	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/categoryattribute"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
)

type categoryAttributeHandler struct {
	assignHandler   command.AssignAttributeToCategoryCommandHandler
	updateHandler   command.UpdateCategoryAttributeCommandHandler
	unassignHandler command.UnassignAttributeFromCategoryCommandHandler
	getListHandler  query.GetCategoryAttributeListQueryHandler
}

func newCategoryAttributeHandler(
	assignHandler command.AssignAttributeToCategoryCommandHandler,
	updateHandler command.UpdateCategoryAttributeCommandHandler,
	unassignHandler command.UnassignAttributeFromCategoryCommandHandler,
	getListHandler query.GetCategoryAttributeListQueryHandler,
) *categoryAttributeHandler {
	return &categoryAttributeHandler{
		assignHandler:   assignHandler,
		updateHandler:   updateHandler,
		unassignHandler: unassignHandler,
		getListHandler:  getListHandler,
	}
}

func toCategoryAttributeResponse(ca *categoryattribute.CategoryAttribute) httpapi.CategoryAttributeResponse {
	return httpapi.CategoryAttributeResponse{
		Id:          ca.ID,
		Version:     ca.Version,
		CategoryId:  ca.CategoryID,
		AttributeId: ca.AttributeID,
		Required:    ca.Required,
		SortOrder:   ca.SortOrder,
		Filterable:  ca.Filterable,
		Searchable:  ca.Searchable,
		Enabled:     ca.Enabled,
		CreatedAt:   ca.CreatedAt,
		ModifiedAt:  ca.ModifiedAt,
	}
}

func (h *categoryAttributeHandler) AssignAttributeToCategory(c context.Context, request httpapi.AssignAttributeToCategoryRequestObject) (httpapi.AssignAttributeToCategoryResponseObject, error) {
	cmd := command.AssignAttributeToCategoryCommand{
		ID:          uuidPtrToStringPtr(request.Body.Id),
		CategoryID:  request.CategoryId,
		AttributeID: request.Body.AttributeId.String(),
		Required:    lo.FromPtrOr(request.Body.Required, false),
		SortOrder:   lo.FromPtrOr(request.Body.SortOrder, 0),
		Filterable:  request.Body.Filterable,
		Searchable:  request.Body.Searchable,
		Enabled:     request.Body.Enabled,
	}

	created, err := h.assignHandler.Handle(c, cmd)
	if err != nil {
		if errors.Is(err, categoryattribute.ErrAlreadyAssigned) {
			return httpapi.AssignAttributeToCategory409ApplicationProblemPlusJSONResponse{
				Status: 409,
				Type:   "about:blank",
				Title:  "Attribute is already assigned to this category",
			}, nil
		}
		if errors.Is(err, persistence.ErrEntityNotFound) {
			return httpapi.AssignAttributeToCategory404ApplicationProblemPlusJSONResponse{
				Status: 404,
				Type:   "about:blank",
				Title:  "Attribute not found",
			}, nil
		}
		return nil, err
	}

	return httpapi.AssignAttributeToCategory200JSONResponse(toCategoryAttributeResponse(created)), nil
}

func (h *categoryAttributeHandler) UpdateCategoryAttribute(c context.Context, request httpapi.UpdateCategoryAttributeRequestObject) (httpapi.UpdateCategoryAttributeResponseObject, error) {
	cmd := command.UpdateCategoryAttributeCommand{
		ID:         request.Body.Id.String(),
		CategoryID: request.CategoryId,
		Version:    request.Body.Version,
		Required:   lo.FromPtrOr(request.Body.Required, false),
		SortOrder:  lo.FromPtrOr(request.Body.SortOrder, 0),
		Filterable: request.Body.Filterable,
		Searchable: request.Body.Searchable,
		Enabled:    request.Body.Enabled,
	}

	updated, err := h.updateHandler.Handle(c, cmd)
	if err != nil {
		if errors.Is(err, persistence.ErrEntityNotFound) {
			return httpapi.UpdateCategoryAttribute404ApplicationProblemPlusJSONResponse{
				Status: 404,
				Type:   "about:blank",
				Title:  "Category attribute assignment not found",
			}, nil
		}
		if errors.Is(err, persistence.ErrOptimisticLocking) {
			return httpapi.UpdateCategoryAttribute412ApplicationProblemPlusJSONResponse{
				Status: 412,
				Type:   "about:blank",
				Title:  "Version mismatch",
			}, nil
		}
		return nil, err
	}

	return httpapi.UpdateCategoryAttribute200JSONResponse(toCategoryAttributeResponse(updated)), nil
}

func (h *categoryAttributeHandler) GetCategoryAttributeList(c context.Context, request httpapi.GetCategoryAttributeListRequestObject) (httpapi.GetCategoryAttributeListResponseObject, error) {
	sort := "sortOrder"
	order := "asc"

	if request.Params.Sort != nil {
		sort = string(*request.Params.Sort)
	}
	if request.Params.Order != nil {
		order = string(*request.Params.Order)
	}

	q := query.GetCategoryAttributeListQuery{
		CategoryID: request.CategoryId,
		Page:       request.Params.Page,
		Size:       request.Params.Size,
		Enabled:    request.Params.Enabled,
		Filterable: request.Params.Filterable,
		Sort:       sort,
		Order:      order,
	}

	result, err := h.getListHandler.Handle(c, q)
	if err != nil {
		return nil, err
	}

	response := httpapi.GetCategoryAttributeList200JSONResponse{
		Items: lo.Map(result.Items, func(ca *categoryattribute.CategoryAttribute, _ int) httpapi.CategoryAttributeResponse {
			return toCategoryAttributeResponse(ca)
		}),
		Page:  result.Page,
		Size:  result.Size,
		Total: int(result.Total),
	}

	return response, nil
}

func (h *categoryAttributeHandler) UnassignAttributeFromCategory(c context.Context, request httpapi.UnassignAttributeFromCategoryRequestObject) (httpapi.UnassignAttributeFromCategoryResponseObject, error) {
	cmd := command.UnassignAttributeFromCategoryCommand{
		ID:         request.AssignmentId,
		CategoryID: request.CategoryId,
	}

	err := h.unassignHandler.Handle(c, cmd)
	if err != nil {
		if errors.Is(err, persistence.ErrEntityNotFound) {
			return httpapi.UnassignAttributeFromCategory404ApplicationProblemPlusJSONResponse{
				Status: 404,
				Type:   "about:blank",
				Title:  "Category attribute assignment not found",
			}, nil
		}
		return nil, err
	}

	return httpapi.UnassignAttributeFromCategory204Response{}, nil
}
