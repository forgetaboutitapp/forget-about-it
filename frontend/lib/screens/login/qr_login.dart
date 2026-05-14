import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import '../../fn/fn.dart';
import '../../state/login.dart';
import 'submit_type.dart';
import 'qr_scanner.dart';

class QrLogin extends HookConsumerWidget {
  final String? remoteURL;

  const QrLogin({super.key, required this.remoteURL});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Padding(
      padding: const EdgeInsets.all(32.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          const Text(
            'Scan the QR code provided by the provision command to automatically log in.',
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 32),
          ElevatedButton.icon(
            style: ElevatedButton.styleFrom(
              padding: const EdgeInsets.all(24.0),
            ),
            icon: const Icon(Icons.qr_code_scanner, size: 32),
            label: const Text('Open QR Scanner', style: TextStyle(fontSize: 20)),
            onPressed: () async {
              final code = await Navigator.push<String?>(
                context,
                MaterialPageRoute(builder: (context) => const QrScannerPage()),
              );
              if (!context.mounted) return;
              if (code != null) {
                final parts = code.split(';');
                String uri = remoteURL ?? '';
                String token = code;
                if (parts.length >= 2) {
                  uri = parts[0];
                  token = parts[1];
                }
                if (uri.isEmpty) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Invalid QR code: missing server host.')),
                  );
                  return;
                }
                Uri parsedUri;
                try {
                  final normalized = uri.contains('://')
                      ? uri
                      : 'http://$uri';
                  parsedUri = Uri.parse(normalized);
                } catch (_) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Invalid QR code: malformed server host.')),
                  );
                  return;
                }
                
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Logging in...')),
                );
                
                final v = await update(
                  parsedUri,
                  Token(token.trim()),
                );
                if (!context.mounted) return;
                switch (v) {
                  case Ok():
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Login successful!')),
                    );
                    break;
                  case Err():
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('Login failed: ${v.value}')),
                    );
                    break;
                }
              }
            },
          ),
        ],
      ),
    );
  }
}
