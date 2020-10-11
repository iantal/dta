package service

import (
	"context"
	b64 "encoding/base64"

	gpprotos "github.com/iantal/dta/protos/gradle-parser"
)

func (e *Explorer) parseGradle(data string) ([]*gpprotos.Project, error) {
	r := &gpprotos.ParseRequest{
		Data: encodeB64(data),
	}

	resp, err := e.gradleClient.Parse(context.Background(), r)
	if err != nil {
		e.log.Info("Error gradle parser", "error", err)
		return nil, err
	}

	return resp.GetProjects(), nil
}

func encodeB64(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}
