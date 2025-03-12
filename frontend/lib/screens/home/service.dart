import 'dart:convert';

import 'package:app/network/interfaces.dart';
import 'package:app/screens/home/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

Future<IList<Tag>> getAllTags(FetchData fd) async {
  final List<dynamic> parsedVal = jsonDecode(await fd.getAllTags());
  return parsedVal.map((e) => Tag.fromJson(e)).toIList();
}
