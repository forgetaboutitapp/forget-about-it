import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:http/http.dart' as http;

import 'parse.dart';
import 'service.dart';

class BulkEditScreen extends HookWidget {
  static String location = '/bulk-edit';
  final http.Client client;
  final String token;
  final String remoteHost;
  const BulkEditScreen({
    super.key,
    required this.client,
    required this.token,
    required this.remoteHost,
  });

  @override
  Widget build(BuildContext context) {
    final future = useState(getAllQuestions(client, token, remoteHost));
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
                await postAllQuestions(client, token, remoteHost, flashcards);
                return (await getAllQuestions(client, token, remoteHost));
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
