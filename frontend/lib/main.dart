import 'package:app/state/login.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

import 'screens/login/view.dart';
import 'screens/home/view.dart';
import 'screens/settings/view.dart';

late final SharedPreferences sharedPreferences;

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  sharedPreferences = await SharedPreferences.getInstance();
  runApp(
    ProviderScope(
      child: MainApp(),
    ),
  );
}

class MainApp extends HookConsumerWidget {
  const MainApp({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = GoRouter(
      redirect: (BuildContext context, GoRouterState state) async {
        final logginState = ref.watch(loginProvider);
        switch (logginState) {
          case NoLogin():
            return LoginScreen.location;
          case LoggedIn():
            return state.fullPath;
        }
      },
      routes: [
        GoRoute(
          path: HomeScreen.location,
          builder: (context, state) => HomeScreen(),
        ),
        GoRoute(
          path: LoginScreen.location,
          builder: (context, state) => LoginScreen(client: http.Client()),
        ),
        GoRoute(
          path: SettingsScreen.location,
          builder: (context, state) => SettingsScreen(client: http.Client()),
        ),
      ],
    );

    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      routerConfig: router,
    );
  }
}
