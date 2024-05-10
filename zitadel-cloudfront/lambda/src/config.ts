import { readFileSync } from "fs";

type OIDCConfig = {
    issuer: string;
    clientId: string;
    redirectUri: string;
    jwksUri: string;
    tokenEndpoint: string;
    authorizationEndpoint: string;
    requiredRole: string;
    roleKey: string;
    cacheMaxAgeMilliseconds: number;
};

export function getConfig(): OIDCConfig {
    try {
        return getConfigFromFile("config.json");
    } catch (err) {
        return getConfigFromEnv();
    }
}

function getConfigFromFile(filePath: string): OIDCConfig {
    try {
        const data: string = readFileSync(filePath, "utf8");
        const userData: OIDCConfig = JSON.parse(data);
        return userData;
    } catch (err) {
        throw new Error(`Error reading or parsing config JSON: ${err}`);
    }
}

function getConfigFromEnv(): OIDCConfig {
    const issuer = process.env.ISSUER;
    const clientId = process.env.CLIENT_ID;
    const redirectUri = process.env.REDIRECT_URI;
    const jwksUri = process.env.JWKS_URI;
    const tokenEndpoint = process.env.TOKEN_ENDPOINT;
    const authorizationEndpoint = process.env.AUTHORIZATION_ENDPOINT;
    const requiredRole = process.env.REQUIRED_ROLE;
    const roleKey = process.env.ROLE_KEY;
    const cacheMaxAgeMilliseconds = parseInt(
        process.env.CACHE_MAX_AGE_MILLISECONDS ?? "600000"
    );

    if (
        !issuer ||
        !clientId ||
        !redirectUri ||
        !jwksUri ||
        !tokenEndpoint ||
        !authorizationEndpoint ||
        !requiredRole ||
        !roleKey ||
        !cacheMaxAgeMilliseconds
    ) {
        throw new Error(
            "Failed to get config. Missing required environment variables."
        );
    }

    return {
        issuer,
        clientId,
        redirectUri,
        jwksUri,
        tokenEndpoint,
        authorizationEndpoint,
        requiredRole,
        roleKey,
        cacheMaxAgeMilliseconds,
    };
}
