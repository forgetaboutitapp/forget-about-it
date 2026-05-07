import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';

import '../../fn/fn.dart';
import '../../interop/grpc_channel.dart';
import '../../protobufs-build/client_server/v1/server_to_client.pb.dart'
    as server_to_client;
import '../../screens/bulk-edit/parse.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

import 'model.dart' as model;

const String _originalText =
    '''% Enter flashcards in CSV format, deliminated by the | char.
% The questions will be formatted using markdown, and the \\ escapes chars.
% Comments start with a `%`, \\n denotes a newline, and tables can be written with \\|
% For example, a question can be written as:
% | _Question 1_ | *Answer 1*\\n*Answer 2* | Explanation 1 | Way to remember 1| Tag-a Tag-b Tag-c|
% In the above question, Question 1 will be italic, Answer 1 and Answer 2 will be bold on two seperate lines.
% Once you submit the question, the server will give it an ID, so the next time you see it, it will look like:
% 12345| _Question 1_ | *Answer 1*\\n*Answer 2*  | Explanation 1 | Way to remember 1| Tag-a Tag-b Tag-c|
% Questions and IDs must be unique per user\n''';

Future<Result<String>> getAllQuestions(
    {required String token,
    required String remoteHost,
    required Function() logOut}) async {
  final client = await ForgetAboutItServiceClient(
          createGrpcChannel(Uri.parse(remoteHost)))
      .getAllQuestions(GetAllQuestionsRequest(token: token));
  if (client.hasError()) {
    if (client.error.shouldLogOut) {
      logOut();
    }
    return Err(Exception(client.error.error));
  }

  if (!client.hasOk()) {
    return Err(Exception('Server Error'));
  }

  return Result.safe(() => '$_originalText\n${unparse(
        client.ok.flashcards
            .map(
              (e) => model.Flashcard(
                id: e.id,
                question: e.question,
                answer: e.answer,
                explanation: e.explanation,
                memoHint: e.memoHint,
                tags: e.tags.toIList(),
              ),
            )
            .toIList(),
      )}');
}

Future<Result<void>> postAllQuestions({
  required String remoteHost,
  required String token,
  required Function() logOut,
  required IList<model.Flashcard> flashcards,
}) async {
  final client = await ForgetAboutItServiceClient(
          createGrpcChannel(Uri.parse(remoteHost)))
      .postAllQuestions(PostAllQuestionsRequest(
          token: token,
          flashcards: flashcards
              .map(
                (e) => server_to_client.Flashcard(
                  id: e.id,
                  question: e.question,
                  answer: e.answer,
                  memoHint: e.memoHint,
                  explanation: e.explanation,
                  tags: e.tags,
                ),
              )
              .toList()));
  if (client.hasError()) {
    if (client.error.shouldLogOut) {
      logOut();
    }
    return Err(Exception(client.error.error));
  }

  if (!client.hasOk()) {
    return Err(Exception('Server Error'));
  }

  return Ok(null);
}
