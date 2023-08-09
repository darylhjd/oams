import 'package:flutter/material.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:url_launcher/url_launcher.dart';

// Shows the login screen.
class LoginScreen extends StatelessWidget {
  static const double _loginButtonMaxHeight = 200;
  static const String _microsoftLogoPath = "assets/microsoft_logo.png";
  static const double _logoHeight = 150;
  static const String _buttonText = "SSO with Microsoft";
  static const String _urlLaunchMode = "_self";

  final String _redirectUrl;

  const LoginScreen({Key? key, String? redirectUrl = ""})
      : _redirectUrl = redirectUrl ?? "",
        super(key: key);

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(
            maxHeight: _loginButtonMaxHeight,
          ),
          child: ElevatedButton(
            style: ButtonStyle(
              shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                const RoundedRectangleBorder(
                  borderRadius: BorderRadius.zero,
                ),
              ),
            ),
            onPressed: () async {
              final url = await APIClient.getLoginURL(_redirectUrl);
              launchUrl(Uri.parse(url), webOnlyWindowName: _urlLaunchMode);
            },
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Image.asset(_microsoftLogoPath, height: _logoHeight),
                const Text(_buttonText),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
