import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';

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
  if (client.hasError()) {
    if (client.error.shouldLogOut) {
      logOut();
    }
    return Err(Exception(client.error.error));
  }

  if (!client.hasOk()) {
    return Err(Exception('Server Error'));
  }

  final question = client.ok;
  return Ok(QuizQuestionStateData(
    id: question.flashcard.id,
    question: question.flashcard.question,
    answer: question.flashcard.answer,
    questionType: switch (question.typeOfQuestion) {
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
    newCards: question.newQuestions,
    dueCards: question.dueQuestions,
    nonDueCards: question.nonDueQuestions,
  ));
}
