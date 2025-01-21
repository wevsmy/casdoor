// Copyright 2025 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package routers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/beego/beego/context"
	"github.com/casdoor/casdoor/object"
)

type OrganizationThemeCookie struct {
	ThemeData   *object.ThemeData
	LogoUrl     string
	FooterHtml  string
	Favicon     string
	DisplayName string
}

func appendThemeCookie(ctx *context.Context, urlPath string) (*OrganizationThemeCookie, error) {
	if urlPath == "/login" {
		application, err := object.GetDefaultApplication(fmt.Sprintf("admin/built-in"))
		if err != nil {
			return nil, err
		}
		organization := application.OrganizationObj
		if organization == nil {
			organization, err = object.GetOrganization(fmt.Sprintf("admin/built-in"))
			if err != nil {
				return nil, err
			}
		}
		if organization != nil {
			organizationThemeCookie := &OrganizationThemeCookie{
				application.ThemeData,
				application.Logo,
				application.FooterHtml,
				organization.Favicon,
				organization.DisplayName,
			}

			if application.ThemeData != nil {
				organizationThemeCookie.ThemeData = organization.ThemeData
			}
			return organizationThemeCookie, setThemeDataCookie(ctx, organizationThemeCookie)
		}
	} else if strings.HasPrefix(urlPath, "/login/oauth/authorize") {
		clientId := ctx.Input.Query("client_id")
		if clientId == "" {
			return nil, nil
		}
		application, err := object.GetApplicationByClientId(clientId)
		if err != nil {
			return nil, err
		}
		if application != nil {
			organization := application.OrganizationObj
			if organization == nil {
				organization, err = object.GetOrganization(fmt.Sprintf("admin/%s", application.Owner))
				if err != nil {
					return nil, err
				}
			}
			organizationThemeCookie := &OrganizationThemeCookie{
				application.ThemeData,
				application.Logo,
				application.FooterHtml,
				organization.Favicon,
				organization.DisplayName,
			}

			if application.ThemeData != nil {
				organizationThemeCookie.ThemeData = organization.ThemeData
			}
			return organizationThemeCookie, setThemeDataCookie(ctx, organizationThemeCookie)
		}
	} else if strings.HasPrefix(urlPath, "/login/") {
		owner := strings.Replace(urlPath, "/login/", "", -1)
		if owner != "undefined" && owner != "oauth/undefined" {
			application, err := object.GetDefaultApplication(fmt.Sprintf("admin/%s", owner))
			if err != nil {
				return nil, err
			}
			organization := application.OrganizationObj
			if organization == nil {
				organization, err = object.GetOrganization(fmt.Sprintf("admin/%s", owner))
				if err != nil {
					return nil, err
				}
			}
			if organization != nil {
				organizationThemeCookie := &OrganizationThemeCookie{
					application.ThemeData,
					application.Logo,
					application.FooterHtml,
					organization.Favicon,
					organization.DisplayName,
				}

				if application.ThemeData != nil {
					organizationThemeCookie.ThemeData = organization.ThemeData
				}
				return organizationThemeCookie, setThemeDataCookie(ctx, organizationThemeCookie)
			}
		}
	}

	return nil, nil
}

func setThemeDataCookie(ctx *context.Context, organizationThemeCookie *OrganizationThemeCookie) error {
	themeDataString, err := json.Marshal(organizationThemeCookie.ThemeData)
	if err != nil {
		return err
	}
	ctx.SetCookie("organizationTheme", string(themeDataString))
	ctx.SetCookie("organizationLogo", organizationThemeCookie.LogoUrl)
	ctx.SetCookie("organizationFootHtml", organizationThemeCookie.FooterHtml)
	return nil
}
