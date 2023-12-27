package cmdline

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
)

const (
	Host = "db-host"
	Port = "db-port"
	Name = "db-name"

	User     = "db-user"
	Password = "db-password"
)

func (a *App) ConnectDBFromCobra(cmd *cobra.Command) (*pgx.Conn, error) {
	f := cmd.Flags()
	name, err := f.GetString(Name)
	if err != nil {
		return nil, err
	}

	host, err := f.GetString(Host)
	if err != nil {
		return nil, err
	}

	port, err := f.GetUint16(Port)
	if err != nil {
		return nil, err
	}

	user, err := f.GetString(User)
	if err != nil {
		return nil, err
	}

	password, err := f.GetString(Password)
	if err != nil {
		return nil, err
	}

	return a.ConnectDB(host, name, user, password, port)
}

// ConnectDB opens a database, pings it, and returns it if all succeeds.
// Will determine if sslmode is needed based on the host you're connecting to
func (a *App) ConnectDB(host, name, user, password string, port uint16) (*pgx.Conn, error) {
	sslmode := useSSL(host)
	a.logger.Info("Connecting database with no middleware",
		"host", host,
		"port", port,
		"sslmode", sslmode,
		"name", name,
		"user", user,
		"len(password) > 0", len(password) > 0,
	)

	db, err := pgx.Connect(context.Background(), genConnStr(port, host, name, user, password, sslmode))
	if err != nil {
		a.logger.Error("failed connecting to database", "err", err)
		return nil, err
	}

	return db, nil
}

func useSSL(host string) string {
	if host == "localhost" || host == "127.0.0.1" {
		return "disable" // sslmode unnecessary if localhost or 127.0.0.1
	}

	return "require"
}

func genConnStr(port uint16, host, name, user, password, useSSL string) string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%v sslmode=%v user=%s password=%s",
		host,
		port,
		name,
		useSSL,
		user,
		password,
	)
}
