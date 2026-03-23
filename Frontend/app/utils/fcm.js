import * as Notifications from "expo-notifications";
import * as Device from "expo-device";
import { Platform } from "react-native";

export default async function getFCMToken() {
  if (!Device.isDevice) {
    console.warn("FCM radi samo na fizičkom uređaju");
    return null;
  }

  const { status: existingStatus } = await Notifications.getPermissionsAsync();
  let finalStatus = existingStatus;

  if (existingStatus !== "granted") {
    const { status } = await Notifications.requestPermissionsAsync();
    finalStatus = status;
  }

  if (finalStatus !== "granted") {
    console.warn("Notifikacije nisu dozvoljene");
    return null;
  }

  // Expo push token (interno koristi FCM na Androidu)
  // const tokenData = await Notifications.getExpoPushTokenAsync();
  // return tokenData.data; // "ExponentPushToken[xxx]"

  const tokenData = await Notifications.getDevicePushTokenAsync();
  return tokenData.data;
}
