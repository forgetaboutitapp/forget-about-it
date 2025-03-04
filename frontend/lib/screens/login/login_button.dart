import 'package:flutter/material.dart';

class LoginButton extends StatelessWidget {
  final bool shouldEnable;
  final String remoteURLString;
  final Future<bool> Function(Uri) toRun;
  const LoginButton(
      {super.key,
      required this.shouldEnable,
      required this.remoteURLString,
      required this.toRun});

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: shouldEnable
          ? () async {
              Uri remoteURL;
              try {
                remoteURL = Uri.parse(remoteURLString);
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    backgroundColor: Colors.red,
                    content: Text('Cannot parse URL $remoteURLString'),
                  ),
                );
                return;
              }
              try {
                final success = await toRun(remoteURL);
                if (!success && context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      backgroundColor: Colors.red,
                      content: Text('Invalid Login'),
                    ),
                  );
                }
              } catch (e) {
                if (!context.mounted) return;
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    backgroundColor: Colors.red,
                    content: Text('Error Connecting: $e'),
                  ),
                );
              }
            }
          : null,
      child: Text('Submit'),
    );
  }
}
