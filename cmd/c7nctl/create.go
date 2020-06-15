// Copyright © 2018 VinkDong <dong@wenqi.us>
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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var clusterName string
var clusterCode string
var clusterDescription string
var appName string
var appType string
var envCode string
var envName string
var envDescription string
var file string
var devopsEnvGroupId int
var templateAppServiceId int
var templateAppServiceVersionId int
var appServiceId int
var appServiceVersionId int
var instanceName string
var valueFile string
var configMapDescription string
var secretDescription string

// getCmd represents the get command
func newCreateCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "The command to create choerodon resource",
		Long:  `The command to create choerodon resource.such as organization, project, app, instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
			error := c7nclient.Client.CheckIsLogin()
			if error != nil {
				fmt.Println(error)
				return
			}
			if len(args) > 0 {
				fmt.Printf("don't have the resource %s, you can user c7nctl create --help to see the resource you can use!", args[0])
			} else {
				cmd.Help()
			}
		},
	}

	cmd.AddCommand(newCreateClusterCmd(out))
	cmd.AddCommand(newCreateAppCmd(out))
	cmd.AddCommand(newCreateEnvCmd(out))
	cmd.AddCommand(newCreateInstanceCmd(out))
	cmd.AddCommand(newCreateServiceCmd(out))
	cmd.AddCommand(newCreateInstanceCmd(out))
	cmd.AddCommand(newCreateIngressCmd(out))
	cmd.AddCommand(newCreateCertCmd(out))
	cmd.AddCommand(newCreateConfigMapCmd(out))
	cmd.AddCommand(newCreateSecretCmd(out))
	cmd.AddCommand(newCreateCustomCmd(out))
	cmd.AddCommand(newCreatePvcCmd(out))
	cmd.AddCommand(newCreatePvCmd(out))

	return cmd
}

// create cluster command
func newCreateClusterCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "create cluster",
		Long:  `you can use this command to create cluster `,
		Run: func(cmd *cobra.Command, args []string) {
			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
			error := c7nclient.Client.CheckIsLogin()
			if error != nil {
				fmt.Println(error)
				return
			}
			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}
			clusterPostInfo := model.ClusterPostInfo{clusterName, clusterCode, clusterDescription, true}
			c7nclient.Client.CreateCluster(cmd.OutOrStdout(), pro.ID, &clusterPostInfo)
		},
	}

	cmd.Flags().StringVar(&clusterName, "name", "", "cluster name")
	cmd.Flags().StringVar(&clusterCode, "code", "", "cluster code")
	cmd.Flags().StringVar(&clusterDescription, "description", "", "cluster description")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("code")
	cmd.MarkFlagRequired("description")

	return cmd
}

// create app command
func newCreateAppCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "create app",
		Long:  `you can use this command to create app `,
		Run: func(cmd *cobra.Command, args []string) {
			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
			error := c7nclient.Client.CheckIsLogin()
			if error != nil {
				fmt.Println(error)
				return
			}
			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}
			appPostInfo := model.AppPostInfo{appName, appCode, appType, templateAppServiceId, templateAppServiceVersionId}
			c7nclient.Client.CreateApp(cmd.OutOrStdout(), pro.ID, &appPostInfo)
		},
	}

	cmd.Flags().StringVar(&appName, "name", "", "app name")
	cmd.Flags().StringVar(&appCode, "code", "", "app code")
	cmd.Flags().StringVar(&appType, "type", "", "the value can be normal or test")
	cmd.Flags().IntVar(&templateAppServiceId, "templateAppServiceId", 0, "the templateAppServiceId")
	cmd.Flags().IntVar(&templateAppServiceVersionId, "templateAppServiceVersionId", 0, "the templateAppServiceVersionId")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("code")
	cmd.MarkFlagRequired("type")

	return cmd
}

// create Env command
func newCreateEnvCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "create env",
		Long:  `you can use this command to create env `,
		Run: func(cmd *cobra.Command, args []string) {
			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
			error := c7nclient.Client.CheckIsLogin()
			if error != nil {
				fmt.Println(error)
				return
			}
			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}
			err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.ID, clusterCode)
			if err != nil {
				return
			}
			envPostInfo := model.EnvPostInfo{envName, envCode, envDescription, cluster.ID}
			c7nclient.Client.CreateEnv(cmd.OutOrStdout(), pro.ID, &envPostInfo)
		},
	}

	cmd.Flags().StringVar(&envName, "name", "", "env name")
	cmd.Flags().StringVar(&envCode, "code", "", "env code")
	cmd.Flags().StringVarP(&envDescription, "description", "d", "", "env Description ")
	cmd.Flags().StringVarP(&clusterCode, "cluster", "c", "", "the cluster code you want to use")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("code")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newCreateInstanceCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "create instance",
		Long:  `you can use this command to create instance `,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
			if _, err := os.Stat(valueFile); os.IsNotExist(err) {
				fmt.Println(err)
				return
			}
			value, err := ioutil.ReadFile(valueFile)
			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}

			err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
			if err != nil {
				return
			}
			instancePostInfo := model.InstancePostInfo{appServiceId, appServiceVersionId, env.ID, instanceName, "create", string(value)}
			c7nclient.Client.CreateInstance(cmd.OutOrStdout(), pro.ID, &instancePostInfo)
		},
	}

	cmd.Flags().StringVarP(&envCode, "env", "e", "", "the envCode you want to deploy")
	cmd.Flags().IntVarP(&appServiceId, "appServiceId", "a", 0, "the appService's id you want to deploy")
	cmd.Flags().IntVarP(&appServiceVersionId, "appServiceVersionId", "v", 0, "the appServiceVersion's id you want to deploy")
	cmd.Flags().StringVarP(&instanceName, "instanceName", "n", "", "the instance name you want to set")
	cmd.Flags().StringVarP(&valueFile, "valueFile", "f", "", "the deploy value's file")
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("appServiceId")
	cmd.MarkFlagRequired("appServiceVersionId")
	cmd.MarkFlagRequired("instanceName")
	cmd.MarkFlagRequired("valueFile")

	return cmd
}
func newCreateServiceCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "create service",
		Long:  `you can use this command to create service `,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)

			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}
			servicePostInfo := model.ServicePostInfo{}

			err = initService(cmd, &pro, &servicePostInfo)
			if err != nil {
				return
			}

			c7nclient.Client.CreateService(cmd.OutOrStdout(), pro.ID, &servicePostInfo)
		},
	}
	cmd.Flags().StringVarP(&envCode, "env", "e", "", "the envCode you want to deploy")
	cmd.Flags().StringVarP(&file, "file", "f", "", "the service yaml file")
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("file")

	return cmd
}

func newCreateIngressCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "create ingress",
		Long:  `you can use this command to create ingress `,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)

			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}
			ingressPostInfo := model.IngressPostInfo{}

			err = initIngress(cmd, &pro, &ingressPostInfo)
			if err != nil {
				return
			}

			c7nclient.Client.CreateIngress(cmd.OutOrStdout(), pro.ID, &ingressPostInfo)
		},
	}

	cmd.Flags().StringVarP(&envCode, "env", "e", "", "the envCode you want to deploy")
	cmd.Flags().StringVarP(&file, "file", "f", "", "the ingress yaml file")
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("file")

	return cmd
}

func newCreateCertCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cert",
		Short: "create certification",
		Long:  `you can use this command to create certification `,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)

			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}

			data := url.Values{}

			err = initCert(cmd, &pro, &data)
			if err != nil {
				return
			}

			c7nclient.Client.CreateCert(cmd.OutOrStdout(), 999, &data)
		},
	}

	cmd.Flags().StringVarP(&envCode, "env", "e", "", "the envCode you want to deploy")
	cmd.Flags().StringVarP(&file, "file", "f", "", "the cert yaml file")
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("file")
	return cmd
}

func newCreateConfigMapCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cm",
		Short: "create configMap",
		Long:  `you can use this command to create configMap`,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)

			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}
			configMapPostInfo := model.ConfigMapPostInfo{}

			err = initConfigMap(cmd, &pro, configMapDescription, &configMapPostInfo)
			if err != nil {
				return
			}
			c7nclient.Client.CreateConfigMap(cmd.OutOrStdout(), pro.ID, &configMapPostInfo)
		},
	}

	cmd.Flags().StringVarP(&envCode, "env", "e", "", "the envCode you want to deploy")
	cmd.Flags().StringVarP(&file, "file", "f", "", "the cert yaml file")
	cmd.Flags().StringVarP(&configMapDescription, "description", "d", "", "configMap description")
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("file")

	return cmd
}
func newCreateSecretCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "create secret",
		Long:  `you can use this command to create secret`,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)

			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}
			secretPostInfo := model.SecretPostInfo{}

			err = initSecret(cmd, &pro, secretDescription, &secretPostInfo)
			if err != nil {
				return
			}
			c7nclient.Client.CreateSecret(cmd.OutOrStdout(), pro.ID, &secretPostInfo)
		},
	}

	cmd.Flags().StringVarP(&envCode, "env", "e", "", "the envCode you want to deploy")
	cmd.Flags().StringVarP(&file, "file", "f", "", "the secret yaml file")
	cmd.Flags().StringVarP(&secretDescription, "description", "d", "", "secret description")
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("file")

	return cmd
}
func newCreateCustomCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom",
		Short: "create custom resource",
		Long:  `you can use this command to create custom resource`,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)

			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}

			data := url.Values{}

			err = initCustom(cmd, &pro, &data)
			if err != nil {
				return
			}

			c7nclient.Client.CreateCustom(cmd.OutOrStdout(), pro.ID, &data)
		},
	}

	cmd.Flags().StringVarP(&envCode, "env", "e", "", "the envCode you want to deploy")
	cmd.Flags().StringVarP(&file, "file", "f", "", "the custom yaml file")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("env")

	return cmd
}

func newCreatePvcCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pvc",
		Short: "create pvc",
		Long:  `you can use this command to create pvc`,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)

			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}

			pvcPostInfo := model.PvcPostInfo{}

			err = initPvc(cmd, &pro, &pvcPostInfo)
			if err != nil {
				return
			}

			c7nclient.Client.CreatePvc(cmd.OutOrStdout(), pro.ID, &pvcPostInfo)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "the pvc yaml file")
	cmd.Flags().StringVarP(&envCode, "envCode", "e", "", "the envCode you want to deploy")
	cmd.Flags().StringVarP(&clusterCode, "clusterCode", "c", "", "the clusterCode you want to deploy")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("clusterCode")
	cmd.MarkFlagRequired("envCode")

	return cmd
}

func newCreatePvCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pv",
		Short: "create pv",
		Long:  `you can use this command to create pv`,
		Run: func(cmd *cobra.Command, args []string) {

			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)

			err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if err != nil {
				return
			}
			err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if err != nil {
				return
			}
			err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
			if err != nil {
				return
			}

			pvPostInfo := model.PvPostInfo{}

			err = initPv(cmd, &pro, &pvPostInfo)
			if err != nil {
				return
			}

			c7nclient.Client.CreatePv(cmd.OutOrStdout(), pro.ID, &pvPostInfo)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "the pv yaml file")
	cmd.Flags().StringVarP(&clusterCode, "clusterCode", "c", "", "the clusterCode you want to deploy")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("clusterCode")

	return cmd
}

func initService(cmd *cobra.Command, pro *model.Project, servicePostInfo *model.ServicePostInfo) (error error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(file)
	results := strings.Split(string(b), "---")
	var services []v1.Service
	var endPoints []v1.Endpoints
	for _, result := range results {
		if result != "" {
			var data = []byte(result)
			service := v1.Service{}
			endPoint := v1.Endpoints{}
			yaml.Unmarshal(data, &service)
			if service.Kind == "Service" {
				services = append(services, service)
				continue
			}
			yaml.Unmarshal(data, &endPoint)
			if endPoint.Kind == "Endpoints" {
				endPoints = append(endPoints, endPoint)
			}
		}
	}
	if len(services) == 0 {
		return errors.New("The service is empty!")
	}
	service := services[0]
	if len(endPoints) > 0 {
		endPoint := endPoints[0]
		endPointPostInfo := make(map[string][]model.EndPointPortInfo)
		for _, subset := range endPoint.Subsets {
			var addresses string
			for index, address := range subset.Addresses {
				if index == len(subset.Addresses)-1 {
					addresses += address.IP
				} else {
					addresses += address.IP + ","
				}
			}
			var endPointPortInfos []model.EndPointPortInfo
			for _, port := range subset.Ports {
				endPointPortInfo := model.EndPointPortInfo{}
				endPointPortInfo.Port = port.Port
				endPointPortInfo.Name = port.Name
				endPointPortInfos = append(endPointPortInfos, endPointPortInfo)
			}
			endPointPostInfo[addresses] = endPointPortInfos
		}
		servicePostInfo.EndPoints = endPointPostInfo
	}
	if err != nil {
		fmt.Print(err)
		return err
	}
	annotations := service.ObjectMeta.Annotations
	appCode := annotations["choerodon.io/network-service-app"]
	if appCode != "" {
		err, app := c7nclient.Client.GetApp(appCode, pro.ID)
		if err != nil {
			return err
		}
		servicePostInfo.AppID = app.ID
	}
	instanceCode := annotations["choerodon.io/network-service-instances"]
	if instanceCode != "" {
		instances := strings.Split(instanceCode, "+")
		servicePostInfo.Instances = instances
	}
	var servicePorts []model.ServicePort
	for _, port := range service.Spec.Ports {
		servicePost := model.ServicePort{
			Port:       port.Port,
			TargetPort: port.TargetPort,
			NodePort:   port.NodePort,
		}
		servicePorts = append(servicePorts, servicePost)
	}
	servicePostInfo.Ports = servicePorts
	err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
	if err != nil {
		return err
	}
	servicePostInfo.EnvID = env.ID
	servicePostInfo.Name = service.ObjectMeta.Name
	var externalIps string
	for index, externalIp := range service.Spec.ExternalIPs {
		if index == len(service.Spec.ExternalIPs)-1 {
			externalIps += externalIp
		} else {
			externalIps += externalIp + ","
		}
	}
	if externalIps != "" {
		servicePostInfo.ExternalIP = externalIps
	}
	servicePostInfo.Type = string(service.Spec.Type)
	servicePostInfo.Selectors = service.Spec.Selector
	return nil
}

func initIngress(cmd *cobra.Command, pro *model.Project, ingressPostInfo *model.IngressPostInfo) (error error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	ingress := v1beta1.Ingress{}
	yaml.Unmarshal(b, &ingress)
	err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
	if err != nil {
		return err
	}
	var ingressPaths []model.IngressPath

	for _, httpIngressPath := range ingress.Spec.Rules[0].HTTP.Paths {
		err, service := c7nclient.Client.GetService(cmd.OutOrStdout(), pro.ID, env.ID, httpIngressPath.Backend.ServiceName)
		if err != nil {
			return errors.New(" the service in not exist!")
		}
		ingressPath := model.IngressPath{
			Path:        httpIngressPath.Path,
			ServicePort: httpIngressPath.Backend.ServicePort,
			ServiceName: httpIngressPath.Backend.ServiceName,
			ServiceID:   service.ID,
		}
		ingressPaths = append(ingressPaths, ingressPath)
	}
	if ingress.Spec.TLS != nil {
		err, cert := c7nclient.Client.GetCert(cmd.OutOrStdout(), pro.ID, env.ID, ingress.Spec.TLS[0].SecretName)
		if err != nil {
			return
		}
		ingressPostInfo.CertId = cert.ID
	}

	ingressPostInfo.Name = ingress.ObjectMeta.Name
	ingressPostInfo.EnvID = env.ID
	ingressPostInfo.Domain = ingress.Spec.Rules[0].Host
	ingressPostInfo.PathList = ingressPaths

	return nil
}

func initCert(cmd *cobra.Command, pro *model.Project, data *url.Values) (error error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	certificate := model.Certificate{}
	yaml.Unmarshal(b, &certificate)
	err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
	if err != nil {
		return err
	}
	(*data)["envId"] = []string{strconv.Itoa(env.ID)}
	(*data)["certName"] = []string{certificate.Metadata.Name}
	(*data)["certValue"] = []string{certificate.Spec.ExistCert.Cert}
	(*data)["keyValue"] = []string{certificate.Spec.ExistCert.Key}
	if len(certificate.Spec.DnsNames) != 0 {
		(*data)["domains"] = []string{certificate.Spec.CommonName + "," + strings.Join(certificate.Spec.DnsNames, ",")}
	} else {
		(*data)["domains"] = []string{certificate.Spec.CommonName}
	}
	(*data)["type"] = []string{"request"}
	if (*data)["certValue"][0] != "" {
		(*data)["type"] = []string{"upload"}
	}
	return nil
}

func initConfigMap(cmd *cobra.Command, pro *model.Project, description string, configMapPostInfo *model.ConfigMapPostInfo) (error error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	configMap := v1.ConfigMap{}
	yaml.Unmarshal(b, &configMap)
	err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
	if err != nil {
		return err
	}
	configMapPostInfo.EnvID = env.ID
	configMapPostInfo.Type = "create"
	configMapPostInfo.Name = configMap.Name
	configMapPostInfo.Description = description
	configMapPostInfo.Value = configMap.Data
	return nil
}

func initSecret(cmd *cobra.Command, pro *model.Project, secretDescription string, secretPostInfo *model.SecretPostInfo) (error error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	secret := v1.Secret{}
	yaml.Unmarshal(b, &secret)
	err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
	if err != nil {
		return err
	}
	secretPostInfo.EnvID = env.ID
	secretPostInfo.Type = "create"
	secretPostInfo.Name = secret.Name
	secretPostInfo.Description = secretDescription
	secretPostInfo.Value = secret.StringData
	return nil
}

func initCustom(cmd *cobra.Command, pro *model.Project, data *url.Values) (error error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
	if err != nil {
		return err
	}
	(*data)["envId"] = []string{strconv.Itoa(env.ID)}
	(*data)["type"] = []string{"create"}
	(*data)["content"] = []string{string(b)}
	return nil
}

func initPvc(cmd *cobra.Command, pro *model.Project, pvcPostInfo *model.PvcPostInfo) (error error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	pvc := v1.PersistentVolumeClaim{}
	yaml.Unmarshal(b, &pvc)
	err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
	if err != nil {
		return nil
	}
	err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.ID, clusterCode)
	if err != nil {
		return err
	}

	pvcPostInfo.EnvID = env.ID
	pvcPostInfo.Name = pvc.Name
	pvcPostInfo.PvName = pvc.Spec.VolumeName
	pvcPostInfo.ClusterId = cluster.ID

	quantity := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	size, err := quantity.Marshal()
	if err != nil {
		return err
	}

	str := string(size)
	r, _ := regexp.Compile("\\d*[MGT]i")
	result := r.FindString(str)
	pvcPostInfo.RequestResource = result

	if len(pvc.Spec.AccessModes) != 1 {
		return errors.New("only support one accessMode")
	}
	pvcPostInfo.AccessModes = string(pvc.Spec.AccessModes[0])
	return nil
}

func initPv(cmd *cobra.Command, pro *model.Project, pvPostInfo *model.PvPostInfo) (error error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	pv := v1.PersistentVolume{}
	_ = yaml.Unmarshal(b, &pv)
	err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.ID, clusterCode)
	if err != nil {
		return err
	}
	pvPostInfo.Name = pv.ObjectMeta.Name
	pvPostInfo.ClusterId = cluster.ID

	quantity := pv.Spec.Capacity[v1.ResourceStorage]
	size, err := quantity.Marshal()
	if err != nil {
		return err
	}
	str := string(size)
	r, _ := regexp.Compile("\\d*[MGT]i")
	result := r.FindString(str)
	pvPostInfo.RequestResource = result
	if len(pv.Spec.AccessModes) != 1 {
		return errors.New("only support one accessMode")
	}
	pvPostInfo.AccessModes = string(pv.Spec.AccessModes[0])
	err = setValueConfig(pv.Spec.PersistentVolumeSource, pvPostInfo)
	if err != nil {
		return err
	}
	pvPostInfo.SkipCheckProjectPermission = true
	pvPostInfo.ProjectIds = []int{}
	return nil
}

func setValueConfig(persistentVolumeSource v1.PersistentVolumeSource, pvPostInfo *model.PvPostInfo) error {
	if persistentVolumeSource.NFS != nil {
		valueConfigBuf, err := json.Marshal(persistentVolumeSource.NFS)
		if err != nil {
			fmt.Println(err)
			return err
		}
		pvPostInfo.ValueConfig = string(valueConfigBuf)
		pvPostInfo.Type = "NFS"
		return nil
	} else if persistentVolumeSource.HostPath != nil {
		valueConfigBuf, err := json.Marshal(persistentVolumeSource.HostPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		pvPostInfo.ValueConfig = string(valueConfigBuf)
		pvPostInfo.Type = "HostPath"
		return nil
	} else {
		fmt.Println("Only support NFS and HostPath,please check it out")
		return errors.New("type error")
	}
}
