import 'package:app/screens/bulk-edit/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:file_picker/file_picker.dart';

abstract class FetchData {
  Future<String> getAllQuestions();
  Future<void> postAllQuestions(IList<Flashcard> flashcard);
  Future<String> generateNewToken();
  Future<bool> checkNewToken();
  Future<void> deleteNewToken();
  Future<String> getRemoteSettings();
  String getRemoteHost();
  Future<String> getAllTags();
  Future<String> getNextQuestion(ISet<String> tags);
  Future<void> gradeQuestion(int questionID, bool correct);
  Future<void> uploadAlgorithm(String data);
}

abstract class GenericFilepicker {
  Future<FilePickerResult?> pickFile({
    required FileType type,
    required List<String> allowedExtensions,
  });
}
