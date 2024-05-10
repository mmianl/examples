# Using Zitadel as autentication server for a static website hosted on AWS S3 and distributed via Amazon CloudFront

In this tutorial, we'll explore the integration of Zitadel with a static website hosted on AWS S3 and distributed via Amazon CloudFront. By combining these tools, we can establish a seamless authentication flow while leveraging the simplicity and reliability of a static website. 

It's worth noting that static site generators like Hugo, Jekyll, or Gatsby could be used to generate the website content. However, regardless of the tool used for generating the static content, the integration with Zitadel and AWS S3/CloudFront remains consistent.

The tutorial will be structured as follows:

* We'll start by setting up a Zitadel organization including a project and an OIDC app.
* Next, we'll create an AWS S3 bucket to host our static website content. We'll also provide an IAM user that could be used in a deployment pipeline.
* To handle the authentication flow, we'll deploy a Lambda@Edge function. This serverless function will intercept requests to our CloudFront distribution and validate user authentication with Zitadel.
* Finally, we'll configure Amazon CloudFront to distribute our static website content globally, ensuring low-latency access for users worldwide. Additionally, we'll set up AWS Web Application Firewall (WAF) to add an extra layer of security by filtering and monitoring HTTP traffic.

Note that the setup of Zitadel itself is out of scope.

## Zitadel

We define a Zitadel organization, a project within the organization and an OIDC application to handle authentication for our website. This application allows users to authenticate via OpenID Connect (OIDC) protocol, enabling secure access to our resources. 

### Role Grants

We also configure role grants within our Zitadel organization to manage access control for users authenticating via OpenID Connect (OIDC) protocol. Role grants define the permissions granted to users or groups within a project, allowing fine-grained control over resource access. A Zitadel action is defined to automatically assign the role grant to all authenticated users or groups in the organization.

## S3

We define an S3 bucket to store our static website content with server-side encryption enabled, and versioning disabled. Access control is set to private and a bucket policy is defined ensuring that only the CloudFront distribution can access the bucket contents (also see CloudFront / Origin Access Control below). Additionally, an IAM user is created for CI/CD purposes with necessary permissions to manage objects in the S3 bucket and an access key is generated for the CI/CD user to authenticate API requests when interacting with the S3 bucket.

## Lambda@Edge

Now, we'll configure a Lambda@Edge function to handle the authentication flow for our static website. Lambda@Edge allows us to execute custom code in response to CloudFront events, enabling serverless processing at the edge locations closest to the users.

However, it's important to note that Lambda@Edge does not support environment variables, which are commonly used for configuration in Lambda functions. To work around this limitation, we'll set up a configuration file that contains the necessary settings for authentication.

Furthermore, Lambda@Edge functions must be deployed in the us-east-1 (N. Virginia) region which is the resason why we have setup multiple AWS providers and why some resources are explicitly created in us-east-1 (`provider = aws.aws_useast`).

### Function

The Lambda@Edge function is written in TypeScript. It handles the authentication flow for a static website hosted on CloudFront using OpenID Connect (OIDC) authentication with Zitadel as the identity provider.

The isAuthenticated function checks for the presence of the access token in the session and verifies it using the JWK Set. It's important to note that the Issuer, which manages the JWK Set, is created outside of the Lambda handler (`export const handler = async ()`). Clients and variables that are declared outside of the handler method can be reused for subsequent events because Lambda maintains the execution environment for some time [AWS Lambda Developer Guide  / Lambda programming model]. The Issuer itself caches the JWK Set for a configurable amount of time (default: 600000 ms, see `config.json.tftpl`) to reduce the number of requests made against Zitadel's JWK Set endpoint and to speed up requests made to the website.

The isAuthorized function checks for the presence of the required role (see section Zitadel / Role Grants above).

The primary Lambda handler function manages various types of requests: authentication callbacks (/auth/callback) following user login, logout requests (/auth/logout) to clear the session, and all other requests (/*) that require user authentication and authorization.

## CloudFront

Next, we define the CloudFront distribution along with ACM certificate including its validation with Route 53 records. It configures the CloudFront distribution to serve content from an S3 bucket, restricts access based on geographic locations, and enables SSL/TLS encryption using the ACM certificate. Additionally, it sets up the required Route 53 records for CloudFront distribution.

### Origin Access Control

Origin access control is configured to ensure that all requests to the website must go through CloudFront. The S3 bucket policy is configured to allow access only from the distribution's OAC. Also note, that the OAC needs to be added to the distribution's origin in the `origin {}` block of `aws_cloudfront_distribution.website`.

## WAF

Last, we set up a comprehensive Web Application Firewall for the CloudFront distribution, providing protection against various types of threats such as bot traffic and excessive requests, while also enabling logging for monitoring and analysis.

AWS Shield Standard, which is included at no extra cost, can minimize the effects of Distributed Denial of Service (DDoS) attacks.
