package util

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"gopkg.in/ini.v1"
)

var home string

func init() {
	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
}

func GetProfiles() ([]string, error) {

	awsConfig, err := ini.Load(filepath.Join(home, ".aws/config"))
	if err != nil {
		return nil, err
	}

	profiles := []string{}
	for _, s := range awsConfig.SectionStrings() {
		if s == "DEFAULT" {
			continue
		}
		profiles = append(profiles, strings.TrimSpace(strings.TrimPrefix(s, "profile ")))
	}

	return profiles, nil

}

func GetRegions(ctx context.Context, cfg aws.Config) ([]string, error) {

	// region is required for call to get all regions
	cfg.Region = "eu-west-1"
	svc := ec2.NewFromConfig(cfg)
	out, err := svc.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(false), // exclude disabled regions
	})
	if err != nil {
		return nil, err
	}

	regions := []string{}
	for _, r := range out.Regions {
		regions = append(regions, *r.RegionName)
	}

	return regions, nil

}
