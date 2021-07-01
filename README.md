base-of-lambda-container-image
---

This is a base of the AWS Lambda function in Golang and deploys it as a container image.

## Usage

1. Clone repository

    ```bash
    git clone https://github.com/michimani/base-of-lambda-container-image.git \
    && cd base-of-lambda-container-image
    ```

2. Build the image

    ```bash
    docker build -t base-of-lambda-container-image .
    ```

## Run at local

1. Run the image

    ```bash
    docker run \
    --rm \
    -p 9000:8080 \
    base-of-lambda-container-image:latest /main
    ```

3. Invoke function

    ```bash
    curl "http://localhost:9000/2015-03-31/functions/function/invocations"
    ```

## Deploy to Lmabda

1. Build the image

    ```bash
    docker build -t base-of-lambda-container-image .
    ```

2. Create ECR repository

    ```bash
    aws ecr create-repository \
    --repository-name base-of-lambda-container-image \
    --region ap-northeast-1
    ```

3. Login to ECR

    ```bash
    aws ecr get-login-password --region ap-northeast-1 \
    | docker login \
    --username AWS \
    --password-stdin $(aws sts get-caller-identity \
    --query 'Account' \
    --output text).dkr.ecr.ap-northeast-1.amazonaws.com
    ```
    
4. Add image tag

    ```bash
    docker tag base-of-lambda-container-image:latest ************.dkr.ecr.ap-northeast-1.amazonaws.com/base-of-lambda-container-image:latest
    ```
    
5. Push to ECR repository

    ```bash
    docker push ************.dkr.ecr.ap-northeast-1.amazonaws.com/base-of-lambda-container-image:latest
    ```

6. Create Lambda Function from container image

    ```bash
    aws lambda create-function \
    --function-name "base-of-lambda-container-image" \
    --package-type "Image" \
    --code "ImageUri=<your-ecr-repository-uri>" \
    --timeout 30 \
    --role "<your-iam-role-arn>" \
    --region "<your-region-code>"
    ```

7. Invoke function

    ```bash
    aws lambda invoke \
    --function-name base-of-lambda-container-image \
    --invocation-type RequestResponse \
    --region ap-northeast-1 \
    out \
    && cat out | jq .
    ```
