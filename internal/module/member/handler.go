package member

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

// Handler routes member and invitation management under
// /tenants/{tenant_id}/members, /tenants/{tenant_id}/invitations, and the
// tenant-less /invitations/accept and /invitations/decline. GET
// /invitations/{token} is public; every other route requires a valid bearer
// token, with membership/admin checks done in the service.
type Handler struct {
	service        *Service
	userMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, userMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, userMiddleware: userMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/tenants/{tenant_id}/members", func(r chi.Router) {
		r.Use(h.userMiddleware)

		r.Get("/", serve(h.filter))
		r.Post("/", serve(h.create))
		r.Put("/{id}/role", serve(h.updateRole))
		r.Delete("/{id}", serve(h.remove))
	})

	r.Route("/tenants/{tenant_id}/invitations", func(r chi.Router) {
		r.Use(h.userMiddleware)

		r.Get("/", serve(h.filterInvitations))
		r.Post("/", serve(h.invite))
		r.Delete("/{id}", serve(h.cancelInvitation))
	})

	r.Route("/invitations", func(r chi.Router) {
		r.Get("/{token}", serve(h.getInvitationByToken))

		r.Group(func(r chi.Router) {
			r.Use(h.userMiddleware)

			r.Post("/accept", serve(h.acceptInvitation))
			r.Post("/decline", serve(h.declineInvitation))
		})
	})
}

// @Summary		List members
// @Description	Lists a tenant's members. Caller must be a member (any role).
// @Tags			member
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	MembersPage
// @Failure		403			{object}	api.Error
// @Router			/tenants/{tenant_id}/members [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	tenantID := chi.URLParam(r, "tenant_id")
	var query FindMembersFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), tenantID, query)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create member
// @Description	Directly creates a User (with password) and Member in one step. Only available on unlicensed single-tenant instances (no invitation flow). Caller must be a tenant admin.
// @Tags			member
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string			true	"Tenant ID"
// @Param			body		body		CreateMemberDTO	true	"Member data"
// @Success		201			{object}	MemberResponse
// @Failure		403			{object}	api.Error
// @Failure		409			{object}	api.Error
// @Router			/tenants/{tenant_id}/members [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	tenantID := chi.URLParam(r, "tenant_id")
	var dto CreateMemberDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), tenantID, dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update member role
// @Description	Changes a member's role. Caller must be a tenant admin and cannot change their own role.
// @Tags			member
// @Accept			json
// @Security		BearerAuth
// @Param			tenant_id	path	string					true	"Tenant ID"
// @Param			id			path	string					true	"Member ID"
// @Param			body		body	UpdateMemberRoleDTO	true	"New role"
// @Success		204
// @Failure		403	{object}	api.Error
// @Router			/tenants/{tenant_id}/members/{id}/role [put]
func (h *Handler) updateRole(w http.ResponseWriter, r *http.Request) error {
	tenantID := chi.URLParam(r, "tenant_id")
	id := chi.URLParam(r, "id")
	var dto UpdateMemberRoleDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.UpdateRole(r.Context(), tenantID, id, dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Remove member
// @Description	Removes a member from the tenant. Caller must be a tenant admin and cannot remove themselves.
// @Tags			member
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"Member ID"
// @Success		204
// @Failure		403	{object}	api.Error
// @Router			/tenants/{tenant_id}/members/{id} [delete]
func (h *Handler) remove(w http.ResponseWriter, r *http.Request) error {
	tenantID := chi.URLParam(r, "tenant_id")
	id := chi.URLParam(r, "id")
	if err := h.service.Remove(r.Context(), tenantID, id); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		List invitations
// @Description	Lists a tenant's pending and past invitations. Caller must be a tenant admin.
// @Tags			member
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			status		query		string	false	"Invitation status"	Enums(pending, accepted, declined, expired)
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	InvitationsPage
// @Failure		403			{object}	api.Error
// @Router			/tenants/{tenant_id}/invitations [get]
func (h *Handler) filterInvitations(w http.ResponseWriter, r *http.Request) error {
	tenantID := chi.URLParam(r, "tenant_id")
	var query FindInvitationsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.FilterInvitations(r.Context(), tenantID, query)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Invite member
// @Description	Creates a pending invitation and emails it. Caller must be a tenant admin.
// @Tags			member
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string			true	"Tenant ID"
// @Param			body		body		InviteMemberDTO	true	"Invitation data"
// @Success		201			{object}	InvitationResponse
// @Failure		403			{object}	api.Error
// @Failure		409			{object}	api.Error
// @Router			/tenants/{tenant_id}/invitations [post]
func (h *Handler) invite(w http.ResponseWriter, r *http.Request) error {
	tenantID := chi.URLParam(r, "tenant_id")
	var dto InviteMemberDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Invite(r.Context(), tenantID, dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Cancel invitation
// @Description	Deletes a pending invitation. Caller must be a tenant admin.
// @Tags			member
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"Invitation ID"
// @Success		204
// @Failure		403	{object}	api.Error
// @Router			/tenants/{tenant_id}/invitations/{id} [delete]
func (h *Handler) cancelInvitation(w http.ResponseWriter, r *http.Request) error {
	tenantID := chi.URLParam(r, "tenant_id")
	id := chi.URLParam(r, "id")
	if err := h.service.CancelInvitation(r.Context(), tenantID, id); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Get invitation by token
// @Description	Looks up an invitation and its tenant by token. Public — no authentication required. Lets a user preview an invitation (including whether it's expired, accepted, or declined) before accepting or declining it.
// @Tags			member
// @Produce		json
// @Param			token	path		string	true	"Invitation token"
// @Success		200		{object}	InvitationInfoResponse
// @Failure		404		{object}	api.Error
// @Router			/invitations/{token} [get]
func (h *Handler) getInvitationByToken(w http.ResponseWriter, r *http.Request) error {
	token := chi.URLParam(r, "token")
	data, err := h.service.GetInvitationByToken(r.Context(), token)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Accept invitation
// @Description	Joins the authenticated user to the invitation's tenant. The invitation must be pending, unexpired, and addressed to the caller's own email.
// @Tags			member
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			body	body		AcceptInvitationDTO	true	"Invitation token"
// @Success		201		{object}	MemberResponse
// @Failure		403		{object}	api.Error
// @Failure		409		{object}	api.Error
// @Router			/invitations/accept [post]
func (h *Handler) acceptInvitation(w http.ResponseWriter, r *http.Request) error {
	var dto AcceptInvitationDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.AcceptInvitation(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Decline invitation
// @Description	Marks a pending invitation as declined. The token must belong to a pending invitation addressed to the caller's own email.
// @Tags			member
// @Accept			json
// @Security		BearerAuth
// @Param			body	body	DeclineInvitationDTO	true	"Invitation token"
// @Success		204
// @Failure		403		{object}	api.Error
// @Failure		409		{object}	api.Error
// @Router			/invitations/decline [post]
func (h *Handler) declineInvitation(w http.ResponseWriter, r *http.Request) error {
	var dto DeclineInvitationDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.DeclineInvitation(r.Context(), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
