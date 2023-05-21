import 'package:flutter/material.dart';
import 'package:frontend/router.dart';
import 'package:go_router/go_router.dart';

class InvalidSession extends StatelessWidget {
  const InvalidSession({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    Future.delayed(
      Duration.zero,
      () {
        showDialog(
          barrierDismissible: false,
          context: context,
          builder: (context) {
            return AlertDialog(
              content: const Text("You have been logged out. "
                  "Please login again."),
              actions: [
                TextButton(
                  child: const Text("Okay!"),
                  onPressed: () => context.goNamed(Routes.login.name),
                )
              ],
            );
          },
        );
      },
    );
    return const SizedBox.shrink();
  }
}
