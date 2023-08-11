import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/home_screen_guest.dart';
import 'package:frontend/screens/home_screen_logged.dart';

// Shows the home screen.
class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ref.read(sessionProvider)
        ? const HomeScreenLoggedIn()
        : const HomeScreenGuest();
  }
}
