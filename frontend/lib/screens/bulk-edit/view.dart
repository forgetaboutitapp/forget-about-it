import 'dart:developer' as developer;

import 'package:app/future_widget/future_widget.dart';
import 'package:app/network/interfaces.dart';
import 'package:app/screens/general-display/show_error.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';

import '../../fn/fn.dart';
import 'model.dart';
import 'parse.dart';
import 'service.dart';

class BulkEditScreen extends HookWidget {
  static String location = '/bulk-edit';
  final FetchDataWithToken remoteServer;

  const BulkEditScreen({
    super.key,
    required this.remoteServer,
  });

  @override
  Widget build(BuildContext context) {
    final future = useState(() async {
      return await getAllQuestions(remoteServer: remoteServer);
    }());
    developer.log('future: ${future.toString()}');
    return Scaffold(
      appBar: AppBar(title: Text('Edit')),
      body: FutureWidget(
        future: future.value,
        built: (context, Result<String> r) => displayWidget(context, r, future),
        waiting: (context) => Center(child: CircularProgressIndicator()),
      ),
    );
  }

  Widget displayWidget(BuildContext context, Result<String> result,
      ValueNotifier<Future<Result<String>>> future) {
    developer.log('result: $result');
    return MainBulkEditView(
        text: switch (result) {
          Ok(:final value) => value,
          Err(:final value) => showErrorAndReturnEmptyString(context, value),
        },
        postQuestion: (text) async {
          future.value = () async {
            return switch (parse(text)) {
              Ok(:final value) => showErrorAndReturnString(
                  context,
                  await runPostGetAllQuestions(value),
                ),
              Err(:final value) => showErrorAndReturnString(
                  context,
                  (Ok(text), value),
                )
            };
          }();
        });
  }

  Future<(Result<String>, Exception?)> runPostGetAllQuestions(
      IList<Flashcard> r) async {
    final Result<String> paqRes = switch (await (await postAllQuestions(
      remoteServer: remoteServer,
      flashcards: r,
    ))
        .doFlatMap(
      (_) async => await getAllQuestions(
        remoteServer: remoteServer,
      ),
    )) {
      Ok(:final value) => Ok(value),
      Err(:final value) => Err(value),
    };
    developer.log('paqRes: ${paqRes.toString()}');
    return switch (paqRes) {
      Ok(:final value) => (Ok(value), null),
      Err(:final value) => (
          await getAllQuestions(
            remoteServer: remoteServer,
          ),
          value
        ),
    };
  }

  String showErrorAndReturnEmptyString(BuildContext context, Exception value) {
    developer.log('in error: $value');
    showErrorDelayed(context, value.toString());
    return '';
  }

  Result<String> showErrorAndReturnString(
    BuildContext context,
    (Result<String>, Exception?) s,
  ) {
    if (s.$2 != null) {
      showErrorDelayed(context, s.$2.toString());
    }
    developer.log('s1:${s.$1}');
    return s.$1;
  }
}

class MainBulkEditView extends HookWidget {
  final String text;
  final Future<void> Function(String post) postQuestion;

  const MainBulkEditView(
      {super.key, required this.text, required this.postQuestion});

  @override
  Widget build(BuildContext context) {
    developer.log('text: $text');
    final controller = useTextEditingController(text: text);

    useEffect(() {
      controller.text = text;
      return null;
    });
    return Column(
      children: [
        Expanded(
          child: TextField(
            maxLines: 9999,
            controller: controller,
          ),
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: ElevatedButton(
                onPressed: () async {
                  await postQuestion(controller.text);
                },
                child: Text('Save'),
              ),
            ),
          ],
        )
      ],
    );
  }
}
