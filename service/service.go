package service

import (
	"context"
	"fmt"
	"github.com/Knetic/govaluate"
	calculator "github.com/radutopala/calculator/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service is the service dealing with computes
type Service struct {
}

// Compute parses an expression and returns the result
func (s Service) Compute(ctx context.Context, req *calculator.Request) (*calculator.Response, error) {
	expression, err := govaluate.NewEvaluableExpression(req.Expression)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Expression may have invalid syntax: %s", err)
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not evaluate the expression: %s", err)
	}

	return &calculator.Response{Result: fmt.Sprintf("%v", result)}, nil
}
