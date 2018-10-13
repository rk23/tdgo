from invoke import task


@task
def build(c): 
    c.run("GOOS=linux GOARCH=amd64 go build -o bin/main cmd/auth/main.go")
    c.run("zip dist/main.zip -j bin/main")

@task
def package(c, s3_bucket=None):
    c.run("aws cloudformation package --template-file deploy/cfn/auth-sam-prod.yml "
          "--s3-bucket {} --output-template-file deploy/cfn/auth-prod.yml".format(s3_bucket))

@task
def deploy(c):
    c.run("aws cloudformation deploy --template-file deploy/cfn/auth-prod.yml "
          "--stack-name OAuth --capabilities CAPABILITY_IAM")

@task
def ship(c, s3_bucket=None):
    build(c)
    package(c, s3_bucket=s3_bucket)
    deploy(c)
