package member

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/email"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	config         config.AuthConfig
	memberRepo     repository.IMemberRepository
	invitationRepo repository.ITenantInvitationRepository
	tenantRepo     repository.ITenantRepository
	userRepo       repository.IUserRepository
	txManager      repository.ITxManager
	tokenGen       *auth.TokenService
	emailSvc       *email.Service
}

func NewService(
	config config.AuthConfig,
	memberRepo repository.IMemberRepository,
	invitationRepo repository.ITenantInvitationRepository,
	tenantRepo repository.ITenantRepository,
	userRepo repository.IUserRepository,
	txManager repository.ITxManager,
	tokenGen *auth.TokenService,
	emailSvc *email.Service,
) *Service {
	return &Service{
		config:         config,
		memberRepo:     memberRepo,
		invitationRepo: invitationRepo,
		tenantRepo:     tenantRepo,
		userRepo:       userRepo,
		txManager:      txManager,
		tokenGen:       tokenGen,
		emailSvc:       emailSvc,
	}
}

func (s *Service) currentUserID(ctx context.Context) (string, error) {
	return auth.PlatformUserID(ctx)
}

// requireMember ensures the caller belongs to the tenant, any role.
func (s *Service) requireMember(ctx context.Context, tenantID string) (string, error) {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return "", err
	}

	if _, err := s.memberRepo.FindByUserAndTenant(ctx, userID, tenantID); err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return "", apperr.Forbidden("not a member of this tenant")
		}
		return "", err
	}

	return userID, nil
}

// requireAdmin ensures the caller is an active tenant admin, or a platform
// admin acting on any tenant.
func (s *Service) requireAdmin(ctx context.Context, tenantID string) (string, error) {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return "", err
	}

	member, err := s.memberRepo.FindByUserAndTenant(ctx, userID, tenantID)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return "", apperr.Forbidden("not allowed to manage members of this tenant")
		}
		return "", err
	}
	if member.Role != model.RoleAdmin || !member.IsActive {
		return "", apperr.Forbidden("not allowed to manage members of this tenant")
	}

	return userID, nil
}

// Filter lists a tenant's members. Caller must be a member (any role).
func (s *Service) Filter(ctx context.Context, tenantID string, dto FindMembersFilterDTO) (*MembersPage, error) {
	if _, err := s.requireMember(ctx, tenantID); err != nil {
		return nil, err
	}

	paginated, err := s.memberRepo.FilterWithUser(ctx, filter.Merge(dto.ToFilter(), filter.New[model.Member](
		filter.WhereEq("tenant_id", tenantID),
	)))
	if err != nil {
		return nil, err
	}

	return (*MembersPage)(paginated), nil
}

// UpdateRole changes a member's role. Caller must be a tenant admin (or
// platform admin) and cannot change their own role.
func (s *Service) UpdateRole(ctx context.Context, tenantID, memberID string, dto UpdateMemberRoleDTO) error {
	callerID, err := s.requireAdmin(ctx, tenantID)
	if err != nil {
		return err
	}

	target, err := s.memberRepo.FindByID(ctx, memberID)
	if err != nil {
		return err
	}
	if target.TenantID != tenantID {
		return apperr.NotFound("member not found")
	}
	if target.UserID == callerID {
		return apperr.Forbidden("cannot change your own role")
	}

	return s.memberRepo.UpdateByID(ctx, memberID, &model.Member{Role: dto.Role})
}

// Remove deletes a member from the tenant. Caller must be a tenant admin (or
// platform admin) and cannot remove themselves.
func (s *Service) Remove(ctx context.Context, tenantID, memberID string) error {
	callerID, err := s.requireAdmin(ctx, tenantID)
	if err != nil {
		return err
	}

	target, err := s.memberRepo.FindByID(ctx, memberID)
	if err != nil {
		return err
	}
	if target.TenantID != tenantID {
		return apperr.NotFound("member not found")
	}
	if target.UserID == callerID {
		return apperr.Forbidden("cannot remove yourself from the tenant")
	}

	return s.memberRepo.DeleteByID(ctx, memberID)
}

