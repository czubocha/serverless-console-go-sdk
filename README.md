# serverless-console-go-sdk
### Serverless Console instrumentation for Go language.

### Features
* request and response data monitoring

### Setup
#### 1. Register with [Serverless Console](https://console.serverless.com/)
#### 2. In [Serverless Console](https://console.serverless.com/) turn on dev mode integration for your AWS account and chosen Lambdas
#### 3. In lambda main package replace import `github.com/aws/aws-lambda-go/lambda` with `github.com/czubocha/serverless-console-go-sdk/lambda`
#### 4. (optionally) Fine tune default instrumentation behavior with following options

##### `SLS_DISABLE_REQUEST_RESPONSE_MONITORING`
(Dev mode only) Disable monitoring requests and responses (function, AWS SDK requests and HTTP(S) requests)

