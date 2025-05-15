import 'dart:developer' as developer;

import '../../network/network.dart';
import '../../protobufs-build/client_to_server.pb.dart';
import '../../screens/login/submit_type.dart';
import 'package:flutter/foundation.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:hive_ce/hive.dart';
import 'package:http/http.dart' as http;
import '../data/constants.dart';
import '../fn/fn.dart';

part 'login.freezed.dart';

@freezed
sealed class LoginReturn with _$LoginReturn {
  const factory LoginReturn.noLogin() = _LoginReturnNoLogin;
  const factory LoginReturn.loggedIn({required String token}) =
      _LoginReturnLoggedIn;
}

Future<Result<()>> update(
  http.Client client,
  Uri remoteUri,
  SubmitType submitType,
) async {
  InsecureMessage d = switch (submitType) {
    TwelveWords(:final twelveWords) => InsecureMessage(
        getToken: GetToken(
          twelveWords: twelveWords.toList(),
        ),
      ),
    Token(:final token) => InsecureMessage(
        getToken: GetToken(token: token),
      ),
  };
  final remoteServer = RemoteServerWithoutToken(
    remoteHost: remoteUri.toString(),
    client: client,
  );

  final token = await remoteServer.getToken(d);
  final res = switch (token) {
    Ok(:final value) => (value.token, null),
    Err(:final value) => (null, value),
  };
  if (res.$1 != null) {
    final token = res.$1;
    Hive.box(localSettingsHiveBox).put(localSettingsHiveLoginToken, token);
    Hive.box(localSettingsHiveBox).put(localSettingsHiveRemoteHost,
        '${remoteUri.scheme}://${remoteUri.host}:${remoteUri.port}');
    return Ok(());
  }
  developer.log('ret: $res');
  return Err(res.$2!);
}
