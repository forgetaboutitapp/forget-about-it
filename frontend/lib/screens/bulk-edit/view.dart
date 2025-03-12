import 'package:app/network/interfaces.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';

import 'parse.dart';
import 'service.dart';

class BulkEditScreen extends HookWidget {
  static String location = '/bulk-edit';
  final FetchData remoteServer;
  const BulkEditScreen({
    super.key,
    required this.remoteServer,
  });

  @override
  Widget build(BuildContext context) {
    final future = useState(getAllQuestions(remoteServer: remoteServer));
    return Scaffold(
      appBar: AppBar(title: Text('Edit')),
      body: FutureBuilder(
        future: future.value,
        builder: (context, snapshot) => switch (snapshot.connectionState) {
          ConnectionState.none => Center(child: Text('Error')),
          ConnectionState.waiting ||
          ConnectionState.active =>
            Center(child: CircularProgressIndicator()),
          ConnectionState.done => displayWidget(context, snapshot, future),
        },
      ),
    );
  }

  Widget displayWidget(
      BuildContext context, snapshot, ValueNotifier<Future<String>> future) {
    if (snapshot.hasError) {
      return Center(child: Text('Error ${snapshot.error}'));
    } else if (snapshot.hasData) {
      return MainBulkEditView(
          text: snapshot.data ?? '',
          postQuestion: (text) async {
            future.value = () async {
              try {
                final flashcards = parse(text);
                await postAllQuestions(
                  remoteServer: remoteServer,
                  flashcards: flashcards,
                );
                return await getAllQuestions(remoteServer: remoteServer);
              } catch (e) {
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      backgroundColor: Colors.red,
                      content: Text('Error: $e'),
                    ),
                  );
                }
              }
              return text;
            }();
          });
    } else {
      return MainBulkEditView(
        text: 'No Data Or Error',
        postQuestion: (v) async {},
      );
    }
  }
}

class MainBulkEditView extends HookWidget {
  final String text;
  final Future<void> Function(String post) postQuestion;
  const MainBulkEditView(
      {super.key, required this.text, required this.postQuestion});

  @override
  Widget build(BuildContext context) {
    final controller = useTextEditingController(text: text);
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
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: TextButton(
                onPressed: () {},
                child: Text('Cancel'),
              ),
            ),
          ],
        )
      ],
    );
  }
}
