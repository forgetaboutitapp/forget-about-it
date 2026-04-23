import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';
import 'package:grpc/grpc_web.dart';

import '../../screens/settings/models/remote_algorithm.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

import '../../fn/fn.dart';
import 'models/remote_device.dart';
import 'models/remote_settings.dart';

Future<Result<RemoteSettings>> getRemoteSettings(
    String remoteHost, String token, Function logOut) async {
  final client = await ForgetAboutItServiceClient(
          GrpcWebClientChannel.xhr(Uri.parse(remoteHost)))
      .getRemoteSettings(GetRemoteSettingsRequest(
    token: token,
  ));
  if (client.hasError()) {
    return Err(Exception(client.error));
  }

  return Ok(RemoteSettings(
    remoteDevices: client.ok.remoteDevices
        .map((e) => RemoteDevice(
              title: e.title,
              dateAdded: e.dateAdded.toDateTime(toLocal: true),
              lastUsed: !e.hasLastUsed()
                  ? null
                  : e.lastUsed.toDateTime(toLocal: true),
              loginId: e.loginId,
            ))
        .toIList(),
    defaultAlgorithm: client.ok.defaultAlgorithm,
    remoteAlgorithms: client.ok.algorithms
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
  ));
}
