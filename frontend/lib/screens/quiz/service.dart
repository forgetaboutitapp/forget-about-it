import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';
import 'package:forget_about_it/protobufs-build/client_server/v1/server_to_client.pb.dart';

import '../../interop/grpc_channel.dart';
import '../../screens/quiz/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

import '../../fn/fn.dart';
import '../../protobufs-build/client_server/v1/server_to_client.pbenum.dart'
    as server_to_client_enums;

Future<Result<QuizQuestionStateData>> getNextQuestion(
    String remoteHost,
    String token,
    Function logOut,
    ISet<String> tagsQuery,
    bool getNewQuestion) async {
  final client = await ForgetAboutItServiceClient(
          createGrpcChannel(Uri.parse(remoteHost)))
      .getNextQuestion(GetNextQuestionRequest(
          token: token,
          tags: tagsQuery.toList(),
          getNewQuestion: getNewQuestion));
  return switch (client) {
    GetNextQuestion(
      :final flashcard,
      :final typeOfQuestion,
      :final newQuestions,
      :final dueQuestions,
      :final nonDueQuestions
    ) =>
      QuizQuestionStateData(
        id: flashcard.id,
        question: flashcard.question,
        answer: flashcard.answer,
        questionType: switch (typeOfQuestion) {
          server_to_client_enums
                .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_NEW =>
            QuestionType.newQuestion,
          server_to_client_enums
                .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_DUE =>
            QuestionType.dueQuestion,
          server_to_client_enums
                .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_NON_DUE =>
            QuestionType.nonDueQuestion,
          server_to_client_enums
                .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_UNSPECIFIED =>
            throw UnimplementedError(),
          server_to_client_enums.GetNextQuestion_TypeOfQuestion() =>
            throw UnimplementedError(),
        },
        newCards: newQuestions,
        dueCards: dueQuestions,
        nonDueCards: nonDueQuestions,
      ),
    ErrorMessage(:var error, :var shouldLogOut) =>
      shouldLogOut ? logOut() : Err(Exception(error)),
    _ => Err(Exception('Server Error')),
  };
}
