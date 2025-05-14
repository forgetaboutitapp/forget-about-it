import 'package:app/fn/fn.dart';
import 'package:app/screens/bulk-edit/model.dart';
import 'package:app/screens/bulk-edit/parse.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter_test/flutter_test.dart';

void main() async {
  test('Test unparsing', () {
    final cmpText =
        '0000000001 | Q1 | A1 |  |  | T1 T2 |\n0000000002 | Q2 | A2 |  |  | T1 T3 |\n0000000003 | Q3  Q3 | A3 |  |  | T1 T3 |\n0000000004 | Q4 \\n Q4 | A4 \\n A4 |  |  | T1 T3 |\n0000000005 | Q5 \\\\ \\n \\| Q5 | A5 \\\\ A5 \\|  |  |  | T1 T3 |';
    final t = unparse(
      <Flashcard>[
        Flashcard(
          id: 1,
          question: 'Q1',
          answer: 'A1',
          explanation: '',
          memoHint: '',
          tags: ['T1', 'T2'].toIList(),
        ),
        Flashcard(
          id: 2,
          question: 'Q2',
          answer: 'A2',
          explanation: '',
          memoHint: '',
          tags: ['T1', 'T3'].toIList(),
        ),
        Flashcard(
          id: 3,
          question: 'Q3  Q3',
          answer: 'A3',
          explanation: '',
          memoHint: '',
          tags: ['T1', 'T3'].toIList(),
        ),
        Flashcard(
          id: 4,
          question: 'Q4 \n Q4',
          answer: 'A4 \n A4',
          explanation: '',
          memoHint: '',
          tags: ['T1', 'T3'].toIList(),
        ),
        Flashcard(
          id: 5,
          question: 'Q5 \\ \n | Q5',
          answer: 'A5 \\ A5 | ',
          explanation: '',
          memoHint: '',
          tags: ['T1', 'T3'].toIList(),
        ),
      ].toIList(),
    );

    assert(t == cmpText);
  });

  test('Test parsing', () {
    assert(
      parse('|Q1|A1 | | | Tag1 Tag2|') ==
          Ok(
            [
              Flashcard(
                id: null,
                question: 'Q1',
                answer: 'A1',
                explanation: '',
                memoHint: '',
                tags: ['Tag1', 'Tag2'].toIList(),
              )
            ].toIList(),
          ),
    );
    assert(
      parse('|Q1|A1 | | | Tag1 Tag2|\n|Q2|A2 | | | Tag|') ==
          Ok(
            [
              Flashcard(
                  id: null,
                  question: 'Q1',
                  answer: 'A1',
                  explanation: '',
                  memoHint: '',
                  tags: ['Tag1', 'Tag2'].toIList()),
              Flashcard(
                id: null,
                question: 'Q2',
                answer: 'A2',
                explanation: '',
                memoHint: '',
                tags: ['Tag'].toIList(),
              )
            ].toIList(),
          ),
    );

    assert(
      parse('''  % nothing here
1 |Q1|A1 | | | Tag1 Tag2|
 2 |Q2|A2 | | | Tag|
 3|q1\\|q2|a1\\na2 | | | tag1 tag2|
 4| 95\\% is what?|5 less than 100|||t1 t2| % See, this should work
 5|Another q|Another a | | | t321|''') ==
          Ok([
            Flashcard(
                id: 1,
                question: 'Q1',
                answer: 'A1',
                explanation: '',
                memoHint: '',
                tags: ['Tag1', 'Tag2'].toIList()),
            Flashcard(
              id: 2,
              question: 'Q2',
              answer: 'A2',
              explanation: '',
              memoHint: '',
              tags: ['Tag'].toIList(),
            ),
            Flashcard(
              id: 3,
              question: 'q1|q2',
              answer: 'a1\na2',
              explanation: '',
              memoHint: '',
              tags: ['tag1', 'tag2'].toIList(),
            ),
            Flashcard(
              id: 4,
              question: '95% is what?',
              answer: '5 less than 100',
              explanation: '',
              memoHint: '',
              tags: ['t1', 't2'].toIList(),
            ),
            Flashcard(
              id: 5,
              question: 'Another q',
              answer: 'Another a',
              explanation: '',
              memoHint: '',
              tags: ['t321'].toIList(),
            ),
          ].toIList()),
    );

    var caught = false;

    final f = parse('123|q1|a1|tag1 \\');
    Exception? err;
    f.match(onOk: (_) => {err = null}, onErr: (e) => err = e);
    if (err is InvalidEscapeException) {
      InvalidEscapeException err1 = err! as InvalidEscapeException;
      assert(err1.lineNumber == 0);
      assert(err1.stringWithError == '123|q1|a1|tag1 \\');
    } else {
      assert(false);
    }

    final f1 = parse('321|q0|a0|tag1 \n 123|q1|a1|tag1 \\');
    err = null;
    f1.match(onOk: (_) => {err = null}, onErr: (e) => err = e);
    if (err is InvalidEscapeException) {
      InvalidEscapeException e = err! as InvalidEscapeException;
      assert(e.lineNumber == 1);
      assert(e.stringWithError == ' 123|q1|a1|tag1 \\');
    } else {
      assert(false);
    }

    err = null;
    final f2 = parse('321|q0|a0|tag1 \n 123|q1|a1|tag1 \\\n444|q2|a2|tag2');
    f2.match(onOk: (_) => {err = null}, onErr: (e) => err = e);
    if (err is InvalidEscapeException) {
      InvalidEscapeException e = err! as InvalidEscapeException;
      assert(e.lineNumber == 1);
      assert(e.stringWithError == ' 123|q1|a1|tag1 \\');
    } else {
      assert(false);
    }

    err = null;

    final f3 = parse('321|q0|a0');
    f3.match(onOk: (_) => {err = null}, onErr: (e) => err = e);
    if (err is InvalidNumberOfFields) {
      InvalidNumberOfFields e = err! as InvalidNumberOfFields;
      assert(e.lineNumber == 0);
      assert(e.stringWithError == '321|q0|a0');
    } else {
      assert(false);
    }

    final f4 = parse('321|q0|a0|tag2|');
    f4.match(onOk: (_) => {err = null}, onErr: (e) => err = e);
    if (err is InvalidNumberOfFields) {
      InvalidNumberOfFields e = err! as InvalidNumberOfFields;

      assert(e.lineNumber == 0);
      assert(e.stringWithError == '321|q0|a0|tag2|');
      caught = true;
    } else {
      assert(false);
    }
    final f5 = parse('321|q0|a0|a|b|');
    f5.match(onOk: (_) => {err = null}, onErr: (e) => err = e);

    if (err is! NoTagException) {
    } else {
      assert(false);
    }
    assert(caught);
  });
}
