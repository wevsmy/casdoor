// Copyright 2021 The Casdoor Authors. All Rights Reserved.
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

import React from "react";
import {Spin} from "antd";
import {withRouter} from "react-router-dom";
import * as AuthBackend from "./AuthBackend";
import * as Util from "./Util";
import * as Setting from "../Setting";
import i18next from "i18next";
import {authConfig} from "./Auth";
import {renderLoginPanel} from "../Setting";

class SamlCallback extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      msg: null,
    };
  }

  getResponseType(redirectUri) {
    const authServerUrl = authConfig.serverUrl;
    // Casdoor's own login page, so "code" is not necessary
    if (redirectUri === "null") {
      return "login";
    }
    const realRedirectUrl = new URL(redirectUri).origin;
    // For Casdoor itself, we use "login" directly
    if (authServerUrl === realRedirectUrl) {
      return "login";
    } else {
      return "code";
    }
  }

  UNSAFE_componentWillMount() {
    const params = new URLSearchParams(this.props.location.search);
    const relayState = params.get("relayState");
    const samlResponse = params.get("samlResponse");

    const messages = atob(relayState).split("&");
    const clientId = messages[0] === "" ? "" : messages[0];
    const application = messages[0] === "" ? "app-built-in" : "";
    const state = messages[1];
    const providerName = messages[2];
    const redirectUri = messages[3];
    const responseType = this.getResponseType(redirectUri);

    const body = {
      type: responseType,
      clientId: clientId,
      provider: providerName,
      state: state,
      application: application,
      redirectUri: `${window.location.origin}/callback`,
      method: "signup",
      relayState: relayState,
      samlResponse: encodeURIComponent(samlResponse),
    };

    let param;
    if (clientId === null || clientId === "") {
      param = "";
    } else {
      param = `?clientId=${clientId}&responseType=${responseType}&redirectUri=${redirectUri}&scope=read&state=${state}`;
    }

    AuthBackend.loginWithSaml(body, param)
      .then((res) => {
        if (res.status === "ok") {
          const responseType = this.getResponseType(redirectUri);
          const handleLogin = (res2) => {
            if (responseType === "login") {
              Setting.showMessage("success", "Logged in successfully");
              Setting.goToLink("/");
            } else if (responseType === "code") {
              const code = res2.data;
              Setting.goToLink(`${redirectUri}?code=${code}&state=${state}`);
            }
          };
          Setting.checkLoginMfa(res, body, {
            clientId: clientId,
            responseType: responseType,
            redirectUri: messages[3],
            state: state,
            nonce: "",
            scope: "read",
            challengeMethod: "",
            codeChallenge: "",
            type: "code",
          }, handleLogin, this);
        } else {
          this.setState({
            msg: res.msg,
          });
        }
      });
  }

  render() {
    if (this.state.getVerifyTotp !== undefined) {
      const application = Setting.getApplicationObj(this);
      return renderLoginPanel(application, this.state.getVerifyTotp, this, window.location.origin);
    }

    return (
      <div style={{display: "flex", justifyContent: "center", alignItems: "center"}}>
        {
          (this.state.msg === null) ? (
            <Spin size="large" tip={i18next.t("login:Signing in...")} style={{paddingTop: "10%"}} />
          ) : (
            Util.renderMessageLarge(this, this.state.msg)
          )
        }
      </div>
    );
  }
}
export default withRouter(SamlCallback);
