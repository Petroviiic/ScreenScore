import { useEffect, useState } from "react";
import { View, ActivityIndicator } from "react-native";
import * as SecureStore from "expo-secure-store";
import { useRouter } from "expo-router";
import Login from "@/app/components/Login";

export default function Index() {
  const [isLoggedIn, setIsLoggedIn] = useState(null);
  const router = useRouter();

  useEffect(() => {
    const checkAuth = async () => {
      const token = await SecureStore.getItemAsync("jwt_token");
      setIsLoggedIn(!!token);
    };
    checkAuth();
  }, []);

  useEffect(() => {
    if (isLoggedIn === true) {
      router.replace("/(tabs)");
    }
  }, [isLoggedIn]);

  if (isLoggedIn === null || isLoggedIn === true) {
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
