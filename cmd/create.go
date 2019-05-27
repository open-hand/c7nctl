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

package cmd

import (
	"errors"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"os"
	"strings"
)

var appTemplateName string
var appTemplateCode string
var appTemplateDescription string
var clusterName string
var clusterCode string
var clusterDescription string
var copyFrom string
var appName string
var appType string
var appTemplate string
var envCode string
var envName string
var envDescription string
var content string

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createAppTemplateCmd)
	createCmd.AddCommand(createClusterCmd)
	createCmd.AddCommand(createAppCmd)
	createCmd.AddCommand(createEnvCmd)
	createCmd.AddCommand(createInstanceCmd)
	createCmd.AddCommand(createServiceCmd)
	createCmd.AddCommand(createIngressCmd)
	createCmd.AddCommand(createCertCmd)
	createCmd.AddCommand(createConfigMapCmd)
	createCmd.AddCommand(createSecretCmd)

	createAppTemplateCmd.Flags().StringVar(&appTemplateName, "name", "", "appTemplate name")
	createAppTemplateCmd.Flags().StringVar(&appTemplateCode, "code", "", "appTemplate code")
	createAppTemplateCmd.Flags().StringVar(&appTemplateDescription, "description", "", "appTemplate description")
	createAppTemplateCmd.Flags().StringVar(&copyFrom, "copyFrom", "", "appTemplate copy from")
	createClusterCmd.Flags().StringVar(&clusterName, "name", "", "cluster name")
	createClusterCmd.Flags().StringVar(&clusterCode, "code", "", "cluster code")
	createClusterCmd.Flags().StringVar(&clusterDescription, "description", "", "cluster description")
	createAppCmd.Flags().StringVar(&appName, "name", "", "app name")
	createAppCmd.Flags().StringVar(&appCode, "code", "", "app code")
	createAppCmd.Flags().StringVar(&appType, "type", "", "the value can be normal or test")
	createAppCmd.Flags().StringVar(&appTemplate, "appTemplate", "", "the appTemplate code you want to use")
	createEnvCmd.Flags().StringVar(&envName, "name", "", "env name")
	createEnvCmd.Flags().StringVar(&envCode, "code", "", "env code")
	createEnvCmd.Flags().StringVar(&envDescription, "description", "", "env Description ")
	createEnvCmd.Flags().StringVar(&clusterCode, "cluster", "", "the cluster code you want to use")
	createInstanceCmd.Flags().StringVar(&envCode, "env", "", "the envCode you want to deploy")
	createInstanceCmd.Flags().StringVar(&content, "content", "", "the instance  yaml file")
	createServiceCmd.Flags().StringVar(&envCode, "env", "", "the envCode you want to deploy")
	createServiceCmd.Flags().StringVar(&content, "content", "", "the service yaml file")
	createIngressCmd.Flags().StringVar(&envCode, "env", "", "the envCode you want to deploy")
	createIngressCmd.Flags().StringVar(&content, "content", "", "the ingress yaml file")
	createCertCmd.Flags().StringVar(&envCode, "env", "", "the envCode you want to deploy")
	createCertCmd.Flags().StringVar(&content, "content", "", "the cert yaml file")
	createConfigMapCmd.Flags().StringVar(&envCode, "env", "", "the envCode you want to deploy")
	createConfigMapCmd.Flags().StringVar(&content, "content", "", "the configMap yaml file")
	createSecretCmd.Flags().StringVar(&envCode, "env", "", "the envCode you want to deploy")
	createSecretCmd.Flags().StringVar(&content, "content", "", "the secret yaml file")
	createAppTemplateCmd.MarkFlagRequired("name")
	createAppTemplateCmd.MarkFlagRequired("code")
	createAppTemplateCmd.MarkFlagRequired("description")
	createClusterCmd.MarkFlagRequired("name")
	createClusterCmd.MarkFlagRequired("code")
	createClusterCmd.MarkFlagRequired("description")
	createAppCmd.MarkFlagRequired("name")
	createAppCmd.MarkFlagRequired("code")
	createAppCmd.MarkFlagRequired("type")
	createAppCmd.MarkFlagRequired("appTemplate")
	createEnvCmd.MarkFlagRequired("cluster")
	createEnvCmd.MarkFlagRequired("description")
	createEnvCmd.MarkFlagRequired("code")
	createEnvCmd.MarkFlagRequired("name")
	createInstanceCmd.MarkFlagRequired("env")
	createInstanceCmd.MarkFlagRequired("content")
	createServiceCmd.MarkFlagRequired("env")
	createServiceCmd.MarkFlagRequired("content")
	createIngressCmd.MarkFlagRequired("env")
	createIngressCmd.MarkFlagRequired("content")
	createCertCmd.MarkFlagRequired("env")
	createCertCmd.MarkFlagRequired("content")
	createConfigMapCmd.MarkFlagRequired("env")
	createConfigMapCmd.MarkFlagRequired("content")
	createSecretCmd.MarkFlagRequired("env")
	createSecretCmd.MarkFlagRequired("content")
}

