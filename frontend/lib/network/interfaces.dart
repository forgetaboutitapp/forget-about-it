import 'dart:async';

import 'package:file_picker/file_picker.dart';
import '../fn/fn.dart';
import '../protobufs-build/client_to_server.pb.dart' as client_to_server;
import '../protobufs-build/server_to_client.pb.dart' as server_to_client;

part 'fetch_data_with_token.g.dart';

abstract class FetchDataWithToken with _$FetchDataWithToken {
  String getRemoteHost();
}

abstract class FetchDataWithoutToken {
  String getRemoteHost();
  Future<Result<server_to_client.GetToken>> getToken(
      client_to_server.InsecureMessage msg);
}

abstract class GenericFilepicker {
  Future<FilePickerResult?> pickFile({
    required FileType type,
    required List<String> allowedExtensions,
  });
}
