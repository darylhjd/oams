import 'package:flutter/material.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:url_launcher/url_launcher.dart';

// LoginScreen shows the login screen.
class LoginScreen extends StatelessWidget {
  static const double loginButtonMaxHeight = 200;
  static const String microsoftLogoPath = "assets/microsoft_logo.png";
  static const double logoHeight = 150;
  static const String buttonText = "SSO with Microsoft";
  static const String urlLaunchMode = "_self";

  final String redirectUrl;

  const LoginScreen({Key? key, this.redirectUrl = ""}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(
            maxHeight: loginButtonMaxHeight,
          ),
          child: ElevatedButton(
            style: ButtonStyle(
              shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                const RoundedRectangleBorder(
                  borderRadius: BorderRadius.zero,
                ),
              ),
            ),
            onPressed: loginAction,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Image.asset(microsoftLogoPath, height: logoHeight),
                const Text(buttonText),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Future<void> loginAction() async {
    final url = await APIClient.getLoginURL(redirectUrl);
    if (!await launchUrl(Uri.parse(url), webOnlyWindowName: urlLaunchMode)) {
      // TODO: Add handling here.
    }
  }
}
