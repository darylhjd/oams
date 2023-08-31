export function buildRedirectUrlQueryParamsString(path: string): string {
  const redirectUrl = `${process.env.WEB_SERVER_HOST}:${process.env.WEB_SERVER_PORT}${path}`;
  const queryParams = new URLSearchParams({
    redirect_url: redirectUrl,
  });
  return queryParams.toString();
}
