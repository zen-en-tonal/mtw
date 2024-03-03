package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"time"

	"github.com/zen-en-tonal/mtw/database"
	"github.com/zen-en-tonal/mtw/database/address"
	"github.com/zen-en-tonal/mtw/database/webhook"
	"github.com/zen-en-tonal/mtw/http"
	"github.com/zen-en-tonal/mtw/mailbox"
	"github.com/zen-en-tonal/mtw/session"
	"github.com/zen-en-tonal/mtw/smtp"
	"github.com/zen-en-tonal/mtw/spam"
)

var (
	domain string = ""
)

func init() {
	flag.StringVar(&domain, "d", domain, "Domain")
}

func main() {
	flag.Parse()

	logger := slog.Default()

	db, err := sql.Open(
		"postgres",
		"postgres://postgres:postgres@db:5432/postgres?sslmode=disable",
	)
	if err != nil {
		logger.Error("failed to connect to db", "inner", err.Error())
		return
	}

	if err := database.Migrate(db); err != nil {
		logger.Warn("migration failure", "inner", err.Error())
	}

	smtp := smtp.New(
		smtp.WithSessionOptions(
			session.WithFilters(
				spam.RcptMismatchFilter(),
				address.Find(db),
			),
			session.WithHooksSome(
				mailbox.AsHook(webhook.NewFind(db)),
			),
			session.WithLogger(logger),
			session.WithTimeout(time.Second*5),
		),
		smtp.WithLogger(logger),
	)
	smtp.Addr = "0.0.0.0:25"
	smtp.Domain = domain
	smtp.AllowInsecureAuth = false

	rest := http.NewWithDB(db, domain, logger)

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx *context.Context) {
		logger.Info("Listening and serving SMTP on 0.0.0.0:25")
		if err := smtp.ListenAndServe(); err != nil {
			logger.Error("smtp", "inner", err.Error())
		}
		cancel()
	}(&ctx)
	go func(ctx *context.Context) {
		if err := rest.Run("0.0.0.0:8080"); err != nil {
			logger.Error("rest", "inner", err.Error())
		}
		cancel()
	}(&ctx)
	<-ctx.Done()
}
