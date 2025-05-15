import 'dart:async';
import 'dart:developer' as dev;

import '../../network/interfaces.dart';
import '../../protobufs-build/client_to_server.pb.dart' as client_to_server;
import '../../screens/general-display/show_error.dart';
import '../../screens/settings/add_algorithm.dart';
import '../../screens/settings/models/model.dart';
import '../../screens/settings/qr_dialog.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:flutter_settings_ui/flutter_settings_ui.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:permission_handler/permission_handler.dart';

class SettingsScreen extends HookConsumerWidget {
  static const location = '/settings';
  final bool curDarkMode;
  final FetchDataWithToken remoteServer;
  final GenericFilepicker filepicker;
  final FutureOr<void> Function(bool) switchDarkMode;
  const SettingsScreen({
    super.key,
    required this.curDarkMode,
    required this.remoteServer,
    required this.switchDarkMode,
    required this.filepicker,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final localCurDarkMode = useState(curDarkMode);
    final stillWaitingForRemoteServer = useState(true);

    final remoteSettings = ref.watch(remoteSettingsNotifierProvider);
    final ValueNotifier<bool?> hasNotificationPermission = useState(null);
    final ValueNotifier<bool?> permDenied = useState(null);

    useEffect(() {
      () async {
        permDenied.value = await Permission.notification.isPermanentlyDenied;
        hasNotificationPermission.value =
            await Permission.notification.isGranted;
      }();
      return null;
    }, []);
    useEffect(() {
      refresh(stillWaitingForRemoteServer, ref, context);
      return null;
    }, []);
    return Scaffold(
      appBar: AppBar(
        title: Text('Settings'),
      ),
      body: stillWaitingForRemoteServer.value
          ? Center(child: CircularProgressIndicator())
          : SettingsList(
              sections: [
                SettingsSection(
                  title: Text('Local Setting'),
                  tiles: [
                    if (!kIsWeb)
                      SettingsTile.switchTile(
                        initialValue: hasNotificationPermission.value == true,
                        onToggle: permDenied.value == null ||
                                permDenied.value == true
                            ? null
                            : (v) async {
                                if (hasNotificationPermission.value != true) {
                                  await Permission.notification.request();
                                  permDenied.value = await Permission
                                      .notification.isPermanentlyDenied;
                                  hasNotificationPermission.value =
                                      await Permission.notification.isGranted;
                                }
                              },
                        title: Text('Enable Notifications Permissions'),
                      ),
                    SettingsTile.switchTile(
                      initialValue: localCurDarkMode.value,
                      onToggle: (v) {
                        localCurDarkMode.value = v;
                        switchDarkMode(v);
                      },
                      title: Text('Dark Theme'),
                    ),
                  ],
                ),
                SettingsSection(
                  title: Text('General Setting'),
                  tiles: [
                    ...?remoteSettings?.remoteAlgorithms?.map(
                      (e) {
                        dev.log(
                            'algoId: ${e.algorithmID}, remoteSettings.defaultAlgorithm: ${remoteSettings.defaultAlgorithm}');
                        return SettingsTile(
                          title: Text(e.algorithmName),
                          description: Text(
                              'Author: ${e.authorName}\nLicense: ${e.license}\nVersion: ${e.version}\nAdded: ${e.timeAdded}'),
                          leading: Radio(
                            value: e.algorithmID,
                            groupValue: remoteSettings.defaultAlgorithm ?? 0,
                            onChanged: (int? v) async {
                              stillWaitingForRemoteServer.value = true;

                              (await remoteServer.setDefaultAlgorithm(
                                client_to_server.SetDefaultAlgorithm(
                                    algorithmId: e.algorithmID),
                              ))
                                  .match(
                                      onOk: (_) {
                                        if (context.mounted) {
                                          refresh(stillWaitingForRemoteServer,
                                              ref, context);
                                        }
                                      },
                                      onErr: (e) =>
                                          showError(context, e.toString()));
                            },
                          ),
                          trailing: IconButton(
                            icon: Icon(Icons.delete),
                            onPressed: remoteSettings.defaultAlgorithm ==
                                    e.algorithmID
                                ? null
                                : () async {
                                    (await remoteServer.removeAlgorithm(
                                      client_to_server.RemoveAlgorithm(
                                          algorithmId: e.algorithmName),
                                    ))
                                        .match(
                                      onErr: (err) =>
                                          showError(context, err.toString()),
                                      onOk: (_) {
                                        if (context.mounted) {
                                          refresh(
                                            stillWaitingForRemoteServer,
                                            ref,
                                            context,
                                          );
                                        }
                                      },
                                    );
                                  },
                          ),
                        );
                      },
                    ),
                    SettingsTile(
                      title: Text('Add Algorithm'),
                      trailing: Icon(Icons.add),
                      onPressed: (context) async {
                        final needsRefresh = await showDialog(
                          context: context,
                          builder: (context) => AddAlgorithm(
                              algos: remoteSettings?.remoteAlgorithms
                                  ?.toSet()
                                  .toISet(),
                              remoteServer: remoteServer,
                              filepicker: filepicker),
                        );
                        if (context.mounted && needsRefresh == true) {
                          refresh(stillWaitingForRemoteServer, ref, context);
                        }
                      },
                    )
                  ],
                ),
                remoteSettings == null
                    ? SettingsSection(
                        title: Text('Other Devices'),
                        tiles: [
                          SettingsTile(
                            title: Text('Cannot reach remote server... '),
                          ),
                        ],
                      )
                    : SettingsSection(
                        title: Text('Other Devices'),
                        tiles: [
                          ...?remoteSettings.remoteDevices?.indexed.map(
                            (tuple) {
                              final e = tuple.$2;
                              final i = tuple.$1;
                              final lastUsed = e.lastUsed;
                              return SettingsTile(
                                title: Text(e.title),
                                value: Column(
                                  children: [
                                    Padding(
                                      padding:
                                          const EdgeInsets.fromLTRB(8, 8, 8, 8),
                                      child: Text(
                                          'Date Added: ${DateFormat.yMd().format(e.dateAdded)}'),
                                    ),
                                    if (lastUsed != null)
                                      Text(
                                          'Last Used: ${DateFormat.yMd().format(lastUsed)}'),
                                  ],
                                ),
                                trailing: IconButton(
                                  icon: Icon(Icons.delete),
                                  onPressed: i == 0
                                      ? null
                                      : () async {
                                          final res =
                                              await remoteServer.removeLogin(
                                            client_to_server.RemoveLogin(
                                                loginId: e.loginId),
                                          );
                                          res.match(
                                            onOk: (o) {
                                              if (context.mounted) {
                                                refresh(
                                                    stillWaitingForRemoteServer,
                                                    ref,
                                                    context);
                                              }
                                            },
                                            onErr: (e) {
                                              showError(
                                                context,
                                                e.toString(),
                                              );
                                            },
                                          );
                                        },
                                ),
                              );
                            },
                          ),
                          SettingsTile(
                            title: Text('Connect Other Devices'),
                            trailing: Icon(Icons.add),
                            onPressed: (context) async {
                              await showDialog(
                                context: context,
                                builder: (context) =>
                                    QRDialog(remoteServer: remoteServer),
                              );
                              if (context.mounted) {
                                refresh(
                                    stillWaitingForRemoteServer, ref, context);
                              }
                            },
                          ),
                        ],
                      ),
                SettingsSection(
                  title: Text('Licensess'),
                  tiles: [
                    SettingsTile(
                      title: Text('Show Licenses'),
                      onPressed: (context) => showLicensePage(
                          applicationLegalese: 'AGPL-V3',
                          context: context,
                          applicationName: 'Forget About It'),
                    )
                  ],
                ),
              ],
            ),
    );
  }

  void refresh(ValueNotifier<bool> stillWaitingForRemoteServer, WidgetRef ref,
      BuildContext context) {
    () async {
      if (stillWaitingForRemoteServer.value != true) {
        stillWaitingForRemoteServer.value = true;
      }
      final p = await ref
          .read(remoteSettingsNotifierProvider.notifier)
          .getData(remoteServer);
      stillWaitingForRemoteServer.value = false;
      if (p != null) {
        stillWaitingForRemoteServer.value = false;
        if (context.mounted) {
          ScaffoldMessenger.of(context)
              .showSnackBar(SnackBar(content: Text(p.toString())));
        }
      }
    }();
  }
}
