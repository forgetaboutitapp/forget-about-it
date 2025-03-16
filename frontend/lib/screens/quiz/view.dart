import 'package:app/data/errors.dart';
import 'package:app/network/interfaces.dart';
import 'package:app/screens/quiz/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class QuizView extends HookConsumerWidget {
  static String location = '/quiz';
  final FetchData remoteServer;
  final Map<String, List<String>> tags;
  const QuizView({super.key, required this.remoteServer, required this.tags});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final question = ref.watch(quizQuestionsProvider);
    final initialized = useState(false);
    final tagsSet = useMemoized(() {
      return tags['tags']?.toISet();
    });
    useEffect(() {
      if (initialized.value == false) {
        try {
          ref
              .read(quizQuestionsProvider.notifier)
              .getNextQuestion(remoteServer, tagsSet);
        } on ServerException catch (e) {
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                backgroundColor: Colors.red,
                content: Text('$e'),
              ),
            );
          }
        }
        initialized.value = true;
      }
      return () {};
    }(), [key]);
    return Scaffold(
      appBar: AppBar(
        title: Text('Quiz'),
      ),
      body: Center(
        child: switch (question) {
          QuizQuestionStateWaiting() => CircularProgressIndicator(),
          QuizQuestionStateNone() => Center(
              child: Text('No questions available'),
            ),
          QuizQuestionStateError(:final exception) =>
            ErrorScreen(exception: exception),
          QuizQuestionStateData(:final question, :final answer, :final id) =>
            DisplayQuestion(
              remoteServer: remoteServer,
              question: question,
              answer: answer,
              id: id,
              tagsSet: tagsSet,
            ),
        },
      ),
    );
  }
}

class ErrorScreen extends StatelessWidget {
  const ErrorScreen({
    super.key,
    required this.exception,
  });

  final Exception exception;

  @override
  Widget build(BuildContext context) {
    WidgetsBinding.instance.addPostFrameCallback(
      (_) => ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          backgroundColor: Colors.red,
          content: Text(
            'Error: $exception',
          ),
        ),
      ),
    );
    return Center(child: Text(''));
  }
}

class DisplayQuestion extends HookConsumerWidget {
  final String question;
  final String answer;
  final int id;
  final ISet<String>? tagsSet;
  final FetchData remoteServer;
  const DisplayQuestion({
    super.key,
    required this.question,
    required this.answer,
    required this.id,
    required this.tagsSet,
    required this.remoteServer,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final shouldShowQuestion = useState(true);
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Expanded(
          child: Center(
            child: shouldShowQuestion.value
                ? Markdown(
                    data: question,
                    selectable: true,
                  )
                : Markdown(
                    data: answer,
                    selectable: true,
                  ),
          ),
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: shouldShowQuestion.value
              ? [
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: ElevatedButton(
                      onPressed: () {
                        shouldShowQuestion.value = false;
                      },
                      child: Text('Check'),
                    ),
                  )
                ]
              : [
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: ElevatedButton(
                      onPressed: () {
                        try {
                          ref
                              .read(quizQuestionsProvider.notifier)
                              .gradeQuestion(remoteServer, tagsSet, id, true);
                        } on ServerException catch (e) {
                          if (context.mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(
                                backgroundColor: Colors.red,
                                content: Text('$e'),
                              ),
                            );
                          }
                        }
                      },
                      style: ButtonStyle(
                        backgroundColor: WidgetStatePropertyAll(Colors.green),
                      ),
                      child: Text('Correct'),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: ElevatedButton(
                      onPressed: () {
                        ref
                            .read(quizQuestionsProvider.notifier)
                            .gradeQuestion(remoteServer, tagsSet, id, false);
                      },
                      style: ButtonStyle(
                        backgroundColor: WidgetStatePropertyAll(Colors.red),
                      ),
                      child: Text('Incorrect'),
                    ),
                  )
                ],
        ),
      ],
    );
  }
}
