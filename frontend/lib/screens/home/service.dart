import '../../network/interfaces.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

import '../../fn/fn.dart';
import '../../protobufs-build/client_to_server.pb.dart' as client_to_server;
import 'model.dart' as model;

Future<Result<(IList<model.Tag>, bool)>> getAllTags(
    FetchDataWithToken fd) async {
  return (await fd.getAllTags(client_to_server.GetAllTags())).map((e) => (
        e.tags
            .map((t) => model.Tag(tag: t.tag, totalQuestions: t.totalQuestions))
            .toIList(),
        e.canRun
      ));
}
