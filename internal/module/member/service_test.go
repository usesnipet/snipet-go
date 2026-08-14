package member_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/email"
	"github.com/usesnipet/snipet/internal/module/member"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
)

const (
	tenantID = "11111111-1111-1111-1111-111111111111"
	adminID  = "22222222-2222-2222-2222-222222222222"
	userID   = "33333333-3333-3333-3333-333333333333"
)

func newTestService(
	t *testing.T,
	memberRepo repository.IMemberRepository,
	invitationRepo repository.ITenantInvitationRepository,
	tenantRepo repository.ITenantRepository,
	userRepo repository.IUserRepository,
) *member.Service {
	t.Helper()

	txManager := mocks.NewMockITxManager(t)
	txManager.EXPECT().
		WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()

	return member.NewService(
		config.AuthConfig{AppURL: "https://snipet.dev", TenantInvitationExpiration: 168 * time.Hour},
		memberRepo,
		invitationRepo,
		tenantRepo,
		userRepo,
		txManager,
		auth.NewTokenService(),
		email.NewService(config.SMTPConfig{Enable: false}, logger.NewLogger(logger.LevelError)),
	)
}

func ctxFor(userID string) context.Context {
	claims := &auth.PlatformUserClaims{BaseClaims: auth.NewBaseClaims(config.AuthConfig{}, userID)}
	return auth.SetPrincipal(context.Background(), auth.NewPrincipal(auth.PrincipalTypePlatformJWT, nil, claims))
}

func assertAppError(t *testing.T, err error, statusCode int) {
	t.Helper()
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, statusCode, appErr.StatusCode)
}

func TestUpdateRoleRejectsNonAdmin(t *testing.T) {
	t.Parallel()

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, userID, tenantID).
		Return(&model.Member{ID: "m1", UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true}, nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	err := svc.UpdateRole(ctxFor(userID), tenantID, "target-member", member.UpdateMemberRoleDTO{Role: model.RoleAdmin})
	assertAppError(t, err, 403)
}

