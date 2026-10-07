package com.nerus.httpsms.ui.main

import android.telephony.PhoneNumberUtils
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.IntrinsicSize
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.BatteryAlert
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.LinearProgressIndicator
import androidx.compose.material3.LocalContentColor
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.nerus.httpsms.R
import com.nerus.httpsms.ui.theme.Blue500
import com.nerus.httpsms.ui.theme.Pink500
import java.util.Locale

@Composable
fun MainScreen(
    viewModel: MainViewModel,
    onSettingsClick: () -> Unit,
    onSmsPermissionClick: () -> Unit,
    onBatteryOptimizationClick: () -> Unit,
    onBatteryGuideClick: () -> Unit,
    onBatteryGuideDismiss: () -> Unit,
    onHeartbeatClick: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Spacer(modifier = Modifier.height(64.dp))

        Image(
            painter = painterResource(id = R.drawable.logo_cropped),
            contentDescription = stringResource(id = R.string.img_http_sms_logo),
            modifier = Modifier
                .width(147.dp)
                .height(92.dp)
        )

        Spacer(modifier = Modifier.height(24.dp))

        PhoneCard(
            phoneNumber = uiState.phoneNumberSIM1,
            isActive = uiState.isActiveSIM1,
            refreshTime = uiState.lastHeartbeatTime
        )

        if (uiState.isDualSim) {
            Spacer(modifier = Modifier.height(24.dp))
            PhoneCard(
                phoneNumber = uiState.phoneNumberSIM2,
                isActive = uiState.isActiveSIM2,
                refreshTime = uiState.lastHeartbeatTime
            )
        }

        Spacer(modifier = Modifier.height(16.dp))

        if (!uiState.isSmsPermissionGranted || !uiState.isBatteryOptimizationDisabled) {
            Column(modifier = Modifier.fillMaxWidth()) {
                if (!uiState.isSmsPermissionGranted) {
                    Button(
                        onClick = onSmsPermissionClick,
                        modifier = Modifier.align(Alignment.CenterHorizontally),
                        colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF4CAF50)),
                        contentPadding = PaddingValues(horizontal = 32.dp, vertical = 16.dp)
                    ) {
                        Text(
                            stringResource(id = R.string.enable_sms_permission),
                            color = Color.White,
                            fontSize = 18.sp
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Icon(
                            painter = painterResource(id = R.drawable.open_in_new_24),
                            contentDescription = null,
                            tint = Color.White
                        )
                    }
                    Spacer(modifier = Modifier.height(8.dp))
                }

                if (!uiState.isBatteryOptimizationDisabled) {
                    Button(
                        onClick = onBatteryOptimizationClick,
                        modifier = Modifier.align(Alignment.CenterHorizontally),
                        colors = ButtonDefaults.buttonColors(containerColor = Pink500),
                        contentPadding = PaddingValues(horizontal = 32.dp, vertical = 16.dp)
                    ) {
                        Icon(
                            imageVector = Icons.Default.BatteryAlert,
                            contentDescription = null,
                            tint = Color.White
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(
                            stringResource(id = R.string.disable_battery_optimization),
                            color = Color.White,
                            fontSize = 18.sp
                        )
                    }
                }
            }
        }

        if (uiState.showBatteryGuide) {
            Spacer(modifier = Modifier.height(16.dp))
            Card(
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.errorContainer),
                modifier = Modifier.fillMaxWidth()
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "Warning: ${uiState.manufacturer} Devices",
                        color = MaterialTheme.colorScheme.onErrorContainer,
                        fontWeight = FontWeight.Bold,
                        fontSize = 18.sp
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(
                        text = "Your device has aggressive background restrictions that may kill this app. Please disable them.",
                        color = MaterialTheme.colorScheme.onErrorContainer,
                        fontSize = 14.sp
                    )
                    Spacer(modifier = Modifier.height(12.dp))
                    Row {
                        Button(
                            onClick = onBatteryGuideClick,
                            colors = ButtonDefaults.buttonColors(containerColor = MaterialTheme.colorScheme.error)
                        ) {
                            Text("Fix It")
                        }
                        Spacer(modifier = Modifier.width(8.dp))
                        Button(
                            onClick = onBatteryGuideDismiss,
                            colors = ButtonDefaults.buttonColors(containerColor = Color.Gray)
                        ) {
                            Text("Dismiss")
                        }
                    }
                }
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        // Connection Status
        Card(
            modifier = Modifier.fillMaxWidth(),
            elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
        ) {
            Column(modifier = Modifier.padding(16.dp)) {
                Text("Connection Status", fontWeight = FontWeight.Bold, fontSize = 18.sp)
                Spacer(modifier = Modifier.height(8.dp))
                val serverStatus = when (uiState.isServerReachable) {
                    true -> "Online"
                    false -> "Offline"
                    null -> "Unknown"
                }
                Text("Server Reachability: $serverStatus", fontSize = 14.sp)
                Spacer(modifier = Modifier.height(4.dp))
                Text("Last Sent: ${uiState.lastSentTime}", fontSize = 14.sp)
                Spacer(modifier = Modifier.height(4.dp))
                Text("Last Received: ${uiState.lastReceivedTime}", fontSize = 14.sp)
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        Column(
            modifier = Modifier.width(IntrinsicSize.Max),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Button(
                onClick = onHeartbeatClick,
                modifier = Modifier.fillMaxWidth(),
                enabled = !uiState.isHeartbeatLoading,
                colors = ButtonDefaults.buttonColors(containerColor = Blue500),
                contentPadding = PaddingValues(horizontal = 32.dp, vertical = 16.dp)
            ) {
                Icon(
                    imageVector = Icons.Default.Favorite,
                    contentDescription = null,
                    tint = if (uiState.isHeartbeatLoading) LocalContentColor.current else Pink500
                )
                Spacer(modifier = Modifier.width(8.dp))
                Text(
                    stringResource(id = R.string.send_heartbeat),
                    color = Color.White,
                    fontSize = 18.sp
                )
            }

            if (uiState.isHeartbeatLoading) {
                LinearProgressIndicator(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(top = 4.dp),
                    color = Pink500
                )
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        Text(
            text = uiState.appVersion,
            fontSize = 14.sp,
            color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.6f)
        )

        Spacer(modifier = Modifier.weight(1f))
        Spacer(modifier = Modifier.height(16.dp))

        Button(
            onClick = onSettingsClick,
            colors = ButtonDefaults.buttonColors(containerColor = Color.Black),
            contentPadding = PaddingValues(horizontal = 32.dp, vertical = 16.dp)
        ) {
            Icon(Icons.Default.Settings, contentDescription = null, tint = Color.White)
            Spacer(modifier = Modifier.width(8.dp))
            Text(
                stringResource(id = R.string.main_app_settings),
                color = Color.White,
                fontSize = 18.sp
            )
        }

        Spacer(modifier = Modifier.height(16.dp))
    }
}

@Composable
fun PhoneCard(
    phoneNumber: String,
    isActive: Boolean,
    refreshTime: String
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(
                    text = PhoneNumberUtils.formatNumber(phoneNumber, Locale.getDefault().country) ?: phoneNumber,
                    fontSize = 28.sp,
                    fontWeight = FontWeight.Medium,
                    modifier = Modifier.weight(1f),
                    color = MaterialTheme.colorScheme.onSurface
                )
                if (isActive) {
                    Icon(
                        imageVector = Icons.Default.CheckCircle,
                        contentDescription = "Active",
                        tint = Color(0xFF70AB5C),
                        modifier = Modifier.size(24.dp)
                    )
                }
            }
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = refreshTime,
                fontSize = 16.sp,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
            )
        }
    }
}
