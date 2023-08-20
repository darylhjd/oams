import 'package:flutter/material.dart';
import 'package:frontend/screens/screen_template.dart';

class UserScreen extends StatelessWidget {
  final String _userId;

  const UserScreen(this._userId, {super.key});

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      Text(_userId),
    );
  }
}
