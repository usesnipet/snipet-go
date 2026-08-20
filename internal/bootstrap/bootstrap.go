package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/config"
	_ "github.com/usesnipet/snipet/docs/swagger"
	"github.com/usesnipet/snipet/drivers/index"
	"github.com/usesnipet/snipet/drivers/llm"
	"github.com/usesnipet/snipet/drivers/source"
	"github.com/usesnipet/snipet/drivers/tool"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/guard"
	"github.com/usesnipet/snipet/internal/infra/cache"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/license"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/module/agent"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
	appmodule "github.com/usesnipet/snipet/internal/module/app"
	auth_module "github.com/usesnipet/snipet/internal/module/auth"
	"github.com/usesnipet/snipet/internal/module/clientauth"
	clientauth_provider "github.com/usesnipet/snipet/internal/module/clientauth/auth-provider"
	"github.com/usesnipet/snipet/internal/module/clientuser"
	"github.com/usesnipet/snipet/internal/module/email"
	"github.com/usesnipet/snipet/internal/module/knowledge"
	knowledgeindex "github.com/usesnipet/snipet/internal/module/knowledge-index"
	llmmodule "github.com/usesnipet/snipet/internal/module/llm"
	"github.com/usesnipet/snipet/internal/module/member"
	"github.com/usesnipet/snipet/internal/module/session"
	systemmodule "github.com/usesnipet/snipet/internal/module/system"
	"github.com/usesnipet/snipet/internal/module/tenant"
	"github.com/usesnipet/snipet/internal/module/user"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/web"
)

