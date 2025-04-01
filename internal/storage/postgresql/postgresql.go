package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"order-service/pkg/api/test"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     uint16 `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-required:"true"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-required:"true"`
	MaxConn  int32  `yaml:"POSTGRES_MAX_CONN" env:"POSTGRES_MAX_CONN"`
	MinConn  int32  `yaml:"POSTGRES_MIN_CONN" env:"POSTGRES_MIN_CONN"`
}

type Storage struct {
	db *pgxpool.Pool
}

func New(cfg Config) (*Storage, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.MaxConn,
		cfg.MinConn,
	)
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := Migrate(conn); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return &Storage{db: conn}, nil
}

func Migrate(s *pgxpool.Pool) error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		s.Config().ConnConfig.User,
		s.Config().ConnConfig.Password,
		s.Config().ConnConfig.Host,
		s.Config().ConnConfig.Port,
		s.Config().ConnConfig.Database,
	)
	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
func (s *Storage) Stop() {
	s.db.Close()
}

func (s *Storage) Add(ctx context.Context, order *test.Order) error {
	query := `INSERT INTO orders (id, item, quantity) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(ctx, query, order.GetId(), order.GetItem(), order.GetQuantity())
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, id string) (*test.Order, error) {
	query := `SELECT id, item, quantity FROM orders WHERE id = $1`
	row := s.db.QueryRow(ctx, query, id)

	var order test.Order
	err := row.Scan(&order.Id, &order.Item, &order.Quantity)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *Storage) Update(ctx context.Context, order *test.Order) (*test.Order, error) {
	query := `UPDATE orders SET item = $2, quantity = $3 WHERE id = $1`
	_, err := s.db.Exec(ctx, query, order.GetId(), order.GetItem(), order.GetQuantity())
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	updatedOrder, err := s.Get(ctx, order.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %w", err)
	}

	return updatedOrder, nil
}
func (s *Storage) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetAll(ctx context.Context) ([]*test.Order, error) {
	query := `SELECT id, item, quantity FROM orders`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*test.Order
	for rows.Next() {
		var order test.Order
		if err := rows.Scan(&order.Id, &order.Item, &order.Quantity); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
