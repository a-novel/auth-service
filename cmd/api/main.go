package main

import (
	"context"
	"fmt"
	"github.com/a-novel/auth-service/config"
	"github.com/a-novel/auth-service/migrations"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/bunovel"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
	sendgridproxy "github.com/a-novel/sendgrid-proxy"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"io/fs"
)

func getFrontendURL(value string) string {
	return config.App.Frontend.URLs[0] + value
}

func main() {
	ctx := context.Background()
	logger := config.GetLogger()
	permissionsClient := config.GetPermissionsClient(logger)

	postgres, sql, err := bunovel.NewClient(ctx, bunovel.Config{
		Driver:                &bunovel.PGDriver{DSN: config.Postgres.DSN, AppName: config.App.Name},
		Migrations:            &bunovel.MigrateConfig{Files: []fs.FS{migrations.Migrations}},
		DiscardUnknownColumns: true,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to postgres")
	}
	defer func() {
		_ = postgres.Close()
		_ = sql.Close()
	}()

	mailSender := mail.NewEmail(config.Mailer.Sender.Name, config.Mailer.Sender.Email)
	mailClient := sendgridproxy.NewMailer(config.Mailer.APIKey, mailSender, config.Mailer.Sandbox, logger)

	secretKeysDAO, logger := config.GetSecretsRepository(logger)
	credentialsDAO := dao.NewCredentialsRepository(postgres)
	identityDAO := dao.NewIdentityRepository(postgres)
	profileDAO := dao.NewProfileRepository(postgres)
	userDAO := dao.NewUserRepository(postgres)

	generateTokenService := services.NewGenerateTokenService(secretKeysDAO, config.Tokens.TTL)
	getTokenService := services.NewGetTokenStatusService(secretKeysDAO)
	introspectTokenService := services.NewIntrospectTokenService(generateTokenService, getTokenService, config.Tokens.RenewDelta)

	cancelNewEmailService := services.NewCancelNewEmailService(credentialsDAO, introspectTokenService)
	emailExistsService := services.NewEmailExistsService(credentialsDAO)
	listService := services.NewListService(userDAO)
	loginService := services.NewLoginService(credentialsDAO, generateTokenService)
	previewService := services.NewPreviewService(profileDAO, identityDAO)
	previewPrivateService := services.NewPreviewPrivateService(credentialsDAO, profileDAO, identityDAO, introspectTokenService)
	registerService := services.NewRegisterService(credentialsDAO, profileDAO, userDAO, mailClient, goframework.GenerateCode, generateTokenService, getFrontendURL(config.App.Frontend.Routes.ValidateEmail), config.Mailer.Templates.EmailValidation)
	resendEmailValidationService := services.NewResendEmailValidationService(credentialsDAO, identityDAO, mailClient, goframework.GenerateCode, introspectTokenService, getFrontendURL(config.App.Frontend.Routes.ValidateEmail), config.Mailer.Templates.EmailValidation)
	resendNewEmailValidationService := services.NewResendNewEmailValidationService(credentialsDAO, identityDAO, mailClient, goframework.GenerateCode, introspectTokenService, getFrontendURL(config.App.Frontend.Routes.ValidateNewEmail), config.Mailer.Templates.EmailUpdate)
	resetPasswordService := services.NewResetPasswordService(credentialsDAO, identityDAO, mailClient, goframework.GenerateCode, getFrontendURL(config.App.Frontend.Routes.ResetPassword), config.Mailer.Templates.PasswordReset)
	searchService := services.NewSearchService(userDAO)
	slugExistsService := services.NewSlugExistsService(profileDAO)
	updateEmailService := services.NewUpdateEmailService(credentialsDAO, identityDAO, mailClient, goframework.GenerateCode, introspectTokenService, getFrontendURL(config.App.Frontend.Routes.ValidateNewEmail), config.Mailer.Templates.EmailUpdate)
	updateIdentityService := services.NewUpdateIdentityService(identityDAO, introspectTokenService)
	updatePasswordService := services.NewUpdatePasswordService(credentialsDAO)
	updateProfileService := services.NewUpdateProfileService(profileDAO, introspectTokenService)
	validateEmailService := services.NewValidateEmailService(credentialsDAO, permissionsClient)
	validateNewEmailService := services.NewValidateNewEmailService(credentialsDAO, permissionsClient)
	getCredentialsService := services.NewGetCredentialsService(credentialsDAO, introspectTokenService)
	getIdentityService := services.NewGetIdentityService(identityDAO, introspectTokenService)
	getProfileService := services.NewGetProfileService(profileDAO, introspectTokenService)

	introspectTokenHandler := handlers.NewIntrospectTokenHandler(introspectTokenService)
	cancelNewEmailHandler := handlers.NewCancelNewEmailHandler(cancelNewEmailService)
	emailExistsHandler := handlers.NewEmailExistsHandler(emailExistsService)
	listHandler := handlers.NewListHandler(listService)
	loginHandler := handlers.NewLoginHandler(loginService)
	previewHandler := handlers.NewPreviewHandler(previewService)
	previewPrivateHandler := handlers.NewPreviewPrivateHandler(previewPrivateService)
	registerHandler := handlers.NewRegisterHandler(registerService)
	resendEmailValidationHandler := handlers.NewResendEmailValidationHandler(resendEmailValidationService)
	resendNewEmailValidationHandler := handlers.NewResendNewEmailValidationHandler(resendNewEmailValidationService)
	resetPasswordHandler := handlers.NewResetPasswordHandler(resetPasswordService)
	searchHandler := handlers.NewSearchHandler(searchService)
	slugExistsHandler := handlers.NewSlugExistsHandler(slugExistsService)
	updateEmailHandler := handlers.NewUpdateEmailHandler(updateEmailService)
	updateIdentityHandler := handlers.NewUpdateIdentityHandler(updateIdentityService)
	updatePasswordHandler := handlers.NewUpdatePasswordHandler(updatePasswordService)
	updateProfileHandler := handlers.NewUpdateProfileHandler(updateProfileService)
	validateEmailHandler := handlers.NewValidateEmailHandler(validateEmailService)
	validateNewEmailHandler := handlers.NewValidateNewEmailHandler(validateNewEmailService)
	getCredentialsHandler := handlers.NewGetCredentialsHandler(getCredentialsService)
	getIdentityHandler := handlers.NewGetIdentityHandler(getIdentityService)
	getProfileHandler := handlers.NewGetProfileHandler(getProfileService)

	router := apis.GetRouter(apis.RouterConfig{
		Logger:    logger,
		ProjectID: config.Deploy.ProjectID,
		CORS:      apis.GetCORS(config.App.Frontend.URLs),
		Prod:      config.ENV == config.ProdENV,
		Health: map[string]apis.HealthChecker{
			"postgres": func() error {
				return postgres.PingContext(ctx)
			},
			"permissions-client": func() error {
				return permissionsClient.Ping(ctx)
			},
		},
	})

	// /auth
	router.GET("/auth", introspectTokenHandler.Handle)
	router.POST("/auth", loginHandler.Handle)
	router.PUT("/auth", registerHandler.Handle)
	// /email
	router.DELETE("/email", cancelNewEmailHandler.Handle)
	router.PATCH("/email", updateEmailHandler.Handle)
	// /password
	router.DELETE("/password", resetPasswordHandler.Handle)
	router.PATCH("/password", updatePasswordHandler.Handle)
	// /credentials
	router.GET("/credentials", getCredentialsHandler.Handle)
	// /identity
	router.PATCH("/identity", updateIdentityHandler.Handle)
	router.GET("/identity", getIdentityHandler.Handle)
	// /profile
	router.PATCH("/profile", updateProfileHandler.Handle)
	router.GET("/profile", getProfileHandler.Handle)
	// /email/validation
	router.PATCH("/email/validation", resendEmailValidationHandler.Handle)
	router.GET("/email/validation", validateEmailHandler.Handle)
	// //email/pending/validation
	router.PATCH("/email/pending/validation", resendNewEmailValidationHandler.Handle)
	router.GET("/email/pending/validation", validateNewEmailHandler.Handle)
	// /email/exists
	router.GET("/email/exists", emailExistsHandler.Handle)
	// /slug/exists
	router.GET("/slug/exists", slugExistsHandler.Handle)
	// /uses
	router.GET("/users", listHandler.Handle)
	// /uses/search
	router.GET("/users/search", searchHandler.Handle)
	// /user
	router.GET("/user", previewHandler.Handle)
	// /user/me
	router.GET("/user/me", previewPrivateHandler.Handle)

	if err := router.Run(fmt.Sprintf(":%d", config.API.Port)); err != nil {
		logger.Fatal().Err(err).Msg("a fatal error occurred while running the API, and the server had to shut down")
	}
}
