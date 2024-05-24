import {
    Callback,
    CloudFrontHeaders,
    CloudFrontRequestEvent,
    Context,
} from "aws-lambda";
import {
    createRemoteJWKSet,
    jwtVerify,
    type JWTVerifyGetKey,
    type JWTVerifyResult,
} from "jose";
import type { CallbackParamsType } from "openid-client";
import { Issuer, generators } from "openid-client";
import { getConfig } from "./config";

type Cookies = Record<string, string>;

type Scope = Record<string, string>;

type Scopes = Record<string, Scope>;

const oidcConfig = getConfig();

export const iss = new Issuer({
    issuer: oidcConfig.issuer,
    authorization_endpoint: oidcConfig.authorizationEndpoint,
    token_endpoint: oidcConfig.tokenEndpoint,
    jwks_uri: oidcConfig.jwksUri,
});

const client = new iss.Client({
    client_id: oidcConfig.clientId,
    redirect_uris: [oidcConfig.redirectUri],
    response_types: ["code"],
    token_endpoint_auth_method: "none",
});

// Initialize the JWK Set to verify the access token
const options = {
    cacheMaxAge: oidcConfig.cacheMaxAgeMilliseconds,
};
const jwkSet: JWTVerifyGetKey = createRemoteJWKSet(
    new URL(oidcConfig.jwksUri),
    options
);

// isAuthenticated checks for the presence of the access token in the session
async function isAuthenticated(
    parsedCookies: Cookies
): Promise<JWTVerifyResult> {
    const accessToken = parsedCookies["zitadel.session"];
    if (!accessToken) {
        throw new Error("No access token found");
    }

    return await jwtVerify(accessToken, jwkSet);
}

// isAuthorized checks for the presence of the required role in the session
async function isAuthorized(result: JWTVerifyResult): Promise<boolean> {
    const roles = result.payload[oidcConfig.roleKey];
    if (!roles || typeof roles !== "object") {
        return false;
    }

    const r = roles as Scopes;
    if (!r[oidcConfig.requiredRole]) {
        return false;
    }

    return true;
}

// parseCookies parses cookies from the request headers
function parseCookies(headers: CloudFrontHeaders): Cookies {
    const parsedCookie: Cookies = {};
    if (headers.cookie) {
        headers.cookie[0].value.split(";").forEach((cookie: string) => {
            if (cookie) {
                const parts = cookie.split("=");
                const trimmedKey = parts[0].trim();
                const trimmedVal = parts[1].trim();
                parsedCookie[trimmedKey] = trimmedVal;
            }
        });
    }

    return parsedCookie;
}

// The handler function is the entry point for the Lambda function
export const handler = async (
    event: CloudFrontRequestEvent,
    _context: Context,
    callback: Callback
) => {
    const request = event.Records[0].cf.request;
    const headers = request.headers;
    const parsedCookies = parseCookies(headers);
    const url = request.uri;
    const querystring = request.querystring;

    if (url.startsWith("/auth/callback")) {
        console.log("Running handler for /auth/callback");

        // Handle callback after authentication
        const code_verifier = parsedCookies["zitadel.verifier"];
        const code = new URLSearchParams(querystring).get("code") ?? undefined;

        var callbackParams: CallbackParamsType = {
            code: code as string,
            redirect_uri: oidcConfig.redirectUri,
        };

        try {
            // Exchange code for tokens using PKCE code_verifier
            const tokenSet = await client.callback(
                oidcConfig.redirectUri,
                callbackParams,
                { code_verifier }
            );

            // Redirect back to the root URL after authentication
            console.log("Redirecting back after authentication");
            const response = {
                status: "302",
                statusDescription: "Found",
                headers: {
                    location: [
                        {
                            key: "location",
                            value: `/`,
                        },
                    ],
                    "set-cookie": [
                        {
                            key: "set-cookie",
                            value: `zitadel.session=${tokenSet.id_token}; Path=/; HttpOnly; Max-Age=3600`,
                        },
                    ],
                },
            };
            callback(null, response);
            return;
        } catch (e) {
            console.error(e);

            // Display error message
            const data = {
                status: 500,
                body: "Error during authentication",
            };
            callback(null, data);
            return;
        }
    } else if (url.startsWith("/auth/logout")) {
        console.log("Running handler for /auth/logout");

        // Clear the session and redirect to the post logout url
        const response = {
            status: "302",
            statusDescription: "Found",
            headers: {
                location: [
                    {
                        key: "location",
                        value: `/logout`,
                    },
                ],
                "set-cookie": [
                    {
                        key: "set-cookie",
                        value: "zitadel.session=; Path=/; HttpOnly; Max-Age=0",
                    },
                ],
            },
        };
        callback(null, response);
        return;
    } else if (url.startsWith("/logout")) {
        console.log("Running handler for /logout");
        const data = {
            status: 200,
            body: "You are logged out",
        };
        callback(null, data);
        return;
    } else {
        console.log("Running handler for /*");

        // Check if the user is authenticated
        var authN: JWTVerifyResult | undefined = undefined;
        try {
            authN = await isAuthenticated(parsedCookies);
        } catch (e) {
            console.log("User is not authenticated");
            console.error(e);

            console.log("Generating authorization URL");

            // If the user is not authenticated...
            // generate PKCE code_challenge
            const code_verifier = generators.codeVerifier();
            const code_challenge = generators.codeChallenge(code_verifier);

            // construct authorization URL with PKCE parameters
            const authorizationUrl = client.authorizationUrl({
                scope: "openid profile email", // adjust scopes as needed
                code_challenge,
                code_challenge_method: "S256", // use SHA-256 for PKCE
            });

            // write code_verifier to session and redirect to the authorization URL
            const response = {
                status: "302",
                statusDescription: "Found",
                headers: {
                    location: [
                        {
                            key: "location",
                            value: authorizationUrl,
                        },
                    ],
                    "set-cookie": [
                        {
                            key: "set-cookie",
                            value: `zitadel.verifier=${code_verifier}; Path=/; HttpOnly; Max-Age=1800`,
                        },
                    ],
                },
            };
            callback(null, response);
            return;
        }

        // Check if the user is authorized
        const authZ = await isAuthorized(authN);

        if (authZ) {
            console.log("User is authorized. Serving content...");

            callback(null, request);
            return;
        }

        console.log("User is not authorized");
        const data = {
            status: 403,
            body: "Forbidden",
        };
        callback(null, data);
        return;
    }
};
