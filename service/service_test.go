package service

import (
	"context"
	calculator "github.com/radutopala/grpc-calculator/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompute(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			"1+2",
			"3",
		},
		{
			"1+2-5",
			"-2",
		},
		{
			"100.5*2/2",
			"100.5",
		},
		{
			"5.55*(10-1)",
			"49.949999999999996",
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			s := &Service{}
			response, err := s.Compute(
				context.Background(),
				&calculator.Request{
					Expression: c.input,
				},
			)

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, response.Result, c.expected)
		})
	}
}
