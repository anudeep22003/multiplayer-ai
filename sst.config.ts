/// <reference path="./.sst/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "multiplayer-ai",
      home: "aws",
      providers: {
        aws: {
          region: "us-east-2",
          profile: process.env.CI ? undefined : "sst-dev",
        }
      },
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
            allowOrigins: ["*"],
            allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
            allowHeaders: ["Content-Type", "Authorization", "X-Dev-User-Id"],
          },
    });
    httpApi.route("ANY /{proxy+}", apiHandler.arn);

    const staticSite = new sst.aws.StaticSite("Web", {
      path: "frontend",
      build: {
        command: "bun run build",
        output: "dist"
      },
      environment: {
        VITE_API_URL: httpApi.url,
      },
    })

    return { apiUrl: httpApi.url, siteUrl: staticSite.url };
  },
});

