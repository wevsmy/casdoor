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

package sync

import (
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/xorm-io/xorm"
)

func getUpdateSql(schemaName string, tableName string, columnNames []string, newColumnVal []interface{}, pkColumnNames []string, pkColumnValue []interface{}) (string, []interface{}, error) {
	updateSql := squirrel.Update(schemaName + "." + tableName)
	for i, columnName := range columnNames {
		updateSql = updateSql.Set(columnName, newColumnVal[i])
	}

	for i, pkColumnName := range pkColumnNames {
		updateSql = updateSql.Where(squirrel.Eq{pkColumnName: pkColumnValue[i]})
	}

	sql, args, err := updateSql.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

func getInsertSql(schemaName string, tableName string, columnNames []string, columnValue []interface{}) (string, []interface{}, error) {
	insertSql := squirrel.Insert(schemaName + "." + tableName).Columns(columnNames...).Values(columnValue...)

	return insertSql.ToSql()
}

func getDeleteSql(schemaName string, tableName string, pkColumnNames []string, pkColumnValue []interface{}) (string, []interface{}, error) {
	deleteSql := squirrel.Delete(schemaName + "." + tableName)

	for i, columnName := range pkColumnNames {
		deleteSql = deleteSql.Where(squirrel.Eq{columnName: pkColumnValue[i]})
	}

	return deleteSql.ToSql()
}

func createEngine(dataSourceName string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	// ping mysql
	err = engine.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("mysql connection success……")
	return engine, nil
}

func getServerId(engin *xorm.Engine) (uint32, error) {
	record, err := engin.QueryInterface("SELECT @@server_id")
	if err != nil {
		return 0, err
	}

	res := uint32(record[0]["@@server_id"].(int64))
	return res, nil
}

func getServerUuid(engin *xorm.Engine) (string, error) {
	record, err := engin.QueryString("show variables like 'server_uuid'")
	if err != nil {
		return "", err
	}

	res := record[0]["Value"]
	return res, err
}

func getPkColumnNames(columnNames []string, PKColumns []int) []string {
	pkColumnNames := make([]string, len(PKColumns))
	for i, index := range PKColumns {
		pkColumnNames[i] = columnNames[index]
	}
	return pkColumnNames
}

func getPkColumnValues(columnValues []interface{}, PKColumns []int) []interface{} {
	pkColumnNames := make([]interface{}, len(PKColumns))
	for i, index := range PKColumns {
		pkColumnNames[i] = columnValues[index]
	}
	return pkColumnNames
}
