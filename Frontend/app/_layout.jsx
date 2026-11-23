import { Stack } from "expo-router";
import AppDataProvider from "@/contexts/AppDataContext";
import ShopProvider from "@/contexts/ShopContext";

export default function RootLayout() {
  return (
    <AppDataProvider>
      <ShopProvider>
        <Stack screenOptions={{ headerShown: false }} />
      </ShopProvider>
    </AppDataProvider>
  );
}
