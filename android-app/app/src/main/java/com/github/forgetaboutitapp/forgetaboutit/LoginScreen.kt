package com.github.forgetaboutitapp.forgetaboutit

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.RadioButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.runtime.remember
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.focus.focusRequester
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.PopupProperties
import com.github.forgetaboutitapp.forgetaboutit.ui.theme.ForgetAboutItTheme
import java.util.Locale
import java.util.UUID

@Composable
fun App(
    hasUsageKey: Boolean,
    onLogin: (String) -> Unit,
    onLogout: () -> Unit,
    bitcoinWords: Set<String>,
    modifier: Modifier = Modifier
) {
    if (hasUsageKey) {
        MainScreen(modifier = modifier, onLogout = onLogout)
    } else {
        LoginScreen(modifier = modifier, onLogin = onLogin, bitcoinWords = bitcoinWords)
    }
}

@Composable
private fun LoginScreen(
    onLogin: (String) -> Unit,
    bitcoinWords: Set<String>,
    modifier: Modifier = Modifier
) {
    var selectedMethod by rememberSaveable { mutableStateOf(LoginMethod.WORDS) }
    var words by rememberSaveable { mutableStateOf(List(12) { "" }) }
    var uuid by rememberSaveable { mutableStateOf("") }
    var errorMessage by rememberSaveable { mutableStateOf<String?>(null) }
    val focusRequesters = remember { List(12) { FocusRequester() } }
    val wordsAreValid = words.all { it.isNotBlank() && bitcoinWords.contains(it) }
    val uuidIsValid = uuid.trim().isValidUuid()
    val canSubmit = if (selectedMethod == LoginMethod.WORDS) wordsAreValid else uuidIsValid

    Column(
        modifier = modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(horizontal = 24.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = androidx.compose.ui.Alignment.CenterHorizontally
    ) {
        Text("Welcome back", style = MaterialTheme.typography.headlineMedium)
        Text(
            text = "Choose how you want to continue",
            modifier = Modifier.padding(top = 8.dp, bottom = 24.dp),
            style = MaterialTheme.typography.bodyLarge
        )

        LoginMethodOption(
            title = "12 words",
            selected = selectedMethod == LoginMethod.WORDS,
            onClick = {
                selectedMethod = LoginMethod.WORDS
                errorMessage = null
            }
        )
        LoginMethodOption(
            title = "UUID",
            selected = selectedMethod == LoginMethod.UUID,
            onClick = {
                selectedMethod = LoginMethod.UUID
                errorMessage = null
            }
        )

        if (selectedMethod == LoginMethod.WORDS) {
            BitcoinWordsGrid(
                words = words,
                onWordChange = { index, word ->
                    words = words.toMutableList().also { it[index] = word.lowercase(Locale.ROOT) }
                    errorMessage = null
                },
                bitcoinWords = bitcoinWords,
                focusRequesters = focusRequesters,
                modifier = Modifier.padding(top = 16.dp)
            )
        } else {
            OutlinedTextField(
                value = uuid,
                onValueChange = {
                    uuid = it
                    errorMessage = null
                },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(top = 16.dp),
                label = { Text("UUID") },
                isError = uuid.isNotBlank() && !uuidIsValid,
                supportingText = if (uuid.isNotBlank() && !uuidIsValid) {
                    { Text("Enter a valid UUID.") }
                } else {
                    null
                },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Ascii),
                singleLine = true
            )
        }

        errorMessage?.let { message ->
            Text(
                text = message,
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(top = 8.dp),
                color = MaterialTheme.colorScheme.error
            )
        }

        Button(
            enabled = canSubmit,
            onClick = {
                val value = if (selectedMethod == LoginMethod.WORDS) {
                    words.joinToString(" ")
                } else {
                    uuid.trim()
                }
                val isValid = when (selectedMethod) {
                    LoginMethod.WORDS -> {
                        words.size == 12 && words.all(bitcoinWords::contains)
                    }
                    LoginMethod.UUID -> value.isValidUuid()
                }
                if (isValid) {
                    onLogin(value)
                } else {
                    errorMessage = when (selectedMethod) {
                        LoginMethod.WORDS -> "Enter exactly 12 Bitcoin words."
                        LoginMethod.UUID -> "Enter a valid UUID."
                    }
                }
            },
            modifier = Modifier
                .fillMaxWidth()
                .padding(top = 24.dp)
        ) {
            Text("Continue")
        }
    }
}

