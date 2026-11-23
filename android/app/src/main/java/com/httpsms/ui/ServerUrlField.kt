package com.nerus.httpsms.ui

import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import com.nerus.httpsms.R
import com.nerus.httpsms.ServerUrl

private val WarningColor = Color(0xFFF57C00)

/**
 * Text field for the URL of the user's own httpSMS server. Shows the validation error if
 * any, otherwise a warning when the URL uses plain http:// or points to the phone itself.
 */
@Composable
fun ServerUrlField(
    value: String,
    onValueChange: (String) -> Unit,
    error: String?,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    onDone: (() -> Unit)? = null
) {
    val warning = when {
        error != null -> null
        ServerUrl.isLoopback(value) -> stringResource(id = R.string.server_url_warning_localhost)
        ServerUrl.isInsecure(value) -> stringResource(id = R.string.server_url_warning_http)
        else -> null
    }

    OutlinedTextField(
        value = value,
        onValueChange = onValueChange,
        label = { Text(stringResource(id = R.string.server_url)) },
        placeholder = { Text(stringResource(id = R.string.login_server_url_hint)) },
        modifier = modifier,
        singleLine = true,
        isError = error != null,
        supportingText = when {
            error != null -> {
                { Text(error) }
            }

            warning != null -> {
                { Text(warning, color = WarningColor) }
            }

            else -> null
        },
        keyboardOptions = KeyboardOptions(
            keyboardType = KeyboardType.Uri,
            autoCorrectEnabled = false,
            imeAction = ImeAction.Done
        ),
        keyboardActions = KeyboardActions(onDone = { onDone?.invoke() }),
        enabled = enabled
    )
}
