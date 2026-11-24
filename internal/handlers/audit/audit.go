package audit

import (
	"context"
	"encoding/json"
	"net/http"

	auditentity "github.com/Dmitro44/relper/internal/entities/audit"
	apperr "github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	auditsvc "github.com/Dmitro44/relper/internal/services/audit"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type AuditHandler struct {
	svc *auditsvc.AuditService
}

func NewAuditHandler(svc *auditsvc.AuditService) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) Register(admin huma.API) {
	huma.Register(admin, huma.Operation{
		OperationID:   "list-audit",
		Method:        http.MethodGet,
		Path:          "/audit",
		Summary:       "List audit events",
		Tags:          []string{"Audit"},
		DefaultStatus: http.StatusOK,
	}, h.HandleList)
}

func (h *AuditHandler) HandleList(ctx context.Context, in *ListAuditRequest) (*AuditResponse, error) {
	if in.Source != "" && in.Source != "app" && in.Source != "trigger" {
		return nil, httperr.Map(apperr.NewBadRequestError("source must be app or trigger"))
	}

	filter := auditentity.Filter{
		Source:     in.Source,
		Action:     in.Action,
		EntityType: in.EntityType,
		Limit:      in.Limit,
	}
	if in.ActorID != uuid.Nil {
		filter.ActorID = &in.ActorID
	}

	events, err := h.svc.List(ctx, filter)
	if err != nil {
		return nil, httperr.Map(err)
	}

	items := make([]AuditItem, 0, len(events))
	for _, event := range events {
		items = append(items, toAuditItem(event))
	}

	return &AuditResponse{Body: items}, nil
}

func toAuditItem(event auditentity.Event) AuditItem {
	var details any
	if len(event.Details) > 0 {
		_ = json.Unmarshal(event.Details, &details)
	}

	return AuditItem{
		ID:            event.ID,
		Source:        event.Source,
		Action:        event.Action,
		EntityType:    event.EntityType,
		EntityID:      event.EntityID,
		ActorID:       event.ActorID,
		ActorNickname: event.ActorNickname,
		Details:       details,
		CreatedAt:     event.CreatedAt,
	}
}
