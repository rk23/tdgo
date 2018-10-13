package cache

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rs/zerolog"
)

// SSM base type for ssm
type SSM struct {
	Client *ssm.SSM
	Log    zerolog.Logger
}

// PutParamInput type for saving params
type PutParamInput struct {
	Name        string
	Description string
	Value       string
}

// Put params to store
func (s *SSM) Put(opts *PutParamInput) error {
	_, err := s.Client.PutParameter(&ssm.PutParameterInput{
		Description: aws.String(opts.Description),
		Name:        aws.String(opts.Name),
		Overwrite:   aws.Bool(true),
		Type:        aws.String("SecureString"),
		Value:       aws.String(opts.Value),
	})

	if err != nil {
		s.Log.Error().Msgf("Error putting access token to SSM: %s", err.Error())
		return err
	}

	return nil
}

// BatchPut simply iterates over a list of PutParamInputs
func (s *SSM) BatchPut(opts []*PutParamInput) error {
	for _, opt := range opts {
		err := s.Put(opt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get params
func (s *SSM) Get(params []*string, decrypt bool) (map[string]string, error) {
	o, err := s.Client.GetParameters(&ssm.GetParametersInput{
		Names:          params,
		WithDecryption: aws.Bool(decrypt),
	})

	if err != nil {
		return nil, err
	}

	p := make(map[string]string)
	for _, param := range o.Parameters {
		p[*param.Name] = *param.Value
	}

	return p, nil
}
