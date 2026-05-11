import '../../screens/login/login_button.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import '../../data/keys.dart';
import '../../fn/fn.dart';
import '../../state/login.dart';
import 'submit_type.dart';
import 'mdns_lookup.dart';

class TokenLogin extends HookConsumerWidget {
  final String? remoteURL;

  const TokenLogin({super.key, required this.remoteURL});

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
                suffixIcon: defaultTargetPlatform == TargetPlatform.android
                    ? IconButton(
                        icon: const Icon(Icons.search),
                        tooltip: 'Discover via mDNS',
                        onPressed: () async {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('Searching for server on local network...')),
                          );
                          final result = await discoverServerViaMdns();
                          if (!context.mounted) return;
                          if (result != null) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(content: Text('Found server: $result')),
                            );
                            urlController.text = result;
                            urlText.value = result;
                          } else {
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('No server found on local network.')),
                            );
                          }
                        },
                      )
                    : null,
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
                  final v = await update(
                    uri,
                    Token(tokenController.text.trim()),
                  );
                  return switch (v) {
                    Ok() => true,
                    Err() => false,
                  };
                }),
          )
        ],
      ),
    );
  }
}
