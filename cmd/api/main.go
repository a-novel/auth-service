package main

import (
	"fmt"
	"github.com/a-novel/auth-service/config"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/mailer"
	"github.com/a-novel/go-framework/security"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func getFrontendURL(value string) string {
	return config.App.Frontend.URLs[0] + value
}

func main() {
	logger := config.GetLogger()

	postgres, closer := config.GetPostgres(logger)
	defer closer()

	mailSender := mail.NewEmail(config.Mailer.Sender.Name, config.Mailer.Sender.Email)
	mailClient := mailer.NewMailer(config.Mailer.APIKey, mailSender, config.Mailer.Sandbox, logger)

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
	registerService := services.NewRegisterService(credentialsDAO, profileDAO, userDAO, mailClient, security.GenerateCode, generateTokenService, config.Mailer.Templates.EmailValidation, getFrontendURL(config.App.Frontend.Routes.ValidateEmail))
	resendEmailValidationService := services.NewResendEmailValidationService(credentialsDAO, identityDAO, mailClient, security.GenerateCode, introspectTokenService, config.Mailer.Templates.EmailValidation, getFrontendURL(config.App.Frontend.Routes.ValidateEmail))
	resendNewEmailValidationService := services.NewResendNewEmailValidationService(credentialsDAO, identityDAO, mailClient, security.GenerateCode, introspectTokenService, config.Mailer.Templates.EmailUpdate, getFrontendURL(config.App.Frontend.Routes.ValidateNewEmail))
	resetPasswordService := services.NewResetPasswordService(credentialsDAO, identityDAO, mailClient, security.GenerateCode, config.Mailer.Templates.PasswordReset, getFrontendURL(config.App.Frontend.Routes.ResetPassword))
	searchService := services.NewSearchService(userDAO)
	slugExistsService := services.NewSlugExistsService(profileDAO)
	updateEmailService := services.NewUpdateEmailService(credentialsDAO, identityDAO, mailClient, security.GenerateCode, introspectTokenService, config.Mailer.Templates.EmailUpdate, getFrontendURL(config.App.Frontend.Routes.ValidateNewEmail))
	updateIdentityService := services.NewUpdateIdentityService(identityDAO, introspectTokenService)
	updatePasswordService := services.NewUpdatePasswordService(credentialsDAO)
	updateProfileService := services.NewUpdateProfileService(profileDAO, introspectTokenService)
	validateEmailService := services.NewValidateEmailService(credentialsDAO)
	validateNewEmailService := services.NewValidateNewEmailService(credentialsDAO)

	pingHandler := handlers.NewPingHandler()
	healthCheckHandler := handlers.NewHealthCheckHandler(postgres)
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

	router := config.GetRouter(logger)

	router.GET("/ping", pingHandler.Handle)
	router.GET("/healthcheck", healthCheckHandler.Handle)
	router.GET("/auth", introspectTokenHandler.Handle)
	router.POST("/auth", loginHandler.Handle)
	router.PUT("/auth", registerHandler.Handle)
	router.DELETE("/email", cancelNewEmailHandler.Handle)
	router.DELETE("/password", resetPasswordHandler.Handle)
	router.PATCH("/email", updateEmailHandler.Handle)
	router.PATCH("/identity", updateIdentityHandler.Handle)
	router.PATCH("/profile", updateProfileHandler.Handle)
	router.PATCH("/password", updatePasswordHandler.Handle)
	router.PATCH("/email/validation", resendEmailValidationHandler.Handle)
	router.PATCH("/email/pending/validation", resendNewEmailValidationHandler.Handle)
	router.GET("/email/validation", validateEmailHandler.Handle)
	router.GET("/email/pending/validation", validateNewEmailHandler.Handle)
	router.GET("/email/exists", emailExistsHandler.Handle)
	router.GET("/slug/exists", slugExistsHandler.Handle)
	router.GET("/users", listHandler.Handle)
	router.GET("/users/search", searchHandler.Handle)
	router.GET("/user", previewHandler.Handle)
	router.GET("/user/me", previewPrivateHandler.Handle)

	if err := router.Run(fmt.Sprintf(":%d", config.API.Port)); err != nil {
		logger.Fatal().Err(err).Msg("a fatal error occurred while running the API, and the server had to shut down")
	}
}
