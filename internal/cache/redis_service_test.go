package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type redisContainer struct {
	testcontainers.Container
	host string
	port int
}

func setupRedis(ctx context.Context) (*redisContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &redisContainer{Container: container, host: hostIP, port: mappedPort.Int()}, nil
}

func TestIntegrationSetGetDel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	redisContainer, err := setupRedis(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = redisContainer.Terminate(ctx) }()

	client, err := NewRedisCache(ctx, redisContainer.host, redisContainer.port)
	if err != nil {
		t.Fatal(err)
	}

	// Set data
	key := "key"
	value := "value"
	ttl, _ := time.ParseDuration("2h")

	assert.NoError(t, client.Set(key, value, ttl))

	// Get data
	savedValue, err := client.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, savedValue)

	assert.NoError(t, client.Delete(key))

	_, err = client.Get(key)
	assert.ErrorIs(t, err, ErrNotFound)
}
