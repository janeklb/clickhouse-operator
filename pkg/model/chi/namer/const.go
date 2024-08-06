// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package namer

import "github.com/altinity/clickhouse-operator/pkg/model/common/namer/macro"

const (
	// configMapNamePatternCommon is a template of common settings for the CHI ConfigMap. "chi-{chi}-common-configd"
	configMapNamePatternCommon = "chi-" + macro.MacrosChiName + "-common-configd"

	// configMapNamePatternCommonUsers is a template of common users settings for the CHI ConfigMap. "chi-{chi}-common-usersd"
	configMapNamePatternCommonUsers = "chi-" + macro.MacrosChiName + "-common-usersd"

	// configMapNamePatternHost is a template of macros ConfigMap. "chi-{chi}-deploy-confd-{cluster}-{shard}-{host}"
	configMapNamePatternHost = "chi-" + macro.MacrosChiName + "-deploy-confd-" + macro.MacrosClusterName + "-" + macro.MacrosHostName
)