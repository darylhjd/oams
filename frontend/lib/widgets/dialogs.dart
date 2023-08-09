import 'package:flutter/material.dart';
import 'package:frontend/router.dart';
import 'package:go_router/go_router.dart';

// Shows a dialog on the screen to inform a logged session user that the session
// has expired.
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
              content: const Text(
                "Your session has been invalidated. Please login again.",
              ),
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
