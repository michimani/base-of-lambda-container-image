base-of-lambda-container-image
---

This is a base of the AWS Lambda function in Golang and deploys it as a container image.

## Usage

1. Clone repository

    ```bash
    git clone https://github.com/michimani/base-of-lambda-container-image.git \
    && cd base-of-lambda-container-image
    ```

2. Build the image for local

    ```bash
    docker build -t base-of-lambda-container-image:local -f Dockerfile_local .
    ```

## Run at local

1. Run the image

    ```bash
    docker run \
    --rm \
    -p 9000:8080 \
    base-of-lambda-container-image:local
    ```

3. Invoke function

    ```bash
    curl "http://localhost:9000/2015-03-31/functions/function/invocations"
    ```

## Deploy to Lmabda

1. Build the image for production

    ```bash
    docker build -t base-of-lambda-container-image:prod .
    ```

2. Create ECR repository

    ```bash
    REGION='ap-northeast-1' \
    && aws ecr create-repository \
    --repository-name base-of-lambda-container-image \
    --region ${REGION}
    ```

3. Login to ECR

    ```bash
    AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text) \
    && aws ecr get-login-password --region ${REGION} \
    | docker login \
    --username AWS \
    --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com
    ```
    
4. Add image tag

    ```bash
    docker tag base-of-lambda-container-image:prod ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/base-of-lambda-container-image:latest
    ```
    
5. Push to ECR repository

    ```bash
    docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/base-of-lambda-container-image:latest
    ```

6. Create IAM Role for the function

    ```bash
    ROLE_ARN=$(
        aws iam create-role \
        --role-name role-for-base-of-lambda-container-image \
        --assume-role-policy-document file://trust-policy.json \
        --query 'Role.Arn' \
        --output text
    )
    ```

7. Create Lambda Function from container image

    ```bash
    aws lambda create-function \
    --function-name "base-of-lambda-container-image" \
    --package-type "Image" \
    --code "ImageUri=${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/base-of-lambda-container-image:latest" \
    --timeout 30 \
    --role "${ROLE_ARN}" \
    --region ${REGION}
    ```

8. Invoke function

    ```bash
    aws lambda invoke \
    --function-name base-of-lambda-container-image \
    --invocation-type RequestResponse \
    --region ${REGION} \
    out \
    && cat out | jq .
    ```
    
9.  (Update function code)

    ```bash
    aws lambda update-function-code \
    --function-name base-of-lambda-container-image \
    --image-uri "${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/base-of-lambda-container-image:latest"
    ```

