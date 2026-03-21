import { useEffect, useState } from "react";
import { View, ActivityIndicator } from "react-native";
import * as SecureStore from "expo-secure-store";
import Home from "@/app/components/Home";
import Login from "@/app/components/Login";

export default function Index() {
  const [isLoggedIn, setIsLoggedIn] = useState(null);

  useEffect(() => {
    const checkAuth = async () => {
      const token = await SecureStore.getItemAsync("jwt_token");
      setIsLoggedIn(!!token);
    };
    checkAuth();
  }, []);

  // dok čeka provjeru tokena
  if (isLoggedIn === null) {
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

  return isLoggedIn ? <Home /> : <Login />;
}
