package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	ns "net/smtp"

	"github.com/gin-gonic/gin"
	"github.com/zen-en-tonal/mtw/database"
	"github.com/zen-en-tonal/mtw/database/address"
	"github.com/zen-en-tonal/mtw/database/webhook"
	"github.com/zen-en-tonal/mtw/forward"
	"github.com/zen-en-tonal/mtw/http"
	"github.com/zen-en-tonal/mtw/session"
	"github.com/zen-en-tonal/mtw/smtp"
)

var (
	domain string = ""

	smtpUser  string = ""
	smtpPass  string = ""
	smtpHost  string = ""
	forwardTo string = ""

	dbconn string = ""

	secret string = ""
)

func init() {
	domain, _ = os.LookupEnv("DOMAIN")
	smtpUser, _ = os.LookupEnv("SMTP_USER")
	smtpPass, _ = os.LookupEnv("SMTP_PASS")
	smtpHost, _ = os.LookupEnv("SMTP_HOST")
	forwardTo, _ = os.LookupEnv("FORWARD_TO")
	dbconn, _ = os.LookupEnv("DB_CONN")
	secret, _ = os.LookupEnv("SECRET")
}

func main() {
	logger := slog.Default()

	db, err := sql.Open("postgres", dbconn)
	if err != nil {
		logger.Error("failed to connect to db", "inner", err.Error())
		return
	}
	if err := database.Migrate(db); err != nil {
		logger.Warn("migration failure", "inner", err.Error())
	}

	if secret == "" {
		logger.Error("secret must be set")
		return
	}

	auth := ns.PlainAuth("", smtpUser, smtpPass, smtpHost)

	smtp := smtp.New(
		smtp.WithSessionOptions(
			session.AppendFilters(address.Find(db)),
			session.AppendHookProviders(webhook.NewFind(db)),
			session.AppendHooks(forward.NewSmtp(smtpHost, auth, forwardTo)),
			session.WithLogger(logger),
			session.WithTimeout(time.Second*5),
		),
		smtp.WithLogger(logger),
	)
	smtp.Addr = "0.0.0.0:25"
	smtp.Domain = domain
	smtp.AllowInsecureAuth = false

	rest := gin.New()
	rest.Use(authMiddle)
	http.SetRoutes(rest, db, domain, logger)

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx *context.Context) {
		logger.Info("Listening and serving SMTP on 0.0.0.0:25")
		if err := smtp.ListenAndServe(); err != nil {
			logger.Error("smtp", "inner", err.Error())
		}
		cancel()
	}(&ctx)
	go func(ctx *context.Context) {
		logger.Info("Listening and serving HTTP on 0.0.0.0:8080")
		if err := rest.Run("0.0.0.0:8080"); err != nil {
			logger.Error("http", "inner", err.Error())
		}
		cancel()
	}(&ctx)
	<-ctx.Done()
}

func authMiddle(c *gin.Context) {
	if c.GetHeader("Authorization") != "Bearer "+secret {
		c.Status(401)
		c.Abort()
	} else {
		c.Next()
	}
}
