import 'dart:developer' as developer;

import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';

import '../../screens/login/submit_type.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:hive_ce/hive.dart';
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
  Uri remoteUri,
  SubmitType submitType,
) async {
  final response =
      await ForgetAboutItServiceClient(createGrpcChannel(remoteUri))
          .getToken(switch (submitType) {
    TwelveWords(:final twelveWords) => GetTokenRequest(
        twelveWords: twelveWords.toList(),
      ),
    Token(:final token) => GetTokenRequest(token: token),
  });

  if (response.hasOk()) {
    final token = response.ok.token;
    Hive.box(localSettingsHiveBox).put(localSettingsHiveLoginToken, token);
    Hive.box(localSettingsHiveBox).put(localSettingsHiveRemoteHost,
        '${remoteUri.scheme}://${remoteUri.host}:${remoteUri.port}');
    return Ok(());
  }

  if (response.hasError()) {
    developer.log('getToken error: ${response.error}');
    return Err(Exception(response.error.error));
  }

  developer.log('getToken returned no result: $response');
  return Err(Exception('Server Error'));
}
