package utilities

import (
	"errors"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.ipify.org").
		Get("/").
		MatchParams(map[string]string{"format": "json"}).
		Reply(200).
		JSON(map[string]string{"ip": "103.186.202.7"})

	ip, err := GetPublicIP()
	assert.NoError(t, err)
	assert.Equal(t, "103.186.202.7", ip)
}

func TestGetPublicIP_Error(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.ipify.org").
		Get("/").
		MatchParams(map[string]string{"format": "json"}).
		Reply(500)

	ip, err := GetPublicIP()
	assert.Error(t, err)
	assert.Empty(t, ip)
}

func TestGetPublicIP_InvalidResponse(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.ipify.org").
		Get("/").
		MatchParams(map[string]string{"format": "json"}).
		Reply(200).
		JSON(map[string]int{"ip": 123})

	ip, err := GetPublicIP()
	assert.Error(t, err)
	assert.Empty(t, ip)
}

func TestGetPublicIP_ErrorClientDo(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.ipify.org").
		Get("/").
		MatchParams(map[string]string{"format": "json"}).
		ReplyError(errors.New("error"))

	ip, err := GetPublicIP()
	assert.Error(t, err)
	assert.Empty(t, ip)
}

func TestGetPublicIP_NewRequestError(t *testing.T) {
	gock.New("https://api.ipify.org").
		Get("/").
		ReplyError(errors.New("error"))

	ip, err := GetPublicIP()
	assert.Error(t, err)
	assert.Empty(t, ip)
}
