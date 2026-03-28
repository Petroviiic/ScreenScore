import {
  getMessaging,
  setBackgroundMessageHandler,
} from "@react-native-firebase/messaging";
import * as SecureStore from "expo-secure-store";
import { NativeModules } from "react-native";
import * as Application from "expo-application";

const { ScreenTimeModule } = NativeModules;
const API_URL = "https://shenika-ovarian-unpiratically.ngrok-free.dev";

// sync request se salje kada je aplikacija zatvorena
const messagingInstance = getMessaging();
setBackgroundMessageHandler(messagingInstance, async (remoteMessage) => {
  if (remoteMessage.data?.type === "sync") {
    const token = await SecureStore.getItemAsync("jwt_token");
    const stats = await ScreenTimeModule.getWeeklyStats();
    if (!stats) return;
    await fetch(`${API_URL}/v1/stats/sync-stats`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        screen_time: Math.floor(stats.totalTimeMs / 60000),
        recorded_at: new Date().toISOString(),
        device_id: Application.getAndroidId(),
      }),
    });
  }
});

// Normalni Expo Router entry
import "expo-router/entry";
