package datasync_workflow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dyntypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	sqlmanager_shared "github.com/nucleuscloud/neosync/backend/pkg/sqlmanager/shared"
	awsmanager "github.com/nucleuscloud/neosync/internal/aws"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	testmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	testpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/sync/errgroup"
)

type postgresTestContainer struct {
	pool *pgxpool.Pool
	url  string
}
type postgresTest struct {
	pool          *pgxpool.Pool
	testcontainer *testpg.PostgresContainer

	source *postgresTestContainer
	target *postgresTestContainer

	databases []string
}

type mysqlTestContainer struct {
	pool      *sql.DB
	container *testmysql.MySQLContainer
	url       string
	close     func()
}

type mysqlTest struct {
	source *mysqlTestContainer
	target *mysqlTestContainer
}

type redisTest struct {
	url           string
	testcontainer *redis.RedisContainer
}

type IntegrationTestSuite struct {
	suite.Suite

	ctx context.Context

	mysql    *mysqlTest
	postgres *postgresTest
	redis    *redisTest
	dynamo   *dynamodbTest
}

func (s *IntegrationTestSuite) SetupPostgres() (*postgresTest, error) {
	pgcontainer, err := testpg.Run(
		s.ctx,
		"postgres:15",
		postgres.WithDatabase("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(20*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}
	postgresTest := &postgresTest{
		testcontainer: pgcontainer,
	}
	connstr, err := pgcontainer.ConnectionString(s.ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	postgresTest.databases = []string{"datasync_source", "datasync_target"}
	pool, err := pgxpool.New(s.ctx, connstr)
	if err != nil {
		return nil, err
	}
	postgresTest.pool = pool

	s.T().Logf("creating databases. %+v \n", postgresTest.databases)
	for _, db := range postgresTest.databases {
		_, err = postgresTest.pool.Exec(s.ctx, fmt.Sprintf("CREATE DATABASE %s;", db))
		if err != nil {
			return nil, err
		}
	}

	srcUrl, err := getDbPgUrl(connstr, "datasync_source", "disable")
	if err != nil {
		return nil, err
	}
	postgresTest.source = &postgresTestContainer{
		url: srcUrl,
	}
	sourceConn, err := pgxpool.New(s.ctx, postgresTest.source.url)
	if err != nil {
		return nil, err
	}
	postgresTest.source.pool = sourceConn

	targetUrl, err := getDbPgUrl(connstr, "datasync_target", "disable")
	if err != nil {
		return nil, err
	}
	postgresTest.target = &postgresTestContainer{
		url: targetUrl,
	}
	targetConn, err := pgxpool.New(s.ctx, postgresTest.target.url)
	if err != nil {
		return nil, err
	}
	postgresTest.target.pool = targetConn
	return postgresTest, nil
}

func (s *IntegrationTestSuite) SetupMysql() (*mysqlTest, error) {
	var source *mysqlTestContainer
	var target *mysqlTestContainer

	errgrp := errgroup.Group{}
	errgrp.Go(func() error {
		sourcecontainer, err := createMysqlTestContainer(s.ctx, "datasync", "root", "pass-source")
		if err != nil {
			return err
		}
		source = sourcecontainer
		return nil
	})

	errgrp.Go(func() error {
		targetcontainer, err := createMysqlTestContainer(s.ctx, "datasync", "root", "pass-target")
		if err != nil {
			return err
		}
		target = targetcontainer
		return nil
	})

	err := errgrp.Wait()
	if err != nil {
		return nil, err
	}

	return &mysqlTest{
		source: source,
		target: target,
	}, nil
}

func createMysqlTestContainer(
	ctx context.Context,
	database, username, password string,
) (*mysqlTestContainer, error) {
	container, err := testmysql.Run(ctx,
		"mysql:8.0.36",
		testmysql.WithDatabase(database),
		testmysql.WithUsername(username),
		testmysql.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server").
				WithOccurrence(1).WithStartupTimeout(20*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}
	connstr, err := container.ConnectionString(ctx, "multiStatements=true")
	if err != nil {
		panic(err)
	}
	pool, err := sql.Open(sqlmanager_shared.MysqlDriver, connstr)
	if err != nil {
		panic(err)
	}
	containerPort, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return nil, err
	}
	containerHost, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	connUrl := fmt.Sprintf("mysql://%s:%s@%s:%s/%s?multiStatements=true", username, password, containerHost, containerPort.Port(), database)
	return &mysqlTestContainer{
		pool:      pool,
		url:       connUrl,
		container: container,
		close: func() {
			if pool != nil {
				pool.Close()
			}
		},
	}, nil
}

func (s *IntegrationTestSuite) SetupRedis() (*redisTest, error) {
	redisContainer, err := redis.Run(
		s.ctx,
		"docker.io/redis:7",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, err
	}
	redisUrl, err := redisContainer.ConnectionString(s.ctx)
	if err != nil {
		return nil, err
	}
	return &redisTest{
		testcontainer: redisContainer,
		url:           redisUrl,
	}, nil
}

type dynamodbTest struct {
	container testcontainers.Container
	endpoint  string

	// Used by plugging in to Neosync resources so Benthos can wire up its aws config
	dtoAwsCreds *mgmtv1alpha1.AwsS3Credentials

	dynamoclient *dynamodb.Client
}

func (s *IntegrationTestSuite) SetupDynamoDB() (*dynamodbTest, error) {
	port := nat.Port("8000/tcp")
	container, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "amazon/dynamodb-local:2.5.2",
			ExposedPorts: []string{string(port)},
			WaitingFor:   wait.ForListeningPort(port),
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}

	mappedport, err := container.MappedPort(s.ctx, port)
	if err != nil {
		return nil, err
	}
	host, err := container.Host(s.ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("http://%s:%d", host, mappedport.Int())
	fakeId := "fakeid"
	fakeSecret := "fakesecret"
	fakeToken := "faketoken"

	awscfg, err := awsmanager.GetAwsConfig(s.ctx, &awsmanager.AwsCredentialsConfig{
		Endpoint: endpoint,
		Id:       fakeId,
		Secret:   fakeSecret,
		Token:    fakeToken,
	})
	if err != nil {
		return nil, err
	}

	dtoAwsCreds := &mgmtv1alpha1.AwsS3Credentials{
		AccessKeyId:     &fakeId,
		SecretAccessKey: &fakeSecret,
		SessionToken:    &fakeToken,
	}

	return &dynamodbTest{
		container:   container,
		endpoint:    endpoint,
		dtoAwsCreds: dtoAwsCreds,
		dynamoclient: dynamodb.NewFromConfig(*awscfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = &endpoint
		}),
	}, nil
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()

	var postgresTest *postgresTest
	var mysqlTest *mysqlTest
	var redisTest *redisTest
	var dynamoTest *dynamodbTest

	errgrp := errgroup.Group{}
	errgrp.Go(func() error {
		p, err := s.SetupPostgres()
		if err != nil {
			return err
		}
		postgresTest = p
		return nil
	})

	errgrp.Go(func() error {
		m, err := s.SetupMysql()
		if err != nil {
			return err
		}
		mysqlTest = m
		return nil
	})

	errgrp.Go(func() error {
		r, err := s.SetupRedis()
		if err != nil {
			return err
		}
		redisTest = r
		return nil
	})

	errgrp.Go(func() error {
		d, err := s.SetupDynamoDB()
		if err != nil {
			return err
		}
		dynamoTest = d
		return nil
	})

	err := errgrp.Wait()
	if err != nil {
		panic(err)
	}

	s.postgres = postgresTest
	s.mysql = mysqlTest
	s.redis = redisTest
	s.dynamo = dynamoTest
}

func (s *IntegrationTestSuite) RunPostgresSqlFiles(pool *pgxpool.Pool, testFolder string, files []string) {
	s.T().Logf("running postgres sql file. folder: %s \n", testFolder)
	for _, file := range files {
		sqlStr, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", testFolder, file))
		if err != nil {
			panic(err)
		}
		_, err = pool.Exec(s.ctx, string(sqlStr))
		if err != nil {
			panic(err)
		}
	}
}

func (s *IntegrationTestSuite) RunMysqlSqlFiles(pool *sql.DB, testFolder string, files []string) {
	s.T().Logf("running mysql sql file. folder: %s \n", testFolder)
	for _, file := range files {
		sqlStr, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", testFolder, file))
		if err != nil {
			panic(err)
		}

		_, err = pool.ExecContext(s.ctx, string(sqlStr))
		if err != nil {
			panic(err)
		}
	}
}

func (s *IntegrationTestSuite) SetupDynamoDbTable(ctx context.Context, tableName, primaryKey string) error {
	s.T().Logf("Creating DynamoDB table: %s\n", tableName)

	out, err := s.dynamo.dynamoclient.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:            &tableName,
		KeySchema:            []dyntypes.KeySchemaElement{{KeyType: dyntypes.KeyTypeHash, AttributeName: &primaryKey}},
		AttributeDefinitions: []dyntypes.AttributeDefinition{{AttributeName: &primaryKey, AttributeType: dyntypes.ScalarAttributeTypeS}},
		BillingMode:          dyntypes.BillingModePayPerRequest,
	})
	if err != nil {
		return err
	}
	if out.TableDescription.TableStatus == dyntypes.TableStatusActive {
		return nil
	}
	if out.TableDescription.TableStatus == dyntypes.TableStatusCreating {
		return s.waitUntilDynamoTableExists(ctx, tableName)
	}
	return fmt.Errorf("%s dynamo table created but unexpected table status: %s", tableName, out.TableDescription.TableStatus)
}

