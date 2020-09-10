/*
 * Copyright 2020 The Knative Authors
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

package testing

import (
	"encoding/json"

	"github.com/gogo/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgotesting "k8s.io/client-go/testing"
	reconcilertesting "knative.dev/eventing/pkg/reconciler/testing/v1"

	coreconfig "knative.dev/eventing-kafka-broker/control-plane/pkg/core/config"
	"knative.dev/eventing-kafka-broker/control-plane/pkg/reconciler/base"
	. "knative.dev/eventing-kafka-broker/control-plane/pkg/reconciler/broker"
)

const (
	ConfigMapNamespace = "test-namespace-config-map"
	ConfigMapName      = "test-config-cm"

	serviceNamespace = "test-service-namespace"
	serviceName      = "test-service"
	ServiceURL       = "http://test-service.test-service-namespace.svc.cluster.local/"

	TriggerUUID = "e7185016-5d98-4b54-84e8-3b1cd4acc6b5"
)

var (
	Formats = []string{base.Protobuf, base.Json}
)

func NewService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: serviceNamespace,
		},
	}
}

func NewConfigMap(configs *Configs, data []byte) runtime.Object {
	return reconcilertesting.NewConfigMap(
		configs.DataPlaneConfigMapName,
		configs.DataPlaneConfigMapNamespace,
		func(configMap *corev1.ConfigMap) {
			if configMap.BinaryData == nil {
				configMap.BinaryData = make(map[string][]byte, 1)
			}
			if data == nil {
				data = []byte("")
			}
			configMap.BinaryData[base.ConfigMapDataKey] = data
		},
	)
}

func NewConfigMapFromBrokers(brokers *coreconfig.Brokers, configs *Configs) runtime.Object {
	var data []byte
	var err error
	if configs.DataPlaneConfigFormat == base.Protobuf {
		data, err = proto.Marshal(brokers)
	} else {
		data, err = json.Marshal(brokers)
	}
	if err != nil {
		panic(err)
	}

	return NewConfigMap(configs, data)
}

func ConfigMapUpdate(configs *Configs, brokers *coreconfig.Brokers) clientgotesting.UpdateActionImpl {
	return clientgotesting.NewUpdateAction(
		schema.GroupVersionResource{
			Group:    "*",
			Version:  "v1",
			Resource: "ConfigMap",
		},
		configs.DataPlaneConfigMapNamespace,
		NewConfigMapFromBrokers(brokers, configs),
	)
}
