/// <reference path="./.sst/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "conman",
      home: "aws",
      providers: { aws: { region: "us-east-2" } },
      removal: input?.stage === "prod" ? "retain" : "remove",
    };
  },
  async run() {
    const databaseUrl = new sst.Secret("DatabaseUrl");
    const apiHandler = new sst.aws.Function("ApiHandler", {
      runtime: "go",
      architecture: "arm64",
      handler: "./backend/cmd/api/main.go",
      timeout: "29 seconds",
      environment: {
        DATABASE_URL: databaseUrl.value,
        ENVIRONMENT: $app.stage,
        APP_ENV: $app.stage,
        DB_MIN_CONNS: "0",
        DB_MAX_CONNS: "4",
      },
    });

    const httpApi = new sst.aws.ApiGatewayV2("HttpApi", {
      cors:
        $app.stage === "prod"
          ? false
          : {
              allowOrigins: ["http://localhost:5173"],
              allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
              allowHeaders: ["Content-Type", "Authorization", "X-Dev-User-Id"],
            },
    });
    httpApi.route("ANY /{proxy+}", apiHandler.arn);

    return { apiUrl: httpApi.url };
  },
});
