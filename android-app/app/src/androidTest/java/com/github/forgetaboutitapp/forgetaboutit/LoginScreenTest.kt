package com.github.forgetaboutitapp.forgetaboutit

import androidx.compose.ui.test.assertCountEquals
import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.hasSetTextAction
import androidx.compose.ui.test.onAllNodes
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import androidx.compose.ui.test.performTextInput
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.github.forgetaboutitapp.forgetaboutit.ui.theme.ForgetAboutItTheme
import org.junit.Assert.assertEquals
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class LoginScreenTest {
    @get:Rule
    val composeTestRule = createComposeRule()

    private val bip39Words = setOf(
        "abandon",
        "ability",
        "able",
        "about",
        "above",
        "absent",
        "absorb",
        "abstract",
        "absurd",
        "abuse",
        "access",
        "accident"
    )

    @Test
    fun loginScreenShowsTwelveWordFields() {
        setApp(hasUsageKey = false)

        onNodeWithText("Welcome back").assertIsDisplayed()
        onAllNodes(hasSetTextAction()).assertCountEquals(12)
        onNodeWithText("Continue").assertIsDisplayed()
    }

    @Test
    fun uuidOptionShowsUuidField() {
        setApp(hasUsageKey = false)

        onNodeWithText("UUID").performClick()

        onNodeWithText("UUID").assertIsDisplayed()
        onAllNodes(hasSetTextAction()).assertCountEquals(1)
    }

    @Test
    fun validTwelveWordsCallsLogin() {
        var loggedInValue: String? = null
        setApp(
            hasUsageKey = false,
            onLogin = { loggedInValue = it }
        )

        val fields = onAllNodes(hasSetTextAction())
        bip39Words.toList().forEachIndexed { index, word ->
            fields[index].performTextInput(word)
        }
        onNodeWithText("Continue").performClick()

        assertEquals(bip39Words.joinToString(" "), loggedInValue)
    }

    @Test
    fun usageKeyShowsMainScreen() {
        setApp(hasUsageKey = true)

        onNodeWithText("Main screen").assertIsDisplayed()
        onNodeWithText("You are signed in.").assertIsDisplayed()
    }

    private fun setApp(
        hasUsageKey: Boolean,
        onLogin: (String) -> Unit = {}
    ) {
        composeTestRule.setContent {
            ForgetAboutItTheme {
                App(
                    hasUsageKey = hasUsageKey,
                    onLogin = onLogin,
                    onLogout = {},
                    bitcoinWords = bip39Words
                )
            }
        }
    }
}
