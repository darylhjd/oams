import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/about_screen.dart';
import 'package:frontend/screens/home_screen.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:frontend/screens/profile_screen.dart';
import 'package:go_router/go_router.dart';

typedef RouteInfo = ({String name, String path});

// Routes stores the different routes in the frontend.
class Routes {
  static const RouteInfo index = (name: "index", path: "/");
  static const RouteInfo about = (name: "about", path: "about");
  static const RouteInfo login = (name: "login", path: "login");
  static const RouteInfo profile = (name: "profile", path: "profile");
}

final routerProvider = Provider<GoRouter>(
  (ref) {
    var isLoggedIn = ref.watch(sessionProvider);
    return GoRouter(
      routes: [
        GoRoute(
          name: Routes.index.name,
          path: Routes.index.path,
          pageBuilder: (context, state) => NoTransitionPage(
            key: state.pageKey,
            child: const HomeScreen(),
          ),
          routes: [
            GoRoute(
              name: Routes.about.name,
              path: Routes.about.path,
              pageBuilder: (context, state) => NoTransitionPage(
                key: state.pageKey,
                child: const AboutScreen(),
              ),
            ),
            GoRoute(
              name: Routes.login.name,
              path: Routes.login.path,
              pageBuilder: (context, state) {
                return NoTransitionPage(
                  key: state.pageKey,
                  child: LoginScreen(
                      returnTo: state.queryParameters["return_to"] ?? ""),
                );
              },
              redirect: (context, state) =>
                  isLoggedIn ? Routes.index.name : null,
            ),
            GoRoute(
              name: Routes.profile.name,
              path: Routes.profile.path,
              pageBuilder: (context, state) {
                return NoTransitionPage(
                  key: state.pageKey,
                  child: const ProfileScreen(),
                );
              },
              redirect: (context, state) =>
                  isLoggedIn ? null : Routes.index.path + Routes.login.path,
            ),
          ],
        ),
      ],
    );
  },
);
