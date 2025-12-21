package http

import (
	"github.com/Sokol111/ecommerce-attribute-service-api/gen/httpapi"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func uuidPtrToStringPtr(u *openapi_types.UUID) *string {
	if u == nil {
		return nil
	}
	s := u.String()
	return &s
}

// combinedHandler combines attribute and category attribute handlers to implement the full StrictServerInterface
type combinedHandler struct {
	*attributeHandler
	*categoryAttributeHandler
}

func newCombinedHandler(
	attrHandler *attributeHandler,
	caHandler *categoryAttributeHandler,
) httpapi.StrictServerInterface {
	return &combinedHandler{
		attributeHandler:         attrHandler,
		categoryAttributeHandler: caHandler,
	}
}