func (s *IntegrationTestSuite) waitUntilDynamoTableExists(ctx context.Context, tableName string) error {
	input := &dynamodb.DescribeTableInput{TableName: &tableName}
	for {
		out, err := s.dynamo.dynamoclient.DescribeTable(ctx, input)
		if err != nil && !awsmanager.IsNotFound(err) {
			return err
		}
		if err != nil && awsmanager.IsNotFound(err) {
			continue
		}
		if out.Table.TableStatus == dyntypes.TableStatusActive {
			return nil
		}
	}
}

func (s *IntegrationTestSuite) DestroyDynamoDbTable(ctx context.Context, tableName string) error {
	s.T().Logf("Destroying DynamoDB table: %s\n", tableName)

	_, err := s.dynamo.dynamoclient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: &tableName,
	})
	if err != nil {
		return err
	}
	return s.waitUntilDynamoTableDestroy(ctx, tableName)
}

func (s *IntegrationTestSuite) waitUntilDynamoTableDestroy(ctx context.Context, tableName string) error {
	input := &dynamodb.DescribeTableInput{TableName: &tableName}
	for {
		_, err := s.dynamo.dynamoclient.DescribeTable(ctx, input)
		if err != nil && !awsmanager.IsNotFound(err) {
			return err
		}
		if err != nil && awsmanager.IsNotFound(err) {
			return nil
		}
	}
}

