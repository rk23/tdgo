# TDGO

Go library for interacting with the TD Ameritrade API. Includes 

Under heavy development as a personal project but free to use if you find it helpful. 

## Prerequisites

- `~/.aws/config` && `~/.aws/credentials` populated with your info
- S3 Bucket for storing lambda code

## How To

Check out the tasks folder to see how to build and deploy the lambda OAuth handler. The basic commands are:

- `inv build  # Builds the main binary and stores it in a consistent location`
- `inv package  # Upload code to S3 bucket, generate cloudformation`
- `inv deploy # Update live stack`
- `inv ship --s3-bucket BUCKET_NAME # All the above`

You'll need to have this deployed before you create your app, TD doesn't let you update 
your app and it's very specific about using _exactly_ the uri you specified. 

I've included a cloudformation example which is 
simply a replica of the cloudformation I'm using with obfuscated env 
variables. From root:

- `mv /deploy/cfn/auth-sam-example.yml /deploy/cfn/auth-sam-prod.yml` 

Add env variables when you get them and customize as you see fit. You can deploy without updating the env vars but you _will_ need to update these env variables for the handler to work. 

Check out the .png in docs to get a better understanding of how the 
authentication flow works.