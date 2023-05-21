import 'package:flutter/material.dart';
import 'package:frontend/screens/screen_template.dart';

// HomeScreen shows the home screen.
class HomeScreen extends StatelessWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return const ScreenTemplate(
      Center(
        child: Text("Welcome to Online Attendance Management System!"),
      ),
    );
  }
}
