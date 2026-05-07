import 'dart:developer';
import 'dart:io';

import 'package:android_alarm_manager_plus/android_alarm_manager_plus.dart';
import 'package:toastification/toastification.dart';
import '../../data/constants.dart';
import '../../interop/get_url.dart';
import '../../screens/bulk-edit/view.dart';
import '../../screens/quiz/view.dart';
import '../../service.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hive_ce_flutter/adapters.dart';

import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:http/http.dart' as http;
import 'interop/notification_service_stub.dart';
import 'screens/login/view.dart';
import 'screens/home/view.dart';
import 'screens/settings/view.dart';
import 'screens/stats/view.dart';

void main() async {
  log('starting');
  WidgetsFlutterBinding.ensureInitialized();
  await Hive.initFlutter();
  await Hive.openBox<dynamic>(localSettingsHiveBox);
  usePathUrlStrategySafe();
  if (!kIsWeb && Platform.isAndroid) {
    await NotificationService().initNotification();
    await AndroidAlarmManager.initialize();
    await AndroidAlarmManager.periodic(
      const Duration(minutes: 5),
      0,
      printHello,
      wakeup: true,
      rescheduleOnReboot: true,
      allowWhileIdle: true,
    );
  }

  runApp(ProviderScope(child: MainApp()));
}

class MainApp extends HookConsumerWidget {
  MainApp({super.key});
  final router = GoRouter(
    refreshListenable: Hive.box(localSettingsHiveBox).listenable(),
    redirect: (BuildContext context, GoRouterState state) async {
      final settingsBox = Hive.box(localSettingsHiveBox);
      if (settingsBox.get(localSettingsHiveLoginToken) == null) {
        return LoginScreen.location;
      } else if (state.fullPath != LoginScreen.location) {
        return state.uri.toString();
      } else {
        return HomeScreen.location;
      }
    },
    routes: [
      GoRoute(
        path: HomeScreen.location,
        builder: (context, state) => HomeScreen(
            logOut: () {},
            remoteServer: _getRemoteServer().remoteHost,
            token: _getRemoteServer().token),
      ),
      GoRoute(
        path: BulkEditScreen.location,
        builder: (context, state) => BulkEditScreen(
          remoteServer: _getRemoteServer().remoteHost,
          token: _getRemoteServer().token,
          logOut: () {},
        ),
      ),
      GoRoute(
        path: LoginScreen.location,
        builder: (context, state) => LoginScreen(),
      ),
      GoRoute(
        path: Stats.location,
        builder: (context, state) => Stats(
            remoteServer: _getRemoteServer().remoteHost,
            token: _getRemoteServer().token,
            logOut: () {}),
      ),
      GoRoute(
        path: QuizView.location,
        builder: (context, state) => QuizView(
          remoteServer: _getRemoteServer().remoteHost,
          token: _getRemoteServer().token,
          logOut: () {},
          tags: state.uri.queryParametersAll,
          isDarkMode: Hive.box(
            localSettingsHiveBox,
          ).get(localSettingsHiveDarkTheme, defaultValue: false),
        ),
      ),
      GoRoute(
        path: SettingsScreen.location,
        builder: (context, state) => SettingsScreen(
          token: _getRemoteServer().token,
          remoteServer: _getRemoteServer().remoteHost,
          logout: () {},
          curDarkMode: Hive.box(
            localSettingsHiveBox,
          ).get(localSettingsHiveDarkTheme, defaultValue: false),
          switchDarkMode: (_) async {
            final v = Hive.box(
              localSettingsHiveBox,
            ).get(localSettingsHiveDarkTheme, defaultValue: false);
            await Hive.box(
              localSettingsHiveBox,
            ).put(localSettingsHiveDarkTheme, !v);
          },
        ),
      ),
    ],
  );

  static RemoteServer _getRemoteServer() {
    return RemoteServer(
      client: http.Client(),
      token: Hive.box(localSettingsHiveBox).get(localSettingsHiveLoginToken),
      remoteHost: Hive.box(
        localSettingsHiveBox,
      ).get(localSettingsHiveRemoteHost),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ValueListenableBuilder(
      valueListenable: Hive.box(localSettingsHiveBox).listenable(),
      builder: (context, settingsBox, w) {
        return ToastificationWrapper(
          child: MaterialApp.router(
            debugShowCheckedModeBanner: false,
            routerConfig: router,
            theme: settingsBox.get(localSettingsHiveDarkTheme) == true
                ? ThemeData.dark(useMaterial3: true)
                : ThemeData.light(useMaterial3: true),
          ),
        );
      },
    );
  }
}

class RemoteServer {
  final http.Client client;
  final String remoteHost;
  final String token;
  const RemoteServer(
      {required this.remoteHost, required this.token, required this.client});
}
