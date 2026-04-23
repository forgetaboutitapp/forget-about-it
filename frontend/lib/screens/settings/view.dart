import 'dart:async';
import 'dart:developer' as dev;

import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';
import 'package:grpc/grpc_web.dart';

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
  final String remoteServer;
  final String token;
  final Function() logout;
  final FutureOr<void> Function(bool) switchDarkMode;
  const SettingsScreen({
    super.key,
    required this.remoteServer,
    required this.token,
    required this.logout,
    required this.curDarkMode,
    required this.switchDarkMode,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final localCurDarkMode = useState(curDarkMode);
    final stillWaitingForRemoteServer = useState(true);

    final remoteSettings = ref.watch(remoteSettingsProvider);
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
                              final client = await ForgetAboutItServiceClient(
                                      GrpcWebClientChannel.xhr(
                                          Uri.parse(remoteServer)))
                                  .getRemoteSettings(
                                      GetRemoteSettingsRequest(token: token));
                              if (client.hasError()) {
                                if (context.mounted) {
                                  showError(context, e.toString());
                                }
                              } else {
                                if (context.mounted) {
                                  refresh(stillWaitingForRemoteServer, ref,
                                      context);
                                }
                              }
                            },
                          ),
                          trailing: IconButton(
                            icon: Icon(Icons.delete),
                            onPressed: remoteSettings.defaultAlgorithm ==
                                    e.algorithmID
                                ? null
                                : () async {
                                    final client =
                                        await ForgetAboutItServiceClient(
                                                GrpcWebClientChannel.xhr(
                                                    Uri.parse(remoteServer)))
                                            .removeAlgorithm(
                                      RemoveAlgorithmRequest(
                                          token: token,
                                          algorithmId: e.algorithmName),
                                    );

                                    if (client.hasError()) {
                                      if (context.mounted) {
                                        showError(context, client.error.error);
                                      }
                                    } else {
                                      if (context.mounted) {
                                        refresh(
                                          stillWaitingForRemoteServer,
                                          ref,
                                          context,
                                        );
                                      }
                                    }
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
                            token: token,
                            algos: remoteSettings?.remoteAlgorithms
                                ?.toSet()
                                .toISet(),
                            remoteServer: remoteServer,
                          ),
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
                                              await ForgetAboutItServiceClient(
                                                      GrpcWebClientChannel.xhr(
                                                          Uri.parse(
                                                              remoteServer)))
                                                  .removeLogin(
                                            RemoveLoginRequest(
                                                token: token,
                                                loginId: e.loginId),
                                          );

                                          if (res.hasOk()) {
                                            if (context.mounted) {
                                              refresh(
                                                  stillWaitingForRemoteServer,
                                                  ref,
                                                  context);
                                            }
                                          } else if (res.hasError()) {
                                            if (context.mounted) {
                                              showError(
                                                context,
                                                e.toString(),
                                              );
                                            }
                                          }
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
                                builder: (context) => QRDialog(
                                    remoteHost: remoteServer,
                                    token: token,
                                    logOut: logout),
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
          .read(remoteSettingsProvider.notifier)
          .getData(remoteServer, token, logout);
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
