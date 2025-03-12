import 'package:app/screens/bulk-edit/model.dart';
import 'package:app/screens/bulk-edit/parse.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter_test/flutter_test.dart';

void main() async {
  test('Test parsing', () {
    assert(
      parse('|Q1|A1|Tag1 Tag2') ==
          [
            Flashcard(
                id: null,
                question: 'Q1',
                answer: 'A1',
                tags: ['Tag1', 'Tag2'].toIList())
          ].toIList(),
    );
    assert(
      parse('|Q1|A1|Tag1 Tag2\n|Q2|A2|Tag') ==
          [
            Flashcard(
                id: null,
                question: 'Q1',
                answer: 'A1',
                tags: ['Tag1', 'Tag2'].toIList()),
            Flashcard(
                id: null, question: 'Q2', answer: 'A2', tags: ['Tag'].toIList())
          ].toIList(),
    );

    assert(
      parse('''  % nothing here
1 |Q1|A1|Tag1 Tag2
 2 |Q2|A2|Tag
 3|q1\\|q2|a1\\na2|tag1 tag2
 4| 95\\% is what?|5 less than 100|t1 t2 % See, this should work
 5|Another q|Another a|t321''') ==
          [
            Flashcard(
                id: 1,
                question: 'Q1',
                answer: 'A1',
                tags: ['Tag1', 'Tag2'].toIList()),
            Flashcard(
              id: 2,
              question: 'Q2',
              answer: 'A2',
              tags: ['Tag'].toIList(),
            ),
            Flashcard(
              id: 3,
              question: 'q1|q2',
              answer: 'a1\na2',
              tags: ['tag1', 'tag2'].toIList(),
            ),
            Flashcard(
              id: 4,
              question: '95% is what?',
              answer: '5 less than 100',
              tags: ['t1', 't2'].toIList(),
            ),
            Flashcard(
              id: 5,
              question: 'Another q',
              answer: 'Another a',
              tags: ['t321'].toIList(),
            ),
          ].toIList(),
    );

    var caught = false;
    try {
      parse('123|q1|a1|tag1 \\');
    } on InvalidEscapeException catch (e) {
      assert(e.lineNumber == 0);
      assert(e.stringWithError == '123|q1|a1|tag1 \\');
      caught = true;
    }
    assert(caught);

    caught = false;
    try {
      parse('321|q0|a0|tag1 \n 123|q1|a1|tag1 \\');
    } on InvalidEscapeException catch (e) {
      assert(e.lineNumber == 1);
      assert(e.stringWithError == ' 123|q1|a1|tag1 \\');
      caught = true;
    }
    assert(caught);

    caught = false;
    try {
      parse('321|q0|a0|tag1 \n 123|q1|a1|tag1 \\\n444|q2|a2|tag2');
    } on InvalidEscapeException catch (e) {
      assert(e.lineNumber == 1);
      assert(e.stringWithError == ' 123|q1|a1|tag1 \\');
      caught = true;
    }
    assert(caught);

    caught = false;
    try {
      parse('321|q0|a0');
    } on InvalidNumberOfFields catch (e) {
      assert(e.lineNumber == 0);
      assert(e.stringWithError == '321|q0|a0');
      caught = true;
    }
    assert(caught);

    caught = false;
    try {
      parse('321|q0|a0|tag2|');
    } on InvalidNumberOfFields catch (e) {
      assert(e.lineNumber == 0);
      assert(e.stringWithError == '321|q0|a0|tag2|');
      caught = true;
    }
    assert(caught);

    caught = false;
    try {
      parse('321|q0|a0|');
    } on NoTagException catch (e) {
      assert(e.id == 321);
      assert(e.question == 'q0');
      caught = true;
    }
    assert(caught);
  });
}
