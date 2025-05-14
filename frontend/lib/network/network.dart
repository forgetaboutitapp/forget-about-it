import 'dart:core';
import 'dart:developer' as developer;

import 'package:app/data/errors.dart';
import 'package:app/network/interfaces.dart';
import 'package:app/protobufs-build/client_to_server.pb.dart'
    as client_to_server;
import 'package:app/protobufs-build/server_to_client.pb.dart'
    as server_to_client;

import 'package:http/http.dart' as http;
import 'package:http/http.dart' as client;
import '../fn/fn.dart';
part 'network.g.dart';

class RemoteServer with _$RemoteServer implements FetchDataWithToken {
  @override
  final String remoteHost;
  final String token;
  final http.Client client;

  RemoteServer({
    required this.remoteHost,
    required this.token,
    required this.client,
  });

  @override
  String getRemoteHost() => remoteHost;
}

class RemoteServerWithoutToken implements FetchDataWithoutToken {
  final String remoteHost;
  final http.Client client;

  RemoteServerWithoutToken({
    required this.remoteHost,
    required this.client,
  });

  @override
  String getRemoteHost() => remoteHost;

  @override
  Future<Result<server_to_client.GetToken>> getToken(
      client_to_server.InsecureMessage msg) async {
    return (await getData(
      remoteHost,
      client_to_server.Message(insecureMessage: msg).writeToBuffer(),
    ))
        .flatMap(
      (v) => Result.safe(() => server_to_client.Message.fromBuffer(v)).flatMap(
        (msg) {
          if (msg.hasOkMessage() && msg.okMessage.hasGetToken()) {
            return Ok(msg.okMessage.getToken);
          } else if (msg.hasErrorMessage()) {
            return Err(
              MyException(
                shouldLogout: msg.errorMessage.shouldLogOut,
                inner: Exception(msg.errorMessage.error),
              ),
            );
          } else {
            developer.log('msg is not ok or error, is ${msg.writeToJson()}');
            return Err(
              MyException(
                shouldLogout: false,
                inner: Exception('Internal Error'),
              ),
            );
          }
        },
      ),
    );
  }
}