func TestUpdateRoleRejectsChangingOwnRole(t *testing.T) {
	t.Parallel()

	adminMemberID := "admin-member-id"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, adminID, tenantID).
		Return(&model.Member{ID: adminMemberID, UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)
	memberRepo.EXPECT().
		FindByID(mock.Anything, adminMemberID).
		Return(&model.Member{ID: adminMemberID, UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	err := svc.UpdateRole(ctxFor(adminID), tenantID, adminMemberID, member.UpdateMemberRoleDTO{Role: model.RoleUser})
	assertAppError(t, err, 403)
}

func TestUpdateRoleAllowsAdminToChangeOthersRole(t *testing.T) {
	t.Parallel()

	targetMemberID := "target-member-id"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, adminID, tenantID).
		Return(&model.Member{ID: "admin-member-id", UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)
	memberRepo.EXPECT().
		FindByID(mock.Anything, targetMemberID).
		Return(&model.Member{ID: targetMemberID, UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true}, nil)
	memberRepo.EXPECT().
		UpdateByID(mock.Anything, targetMemberID, mock.Anything).
		Return(nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	err := svc.UpdateRole(ctxFor(adminID), tenantID, targetMemberID, member.UpdateMemberRoleDTO{Role: model.RoleAdmin})
	require.NoError(t, err)
}

func TestRemoveRejectsRemovingSelf(t *testing.T) {
	t.Parallel()

	adminMemberID := "admin-member-id"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, adminID, tenantID).
		Return(&model.Member{ID: adminMemberID, UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)
	memberRepo.EXPECT().
		FindByID(mock.Anything, adminMemberID).
		Return(&model.Member{ID: adminMemberID, UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	err := svc.Remove(ctxFor(adminID), tenantID, adminMemberID)
	assertAppError(t, err, 403)
}

func TestRemoveRejectsMemberFromAnotherTenant(t *testing.T) {
	t.Parallel()

	otherTenantMemberID := "other-tenant-member-id"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, adminID, tenantID).
		Return(&model.Member{ID: "admin-member-id", UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)
	memberRepo.EXPECT().
		FindByID(mock.Anything, otherTenantMemberID).
		Return(&model.Member{ID: otherTenantMemberID, UserID: userID, TenantID: "other-tenant-id", Role: model.RoleUser, IsActive: true}, nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	err := svc.Remove(ctxFor(adminID), tenantID, otherTenantMemberID)
	assertAppError(t, err, 404)
}

func TestFilterRejectsNonMember(t *testing.T) {
	t.Parallel()

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, userID, tenantID).
		Return(nil, apperr.NotFound("member not found"))

	svc := newTestService(t, memberRepo, nil, nil, nil)

	_, err := svc.Filter(ctxFor(userID), tenantID, member.FindMembersFilterDTO{})
	assertAppError(t, err, 403)
}

func TestFilterAllowsAnyMemberRole(t *testing.T) {
	t.Parallel()

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, userID, tenantID).
		Return(&model.Member{ID: "m1", UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true}, nil)
	memberRepo.EXPECT().
		FilterWithUser(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.Member{}, 0, 0, 20), nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	result, err := svc.Filter(ctxFor(userID), tenantID, member.FindMembersFilterDTO{})
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestInviteRejectsWhenAlreadyMember(t *testing.T) {
	t.Parallel()

	adminMemberID := "admin-member-id"
	existing := &model.User{ID: userID, Name: "Existing", Email: "existing@example.com"}

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, adminID, tenantID).
		Return(&model.Member{ID: adminMemberID, UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, userID, tenantID).
		Return(&model.Member{ID: "existing-member-id", UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true}, nil)

	tenantRepo := mocks.NewMockITenantRepository(t)
	tenantRepo.EXPECT().FindByID(mock.Anything, tenantID).Return(&model.Tenant{ID: tenantID, Name: "Acme"}, nil)

	userRepo := mocks.NewMockIUserRepository(t)
	userRepo.EXPECT().FindByID(mock.Anything, adminID).Return(&model.User{ID: adminID, Name: "Admin"}, nil)
	userRepo.EXPECT().FindByEmail(mock.Anything, existing.Email).Return(existing, nil)

	svc := newTestService(t, memberRepo, nil, tenantRepo, userRepo)

	_, err := svc.Invite(ctxFor(adminID), tenantID, member.InviteMemberDTO{Email: existing.Email, Role: model.RoleUser})
	assertAppError(t, err, 409)
}

func TestInviteRejectsWhenInvitationAlreadyPending(t *testing.T) {
	t.Parallel()

	adminMemberID := "admin-member-id"
	const email = "new@example.com"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, adminID, tenantID).
		Return(&model.Member{ID: adminMemberID, UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)

	tenantRepo := mocks.NewMockITenantRepository(t)
	tenantRepo.EXPECT().FindByID(mock.Anything, tenantID).Return(&model.Tenant{ID: tenantID, Name: "Acme"}, nil)

	userRepo := mocks.NewMockIUserRepository(t)
	userRepo.EXPECT().FindByID(mock.Anything, adminID).Return(&model.User{ID: adminID, Name: "Admin"}, nil)
	userRepo.EXPECT().FindByEmail(mock.Anything, email).Return(nil, apperr.NotFound("user not found"))

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().
		FindPendingByEmailAndTenant(mock.Anything, email, tenantID).
		Return(&model.TenantInvitation{ID: "pending-invite"}, nil)

	svc := newTestService(t, memberRepo, invitationRepo, tenantRepo, userRepo)

	_, err := svc.Invite(ctxFor(adminID), tenantID, member.InviteMemberDTO{Email: email, Role: model.RoleUser})
	assertAppError(t, err, 409)
}

func TestCancelInvitationRejectsInvitationFromAnotherTenant(t *testing.T) {
	t.Parallel()

	adminMemberID := "admin-member-id"
	invitationID := "invite-1"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, adminID, tenantID).
		Return(&model.Member{ID: adminMemberID, UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}, nil)

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().
		FindByID(mock.Anything, invitationID).
		Return(&model.TenantInvitation{ID: invitationID, TenantID: "other-tenant-id"}, nil)

	svc := newTestService(t, memberRepo, invitationRepo, nil, nil)

	err := svc.CancelInvitation(ctxFor(adminID), tenantID, invitationID)
	assertAppError(t, err, 404)
}

func TestAcceptInvitationJoinsTenantWithInvitationRole(t *testing.T) {
	t.Parallel()

	const token = "invite-token"
	invitation := &model.TenantInvitation{
		ID:        "invite-1",
		TenantID:  tenantID,
		Email:     "user@example.com",
		Token:     token,
		Role:      model.RoleAdmin,
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, userID, tenantID).
		Return(nil, apperr.NotFound("member not found"))
	memberRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, m *model.Member) { m.ID = "new-member-id" }).
		Return(nil)

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().FindByToken(mock.Anything, token).Return(invitation, nil)
	invitationRepo.EXPECT().UpdateByID(mock.Anything, invitation.ID, mock.Anything).Return(nil)

	userRepo := mocks.NewMockIUserRepository(t)
	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(&model.User{ID: userID, Email: "user@example.com"}, nil)

	svc := newTestService(t, memberRepo, invitationRepo, nil, userRepo)

	created, err := svc.AcceptInvitation(ctxFor(userID), member.AcceptInvitationDTO{Token: token})
	require.NoError(t, err)
	assert.Equal(t, tenantID, created.TenantID)
	assert.Equal(t, model.RoleAdmin, created.Role)
}

func TestAcceptInvitationRejectsEmailMismatch(t *testing.T) {
	t.Parallel()

	const token = "invite-token"
	invitation := &model.TenantInvitation{
		ID:        "invite-1",
		TenantID:  tenantID,
		Email:     "invited@example.com",
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().FindByToken(mock.Anything, token).Return(invitation, nil)

	userRepo := mocks.NewMockIUserRepository(t)
	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(&model.User{ID: userID, Email: "someone-else@example.com"}, nil)

	svc := newTestService(t, nil, invitationRepo, nil, userRepo)

	_, err := svc.AcceptInvitation(ctxFor(userID), member.AcceptInvitationDTO{Token: token})
	assertAppError(t, err, 403)
}

func TestAcceptInvitationRejectsExpired(t *testing.T) {
	t.Parallel()

	const token = "invite-token"
	invitation := &model.TenantInvitation{
		ID:        "invite-1",
		TenantID:  tenantID,
		Email:     "user@example.com",
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().FindByToken(mock.Anything, token).Return(invitation, nil)

	userRepo := mocks.NewMockIUserRepository(t)
	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(&model.User{ID: userID, Email: "user@example.com"}, nil)

	svc := newTestService(t, nil, invitationRepo, nil, userRepo)

	_, err := svc.AcceptInvitation(ctxFor(userID), member.AcceptInvitationDTO{Token: token})
	assertAppError(t, err, 409)
}

func TestAcceptInvitationRejectsWhenAlreadyMember(t *testing.T) {
	t.Parallel()

	const token = "invite-token"
	invitation := &model.TenantInvitation{
		ID:        "invite-1",
		TenantID:  tenantID,
		Email:     "user@example.com",
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, userID, tenantID).
		Return(&model.Member{ID: "existing-member-id", UserID: userID, TenantID: tenantID}, nil)

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().FindByToken(mock.Anything, token).Return(invitation, nil)

	userRepo := mocks.NewMockIUserRepository(t)
	userRepo.EXPECT().FindByID(mock.Anything, userID).Return(&model.User{ID: userID, Email: "user@example.com"}, nil)

	svc := newTestService(t, memberRepo, invitationRepo, nil, userRepo)

	_, err := svc.AcceptInvitation(ctxFor(userID), member.AcceptInvitationDTO{Token: token})
	assertAppError(t, err, 409)
}
