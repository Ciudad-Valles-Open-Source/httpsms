package com.nerus.httpsms.ui.settings

import android.content.Context
import androidx.lifecycle.ViewModel
import com.nerus.httpsms.Constants
import com.nerus.httpsms.R
import com.nerus.httpsms.ServerUrl
import com.nerus.httpsms.Settings
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow

data class SettingsUiState(
    val isDebugLogEnabled: Boolean = false,
    val phoneNumberSIM1: String = "",
    val isIncomingSIM1Enabled: Boolean = false,
    val isActiveSIM1: Boolean = false,
    val isIncomingCallEventsSIM1Enabled: Boolean = false,
    val isDualSim: Boolean = false,
    val phoneNumberSIM2: String = "",
    val isIncomingSIM2Enabled: Boolean = false,
    val isActiveSIM2: Boolean = false,
    val isIncomingCallEventsSIM2Enabled: Boolean = false,
    val encryptionKey: String = "",
    val isEncryptReceivedMessagesEnabled: Boolean = false,
    val serverUrl: String = "",
    val savedServerUrl: String = "",
    val serverUrlError: String? = null
)

class SettingsViewModel : ViewModel() {
    private val _uiState = MutableStateFlow(SettingsUiState())
    val uiState = _uiState.asStateFlow()

    fun initialize(context: Context) {
        _uiState.value = SettingsUiState(
            isDebugLogEnabled = Settings.isDebugLogEnabled(context),
            phoneNumberSIM1 = Settings.getSIM1PhoneNumber(context) ?: "",
            isIncomingSIM1Enabled = Settings.isIncomingMessageEnabled(context, Constants.SIM1),
            isActiveSIM1 = Settings.getActiveStatus(context, Constants.SIM1),
            isIncomingCallEventsSIM1Enabled = Settings.isIncomingCallEventsEnabled(context, Constants.SIM1),
            isDualSim = Settings.isDualSIM(context),
            phoneNumberSIM2 = Settings.getSIM2PhoneNumber(context) ?: "",
            isIncomingSIM2Enabled = Settings.isIncomingMessageEnabled(context, Constants.SIM2),
            isActiveSIM2 = Settings.getActiveStatus(context, Constants.SIM2),
            isIncomingCallEventsSIM2Enabled = Settings.isIncomingCallEventsEnabled(context, Constants.SIM2),
            encryptionKey = Settings.getEncryptionKey(context) ?: "",
            isEncryptReceivedMessagesEnabled = Settings.encryptReceivedMessages(context),
            serverUrl = Settings.getServerUrl(context) ?: "",
            savedServerUrl = Settings.getServerUrl(context) ?: ""
        )
    }

    fun setDebugLogEnabled(context: Context, enabled: Boolean) {
        Settings.setDebugLogEnabled(context, enabled)
        _uiState.value = _uiState.value.copy(isDebugLogEnabled = enabled)
    }

    fun setIncomingSIM1Enabled(context: Context, enabled: Boolean) {
        Settings.setIncomingActiveSIM1(context, enabled)
        _uiState.value = _uiState.value.copy(isIncomingSIM1Enabled = enabled)
    }

    fun setActiveSIM1(context: Context, enabled: Boolean) {
        Settings.setActiveStatusAsync(context, enabled, Constants.SIM1)
        _uiState.value = _uiState.value.copy(isActiveSIM1 = enabled)
    }

    fun setIncomingCallEventsSIM1Enabled(context: Context, enabled: Boolean) {
        Settings.setIncomingCallEventsEnabled(context, Constants.SIM1, enabled)
        _uiState.value = _uiState.value.copy(isIncomingCallEventsSIM1Enabled = enabled)
    }

    fun setIncomingSIM2Enabled(context: Context, enabled: Boolean) {
        Settings.setIncomingActiveSIM2(context, enabled)
        _uiState.value = _uiState.value.copy(isIncomingSIM2Enabled = enabled)
    }

    fun setActiveSIM2(context: Context, enabled: Boolean) {
        Settings.setActiveStatusAsync(context, enabled, Constants.SIM2)
        _uiState.value = _uiState.value.copy(isActiveSIM2 = enabled)
    }

    fun setIncomingCallEventsSIM2Enabled(context: Context, enabled: Boolean) {
        Settings.setIncomingCallEventsEnabled(context, Constants.SIM2, enabled)
        _uiState.value = _uiState.value.copy(isIncomingCallEventsSIM2Enabled = enabled)
    }

    fun setEncryptionKey(context: Context, key: String) {
        val trimmedKey = key.trim()
        if (trimmedKey.isEmpty()) {
            Settings.setEncryptionKey(context, null)
            Settings.setEncryptReceivedMessages(context, false)
            _uiState.value = _uiState.value.copy(encryptionKey = "", isEncryptReceivedMessagesEnabled = false)
        } else {
            Settings.setEncryptionKey(context, trimmedKey)
            _uiState.value = _uiState.value.copy(encryptionKey = trimmedKey)
        }
    }

    fun setEncryptReceivedMessagesEnabled(context: Context, enabled: Boolean) {
        Settings.setEncryptReceivedMessages(context, enabled)
        _uiState.value = _uiState.value.copy(isEncryptReceivedMessagesEnabled = enabled)
    }

    fun onServerUrlChange(url: String) {
        _uiState.value = _uiState.value.copy(serverUrl = url, serverUrlError = null)
    }

    /**
     * Validates, normalizes and saves the server URL typed in the settings screen.
     * Returns true when the URL was saved.
     */
    fun saveServerUrl(context: Context): Boolean {
        val serverUrl = ServerUrl.normalize(_uiState.value.serverUrl)
        val error = when (ServerUrl.validate(serverUrl)) {
            ServerUrl.Validation.VALID -> null
            ServerUrl.Validation.EMPTY -> context.getString(R.string.server_url_error_empty)
            ServerUrl.Validation.INVALID -> context.getString(R.string.server_url_error_invalid, serverUrl)
        }
        if (error != null) {
            _uiState.value = _uiState.value.copy(serverUrl = serverUrl, serverUrlError = error)
            return false
        }

        if (serverUrl != _uiState.value.savedServerUrl) {
            Settings.setServerUrlAsync(context, serverUrl)
            // Register this phone's FCM token with the new server the next time the app resumes.
            Settings.setFcmTokenLastUpdateTimestampAsync(context, 0)
        }
        _uiState.value = _uiState.value.copy(
            serverUrl = serverUrl,
            savedServerUrl = serverUrl,
            serverUrlError = null
        )
        return true
    }

    fun logout(context: Context, onLogoutComplete: () -> Unit) {
        Settings.setApiKeyAsync(context, null)
        Settings.setSIM1PhoneNumber(context, null)
        Settings.setSIM2PhoneNumber(context, null)
        Settings.setActiveStatusAsync(context, true, Constants.SIM1)
        Settings.setActiveStatusAsync(context, true, Constants.SIM2)
        Settings.setIncomingActiveSIM1(context, true)
        Settings.setIncomingActiveSIM2(context, true)
        Settings.setUserID(context, null)
        Settings.setEncryptionKey(context, null)
        Settings.setEncryptReceivedMessages(context, false)
        Settings.setFcmTokenLastUpdateTimestampAsync(context, 0)
        Settings.setIncomingCallEventsEnabled(context, Constants.SIM1, false)
        Settings.setIncomingCallEventsEnabled(context, Constants.SIM2, false)
        onLogoutComplete()
    }
}
