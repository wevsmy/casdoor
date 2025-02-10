// Copyright 2023 The Casdoor Authors. All Rights Reserved.
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
import {Link} from "react-router-dom";
import {Button, Table} from "antd";
import {ClockCircleOutlined, CloseCircleOutlined, ExclamationCircleOutlined, MinusCircleOutlined, SyncOutlined} from "@ant-design/icons";
import moment from "moment";
import * as Setting from "./Setting";
import * as SubscriptionBackend from "./backend/SubscriptionBackend";
import i18next from "i18next";
import BaseListPage from "./BaseListPage";
import PopconfirmModal from "./common/modal/PopconfirmModal";

class SubscriptionListPage extends BaseListPage {
  newSubscription() {
    const randomName = Setting.getRandomName();
    const owner = Setting.getRequestOrganization(this.props.account);

    return {
      owner: owner,
      name: `sub_${randomName}`,
      createdTime: moment().format(),
      displayName: `New Subscription - ${randomName}`,
      startTime: moment().format(),
      endTime: moment().add(30, "d").format(),
      period: "Monthly",
      description: "",
      user: "",
      plan: "",
      state: "Active",
    };
  }

  addSubscription() {
    const newSubscription = this.newSubscription();
    SubscriptionBackend.addSubscription(newSubscription)
      .then((res) => {
        if (res.status === "ok") {
          this.props.history.push({pathname: `/subscriptions/${newSubscription.owner}/${newSubscription.name}`, mode: "add"});
          Setting.showMessage("success", i18next.t("general:Successfully added"));
        } else {
          Setting.showMessage("error", `${i18next.t("general:Failed to add")}: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `${i18next.t("general:Failed to connect to server")}: ${error}`);
      });
  }

  deleteSubscription(i) {
    SubscriptionBackend.deleteSubscription(this.state.data[i])
      .then((res) => {
        if (res.status === "ok") {
          Setting.showMessage("success", i18next.t("general:Successfully deleted"));
          this.fetch({
            pagination: {
              ...this.state.pagination,
              current: this.state.pagination.current > 1 && this.state.data.length === 1 ? this.state.pagination.current - 1 : this.state.pagination.current,
            },
          });
        } else {
          Setting.showMessage("error", `${i18next.t("general:Failed to delete")}: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `${i18next.t("general:Failed to connect to server")}: ${error}`);
      });
  }

  renderTable(subscriptions) {
    const columns = [
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "140px",
        fixed: "left",
        sorter: true,
        ...this.getColumnSearchProps("name"),
        render: (text, record, index) => {
          return (
            <Link to={`/subscriptions/${record.owner}/${record.name}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:Organization"),
        dataIndex: "owner",
        key: "owner",
        width: "120px",
        sorter: true,
        ...this.getColumnSearchProps("owner"),
        render: (text, record, index) => {
          return (
            <Link to={`/organizations/${text}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:Created time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "160px",
        sorter: true,
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Display name"),
        dataIndex: "displayName",
        key: "displayName",
        width: "170px",
        sorter: true,
        ...this.getColumnSearchProps("displayName"),
      },
      {
        title: i18next.t("subscription:Period"),
        dataIndex: "period",
        key: "period",
        width: "140px",
        ...this.getColumnSearchProps("period"),
      },
      {
        title: i18next.t("subscription:Start time"),
        dataIndex: "startTime",
        key: "startTime",
        width: "140px",
        ...this.getColumnSearchProps("startTime"),
      },
      {
        title: i18next.t("subscription:End time"),
        dataIndex: "endTime",
        key: "endTime",
        width: "140px",
        ...this.getColumnSearchProps("endTime"),
      },
      {
        title: i18next.t("general:Plan"),
        dataIndex: "plan",
        key: "plan",
        width: "140px",
        ...this.getColumnSearchProps("plan"),
        render: (text, record, index) => {
          return (
            <Link to={`/plans/${record.owner}/${text}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:User"),
        dataIndex: "user",
        key: "user",
        width: "140px",
        ...this.getColumnSearchProps("user"),
        render: (text, record, index) => {
          return (
            <Link to={`/users/${record.owner}/${text}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:Payment"),
        dataIndex: "payment",
        key: "payment",
        width: "140px",
        ...this.getColumnSearchProps("payment"),
        render: (text, record, index) => {
          return (
            <Link to={`/payments/${record.owner}/${text}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:State"),
        dataIndex: "state",
        key: "state",
        width: "120px",
        sorter: true,
        ...this.getColumnSearchProps("state"),
        render: (text, record, index) => {
          switch (text) {
          case "Pending":
            return Setting.getTag("processing", i18next.t("subscription:Pending"), <ExclamationCircleOutlined />);
          case "Active":
            return Setting.getTag("success", i18next.t("subscription:Active"), <SyncOutlined spin />);
          case "Upcoming":
            return Setting.getTag("warning", i18next.t("subscription:Upcoming"), <ClockCircleOutlined />);
          case "Expired":
            return Setting.getTag("warning", i18next.t("subscription:Expired"), <ClockCircleOutlined />);
          case "Error":
            return Setting.getTag("error", i18next.t("subscription:Error"), <CloseCircleOutlined />);
          case "Suspended":
            return Setting.getTag("default", i18next.t("subscription:Suspended"), <MinusCircleOutlined />);
          default:
            return null;
          }
        },
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "",
        key: "op",
        width: "230px",
        fixed: (Setting.isMobile()) ? "false" : "right",
        render: (text, record, index) => {
          return (
            <div>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/subscriptions/${record.owner}/${record.name}`)}>{i18next.t("general:Edit")}</Button>
              <PopconfirmModal
                title={i18next.t("general:Sure to delete") + `: ${record.name} ?`}
                onConfirm={() => this.deleteSubscription(index)}
              >
              </PopconfirmModal>
            </div>
          );
        },
      },
    ];

    const paginationProps = {
      total: this.state.pagination.total,
      showQuickJumper: true,
      showSizeChanger: true,
      showTotal: () => i18next.t("general:{total} in total").replace("{total}", this.state.pagination.total),
    };

    return (
      <div>
        <Table scroll={{x: "max-content"}} columns={columns} dataSource={subscriptions} rowKey={(record) => `${record.owner}/${record.name}`} size="middle" bordered pagination={paginationProps}
          title={() => (
            <div>
              {i18next.t("general:Subscriptions")}&nbsp;&nbsp;&nbsp;&nbsp;
              <Button type="primary" size="small" onClick={this.addSubscription.bind(this)}>{i18next.t("general:Add")}</Button>
            </div>
          )}
          loading={this.state.loading}
          onChange={this.handleTableChange}
        />
      </div>
    );
  }

  fetch = (params = {}) => {
    let field = params.searchedColumn, value = params.searchText;
    const sortField = params.sortField, sortOrder = params.sortOrder;
    if (params.type !== undefined && params.type !== null) {
      field = "type";
      value = params.type;
    }
    this.setState({loading: true});
    SubscriptionBackend.getSubscriptions(Setting.isDefaultOrganizationSelected(this.props.account) ? "" : Setting.getRequestOrganization(this.props.account), params.pagination.current, params.pagination.pageSize, field, value, sortField, sortOrder)
      .then((res) => {
        this.setState({
          loading: false,
        });
        if (res.status === "ok") {
          this.setState({
            data: res.data,
            pagination: {
              ...params.pagination,
              total: res.data2,
            },
            searchText: params.searchText,
            searchedColumn: params.searchedColumn,
          });
        } else {
          if (Setting.isResponseDenied(res)) {
            this.setState({
              isAuthorized: false,
            });
          } else {
            Setting.showMessage("error", res.msg);
          }
        }
      });
  };
}

export default SubscriptionListPage;
