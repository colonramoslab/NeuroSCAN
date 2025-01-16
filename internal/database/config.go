package database

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	Name               string
	User               string
	Host               string
	Port               string
	SSLMode            string
	ConnectionTimeout  int
	Password           string
	SSLCertPath        string
	SSLKeyPath         string
	SSLRootCertPath    string
	PoolMinConnections string
	PoolMaxConnections string
	PoolMaxConnLife    time.Duration
	PoolMaxConnIdle    time.Duration
	PoolHealthCheck    time.Duration
}

func (c *Config) DatabaseConfig() *Config {
	return c
}

func NewConfigFromEnv() *Config {
	config := &Config{}

	err := godotenv.Load()
	if err != nil {
		return config
	}

	dsn := os.Getenv("DB_DSN")
	if dsn != "" {
		config, _ := NewFromDSN(dsn)
		return config
	}

	config.Name = os.Getenv("DB_NAME")
	config.User = os.Getenv("DB_USER")
	config.Host = os.Getenv("DB_HOST")
	config.Port = os.Getenv("DB_PORT")
	config.SSLMode = os.Getenv("DB_SSLMODE")
	config.Password = os.Getenv("DB_PASSWORD")
	config.SSLCertPath = os.Getenv("DB_SSLCERT")
	config.SSLKeyPath = os.Getenv("DB_SSLKEY")
	config.SSLRootCertPath = os.Getenv("DB_SSLROOTCERT")
	config.PoolMinConnections = os.Getenv("DB_POOL_MIN_CONNS")
	config.PoolMaxConnections = os.Getenv("DB_POOL_MAX_CONNS")

	timeout, err := strconv.Atoi(os.Getenv("DB_CONNECT_TIMEOUT"))
	if err != nil {
		config.ConnectionTimeout = 0
	} else {
		config.ConnectionTimeout = timeout
	}

	poolMaxConnLife, err := time.ParseDuration(os.Getenv("DB_POOL_MAX_CONN_LIFETIME"))
	if err != nil {
		config.PoolMaxConnLife = 5 * time.Minute
	} else {
		config.PoolMaxConnLife = poolMaxConnLife
	}

	poolMaxConnIdle, err := time.ParseDuration(os.Getenv("DB_POOL_MAX_CONN_IDLE_TIME"))
	if err != nil {
		config.PoolMaxConnIdle = 1 * time.Minute
	} else {
		config.PoolMaxConnIdle = poolMaxConnIdle
	}

	poolHealthCheck, err := time.ParseDuration(os.Getenv("DB_POOL_HEALTH_CHECK_PERIOD"))
	if err != nil {
		config.PoolHealthCheck = 1 * time.Minute
	} else {
		config.PoolHealthCheck = poolHealthCheck
	}

	return config
}

func NewFromDSN(dsn string) (*Config, error) {
	cfg := &Config{}

	vals, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	port := vals.Config.Port
	// convert to string
	portStr := strconv.FormatUint(uint64(port), 10)

	cfg.Name = vals.Config.Database
	cfg.User = vals.Config.User
	cfg.Host = vals.Config.Host
	cfg.Password = vals.Config.Password
	cfg.Port = portStr

	return cfg, nil
}

func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	host := c.Host
	if v := c.Port; v != "" {
		host = host + ":" + v
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	q := u.Query()
	if v := c.ConnectionTimeout; v > 0 {
		q.Add("connect_timeout", strconv.Itoa(v))
	}
	if v := c.SSLMode; v != "" {
		q.Add("sslmode", v)
	}
	if v := c.SSLCertPath; v != "" {
		q.Add("sslcert", v)
	}
	if v := c.SSLKeyPath; v != "" {
		q.Add("sslkey", v)
	}
	if v := c.SSLRootCertPath; v != "" {
		q.Add("sslrootcert", v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