@Composable
private fun BitcoinWordsGrid(
    words: List<String>,
    onWordChange: (Int, String) -> Unit,
    bitcoinWords: Set<String>,
    focusRequesters: List<FocusRequester>,
    modifier: Modifier = Modifier
) {
    Column(modifier = modifier.fillMaxWidth()) {
        Text(
            text = "Enter your 12 Bitcoin words",
            style = MaterialTheme.typography.labelLarge,
            modifier = Modifier.padding(bottom = 8.dp)
        )
        repeat(6) { rowIndex ->
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(bottom = 8.dp),
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                repeat(2) { columnIndex ->
                    val index = rowIndex * 2 + columnIndex
                    BitcoinWordField(
                        value = words[index],
                        bitcoinWords = bitcoinWords,
                        onWordChange = { onWordChange(index, it) },
                        onWordSelected = {
                            if (index < focusRequesters.lastIndex) {
                                focusRequesters[index + 1].requestFocus()
                            }
                        },
                        modifier = Modifier
                            .weight(1f)
                            .focusRequester(focusRequesters[index])
                    )
                }
            }
        }
    }
}

@Composable
private fun BitcoinWordField(
    value: String,
    bitcoinWords: Set<String>,
    onWordChange: (String) -> Unit,
    onWordSelected: () -> Unit,
    modifier: Modifier = Modifier
) {
    val suggestions = bitcoinWords
        .asSequence()
        .filter { it.startsWith(value) && it != value }
        .sorted()
        .take(5)
        .toList()
    var suggestionsVisible by rememberSaveable { mutableStateOf(false) }

    Column(modifier = modifier) {
        OutlinedTextField(
            value = value,
            onValueChange = { input ->
                val normalized = input.lowercase(Locale.ROOT)
                if (normalized.isEmpty() || bitcoinWords.any { it.startsWith(normalized) }) {
                    onWordChange(normalized)
                    suggestionsVisible = normalized.isNotEmpty()
                }
            },
            modifier = Modifier.fillMaxWidth(),
            label = { Text("Word") },
            singleLine = true,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Ascii)
        )
        DropdownMenu(
            expanded = suggestionsVisible && suggestions.isNotEmpty(),
            onDismissRequest = { suggestionsVisible = false },
            properties = PopupProperties(focusable = false)
        ) {
            suggestions.forEach { suggestion ->
                DropdownMenuItem(
                    text = { Text(suggestion) },
                    onClick = {
                        onWordChange(suggestion)
                        suggestionsVisible = false
                        onWordSelected()
                    }
                )
            }
        }
    }
}

@Composable
private fun LoginMethodOption(
    title: String,
    selected: Boolean,
    onClick: () -> Unit
) {
    Card(onClick = onClick, modifier = Modifier.fillMaxWidth()) {
        androidx.compose.foundation.layout.Row(
            modifier = Modifier.padding(horizontal = 8.dp),
            verticalAlignment = androidx.compose.ui.Alignment.CenterVertically
        ) {
            RadioButton(selected = selected, onClick = onClick)
            Text(title, style = MaterialTheme.typography.titleMedium)
        }
    }
}

private enum class LoginMethod {
    WORDS,
    UUID
}

private fun String.isValidUuid(): Boolean = try {
    UUID.fromString(this)
    true
} catch (_: IllegalArgumentException) {
    false
}

@Preview(showBackground = true)
@Composable
private fun LoginPreview() {
    ForgetAboutItTheme {
        LoginScreen(onLogin = {}, bitcoinWords = emptySet())
    }
}
