import 'package:flutter/material.dart';

void showError(BuildContext context, String error) =>
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        backgroundColor: Colors.red,
        content: Text(error),
      ),
    );
