import 'dart:developer' as developer;

import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';
import 'package:forget_about_it/protobufs-build/client_server/v1/server_to_client.pb.dart';

import '../../screens/login/submit_type.dart';
import 'package:flutter/foundation.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:hive_ce/hive.dart';
import 'package:http/http.dart' as http;
import '../data/constants.dart';
import '../fn/fn.dart';
import '../interop/grpc_channel.dart';

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
  final client = await ForgetAboutItServiceClient(createGrpcChannel(remoteUri))
      .getToken(switch (submitType) {
    TwelveWords(:final twelveWords) => GetTokenRequest(
        twelveWords: twelveWords.toList(),
      ),
    Token(:final token) => GetTokenRequest(token: token),
  });
  final res = switch (client) {
    GetToken(:final token) => (token, null),
    Err(:final value) => (null, value),
    _ => throw Exception('Unknown Error'),
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
