import 'package:app/screens/login/login_button.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:http/http.dart' as http;
import '../../data/keys.dart';
import '../../state/login.dart';
import 'submit_type.dart';

class TokenLogin extends HookConsumerWidget {
  final String? remoteURL;
  final http.Client client;

  const TokenLogin({super.key, required this.remoteURL, required this.client});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final urlController = useTextEditingController(text: remoteURL ?? '');
    final tokenController = useTextEditingController();

    final urlText = useState(remoteURL ?? '');
    final tokenText = useState('');
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: TextField(
              key: ValueKey(loginURLKeyInToken),
              onChanged: (value) => urlText.value = value,
              controller: urlController,
              readOnly: kReleaseMode && remoteURL != null,
              decoration: InputDecoration(
                border: OutlineInputBorder(),
                labelText: 'Server Host',
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: TextField(
              key: ValueKey(loginTokenKeyInToken),
              onChanged: (value) => tokenText.value = value,
              controller: tokenController,
              decoration: InputDecoration(
                border: OutlineInputBorder(),
                labelText: 'Token',
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: LoginButton(
              shouldEnable: tokenController.text.trim() != '' &&
                  urlController.text.trim() != '',
              remoteURLString: urlController.text,
              toRun: (uri) async {
                return await update(
                  client,
                  uri,
                  Token(
                    tokenController.text.trim(),
                  ),
                );
              },
            ),
          )
        ],
      ),
    );
  }
}
