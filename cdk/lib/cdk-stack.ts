import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import path = require("path");

export class CdkStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const lambdaFn = new lambda.Function(this, "contact-api", {
            runtime: lambda.Runtime.GO_1_X,
            code: lambda.Code.fromAsset(
                path.join(__dirname, "..", "build", "app")
            ),
            handler: "main",
            environment: {
                TARGET_EMAIL: "",
            },
        });
    }
}
