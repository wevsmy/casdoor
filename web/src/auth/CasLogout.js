// Copyright 2022 The Casdoor Authors. All Rights Reserved.
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
import {Card, Spin} from "antd";
import {withRouter} from "react-router-dom";
import * as AuthBackend from "./AuthBackend";
import * as Setting from "../Setting";
import i18next from "i18next";

class CasLogout extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      msg: null,
    };
    if (props.match?.params.casApplicationName !== undefined) {
      this.state.owner = props.match?.params.owner;
      this.state.applicationName = props.match?.params.casApplicationName;
    }
  }

  UNSAFE_componentWillMount() {
    const params = new URLSearchParams(this.props.location.search);
    const logoutInterval = 100;

    const logoutTimeOut = (redirectUri) => {
      setTimeout(() => {
        AuthBackend.getAccount().then((accountRes) => {
          if (accountRes.status === "ok") {
            AuthBackend.logout().then((logoutRes) => {
              if (logoutRes.status === "ok") {
                logoutTimeOut(logoutRes.data2);
              } else {
                Setting.showMessage("error", `${i18next.t("login:Failed to log out")}: ${logoutRes.msg}`);
              }
            });
          } else {
            Setting.showMessage("success", i18next.t("application:Logged out successfully"));
            this.props.onUpdateAccount(null);
            if (redirectUri !== null && redirectUri !== undefined && redirectUri !== "") {
              Setting.goToLink(redirectUri);
            } else if (params.has("service")) {
              Setting.goToLink(params.get("service"));
            } else {
              Setting.goToLinkSoft(this, `/cas/${this.state.owner}/${this.state.applicationName}/login`);
            }
          }
        });
      }, logoutInterval);
    };

    AuthBackend.logout()
      .then((res) => {
        if (res.status === "ok") {
          logoutTimeOut(res.data2);
        } else {
          Setting.showMessage("error", `${i18next.t("login:Failed to log out")}: ${res.msg}`);
        }
      });
  }

  render() {
    return (
      <Card>
        <div style={{display: "flex", justifyContent: "center", alignItems: "center"}}>
          {
            <Spin size="large" tip={i18next.t("login:Logging out...")} style={{paddingTop: "10%"}} />
          }
        </div>
      </Card>
    );
  }
}
export default withRouter(CasLogout);
