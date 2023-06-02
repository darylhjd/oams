import 'dart:convert';
import 'dart:io';

import 'package:http/browser_client.dart';

import 'models.dart';

// APIClient helps to interface with the OAMS API.
class APIClient {
  static final client = () {
    var client = BrowserClient();
    client.withCredentials = true;
    return client;
  }();

  // TODO: Use .env for this.
  static const String apiHost = "localhost";
  static const int apiPort = 8080;

  static const String defaultRedirectUrl = "http://localhost:8000/";

  static const String loginPath = "api/v1/login";
  static const String loginReturnToParam = "return_to";
  static const String logoutPath = "api/v1/sign-out";
  static const String userPath = "api/v1/user";

  static Future<String> getLoginURL(String returnTo) async {
    if (returnTo.isEmpty) {
      returnTo = defaultRedirectUrl;
    }

    final uri = Uri(
      scheme: "http",
      host: apiHost,
      port: apiPort,
      path: loginPath,
      queryParameters: {
        loginReturnToParam: returnTo,
      },
    );

    final response = await client.get(uri);
    if (response.statusCode != HttpStatus.ok) {
      return Future.error(const HttpException("cannot get login url"));
    }

    final loginResponse = LoginResponse.fromJson(jsonDecode(response.body));
    return loginResponse.redirectUrl;
  }

  static Future<bool> logout() async {
    final uri = Uri(
      scheme: "http",
      host: apiHost,
      port: apiPort,
      path: logoutPath,
    );

    final response = await client.get(uri);
    return response.statusCode == HttpStatus.ok;
  }

  static Future<User> getUserInfo() async {
    final uri =
        Uri(scheme: "http", host: apiHost, port: apiPort, path: userPath);

    final response = await client.get(uri);
    if (response.statusCode != HttpStatus.ok) {
      return Future.error(const HttpException("cannot get user details"));
    }

    return User.fromJson(jsonDecode(response.body));
  }
}
