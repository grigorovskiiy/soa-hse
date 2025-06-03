package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/fx"

	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/logger"
)

func CreateDbTables(conn *sql.DB) error {
	var queries []string
	tableNames := []string{"comments", "likes", "views"}
	for _, name := range tableNames {
		queries = append(queries, fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s (
                		time DateTime('UTC'),
  						user_id Int32,
  						post_id Int32
        			)  
            		ENGINE = MergeTree() 
 					PARTITION BY toYYYYMM(time)
 					ORDER BY (time)`, name))
	}

	for _, query := range queries {
		_, err := conn.Exec(query)
		if err != nil {
			logger.Logger.Error("error execing query", "error", err)
			return err
		}
	}

	return nil
}

func CreateKafkaTables(conn *sql.DB) error {
	var queries []string
	tableNames := []string{"commentskafka", "likeskafka", "viewskafka"}
	topics := []string{"comments.topic", "likes.topic", "views.topic"}
	for i := range tableNames {
		queries = append(queries, fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s (
  						time DateTime('UTC'),
  						user_id Int32,
  						post_id Int32
 					)
 					ENGINE = Kafka()
 					SETTINGS kafka_broker_list = 'kafka:9092',
    				kafka_topic_list = '%s',
    				kafka_group_name = 'clickhouse_%s_consumer',
					kafka_format = 'JSONEachRow',
					kafka_num_consumers = 1,
					kafka_skip_broken_messages = 1,
					date_time_input_format = 'best_effort'`, tableNames[i], topics[i], topics[i]))

	}

	for _, query := range queries {
		_, err := conn.Exec(query)
		if err != nil {
			logger.Logger.Error("error execing query", "error", err)
			return err
		}
	}

	return nil
}

func CreateMwTables(conn *sql.DB) error {
	var queries []string
	mwTableNames := []string{"commentsmw", "likesmw", "viewsmw"}
	dbTableNames := []string{"comments", "likes", "views"}
	kafkaTableNames := []string{"commentskafka", "likeskafka", "viewskafka"}

	for i := range mwTableNames {
		queries = append(queries, fmt.Sprintf(
			`CREATE MATERIALIZED VIEW IF NOT EXISTS %s
   						TO %s
   						AS 
   						SELECT 
    						time,
    						user_id,
    						post_id
						FROM %s`, mwTableNames[i], dbTableNames[i], kafkaTableNames[i]))
	}

	for _, query := range queries {
		_, err := conn.Exec(query)
		if err != nil {
			logger.Logger.Error("error execing query", "error", err)
			return err
		}
	}

	return nil
}

func InitDb(lc fx.Lifecycle, cfg *config.Config) (*sql.DB, error) {
	logger.Logger.Info("Connecting to ClickHouse")

	conn := sql.OpenDB(clickhouse.Connector(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s%s", cfg.ClickHouseHost, cfg.ClickHousePort)},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouseDb,
			Username: cfg.ClickHouseUser,
			Password: cfg.ClickHousePassword,
		},
	}))

	if err := conn.Ping(); err != nil {
		logger.Logger.Error("ping clickhouse error", "error", err.Error())
		return nil, err
	}

	if err := CreateDbTables(conn); err != nil {
		logger.Logger.Error("create db tables error", "error", err.Error())
		return nil, err
	}

	if err := CreateKafkaTables(conn); err != nil {
		logger.Logger.Error("create kafka tables error", "error", err.Error())
		return nil, err
	}

	if err := CreateMwTables(conn); err != nil {
		logger.Logger.Error("create mw tables error", "error", err.Error())
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return conn, nil
}
