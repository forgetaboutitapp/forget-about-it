import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';
import 'package:forget_about_it/protobufs-build/client_server/v1/server_to_client.pb.dart';
import 'package:grpc/grpc_web.dart';

import '../../fn/fn.dart';
import 'model.dart' as model;

Future<Result<(IList<model.Tag>, bool)>> getAllTags(
    String token, String remoteHost, Function logOut) async {
  final client = await ForgetAboutItServiceClient(
          GrpcWebClientChannel.xhr(Uri.parse(remoteHost)))
      .getAllTags(GetAllTagsRequest(token: token));

  return switch (client) {
    GetAllTags(:var tags, :var canRun) => Ok((
        tags
            .map((t) => model.Tag(tag: t.tag, totalQuestions: t.totalQuestions))
            .toIList(),
        canRun
      )),
    ErrorMessage(:var error, :var shouldLogOut) =>
      shouldLogOut ? logOut() : Err(Exception(error)),
    _ => Err(Exception('Server Error')),
  };
}
