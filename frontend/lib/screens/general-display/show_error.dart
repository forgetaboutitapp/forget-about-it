import 'package:flutter/material.dart';

void showErrorDelayed(BuildContext context, String error) =>
    WidgetsBinding.instance.addPostFrameCallback(
      (_) => showError(
        context,
        error,
      ),
    );
void showError(BuildContext context, String error) => context.mounted
    ? ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          backgroundColor: Colors.red,
          content: Text(error),
        ),
      )
    : () {};
