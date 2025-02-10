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

package pp

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func getPriceString(price float64) string {
	priceString := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", price), "0"), ".")
	return priceString
}

func joinAttachString(tokens []string) string {
	return strings.Join(tokens, "|")
}

func parseAttachString(s string) (string, string, string, error) {
	tokens := strings.Split(s, "|")
	if len(tokens) != 3 {
		return "", "", "", fmt.Errorf("parseAttachString() error: len(tokens) expected 3, got: %d", len(tokens))
	}
	return tokens[0], tokens[1], tokens[2], nil
}

func priceInt64ToFloat64(price int64) float64 {
	return float64(price) / 100
}

func priceFloat64ToInt64(price float64) int64 {
	return int64(math.Round(price * 100))
}

func priceFloat64ToString(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}

func priceStringToFloat64(price string) float64 {
	f, err := strconv.ParseFloat(price, 64)
	if err != nil {
		panic(err)
	}
	return f
}
