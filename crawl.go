package crawl

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/kloudyuk/crawl/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var errorc = make(chan error)
var resultsc = make(chan interface{})

type CrawlFunc func(context.Context, aws.Config) (interface{}, error)

func Exec(fn CrawlFunc) []interface{} {

	var wg sync.WaitGroup
	var results []interface{}

	profiles, err := util.GetProfiles()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case err := <-errorc:
				log.Fatal(err)
			case r := <-resultsc:
				results = append(results, r)
			}
		}
	}()

	for _, p := range profiles {
		wg.Add(1)
		go crawlAccount(&wg, p, fn)
	}

	wg.Wait()

	return results

}

func crawlAccount(wg *sync.WaitGroup, profile string, fn CrawlFunc) {

	defer wg.Done()

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		errorc <- fmt.Errorf("profile: %s, error: %s", profile, err)
		return
	}

	regions, err := util.GetRegions(ctx, cfg)
	if err != nil {
		errorc <- fmt.Errorf("profile: %s, error: %s", profile, err)
		return
	}

	var regionWG sync.WaitGroup

	// For each region, execute the func
	for _, region := range regions {
		regionWG.Add(1)
		go crawlRegion(ctx, cfg, &regionWG, profile, region, fn)
	}

	regionWG.Wait()

}

func crawlRegion(ctx context.Context, cfg aws.Config, wg *sync.WaitGroup, profile string, region string, fn CrawlFunc) {

	defer wg.Done()

	cfg.Region = region
	results, err := fn(ctx, cfg)
	if err != nil {
		errorc <- fmt.Errorf("profile: %s, region: %s, error: %s", profile, region, err)
		return
	}

	resultsc <- results

}
