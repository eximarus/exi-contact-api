import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as path from "path";

export class CdkStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const lambdaFn = new lambda.Function(this, "ContactApiLambdaFn", {
            runtime: lambda.Runtime.GO_1_X,
            code: lambda.Code.fromAsset(
                path.join(__dirname, "..", "..", "build", "main.zip")
            ),
            handler: "main",
            environment: {
                TARGET_EMAIL: "",
                GMAIL_USER: "",
                GMAIL_PASSWORD: "",
            },
        });
        new apigateway.LambdaRestApi(this, "ContactApiLambdaFnEndpoint", {
            handler: lambdaFn,
        });
    }
}
