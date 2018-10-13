package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/rk23/tdgo/pkg/cache"
	"github.com/rk23/tdgo/pkg/oauth"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/namsral/flag"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	clientID    = flag.String("client_id", "", "TD Ameritrade Client ID ending (and including) @AMER.OAUTHAP")
	redirectURI = flag.String("redirect_uri", "", "Redirect URI as set in your TD Ameritrade app")
)

type cacher interface {
	Put(*cache.PutParamInput) error
	BatchPut([]*cache.PutParamInput) error
	Get([]*string, bool) (map[string]string, error)
}

type server struct {
	log   zerolog.Logger
	cache cacher
}

type jsonres struct {
	Data   map[string]string `json:"data"`
	Errors map[string]string `json:"errors"`
}

func jsonhttp(data map[string]string, errors map[string]string) string {
	if data == nil {
		data = map[string]string{}
	}
	if errors == nil {
		errors = map[string]string{}
	}

	res := &jsonres{
		Data:   data,
		Errors: errors,
	}
	r, _ := json.Marshal(res)
	return string(r)
}

func (s *server) RegisterTokens(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	code := req.QueryStringParameters["code"]

	res, err := oauth.AccessToken(&oauth.AccessTokenRequest{
		Code:        code,
		AccessType:  "offline",
		ClientID:    *clientID,
		RedirectURI: *redirectURI,
	})
	if err != nil {
		s.log.Error().Str("error", err.Error()).Msg("err in access token")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	err = s.cache.BatchPut(
		[]*cache.PutParamInput{&cache.PutParamInput{
			Description: "Access Token for the Aurifodina",
			Name:        "accesstoken",
			Value:       res.AccessToken,
		}, &cache.PutParamInput{
			Description: "Refresh Token for the Aurifodina",
			Name:        "refreshtoken",
			Value:       res.RefreshToken,
		}})

	if err != nil {
		s.log.Error().
			Str("error", err.Error()).
			Msg("Error saving params, continuing")
	}

	b, err := json.Marshal(res)
	if err != nil {
		s.log.Error().Msg("Err in marshallng res")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (s *server) router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/":
		return s.RegisterTokens(ctx, req)
	case "/cache":
		p, err := s.cache.Get([]*string{aws.String("accesstoken"), aws.String("refreshtoken")}, true)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       jsonhttp(nil, map[string]string{"description": err.Error()}),
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       jsonhttp(p, nil),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotImplemented,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	debug := flag.Bool("debug", true, "sets log level to debug")
	flag.Parse()

	zerolog.TimeFieldFormat = ""
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Debug().Msg("starting handler")
	log := zerolog.New(os.Stdout)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Error().Msgf("failed to start lambda, err in aws session: %s", err.Error())
		panic(err)
	}

	sc := ssm.New(sess)
	srv := &server{
		log: log,
		cache: &cache.SSM{
			Log:    log,
			Client: sc,
		},
	}
	lambda.Start(srv.router)
}
