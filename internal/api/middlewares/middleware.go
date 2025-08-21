// Package middlewares provides HTTP middleware functions for the API.
// It includes authentication and authorization checks.
package middlewares

import "github.com/uchupx/saceri-chatbot-api/pkg/grpc/client"

type Middleware struct {
	AuthClient *client.AuthClient
}
