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
	Data   interface{}         `json:"data"`
	Errors []map[string]string `json:"errors"`
}

func jsonhttp(data interface{}, errs []string) string {
	errors := []map[string]string{}

	if data == nil {
		data = map[string]string{}
	}
	if len(errs) != 0 {
		for _, err := range errs {
			errors = append(errors, map[string]string{"description": err})
		}
	}

	res := &jsonres{
		Data:   data,
		Errors: errors,
	}
	r, _ := json.Marshal(res)
	return string(r)
}

func httpres(code int, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (s *server) router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/":
		code := req.QueryStringParameters["code"]
		errs := []string{}

		res, err := oauth.AccessToken(&oauth.AccessTokenRequest{
			Code:        code,
			AccessType:  "offline",
			ClientID:    *clientID,
			RedirectURI: *redirectURI,
		})
		if err != nil {
			errs = append(errs, err.Error())
			s.log.Error().Str("error", err.Error()).Msg("err in access token")
			return httpres(http.StatusInternalServerError, jsonhttp(nil, errs))
		}

		err = s.cache.BatchPut(
			[]*cache.PutParamInput{&cache.PutParamInput{
				Description: "TD Ameritrade Access Token",
				Name:        "accesstoken",
				Value:       res.AccessToken,
			}, &cache.PutParamInput{
				Description: "TD Ameritrade Refresh Token",
				Name:        "refreshtoken",
				Value:       res.RefreshToken,
			}})

		if err != nil {
			s.log.Error().
				Str("error", err.Error()).
				Msg("Error saving params, continuing")
			errs = append(errs, err.Error())
		}

		return httpres(http.StatusOK, jsonhttp(res, errs))
	case "/cache":
		p, err := s.cache.Get([]*string{aws.String("accesstoken"), aws.String("refreshtoken")}, true)
		if err != nil {
			return httpres(http.StatusInternalServerError, jsonhttp(nil, []string{err.Error()}))
		}
		return httpres(http.StatusOK, jsonhttp(p, nil))
	}
	return httpres(http.StatusNotFound, jsonhttp(nil, []string{"Route not found"}))
}

func main() {
	debug := flag.Bool("debug", true, "sets log level to debug")
	flag.Parse()

	zerolog.TimeFieldFormat = ""
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
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

	log.Debug().Msg("starting handler")
	lambda.Start(srv.router)
}
