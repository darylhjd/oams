import 'package:flutter_dotenv/flutter_dotenv.dart';

const _apiServerHost = "API_SERVER_HOST";
const _apiServerPort = "API_SERVER_PORT";
const _webServerHost = "WEB_SERVER_HOST";
const _webServerPort = "WEB_SERVER_PORT";

String apiServerHost() {
  return dotenv.get(_apiServerHost);
}

String apiServerPort() {
  return dotenv.get(_apiServerPort);
}

String webServerHost() {
  return dotenv.get(_webServerHost);
}

String webServerPort() {
  return dotenv.get(_webServerPort);
}
