const _apiServerHost = "API_SERVER_HOST";
const _apiServerPort = "API_SERVER_PORT";
const _webServerHost = "WEB_SERVER_HOST";
const _webServerPort = "WEB_SERVER_PORT";

checkEnvVars() {
  if (apiServerHost() == "" ||
      apiServerPort() == "" ||
      webServerHost() == "" ||
      webServerPort() == "") {
    throw Exception("environment variables not set up properly");
  }
}

String apiServerHost() {
  return const String.fromEnvironment(_apiServerHost);
}

String apiServerPort() {
  return const String.fromEnvironment(_apiServerPort);
}

String webServerHost() {
  return const String.fromEnvironment(_webServerHost);
}

String webServerPort() {
  return const String.fromEnvironment(_webServerPort);
}
