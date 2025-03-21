import 'package:freezed_annotation/freezed_annotation.dart';

part 'uploaded_file_state.freezed.dart';

@freezed
class UploadedFileState with _$UploadedFileState {
  @override
  final String filename;
  @override
  final String data;

  UploadedFileState({required this.filename, required this.data});
}
