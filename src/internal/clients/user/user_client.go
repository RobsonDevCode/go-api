package userClient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/RobsonDevCode/go-profile-service/src/internal/caching"
	responses "github.com/RobsonDevCode/go-profile-service/src/internal/clients/user/responses"
	"github.com/RobsonDevCode/go-profile-service/src/internal/config"
	"github.com/google/uuid"
	"github.com/sony/gobreaker"
)

type UserClient struct {
	client  *http.Client
	cb      *gobreaker.CircuitBreaker
	baseUrl *url.URL
	jwt     string
	jwtLock sync.RWMutex
	cache   *caching.Cache
}

func NewUserClient(config config.Config, cache *caching.Cache) (*UserClient, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	cbSettings := gobreaker.Settings{
		Name:        "user-client",
		MaxRequests: 0,
		Interval:    0,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker state changed from %v to %v\n", from, to)
		},
	}

	cb := gobreaker.NewCircuitBreaker(cbSettings)
	baseUrl, err := url.Parse(config.UserClientOptions.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	return &UserClient{
		client:  client,
		baseUrl: baseUrl,
		cb:      cb,
		cache:   cache,
	}, nil
}

func (c *UserClient) SetJwt(token string) {
	c.jwtLock.Lock()
	defer c.jwtLock.Unlock()

	c.jwt = token
}

func (c *UserClient) Get(id uuid.UUID, ctx context.Context) (User, error) {
	url := fmt.Sprintf("%s/%s", c.baseUrl, id)

	key := fmt.Sprintf("user-%v", id)

	result, err := c.cache.GetOrCreate(key, time.Minute*5, func() (interface{}, error) {
		result, err := c.cb.Execute(func() (interface{}, error) {
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return User{}, fmt.Errorf("failed to get user: %w", err)
			}

			response, err := c.client.Do(request)
			if err != nil {
				return User{}, fmt.Errorf("client error: %w", err)
			}
			defer response.Body.Close()

			if response.StatusCode >= 500 {
				return User{}, fmt.Errorf("server error on user client: %d", response.StatusCode)
			} else if response.StatusCode == 400 {
				var problemDetails responses.ProblemDetails
				if err := json.NewDecoder(request.Body).Decode(&problemDetails); err != nil {
					return User{}, fmt.Errorf("failed to decode bad request from user client, %w", err)
				}

				return User{}, fmt.Errorf("bad request: %v", problemDetails)
			}

			var user User

			if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
				return User{}, fmt.Errorf("failed to decode response")
			}

			return user, nil
		})
		if err != nil {
			if err == gobreaker.ErrOpenState {
				return User{}, fmt.Errorf("service unavailable, circuit open: %w", err)
			}

			return User{}, fmt.Errorf("circuit breaker error: %w", err)
		}

		return result, nil
	})

	if err != nil {
		return User{}, err
	}

	user, ok := result.(User)
	if !ok {
		return User{}, fmt.Errorf("unexpected response type")
	}

	return user, nil
}

func (c *UserClient) UserExists(id uuid.UUID, ctx context.Context) (bool, error) {

	user, err := c.Get(id, ctx)
	if err != nil {
		return false, err
	}

	if user.Id != uuid.Nil {
		return true, nil
	}

	return false, nil
}
