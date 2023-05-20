import 'package:flutter/material.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:url_launcher/url_launcher.dart';

import '../api/client.dart';

class LoginScreen extends StatelessWidget {
  final Map<String, String> queryParams;

  const LoginScreen({Key? key, required this.queryParams}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(
            maxHeight: 200,
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
              final redirectUrl = await APIClient.getLoginURL(queryParams);
              if (!await launchUrl(Uri.parse(redirectUrl),
                  webOnlyWindowName: "_self")) {
                // TODO: Add handling here.
              }
            },
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Image.asset("assets/microsoft_logo.png", height: 150),
                const Text("SSO with Microsoft"),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
