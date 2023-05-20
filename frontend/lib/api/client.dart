import 'dart:convert';

import 'package:http/http.dart' as http;

class APIClient {
  static final client = http.Client();
  static const String apiHost = "localhost";
  static const int apiPort = 8080;

  static const String loginPath = "api/v1/login";

  static Future<String> getLoginURL(Map<String, String> queryParams) async {
    String returnUrl = "http://localhost:8000/";
    if (queryParams.containsKey("return_to")) {
      returnUrl = queryParams["return_to"]!;
    }

    final uri = Uri(
      scheme: "http",
      host: apiHost,
      port: apiPort,
      path: loginPath,
      queryParameters: {
        "return_to": returnUrl,
      },
    );

    final response = await client.get(uri);
    final loginResponse = _LoginResponse.fromJson(jsonDecode(response.body));

    return loginResponse.redirectUrl;
  }
}

// _LoginResponse is a data class modelling the response from a login API request.
class _LoginResponse {
  final String redirectUrl;

  _LoginResponse(this.redirectUrl);

  _LoginResponse.fromJson(Map<String, dynamic> json)
      : redirectUrl = json['redirect_url'];
}
