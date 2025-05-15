import '../../protobufs-build/client_to_server.pb.dart' as client_to_server;
import '../../protobufs-build/server_to_client.pb.dart';
import '../../screens/settings/models/remote_algorithm.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

import '../../fn/fn.dart';
import '../../network/interfaces.dart';
import 'models/remote_device.dart';
import 'models/remote_settings.dart';

Future<Result<RemoteSettings>> getRemoteSettings(
        FetchDataWithToken remoteServer) async =>
    (await remoteServer.getRemoteSettings(
      client_to_server.GetRemoteSettings(),
    ))
        .map(
      (GetRemoteSettings s) => RemoteSettings(
        remoteDevices: s.remoteDevices
            .map((e) => RemoteDevice(
                  title: e.title,
                  dateAdded: e.dateAdded.toDateTime(toLocal: true),
                  lastUsed: !e.hasLastUsed()
                      ? null
                      : e.lastUsed.toDateTime(toLocal: true),
                  loginId: e.loginId,
                ))
            .toIList(),
        defaultAlgorithm: s.defaultAlgorithm,
        remoteAlgorithms: s.algorithms
            .map(
              (a) => RemoteAlgorithm(
                algorithmID: a.algorithmId,
                authorName: a.authorName,
                license: a.license,
                remoteURL: a.remoteUrl,
                downloadURL: a.downloadUrl,
                version: a.version,
                algorithmName: a.algorithmName,
                timeAdded: a.dateAdded.toDateTime(toLocal: true).toString(),
              ),
            )
            .toIList(),
      ),
    );
