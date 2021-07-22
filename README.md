# crawl

Crawl AWS accounts & regions executing a given function against each and accumulating the results

## Background

When you have a lot of AWS accounts and use multiple regions it can be hard to operate against these efficiently.
An example could be that AWS have deprecated a given Lambda runtime so you need to find all Lambdas across all of your accounts and regions that use the deprecated runtime.
That would involve time consuming manual searching or the creation of scripts to go through each of your accounts.
Even then, once you get to over 10 accounts and multiply that by the number of AWS regions, these scripts can be quite slow/inefficient to execute.

This package aims to solve that problem by providing a Go package which can be used to execute any given function across all your AWS accounts (profiles in your `~/.aws/config`) and then all enabled regions in those accounts. It uses Go routines & channels so that all requests can be executed concurrently, greatly improving the efficiency and time taken to run, regardless of the number of accounts & regions you need search in.

## Usage

```sh
go get -u github.com/kloudyuk/crawl
```

## Example

The following example shows how to get a list of Lambda functions across all accounts & regions using this library

```go
package main

import (
  "context"
  "fmt"

  "github.com/kloudyuk/crawl"

  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/service/lambda"
  "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// Define the function we want to run in every account & region
// The function signature must match the signature at:
// https://github.com/kloudyuk/crawl/blob/main/crawl.go#L18
func getLambdas(ctx context.Context, profile string, cfg aws.Config) (interface{}, error) {
  svc := lambda.NewFromConfig(cfg)
  in := &lambda.ListFunctionsInput{}
  out, err := svc.ListFunctions(ctx, in)
  if err != nil {
    return nil, err
  }
  return out.Functions, nil
}

func main() {

  // Execute the function across all accounts & regions
  // The function is executed for each profile/region combo concurrently using go routines
  // crawl.Exec blocks until all go routines are complete
  // results contains an interface containing the results from all accounts & regions
  results := crawl.Exec(getLambdas)

  // Cast the result interfaces back to concrete types using a type assertion
  lambdas := []types.FunctionConfiguration{}
  for _, r := range results {
    if l, ok := r.([]types.FunctionConfiguration); ok {
      lambdas = append(lambdas, l...)
    }
  }

  // Do whatever with the results
  for _, lambda := range lambdas {
    fmt.Println(*lambda.FunctionArn)
  }

}
```