func (s *IntegrationTestSuite) InsertDynamoDBRecords(tableName string, data []map[string]dyntypes.AttributeValue) error {
	s.T().Logf("Inserting %d DynamoDB Records into table: %s\n", len(data), tableName)

	writeRequests := make([]dyntypes.WriteRequest, len(data))
	for i, record := range data {
		writeRequests[i] = dyntypes.WriteRequest{
			PutRequest: &dyntypes.PutRequest{
				Item: record,
			},
		}
	}

	_, err := s.dynamo.dynamoclient.BatchWriteItem(s.ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]dyntypes.WriteRequest{
			tableName: writeRequests,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down test suite")
	// postgres
	for _, db := range s.postgres.databases {
		_, err := s.postgres.pool.Exec(s.ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE);", db))
		if err != nil {
			panic(err)
		}
	}
	if s.postgres.source.pool != nil {
		s.postgres.source.pool.Close()
	}
	if s.postgres.target.pool != nil {
		s.postgres.target.pool.Close()
	}
	if s.postgres.pool != nil {
		s.postgres.pool.Close()
	}
	if s.postgres.testcontainer != nil {
		err := s.postgres.testcontainer.Terminate(s.ctx)
		if err != nil {
			panic(err)
		}
	}

	// mysql
	s.mysql.source.close()
	s.mysql.target.close()
	if s.mysql.source.container != nil {
		err := s.mysql.source.container.Terminate(s.ctx)
		if err != nil {
			panic(err)
		}
	}
	if s.mysql.target.container != nil {
		err := s.mysql.target.container.Terminate(s.ctx)
		if err != nil {
			panic(err)
		}
	}

	// redis
	if s.redis.testcontainer != nil {
		if err := s.redis.testcontainer.Terminate(s.ctx); err != nil {
			panic(err)
		}
	}

	// localstack
	if s.dynamo.container != nil {
		if err := s.dynamo.container.Terminate(s.ctx); err != nil {
			panic(err)
		}
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	evkey := "INTEGRATION_TESTS_ENABLED"
	shouldRun := os.Getenv(evkey)
	if shouldRun != "1" {
		slog.Warn(fmt.Sprintf("skipping integration tests, set %s=1 to enable", evkey))
		return
	}
	suite.Run(t, new(IntegrationTestSuite))
}

func getDbPgUrl(dburl, database, sslmode string) (string, error) {
	u, err := url.Parse(dburl)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			return "", fmt.Errorf("unable to parse postgres url [%s]: %w", urlErr.Op, urlErr.Err)
		}
		return "", fmt.Errorf("unable to parse postgres url: %w", err)
	}

	u.Path = database
	query := u.Query()
	query.Add("sslmode", sslmode)
	return u.String(), nil
}
