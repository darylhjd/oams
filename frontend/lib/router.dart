import 'package:frontend/screens/about_screen.dart';
import 'package:frontend/screens/home_screen.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:go_router/go_router.dart';

typedef RouteInfo = ({String name, String path});

// Routes stores the different routes in the frontend.
class Routes {
  static const RouteInfo index = (name: "index", path: "/");
  static const RouteInfo about = (name: "about", path: "about");
  static const RouteInfo login = (name: "login", path: "login");
  static const RouteInfo profile = (name: "profile", path: "profile");
}

// router helps to provide proper URL handling for the frontend.
final router = GoRouter(
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
        )
      ],
    ),
  ],
);
