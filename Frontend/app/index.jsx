import { useEffect } from "react";
import { View, ActivityIndicator } from "react-native";
import { useRouter } from "expo-router";
import { useAppData } from "@/contexts/AppDataContext";
import Login from "@/app/components/Login";

export default function Index() {
  const { isReady, isAuthenticated, initializeApp } = useAppData();
  const router = useRouter();

  useEffect(() => {
    initializeApp();
  }, []);

  useEffect(() => {
    if (isReady && isAuthenticated) {
      router.replace("/(tabs)");
    }
  }, [isReady, isAuthenticated]);

  if (!isReady || isAuthenticated) {
    return (
      <View
        style={{
          flex: 1,
          justifyContent: "center",
          alignItems: "center",
          backgroundColor: "#0A0A0F",
        }}
      >
        <ActivityIndicator color="#00E5A0" />
      </View>
    );
  }

  return <Login />;
}