// Invite creates a pending invitation for an email address and sends the
// invitation email. Caller must be a tenant admin (or platform admin).
func (s *Service) Invite(ctx context.Context, tenantID string, dto InviteMemberDTO) (*model.TenantInvitation, error) {
	callerID, err := s.requireAdmin(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	inviter, err := s.userRepo.FindByID(ctx, callerID)
	if err != nil {
		return nil, err
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(dto.Email))

	existingUser, err := s.userRepo.FindByEmail(ctx, normalizedEmail)
	if err == nil {
		if _, err := s.memberRepo.FindByUserAndTenant(ctx, existingUser.ID, tenantID); err == nil {
			return nil, apperr.Conflict("this user is already a member of the tenant")
		}
	} else {
		var appErr *apperr.Error
		if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	if _, err := s.invitationRepo.FindPendingByEmailAndTenant(ctx, normalizedEmail, tenantID); err == nil {
		return nil, apperr.Conflict("an invitation is already pending for this email")
	} else {
		var appErr *apperr.Error
		if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	token, err := s.tokenGen.GenerateToken()
	if err != nil {
		return nil, err
	}

	invitation := &model.TenantInvitation{
		TenantID:  tenantID,
		Email:     normalizedEmail,
		Token:     token,
		Role:      dto.Role,
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(s.config.TenantInvitationExpiration),
	}
	if err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	link := fmt.Sprintf("%s/invite?token=%s", strings.TrimRight(s.config.AppURL, "/"), token)
	if err := s.emailSvc.SendTemplate(ctx, normalizedEmail, email.TemplateTenantInvitation, email.TenantInvitationData{
		TenantName:  tenant.Name,
		InviterName: inviter.Name,
		Link:        link,
	}); err != nil {
		return nil, err
	}

	return invitation, nil
}

// CancelInvitation deletes a pending invitation. Caller must be a tenant
// admin (or platform admin).
func (s *Service) CancelInvitation(ctx context.Context, tenantID, invitationID string) error {
	if _, err := s.requireAdmin(ctx, tenantID); err != nil {
		return err
	}

	invitation, err := s.invitationRepo.FindByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation.TenantID != tenantID {
		return apperr.NotFound("invitation not found")
	}

	return s.invitationRepo.DeleteByID(ctx, invitationID)
}

// FilterInvitations lists a tenant's invitations. Caller must be a tenant
// admin.
func (s *Service) FilterInvitations(ctx context.Context, tenantID string, dto FindInvitationsFilterDTO) (*InvitationsPage, error) {
	if _, err := s.requireAdmin(ctx, tenantID); err != nil {
		return nil, err
	}

	return s.invitationRepo.Filter(ctx, filter.Merge(dto.ToFilter(), filter.New[model.TenantInvitation](
		filter.WhereEq("tenant_id", tenantID),
	)))
}

// AcceptInvitation joins the authenticated user to the invitation's tenant
// with the role it was issued for. The token must belong to a pending,
// unexpired invitation addressed to the caller's own email.
func (s *Service) AcceptInvitation(ctx context.Context, dto AcceptInvitationDTO) (*model.Member, error) {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	invitation, err := s.invitationRepo.FindByToken(ctx, dto.Token)
	if err != nil {
		return nil, err
	}

	if invitation.Status != model.InvitationStatusPending {
		return nil, apperr.Conflict("invitation is no longer pending")
	}
	if invitation.ExpiresAt.Before(time.Now()) {
		return nil, apperr.Conflict("invitation has expired")
	}
	if !strings.EqualFold(invitation.Email, user.Email) {
		return nil, apperr.Forbidden("this invitation was sent to a different email address")
	}

	if _, err := s.memberRepo.FindByUserAndTenant(ctx, userID, invitation.TenantID); err == nil {
		return nil, apperr.Conflict("you are already a member of this tenant")
	} else {
		var appErr *apperr.Error
		if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	created := &model.Member{
		UserID:   userID,
		TenantID: invitation.TenantID,
		Role:     invitation.Role,
		IsActive: true,
	}

	err = s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.memberRepo.Create(ctx, created); err != nil {
			return err
		}
		return s.invitationRepo.UpdateByID(ctx, invitation.ID, &model.TenantInvitation{Status: model.InvitationStatusAccepted})
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}
