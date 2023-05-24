import 'package:flutter/material.dart';
import 'package:frontend/screens/screen_template.dart';

// HomeScreen shows the home screen.
class HomeScreen extends StatelessWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      ListView(
        padding: const EdgeInsets.all(20),
        children: [
          Container(
            decoration: BoxDecoration(
              border: Border.all(color: Colors.black),
            ),
            height: 300,
            child:
                const Text("Welcome to Online Attendance Management System!"),
          ),
        ],
      ),
    );
  }
}
