import '../../fn/fn.dart';
import '../../network/interfaces.dart';
import '../../protobufs-build/client_to_server.pb.dart' as client_to_server;
import '../../protobufs-build/server_to_client.pb.dart' as server_to_client;
import '../../screens/bulk-edit/model.dart';
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
        {required FetchDataWithToken remoteServer}) async =>
    (await remoteServer.getAllQuestions(client_to_server.GetAllQuestions()))
        .flatMap(
      (a) => Result.safe(
        () => '$_originalText\n${unparse(
          a.flashcards
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
        )}',
      ),
    );

Future<Result<server_to_client.PostAllQuestions>> postAllQuestions({
  required FetchDataWithToken remoteServer,
  required IList<Flashcard> flashcards,
}) async {
  return await remoteServer.postAllQuestions(
    client_to_server.PostAllQuestions(
      flashcards: flashcards
          .map(
            (e) => client_to_server.Flashcard(
              id: e.id,
              question: e.question,
              answer: e.answer,
              memoHint: e.memoHint,
              explanation: e.explanation,
              tags: e.tags,
            ),
          )
          .toList(),
    ),
  );
}
