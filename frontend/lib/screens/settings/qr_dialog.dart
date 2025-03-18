import 'dart:convert';

import 'package:app/network/interfaces.dart';
import 'package:app/screens/login/view.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:qr_flutter/qr_flutter.dart';

class QRDialog extends HookWidget {
  final FetchData remoteServer;

  const QRDialog({
    super.key,
    required this.remoteServer,
  });

  @override
  Widget build(BuildContext context) {
    final ValueNotifier<String?> qrCodeData = useState(null);
    final show12Words = useState(LoginMethod.camera);
    useEffect(() {
      bool isCancelled = false;
      remoteServer.generateNewToken().then((e) async {
        qrCodeData.value = e;
        while (!isCancelled) {
          await remoteServer.checkNewToken().then((cancelled) {
            if (cancelled) {
              isCancelled = true;
              if (context.mounted) {
                Navigator.of(context).pop();
              }
            }
          });
          await Future.delayed(Duration(seconds: 1));
        }
      });
      return () {
        remoteServer.deleteNewToken();
      };
    }, []);
    final qrCodeDataValue = qrCodeData.value;
    if (qrCodeDataValue == null) {
      return Center(child: CircularProgressIndicator());
    }
    final decode = jsonDecode(qrCodeDataValue);
    return Dialog(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(8.0, 8.0, 8.0, 8.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 0, 0, 16.0),
              child: switch (show12Words.value) {
                LoginMethod.twelveWords =>
                  TwelveWordView(twelveWords: decode['mnemonic']),
                LoginMethod.camera => QrImageViewer(
                    uuid: decode['new-uuid'],
                    remoteServer: remoteServer,
                  ),
                LoginMethod.token => TokenView(
                    remoteHost: remoteServer.getRemoteHost(),
                    uuid: decode['new-uuid'],
                  ),
              },
            ),
            SegmentedButton(
                segments: [
                  ButtonSegment(
                      value: 1,
                      label: Text('12 Words'),
                      icon: Icon(Icons.input)),
                  ButtonSegment(
                      enabled: true,
                      value: 2,
                      label: Text('Camera'),
                      icon: Icon(Icons.qr_code)),
                  ButtonSegment(
                      value: 3,
                      label: Text('Token'),
                      icon: Icon(Icons.text_fields)),
                ],
                selected: switch (show12Words.value) {
                  LoginMethod.twelveWords => {1},
                  LoginMethod.camera => {2},
                  LoginMethod.token => {3},
                },
                onSelectionChanged: (v) =>
                    switch (v.map((e) => e).toList()[0]) {
                      1 => show12Words.value = LoginMethod.twelveWords,
                      2 => show12Words.value = LoginMethod.camera,
                      3 => show12Words.value = LoginMethod.token,
                      _ => throw AssertionError('selected an invalid state $v')
                    }),
          ],
        ),
      ),
    );
  }
}

class TokenView extends StatelessWidget {
  final String uuid;
  final String remoteHost;
  const TokenView({super.key, required this.uuid, required this.remoteHost});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Table(
        defaultColumnWidth: IntrinsicColumnWidth(),
        children: [
          TableRow(
            children: [
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Text('URL'),
              ),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: SelectableText(remoteHost),
              )
            ],
          ),
          TableRow(children: [
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: Text('Token'),
            ),
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: SelectableText(uuid),
            )
          ]),
        ],
      ),
    );
  }
}

class TwelveWordView extends StatelessWidget {
  final List<dynamic> twelveWords;

  const TwelveWordView({super.key, required this.twelveWords});
  @override
  Widget build(BuildContext context) {
    return Table(
      children: List.generate(
        6,
        (i) => TableRow(
          children: List.generate(
            2,
            (j) => Padding(
              padding: EdgeInsets.fromLTRB(
                  8,
                  MediaQuery.sizeOf(context).height > 550 ? 4 : 0,
                  8,
                  MediaQuery.sizeOf(context).height > 550 ? 4 : 0),
              child: Center(
                child: Text(
                  twelveWords[i * 2 + j].toString(),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class QrImageViewer extends StatelessWidget {
  final String uuid;
  final FetchData remoteServer;
  const QrImageViewer(
      {super.key, required this.uuid, required this.remoteServer});
  @override
  Widget build(BuildContext context) => QrImageView(
        data: '${remoteServer.getRemoteHost()};$uuid',
        version: QrVersions.auto,
        size: 400,
        gapless: false,
      );
}
