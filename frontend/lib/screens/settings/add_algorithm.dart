import 'dart:convert';

import 'package:app/screens/general-display/show_error.dart';
import 'package:app/screens/settings/models/remote_algorithm.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';

import '../../network/interfaces.dart';
import 'models/uploaded_file_state.dart';

class AddAlgorithm extends HookWidget {
  const AddAlgorithm({
    super.key,
    required this.remoteServer,
    required this.filepicker,
    required this.algos,
  });
  final FetchData remoteServer;
  final GenericFilepicker filepicker;
  final ISet<RemoteAlgorithm>? algos;

  @override
  Widget build(BuildContext context) {
    final ValueNotifier<UploadedFileState?> uploadedFile = useState(null);
    final uploadedFileValue = uploadedFile.value;
    final isUploading = useState(false);
    final data = uploadedFile.value?.data;
    return AlertDialog(
      title: Text('Add Spacing Algorithm'),
      content: ConstrainedBox(
        constraints: BoxConstraints(maxWidth: 340, minWidth: 340),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (isUploading.value) LinearProgressIndicator(),
            Text('Select a file to upload'),
            uploadedFileValue == null
                ? TextButton.icon(
                    onPressed: () async {
                      FilePickerResult? result = await filepicker.pickFile(
                        type: FileType.custom,
                        allowedExtensions: ['json'],
                      );
                      if (result != null && result.files.length == 1) {
                        String? dataAsString;
                        final dataBytes = result.files.first.bytes;
                        if (dataBytes != null) {
                          try {
                            dataAsString = Utf8Decoder().convert(dataBytes);
                          } on FormatException catch (_) {
                            if (context.mounted) {
                              showError(
                                  context, 'Uploaded data is not in UTF8');
                            }
                          }
                          if (dataAsString != null) {
                            uploadedFile.value = UploadedFileState(
                              filename: result.files.first.name,
                              data: dataAsString,
                            );
                          }
                        }
                      }
                    },
                    label: Text('Upload'),
                    icon: Icon(Icons.upload_file),
                  )
                : Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(uploadedFileValue.filename),
                      IconButton(
                        onPressed: isUploading.value
                            ? null
                            : () => uploadedFile.value = null,
                        icon: Icon(Icons.delete),
                      )
                    ],
                  ),
          ],
        ),
      ),
      actions: [
        TextButton(
            onPressed: isUploading.value || data == null
                ? null
                : () async {
                    final parsedJson = jsonDecode(data);
                    bool shouldUpload = true;

                    try {
                      final algoName = parsedJson['algorithm-name'];
                      IMap<String, RemoteAlgorithm> names = IMap();
                      final algosLocal = algos;
                      if (algosLocal != null) {
                        for (var element in algosLocal) {
                          names = names.add(element.algorithmName, element);
                        }
                      }
                      if (algoName is! String) {
                        showError(context,
                            'Algorithm name should be a string, not ${algoName.runtimeType}');
                        return;
                      }
                      final String algoNameString = algoName;
                      if (names.keys.contains(algoNameString)) {
                        shouldUpload = await showDialog(
                          context: context,
                          builder: (context) => AlertDialog(
                            title: Text('Warning'),
                            content: Text(
                              'You already have an algorithm with the name $algoNameString and has a version ${names.get(algoNameString)?.version}\n Do you want to overwright?',
                            ),
                            actions: [
                              TextButton(
                                onPressed: () {
                                  Navigator.of(context).pop(true);
                                },
                                child: Text('Ok'),
                              ),
                              TextButton(
                                onPressed: () {
                                  Navigator.of(context).pop(false);
                                },
                                child: Text('Cancel'),
                              )
                            ],
                          ),
                        );
                      }
                      if (shouldUpload) {
                        isUploading.value = true;
                        final res = await remoteServer.uploadAlgorithm(data);
                        if (res != null) {
                          if (context.mounted) {
                            showError(context, res);
                          }
                        }
                      }
                    } on Exception catch (e) {
                      if (context.mounted) {
                        showError(context, e.toString());
                      }
                    }
                    if (shouldUpload) {
                      isUploading.value = false;
                      if (context.mounted) {
                        Navigator.pop(context, true);
                      }
                    } else {
                      if (context.mounted) {
                        Navigator.pop(context, false);
                      }
                    }
                  },
            child: Text('Ok')),
        TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text('Cancel')),
      ],
    );
  }
}
