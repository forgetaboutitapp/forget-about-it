package com.github.forgetaboutitapp.forgetaboutit

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.core.content.edit
import com.github.forgetaboutitapp.forgetaboutit.ui.theme.ForgetAboutItTheme

private const val PREFERENCES_NAME = "forget_about_it_preferences"
private const val USAGE_KEY = "usage_key"

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        val preferences = getSharedPreferences(PREFERENCES_NAME, MODE_PRIVATE)
        val bitcoinWords = assets.open("bip39_english.txt").bufferedReader().use { it.readLines().toSet() }

        setContent {
            ForgetAboutItTheme {
                var hasUsageKey by remember {
                    mutableStateOf(!preferences.getString(USAGE_KEY, null).isNullOrBlank())
                }
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    App(
                        hasUsageKey = hasUsageKey,
                        onLogin = { usageKey ->
                            preferences.edit { putString(USAGE_KEY, usageKey) }
                            hasUsageKey = true
                        },
                        onLogout = {
                            preferences.edit { remove(USAGE_KEY) }
                            hasUsageKey = false
                        },
                        bitcoinWords = bitcoinWords,
                        modifier = Modifier.padding(innerPadding)
                    )
                }
            }
        }
    }
}

@Composable
internal fun MainScreen(
    onLogout: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .padding(24.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text("Main screen", style = MaterialTheme.typography.headlineMedium)
        Text(
            text = "You are signed in.",
            modifier = Modifier.padding(top = 8.dp)
        )
        TextButton(onClick = onLogout, modifier = Modifier.padding(top = 16.dp)) {
            Text("Sign out")
        }
    }
}
