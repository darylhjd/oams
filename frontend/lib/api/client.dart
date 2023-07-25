import 'dart:convert';
import 'dart:io';

import 'package:frontend/env/env.dart';
import 'package:http/browser_client.dart';

import 'models.dart';

// APIClient helps to interface with the OAMS API.
class APIClient {
  static final client = () {
    var client = BrowserClient();
    client.withCredentials = true;
    return client;
  }();

  static Uri apiUri = Uri.parse("${apiServerHost()}:${apiServerPort()}");

  static String defaultRedirectUrl = "${webServerHost()}:${webServerPort()}";

  static const String loginPath = "api/v1/login";
  static const String loginRedirectUrlParam = "redirect_url";
  static const String logoutPath = "api/v1/sign-out";
  static const String usersPath = "api/v1/users";

  // getLoginURL gets a login URL from the APIServer.
  static Future<String> getLoginURL(String returnTo) async {
    if (returnTo.isEmpty) {
      returnTo = defaultRedirectUrl;
    }

    final uri = apiUri.replace(
      path: loginPath,
      queryParameters: {
        loginRedirectUrlParam: returnTo,
      },
    );

    final response = await client.get(uri);
    if (response.statusCode != HttpStatus.ok) {
      return Future.error(const HttpException("cannot get login url"));
    }

    final loginResponse = LoginResponse.fromJson(jsonDecode(response.body));
    return loginResponse.redirectUrl;
  }

  // logout removes the current logged in session.
  static Future<bool> logout() async {
    final uri = apiUri.replace(
      path: logoutPath,
    );

    final response = await client.get(uri);
    return response.statusCode == HttpStatus.ok;
  }

  // getUserInfo returns the user info of the current logged in user.
  static Future<UsersResponse> getUsers() async {
    final uri = apiUri.replace(
      path: usersPath,
    );

    final response = await client.get(uri);
    if (response.statusCode != HttpStatus.ok) {
      return Future.error(const HttpException("cannot get user details"));
    }

    return UsersResponse.fromJson(jsonDecode(response.body));
  }
}
