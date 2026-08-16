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

func adminUser() *model.User {
	return &model.User{ID: adminID, Name: "Admin", Email: "admin@example.com"}
}

func adminMembership(id string) *model.Member {
	return &model.Member{ID: id, UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}
}

func regularUser() *model.User {
	return &model.User{ID: userID, Name: "Regular", Email: "user@example.com"}
}

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

// ctxFor builds a context carrying the auth.UserIdentity that
// guard.RequireUser would have loaded for user, with the given tenant
// memberships.
func ctxFor(user *model.User, memberships ...*model.Member) context.Context {
	return auth.SetUserIdentity(context.Background(), &auth.UserIdentity{User: user, Memberships: memberships})
}

func assertAppError(t *testing.T, err error, statusCode int) {
	t.Helper()
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, statusCode, appErr.StatusCode)
}

func TestUpdateRoleRejectsNonAdmin(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil, nil, nil, nil)

	ctx := ctxFor(regularUser(), &model.Member{ID: "m1", UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true})
	err := svc.UpdateRole(ctx, tenantID, "target-member", member.UpdateMemberRoleDTO{Role: model.RoleAdmin})
	assertAppError(t, err, 403)
}

func TestUpdateRoleRejectsChangingOwnRole(t *testing.T) {
	t.Parallel()

	adminMemberID := "admin-member-id"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByID(mock.Anything, adminMemberID).
		Return(adminMembership(adminMemberID), nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	ctx := ctxFor(adminUser(), adminMembership(adminMemberID))
	err := svc.UpdateRole(ctx, tenantID, adminMemberID, member.UpdateMemberRoleDTO{Role: model.RoleUser})
	assertAppError(t, err, 403)
}

func TestUpdateRoleAllowsAdminToChangeOthersRole(t *testing.T) {
	t.Parallel()

	targetMemberID := "target-member-id"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByID(mock.Anything, targetMemberID).
		Return(&model.Member{ID: targetMemberID, UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true}, nil)
	memberRepo.EXPECT().
		UpdateByID(mock.Anything, targetMemberID, mock.Anything).
		Return(nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	ctx := ctxFor(adminUser(), adminMembership("admin-member-id"))
	err := svc.UpdateRole(ctx, tenantID, targetMemberID, member.UpdateMemberRoleDTO{Role: model.RoleAdmin})
	require.NoError(t, err)
}

func TestRemoveRejectsRemovingSelf(t *testing.T) {
	t.Parallel()

	adminMemberID := "admin-member-id"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByID(mock.Anything, adminMemberID).
		Return(adminMembership(adminMemberID), nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	ctx := ctxFor(adminUser(), adminMembership(adminMemberID))
	err := svc.Remove(ctx, tenantID, adminMemberID)
	assertAppError(t, err, 403)
}

func TestRemoveRejectsMemberFromAnotherTenant(t *testing.T) {
	t.Parallel()

	otherTenantMemberID := "other-tenant-member-id"

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByID(mock.Anything, otherTenantMemberID).
		Return(&model.Member{ID: otherTenantMemberID, UserID: userID, TenantID: "other-tenant-id", Role: model.RoleUser, IsActive: true}, nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	ctx := ctxFor(adminUser(), adminMembership("admin-member-id"))
	err := svc.Remove(ctx, tenantID, otherTenantMemberID)
	assertAppError(t, err, 404)
}

func TestFilterRejectsNonMember(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil, nil, nil, nil)

	_, err := svc.Filter(ctxFor(regularUser()), tenantID, member.FindMembersFilterDTO{})
	assertAppError(t, err, 403)
}

func TestFilterAllowsAnyMemberRole(t *testing.T) {
	t.Parallel()

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FilterWithUser(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.Member{}, 0, 0, 20), nil)

	svc := newTestService(t, memberRepo, nil, nil, nil)

	ctx := ctxFor(regularUser(), &model.Member{ID: "m1", UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true})
	result, err := svc.Filter(ctx, tenantID, member.FindMembersFilterDTO{})
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestInviteRejectsWhenAlreadyMember(t *testing.T) {
	t.Parallel()

	existing := &model.User{ID: userID, Name: "Existing", Email: "existing@example.com"}

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		FindByUserAndTenant(mock.Anything, userID, tenantID).
		Return(&model.Member{ID: "existing-member-id", UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true}, nil)

	tenantRepo := mocks.NewMockITenantRepository(t)
	tenantRepo.EXPECT().FindByID(mock.Anything, tenantID).Return(&model.Tenant{ID: tenantID, Name: "Acme"}, nil)

	userRepo := mocks.NewMockIUserRepository(t)
	userRepo.EXPECT().FindByEmail(mock.Anything, existing.Email).Return(existing, nil)

	svc := newTestService(t, memberRepo, nil, tenantRepo, userRepo)

	ctx := ctxFor(adminUser(), adminMembership("admin-member-id"))
	_, err := svc.Invite(ctx, tenantID, member.InviteMemberDTO{Email: existing.Email, Role: model.RoleUser})
	assertAppError(t, err, 409)
}

func TestInviteRejectsWhenInvitationAlreadyPending(t *testing.T) {
	t.Parallel()

	const email = "new@example.com"

	tenantRepo := mocks.NewMockITenantRepository(t)
	tenantRepo.EXPECT().FindByID(mock.Anything, tenantID).Return(&model.Tenant{ID: tenantID, Name: "Acme"}, nil)

	userRepo := mocks.NewMockIUserRepository(t)
	userRepo.EXPECT().FindByEmail(mock.Anything, email).Return(nil, apperr.NotFound("user not found"))

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().
		FindPendingByEmailAndTenant(mock.Anything, email, tenantID).
		Return(&model.TenantInvitation{ID: "pending-invite"}, nil)

	svc := newTestService(t, nil, invitationRepo, tenantRepo, userRepo)

	ctx := ctxFor(adminUser(), adminMembership("admin-member-id"))
	_, err := svc.Invite(ctx, tenantID, member.InviteMemberDTO{Email: email, Role: model.RoleUser})
	assertAppError(t, err, 409)
}

func TestCancelInvitationRejectsInvitationFromAnotherTenant(t *testing.T) {
	t.Parallel()

	invitationID := "invite-1"

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().
		FindByID(mock.Anything, invitationID).
		Return(&model.TenantInvitation{ID: invitationID, TenantID: "other-tenant-id"}, nil)

	svc := newTestService(t, nil, invitationRepo, nil, nil)

	ctx := ctxFor(adminUser(), adminMembership("admin-member-id"))
	err := svc.CancelInvitation(ctx, tenantID, invitationID)
	assertAppError(t, err, 404)
}

func TestAcceptInvitationJoinsTenantWithInvitationRole(t *testing.T) {
	t.Parallel()

	const token = "invite-token"
	user := &model.User{ID: userID, Email: "user@example.com"}
	invitation := &model.TenantInvitation{
		ID:        "invite-1",
		TenantID:  tenantID,
		Email:     user.Email,
		Token:     token,
		Role:      model.RoleAdmin,
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	memberRepo := mocks.NewMockIMemberRepository(t)
	memberRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, m *model.Member) { m.ID = "new-member-id" }).
		Return(nil)

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().FindByToken(mock.Anything, token).Return(invitation, nil)
	invitationRepo.EXPECT().UpdateByID(mock.Anything, invitation.ID, mock.Anything).Return(nil)

	svc := newTestService(t, memberRepo, invitationRepo, nil, nil)

	created, err := svc.AcceptInvitation(ctxFor(user), member.AcceptInvitationDTO{Token: token})
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

	svc := newTestService(t, nil, invitationRepo, nil, nil)

	user := &model.User{ID: userID, Email: "someone-else@example.com"}
	_, err := svc.AcceptInvitation(ctxFor(user), member.AcceptInvitationDTO{Token: token})
	assertAppError(t, err, 403)
}

func TestAcceptInvitationRejectsExpired(t *testing.T) {
	t.Parallel()

	const token = "invite-token"
	user := &model.User{ID: userID, Email: "user@example.com"}
	invitation := &model.TenantInvitation{
		ID:        "invite-1",
		TenantID:  tenantID,
		Email:     user.Email,
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().FindByToken(mock.Anything, token).Return(invitation, nil)

	svc := newTestService(t, nil, invitationRepo, nil, nil)

	_, err := svc.AcceptInvitation(ctxFor(user), member.AcceptInvitationDTO{Token: token})
	assertAppError(t, err, 409)
}

func TestAcceptInvitationRejectsWhenAlreadyMember(t *testing.T) {
	t.Parallel()

	const token = "invite-token"
	user := &model.User{ID: userID, Email: "user@example.com"}
	invitation := &model.TenantInvitation{
		ID:        "invite-1",
		TenantID:  tenantID,
		Email:     user.Email,
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	invitationRepo := mocks.NewMockITenantInvitationRepository(t)
	invitationRepo.EXPECT().FindByToken(mock.Anything, token).Return(invitation, nil)

	svc := newTestService(t, nil, invitationRepo, nil, nil)

	ctx := ctxFor(user, &model.Member{ID: "existing-member-id", UserID: userID, TenantID: tenantID})
	_, err := svc.AcceptInvitation(ctx, member.AcceptInvitationDTO{Token: token})
	assertAppError(t, err, 409)
}
