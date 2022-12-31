import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as path from "path";
import * as dotenv from "dotenv";
import { CfnOutput } from "aws-cdk-lib";
import { FunctionUrlAuthType } from "aws-cdk-lib/aws-lambda";
dotenv.config();

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
                TARGET_EMAIL:
                    process.env.TARGET_EMAIL ??
                    (() => {
                        throw new Error("TARGET_EMAIL is not defined");
                    })(),
                SMTP_USER:
                    process.env.SMTP_USER ??
                    (() => {
                        throw new Error("SMTP_USER is not defined");
                    })(),
                SMTP_PASSWORD:
                    process.env.SMTP_PASSWORD ??
                    (() => {
                        throw new Error("SMTP_PASSWORD is not defined");
                    })(),
                SMTP_HOST:
                    process.env.SMTP_HOST ??
                    (() => {
                        throw new Error("SMTP_HOST is not defined");
                    })(),
                SMTP_PORT:
                    process.env.SMTP_PORT ??
                    (() => {
                        throw new Error("SMTP_PORT is not defined");
                    })(),
            },
        });

        const fnUrl = lambdaFn.addFunctionUrl({
            authType: FunctionUrlAuthType.NONE,
        });
        new CfnOutput(this, "ContactApiLambdaFnUrl", {
            value: fnUrl.url,
        });
        new apigateway.LambdaRestApi(this, "ContactApiLambdaFnEndpoint", {
            handler: lambdaFn,
        });
    }
}
