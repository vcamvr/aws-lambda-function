package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

const BaseArn string = "arn:aws-cn:ecs:%s:%s:service/%s/%s"

type eventReq struct {
	Account   string `json:"account"`
	Cluster   string `json:"cluster"`
	Service   string `json:"service"`
	AlarmName string `json:"alarmName"`
}

func UpdateCollector(ctx context.Context, req *eventReq) error {
	log.Printf("aws will to update {%s} service, triggered by {%s} alarm", req.Service, req.AlarmName)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	if cfg.Region == "" {
		return errors.New("default config Region is nil")
	}
	cred, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	} else if cred.AccessKeyID == "" || cred.SecretAccessKey == "" {
		return errors.New("default config AccessKeyID or AccessKey is nil")
	}
	arn := fmt.Sprintf(BaseArn, cfg.Region, req.Account, req.Cluster, req.Service)
	svc := ecs.NewFromConfig(cfg)
	in := &ecs.DescribeServicesInput{
		Services: []string{arn},
		Cluster:  &req.Cluster,
	}
	ins, err := svc.DescribeServices(ctx, in)
	if err != nil {
		return err
	} else if len(ins.Services) == 0 {
		return errors.New(fmt.Sprintf("can't find service by arn: %s", arn))
	}
	u := ecs.UpdateServiceInput{
		Service:                       ins.Services[0].ServiceArn,
		Cluster:                       ins.Services[0].ClusterArn,
		TaskDefinition:                ins.Services[0].TaskDefinition,
		DesiredCount:                  &ins.Services[0].DesiredCount,
		CapacityProviderStrategy:      ins.Services[0].CapacityProviderStrategy,
		DeploymentConfiguration:       ins.Services[0].DeploymentConfiguration,
		PlacementStrategy:             ins.Services[0].PlacementStrategy,
		LoadBalancers:                 ins.Services[0].LoadBalancers,
		HealthCheckGracePeriodSeconds: ins.Services[0].HealthCheckGracePeriodSeconds,
		ForceNewDeployment:            true,
	}
	_, err = svc.UpdateService(ctx, &u)
	if err != nil {
		return err
	}
	log.Printf("update {%s} service success", req.Service)
	return nil
}

func main() {
	lambda.Start(UpdateCollector)
}
