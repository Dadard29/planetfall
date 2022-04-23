package main

import (
	"context"
	"fmt"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func (s *Server) getSecret(secretName string) (string, error) {
	ctx := context.Background()

	secretPath := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", s.projectID, secretName)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	result, err := s.secretManager.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}
	return string(result.Payload.Data), nil
}
