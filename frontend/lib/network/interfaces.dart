import 'package:app/screens/bulk-edit/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

abstract class FetchData {
  Future<String> getAllQuestions();
  Future<void> postAllQuestions(IList<Flashcard> flashcard);
  Future<String> generateNewToken();
  Future<bool> checkNewToken();
  Future<void> deleteNewToken();
  Future<String> getRemoteSettings();
  String getRemoteHost();
  Future<String> getAllTags();
}
