import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';

import '../../fn/fn.dart';
import '../../interop/grpc_channel.dart';
import 'model.dart' as model;

Future<Result<(IList<model.Tag>, bool)>> getAllTags(
    String token, String remoteHost, Function logOut) async {
  final client = await ForgetAboutItServiceClient(
          createGrpcChannel(Uri.parse(remoteHost)))
      .getAllTags(GetAllTagsRequest(token: token));

  if (client.hasError()) {
    if (client.error.shouldLogOut) {
      logOut();
    }
    return Err(Exception(client.error.error));
  }

  if (!client.hasOk()) {
    return Err(Exception('Server Error'));
  }

  return Ok((
    client.ok.tags
        .map((t) => model.Tag(tag: t.tag, totalQuestions: t.totalQuestions))
        .toIList(),
    client.ok.canRun,
  ));
}
