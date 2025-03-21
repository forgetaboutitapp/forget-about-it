import 'package:app/data/constants.dart';
import 'package:app/interop/get_url.dart';
import 'package:app/network/interfaces.dart';
import 'package:app/screens/bulk-edit/view.dart';
import 'package:app/screens/quiz/view.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hive_ce_flutter/adapters.dart';

import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:http/http.dart' as http;
import 'network/network.dart';
import 'screens/login/view.dart';
import 'screens/home/view.dart';
import 'screens/settings/view.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Hive.initFlutter();
  await Hive.openBox<dynamic>(localSettingsHiveBox);
  usePathUrlStrategySafe();

  runApp(
    ProviderScope(
      child: MainApp(),
    ),
  );
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
          remoteServer: _getRemoteServer(),
        ),
      ),
      GoRoute(
        path: BulkEditScreen.location,
        builder: (context, state) => BulkEditScreen(
          remoteServer: _getRemoteServer(),
        ),
      ),
      GoRoute(
        path: LoginScreen.location,
        builder: (context, state) => LoginScreen(
          client: http.Client(),
        ),
      ),
      GoRoute(
        path: QuizView.location,
        builder: (context, state) => QuizView(
          remoteServer: _getRemoteServer(),
          tags: state.uri.queryParametersAll,
        ),
      ),
      GoRoute(
        path: SettingsScreen.location,
        builder: (context, state) => SettingsScreen(
          filepicker: RealFilepicker(),
          remoteServer: _getRemoteServer(),
          curDarkMode: Hive.box(localSettingsHiveBox)
              .get(localSettingsHiveDarkTheme, defaultValue: false),
          switchDarkMode: (_) async {
            final v = Hive.box(localSettingsHiveBox)
                .get(localSettingsHiveDarkTheme, defaultValue: false);
            await Hive.box(localSettingsHiveBox)
                .put(localSettingsHiveDarkTheme, !v);
          },
        ),
      ),
    ],
  );

  static RemoteServer _getRemoteServer() {
    return RemoteServer(
      client: http.Client(),
      token: Hive.box(localSettingsHiveBox).get(localSettingsHiveLoginToken),
      remoteHost:
          Hive.box(localSettingsHiveBox).get(localSettingsHiveRemoteHost),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ValueListenableBuilder(
        valueListenable: Hive.box(localSettingsHiveBox).listenable(),
        builder: (context, settingsBox, w) {
          return MaterialApp.router(
            debugShowCheckedModeBanner: false,
            routerConfig: router,
            theme: settingsBox.get(localSettingsHiveDarkTheme) == true
                ? ThemeData.dark(useMaterial3: true)
                : ThemeData.light(useMaterial3: true),
          );
        });
  }
}

class RealFilepicker extends GenericFilepicker {
  @override
  Future<FilePickerResult?> pickFile(
          {required FileType type, required List<String> allowedExtensions}) =>
      FilePicker.platform
          .pickFiles(type: type, allowedExtensions: allowedExtensions);
}