// getCmd represents the get command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "The command to create choerodon resource",
	Long:  `The command to create choerodon resource.such as organization, project, app, instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
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

// create apptemplate command
var createAppTemplateCmd = &cobra.Command{
	Use:   "app-template",
	Short: "create application template",
	Long:  `you can use this command to create application template `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, organizationId := c7nclient.Client.GetOrganization(cmd.OutOrStdout(), userinfo.ID, orgCode)
		if err != nil {
			return
		}
		var appTemplateId int
		if copyFrom != "" {
			err, appTemplateInfo := c7nclient.Client.GetAppTemplate(cmd.OutOrStdout(), organizationId, copyFrom)
			if err != nil {
				return
			}
			appTemplateId = appTemplateInfo.ID
		}
		appTemplatePostInfo := model.AppTemplatePostInfo{appTemplateName, appTemplateCode, appTemplateDescription, appTemplateId}

		c7nclient.Client.CreateAppTemplate(cmd.OutOrStdout(), organizationId, &appTemplatePostInfo)
	},
}

// create cluster command
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "create cluster",
	Long:  `you can use this command to create cluster `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, organizationId := c7nclient.Client.GetOrganization(cmd.OutOrStdout(), userinfo.ID, orgCode)
		if err != nil {
			return
		}
		clusterPostInfo := model.ClusterPostInfo{clusterName, clusterCode, clusterDescription, true}
		c7nclient.Client.CreateCluster(cmd.OutOrStdout(), organizationId, &clusterPostInfo)
	},
}

// create app command
var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "create app",
	Long:  `you can use this command to create app `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userinfo.ID, proCode)
		if err != nil {
			return
		}
		err, apptemplate := c7nclient.Client.GetAppTemplate(cmd.OutOrStdout(), pro.OrganizationID, appTemplate)
		if err != nil {
			return
		}
		appPostInfo := model.AppPostInfo{appName, appCode, appType, apptemplate.ID, true}
		c7nclient.Client.CreateApp(cmd.OutOrStdout(), pro.ID, &appPostInfo)
	},
}

// create Env command
var createEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "create env",
	Long:  `you can use this command to create env `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userinfo.ID, proCode)
		if err != nil {
			return
		}
		err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.OrganizationID, clusterCode)
		if err != nil {
			return
		}
		envPostInfo := model.EnvPostInfo{envName, envCode, envDescription, cluster.ID}
		c7nclient.Client.CreateEnv(cmd.OutOrStdout(), pro.ID, &envPostInfo)
	},
}

var createInstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "create instance",
	Long:  `you can use this command to create instance `,
	Run: func(cmd *cobra.Command, args []string) {

		c7nclient.InitClient(&clientConfig)

		if _, err := os.Stat(content); os.IsNotExist(err) {
			fmt.Println(err)
			return
		}
		b, err := ioutil.ReadFile(content)
		release := model.Release{}
		yaml.Unmarshal(b, &release)
		if err != nil {
			fmt.Print(err)
			return
		}
		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
		if err != nil {
			return
		}
		err, app := c7nclient.Client.GetApp(release.Spec.ChartName, pro.ID)
		if err != nil {
			return
		}
		err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
		if err != nil {
			return
		}
		err, version := c7nclient.Client.GetAppVersion(cmd.OutOrStdout(), pro.ID, release.Spec.ChartVersion, app.ID)
		if err != nil {
			return
		}
		instancePostInfo := model.InstancePostInfo{version.ID, env.ID, app.ID, release.Metadata.Name, release.Spec.Values, "create", false}
		c7nclient.Client.CreateInstance(cmd.OutOrStdout(), pro.ID, &instancePostInfo)
	},
}

var createServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "create service",
	Long:  `you can use this command to create service `,
	Run: func(cmd *cobra.Command, args []string) {

		c7nclient.InitClient(&clientConfig)

		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
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

var createIngressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "create ingress",
	Long:  `you can use this command to create ingress `,
	Run: func(cmd *cobra.Command, args []string) {

		c7nclient.InitClient(&clientConfig)

		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
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

var createCertCmd = &cobra.Command{
	Use:   "cert",
	Short: "create certification",
	Long:  `you can use this command to create certification `,
	Run: func(cmd *cobra.Command, args []string) {

		c7nclient.InitClient(&clientConfig)

		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
		if err != nil {
			return
		}
		certPostInfo := model.CertificationPostInfo{}

		err = initCert(cmd, &pro, &certPostInfo)
		if err != nil {
			return
		}
		c7nclient.Client.CreateCert(cmd.OutOrStdout(), pro.ID, &certPostInfo)
	},
}

var createConfigMapCmd = &cobra.Command{
	Use:   "configMap",
	Short: "create configMap",
	Long:  `you can use this command to create configMap`,
	Run: func(cmd *cobra.Command, args []string) {

		c7nclient.InitClient(&clientConfig)

		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
		if err != nil {
			return
		}
		configMapPostInfo := model.ConfigMapPostInfo{}

		err = initConfigMap(cmd, &pro, &configMapPostInfo)
		if err != nil {
			return
		}
		c7nclient.Client.CreateConfigMap(cmd.OutOrStdout(), pro.ID, &configMapPostInfo)
	},
}

var createSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "create secret",
	Long:  `you can use this command to create secret`,
	Run: func(cmd *cobra.Command, args []string) {

		c7nclient.InitClient(&clientConfig)

		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
		if err != nil {
			return
		}
		secretPostInfo := model.SecretPostInfo{}

		err = initSecret(cmd, &pro, &secretPostInfo)
		if err != nil {
			return
		}
		c7nclient.Client.CreateSecret(cmd.OutOrStdout(), pro.ID, &secretPostInfo)
	},
}

func initService(cmd *cobra.Command, pro *model.Project, servicePostInfo *model.ServicePostInfo) (error error) {

	if _, err := os.Stat(content); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(content)
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
		servicePostInfo.AppInstance = instances
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
	servicePostInfo.Label = service.Spec.Selector
	return nil
}

func initIngress(cmd *cobra.Command, pro *model.Project, ingressPostInfo *model.IngressPostInfo) (error error) {

	if _, err := os.Stat(content); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(content)
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
			return errors.New("the service in not exist!")
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

func initCert(cmd *cobra.Command, pro *model.Project, certificationPostInfo *model.CertificationPostInfo) (error error) {

	if _, err := os.Stat(content); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(content)
	if err != nil {
		return err
	}
	certificate := model.Certificate{}
	yaml.Unmarshal(b, &certificate)
	err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
	if err != nil {
		return err
	}
	certificationPostInfo.EnvID = env.ID
	certificationPostInfo.CertName = certificate.Metadata.Name
	certificationPostInfo.CertValue = certificate.Spec.ExistCert.Cert
	certificationPostInfo.KeyValue = certificate.Spec.ExistCert.Key
	certificationPostInfo.Domains = []string{certificate.Spec.CommonName}
	certificationPostInfo.Type = "request"
	if certificationPostInfo.CertValue != "" {
		certificationPostInfo.Type = "upload"
	}
	return nil
}

func initConfigMap(cmd *cobra.Command, pro *model.Project, configMapPostInfo *model.ConfigMapPostInfo) (error error) {

	if _, err := os.Stat(content); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(content)
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
	configMapPostInfo.Description = "This is a configMap"
	configMapPostInfo.Value = configMap.Data
	return nil
}

func initSecret(cmd *cobra.Command, pro *model.Project, secretPostInfo *model.SecretPostInfo) (error error) {

	if _, err := os.Stat(content); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	b, err := ioutil.ReadFile(content)
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
	secretPostInfo.Description = "This is a secret"
	secretPostInfo.Value = secret.StringData
	return nil
}
