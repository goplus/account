/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"fmt"
	"os"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var (
	ErrNotExist   = os.ErrNotExist
	ErrPermission = os.ErrPermission
)

type CasdoorConfig struct {
	endPoint         string
	clientId         string
	clientSecret     string
	certificate      string
	organizationName string
	applicationName  string
}

type Account struct {
	casdoorConfig *CasdoorConfig

	zlog *zap.SugaredLogger
}

func New() *Account {
	conf := casdoorConfigInit()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	zlog := logger.Sugar()

	return &Account{casdoorConfig: conf, zlog: zlog}
}

func casdoorConfigInit() *CasdoorConfig {
	endPoint := os.Getenv("GOP_CASDOOR_ENDPOINT")
	clientID := os.Getenv("GOP_CASDOOR_CLIENTID")
	clientSecret := os.Getenv("GOP_CASDOOR_CLIENTSECRET")
	certificate := os.Getenv("GOP_CASDOOR_CERTIFICATE")
	organizationName := os.Getenv("GOP_CASDOOR_ORGANIZATIONNAME")
	applicationName := os.Getenv("GOP_CASDOOR_APPLICATONNAME")

	casdoorsdk.InitConfig(endPoint, clientID, clientSecret, certificate, organizationName, applicationName)

	return &CasdoorConfig{
		endPoint:         endPoint,
		clientId:         clientID,
		clientSecret:     clientSecret,
		certificate:      certificate,
		organizationName: organizationName,
		applicationName:  applicationName,
	}
}

func (a *Account) RedirectToCasdoor(redirect string) (loginURL string) {
	// TODO: Check whitelist from referer
	ResponseType := "code"
	Scope := "read"
	State := "casdoor"
	loginURL = fmt.Sprintf(
		"%s/login/oauth/authorize?client_id=%s&response_type=%s&redirect_uri=%s&scope=%s&state=%s",
		a.casdoorConfig.endPoint,
		a.casdoorConfig.clientId,
		ResponseType,
		redirect,
		Scope,
		State,
	)

	return loginURL
}

func (a *Account) GetAccessToken(code, state string) (token *oauth2.Token, err error) {
	token, err = casdoorsdk.GetOAuthToken(code, state)
	if err != nil {
		a.zlog.Error(err)

		return nil, ErrNotExist
	}

	return token, nil
}
