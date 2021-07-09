# crawl

Crawl AWS accounts & regions executing a given function against each and accumulating the results

## Background

When you have a lot of AWS accounts and use multiple regions it can be hard to operate against these efficiently.
An example could be that AWS have deprecated a given Lambda runtime so you need to find all Lambdas across all of your accounts and regions that use the deprecated runtime.
That would involve time consuming manual searching or the creation of scripts to go through each of your accounts. 
Even then, once you get to over 10 accounts and multiply that by the number of AWS regions, these scripts can be quite slow/inefficient to execute.

This package aims to solve that problem by providing a Go package which can be used to execute any given function across all your AWS accounts (profiles in your `~/.aws/config`) and then all enabled regions in those accounts. It uses Go routines & channels so that all requests can be executed concurrently, greatly improving the efficiency and time taken to run, regardless of the number of accounts & regions you need search in.
