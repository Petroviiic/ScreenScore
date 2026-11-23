import {
  createContext,
  useCallback,
  useContext,
  useState,
  useMemo,
  useEffect,
} from "react";
import * as SecureStore from "expo-secure-store";
import { useAppData } from "@/contexts/AppDataContext";

const API_BASE = "https://shenika-ovarian-unpiratically.ngrok-free.dev";
const ShopContext = createContext(null);

export default function ShopProvider({ children }) {
  const { isAuthenticated } = useAppData();
  const [shopData, setShopData] = useState({
    all: [],
    common: [],
    rare: [],
    legendary: [],
  });

  const fetchShop = useCallback(async () => {
    try {
      const token = await SecureStore.getItemAsync("jwt_token");
      if (!token) return;
      const res = await fetch(
        `${API_BASE}/v1/notifications/get_available_shop`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
      if (res.ok) {
        const data = await res.json();
        // backend vraca null kada nema poruka
        const items = Array.isArray(data.data) ? data.data : [];
        setShopData({
          all: items,
          common: items.filter((i) => i.rarity === "common"),
          rare: items.filter((i) => i.rarity === "rare"),
          legendary: items.filter((i) => i.rarity === "legendary"),
        });
      }
    } catch (err) {
      console.log("fetchShop error:", err);
    }
  }, []);

  // Preload cim je korisnik ulogovan (ne pri otvaranju app-a, tada jos nema tokena)
  useEffect(() => {
    if (isAuthenticated) fetchShop();
  }, [isAuthenticated, fetchShop]);

  const value = useMemo(() => ({ shopData, fetchShop }), [shopData, fetchShop]);

  return <ShopContext.Provider value={value}>{children}</ShopContext.Provider>;
}

export const useShop = () => useContext(ShopContext);
