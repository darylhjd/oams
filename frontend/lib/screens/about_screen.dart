import 'package:flutter/material.dart';
import 'package:frontend/screens/screen_template.dart';

class AboutScreen extends StatelessWidget {
  const AboutScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return const ScreenTemplate(
      Center(
        child: Text("This is the about screen"),
      ),
    );
  }
}