func Bootstrap(cfg *config.Config, logger *logger.Logger) error {
	// database
	db, _, embeddedDB, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Errorf("failed to create database: %v", err)
		return err
	}

	if embeddedDB != nil {
		defer func() {
			logger.Infof("stopping embedded database...")
			if err := embeddedDB.Stop(); err != nil {
				logger.Errorf("failed to stop embedded database: %v", err)
				return
			}
			logger.Infof("embedded database stopped successfully")
		}()
	}

	// repositories
	txManager := repository.NewTxManager(db)
	apiKeyRepo := repository.NewApiKeyRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	llmRepo := repository.NewLLMRepository(db)
	appRepo := repository.NewAppRepository(db)
	sessionRepo := repository.NewSessionRepository(db, appRepo)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	knowledgeIndexRepo := repository.NewKnowledgeIndexRepository(db)
	knowledgeItemRepo := repository.NewKnowledgeItemRepository(db)
	indexedKnowledgeItemRepo := repository.NewIndexedKnowledgeItemRepository(db)
	appUserRepo := repository.NewAppUserRepository(db, appRepo)
	refreshTokenRepo := repository.NewAppUserRefreshTokenRepository(db)
	executionRepo := repository.NewExecutionRepository(db)
	messageRepo := repository.NewExecutionMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	tenantInvitationRepo := repository.NewTenantInvitationRepository(db)

	// runtime
	sourceRegistry := source.Registry(logger)
	sourceManager := manager.NewDriver(sourceRegistry)

	indexRegistry := index.Registry(logger)
	indexManager := manager.NewDriver(indexRegistry)

	llmRegistry := llm.Registry(logger)
	llmManager := manager.NewDriver(llmRegistry)

	toolRegistry := tool.Registry(logger)
	toolManager := manager.NewTool(manager.NewDriver(toolRegistry))

	engine := runtime.NewEngine(llmManager, toolManager, logger)

	syncPool := queue.NewPool(cfg.Sync.Workers, logger)
	syncPool.Start(context.Background())

	indexSyncWorker := knowledgeindex.NewSyncIndexWorker(
		indexManager,
		sourceManager,
		knowledgeRepo,
		knowledgeItemRepo,
		knowledgeIndexRepo,
		indexedKnowledgeItemRepo,
		logger,
	)
	knowledgeSyncWorker := knowledge.NewSyncWorker(
		sourceManager,
		knowledgeRepo,
		knowledgeItemRepo,
		knowledgeIndexRepo,
		indexSyncWorker.Sync,
		100,
		logger,
	)

	// services
	apiKeyGenerator := auth.NewAPIKeyGenerator()
	apiKeyHasher := auth.NewKeyHasher()
	appUserJWTService := auth.NewJWTService(cfg.Auth, func() *auth.AppUserClaims { return &auth.AppUserClaims{} })

	refreshTokenService := auth.NewTokenService()

	authRegistry := clientauth_provider.NewRegistry()

	clientauthService := clientauth.NewService(
		authRegistry,
		appRepo,
		appUserRepo,
		refreshTokenRepo,
		appUserJWTService,
		refreshTokenService,
		cfg.Auth,
	)

	platformUserJWTService := auth.NewJWTService(cfg.Auth, func() *auth.PlatformUserClaims { return &auth.PlatformUserClaims{} })
	userProviderRegistry := auth_module.NewProviderRegistry(
		auth_module.NewGoogleProvider(cfg.Auth),
		auth_module.NewGithubProvider(cfg.Auth),
	)
	emailService := email.NewService(cfg.SMTP, logger)
	licenseService := license.NewService(cfg.License)
	authService := auth_module.NewService(
		cfg.Auth,
		userRepo,
		accountRepo,
		tokenRepo,
		platformUserJWTService,
		refreshTokenService,
		emailService,
		userProviderRegistry,
		licenseService,
	)

	apiKeyService := apikey.NewService(logger, apiKeyRepo, apiKeyGenerator, apiKeyHasher)

	appService := appmodule.NewService(appRepo, apiKeyGenerator, apiKeyHasher, logger)

	agentService := agent.NewService(agentRepo, llmRepo, txManager, engine, executionRepo, messageRepo, logger)
	llmService := llmmodule.NewService(llmRepo, llmManager)

	knowledgeService := knowledge.NewService(
		txManager,
		knowledgeRepo,
		knowledgeItemRepo,
		sourceManager,
		syncPool,
		knowledgeSyncWorker,
	)
	knowledgeIndexService := knowledgeindex.NewService(
		knowledgeIndexRepo,
		knowledgeRepo,
		indexedKnowledgeItemRepo,
		indexManager,
		syncPool,
		indexSyncWorker,
		txManager,
	)

	sessionService := session.NewService(sessionRepo, messageRepo, appService, agentService)

	clientUserService := clientuser.NewService(appUserRepo)
	userService := user.NewService(userRepo, cfg.User)
	systemService := systemmodule.NewService(licenseService)

	tenantService := tenant.NewService(tenantRepo, memberRepo, userRepo, txManager, cfg.Tenant, licenseService, logger)
	memberService := member.NewService(cfg.Auth, memberRepo, tenantInvitationRepo, tenantRepo, userRepo, txManager, refreshTokenService, emailService, licenseService)

	if err := userService.Init(context.Background()); err != nil {
		logger.Errorf("failed to init user service: %v", err)
		return err
	}
	_, err = tenantService.Init(context.Background(), cfg.User.AdminEmail)
	if err != nil {
		logger.Errorf("failed to init tenant service: %v", err)
		return err
	}

	// cache
	apiKeyCache := cache.NewMemoryCache(1000, 1*time.Hour)
	appKeyCache := cache.NewMemoryCache(1000, 1*time.Hour)

	// guards
	requireAPIKey := guard.RequireApiKey(apiKeyService, apiKeyCache)
	requireAppKey := guard.RequireAppKey(appService, appKeyCache)
	requireAppUser := guard.RequireAppUser(appUserJWTService)
	requireUser := guard.RequireUser(platformUserJWTService, userRepo, memberRepo)
	anyClientAuth := guard.Or(requireAppUser, requireAPIKey)
	sessionAuth := guard.Or(requireAppUser, requireAppKey)

	apiKeyMiddleware := requireAPIKey.Handler()
	appKeyMiddleware := requireAppKey.Handler()
	appJWTMiddleware := requireAppUser.Handler()
	userMiddleware := requireUser.Handler()
	anyClientAuthMiddleware := anyClientAuth.Handler()
	sessionAuthMiddleware := sessionAuth.Handler()

	// handlers
	clientAuthHandler := clientauth.NewHandler(clientauthService)
	platformAuthHandler := auth_module.NewHandler(authService, userMiddleware)
	systemHandler := systemmodule.NewHandler(systemService)
	apiKeyHandler := apikey.NewHandler(apiKeyService, userMiddleware, apiKeyMiddleware)
	agentHandler := agent.NewHandler(agentService, userMiddleware, apiKeyMiddleware)
	llmHandler := llmmodule.NewHandler(llmService, userMiddleware)
	appHandler := appmodule.NewHandler(appService, userMiddleware, appKeyMiddleware)
	sessionHandler := session.NewHandler(sessionService, sessionAuthMiddleware)
	knowledgeHandler := knowledge.NewHandler(knowledgeService, userMiddleware)
	knowledgeIndexHandler := knowledgeindex.NewHandler(knowledgeIndexService, userMiddleware)
	clientUserHandler := clientuser.NewHandler(clientUserService, apiKeyMiddleware, anyClientAuthMiddleware, appJWTMiddleware)
	userHandler := user.NewHandler(userService, userMiddleware)
	tenantHandler := tenant.NewHandler(tenantService, userMiddleware)
	memberHandler := member.NewHandler(memberService, userMiddleware)

	// register routes
	api := api.New()
	api.Router.Handle("/*", web.Handler())
	api.Router.Route(config.APIPrefix, func(r chi.Router) {
		clientAuthHandler.RegisterRoutes(r, api.Serve)
		platformAuthHandler.RegisterRoutes(r, api.Serve)
		systemHandler.RegisterRoutes(r, api.Serve)
		apiKeyHandler.RegisterRoutes(r, api.Serve)
		agentHandler.RegisterRoutes(r, api.Serve)
		llmHandler.RegisterRoutes(r, api.Serve)
		appHandler.RegisterRoutes(r, api.Serve)
		sessionHandler.RegisterRoutes(r, api.Serve)
		knowledgeHandler.RegisterRoutes(r, api.Serve)
		knowledgeIndexHandler.RegisterRoutes(r, api.Serve)
		clientUserHandler.RegisterRoutes(r, api.Serve)
		userHandler.RegisterRoutes(r, api.Serve)
		tenantHandler.RegisterRoutes(r, api.Serve)
		memberHandler.RegisterRoutes(r, api.Serve)
	})

	logger.Infof("sync worker pool started with %d workers", cfg.Sync.Workers)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: api.Router,
	}

	go func() {
		logger.Infof("server started on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("failed to start server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Infof("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("failed to shutdown server: %v", err)
	}

	return nil
}
