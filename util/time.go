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

package util

import (
	"strconv"
	"time"
)

func GetCurrentTime() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format(time.RFC3339)
}

func GetCurrentTimeEx(timestamp string) string {
	tm := time.Now()
	inputTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		panic(err)
	}

	if !tm.After(inputTime) {
		tm = inputTime.Add(1 * time.Millisecond)
	}

	return tm.Format("2006-01-02T15:04:05.999Z07:00")
}

func GetCurrentUnixTime() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func String2Time(timestamp string) time.Time {
	if timestamp == "" {
		return time.Now()
	}
	parseTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		panic(err)
	}
	return parseTime
}

func Time2String(timestamp time.Time) string {
	return timestamp.Format(time.RFC3339)
}

func IsTokenExpired(createdTime string, expiresIn int) (bool, string) {
	createdTimeObj, _ := time.Parse(time.RFC3339, createdTime)
	expiresAtObj := createdTimeObj.Add(time.Duration(expiresIn) * time.Second)
	isExpired := time.Now().After(expiresAtObj)
	expireTime := expiresAtObj.Local().Format(time.RFC3339)
	return isExpired, expireTime
}
