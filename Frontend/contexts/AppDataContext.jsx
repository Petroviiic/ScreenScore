import { createContext, useContext, useState } from "react";
import * as SecureStore from "expo-secure-store";
import { NativeModules } from "react-native";
import * as Application from "expo-application";

const { ScreenTimeModule } = NativeModules;
const API_BASE = "https://shenika-ovarian-unpiratically.ngrok-free.dev";
const AppDataContext = createContext(null);

const DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
const DAY_FULL = [
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
  "Sunday",
];

function getCurrentWeekDays() {
  const today = new Date();
  const dayOfWeek = today.getDay();
  const daysFromMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
  const monday = new Date(today);
  monday.setDate(today.getDate() - daysFromMonday);
  monday.setHours(0, 0, 0, 0);
  const days = [];
  for (let i = 0; i < 7; i++) {
    const d = new Date(monday);
    d.setDate(monday.getDate() + i);
    days.push(d);
  }
  return days;
}

export default function AppDataProvider({ children }) {
  const [isReady, setIsReady] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [weekData, setWeekData] = useState([]);
  const [todayIndex, setTodayIndex] = useState(0);
  const [hasPermission, setHasPermission] = useState(false);
  const [groupsData, setGroupsData] = useState([]);

  async function initializeApp() {
    const token = await SecureStore.getItemAsync("jwt_token");

    if (!token) {
      setIsAuthenticated(false);
      setIsReady(true);
      return;
    }

    // validate token i screentime u paralelu
    const [validateRes] = await Promise.allSettled([
      fetch(`${API_BASE}/v1/users/validate_token`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      }),
      fetchScreenTimeData(token, true),
      fetchGroups(token),
    ]);

    if (validateRes.status === "fulfilled") {
      setIsAuthenticated(validateRes.value.ok);
    } else {
      setIsAuthenticated(false);
    }

    setIsReady(true);
  }

  async function fetchScreenTimeData(token, sync = false) {
    try {
      const granted = await ScreenTimeModule.checkPermission();
      setHasPermission(granted);
      if (!granted) return;

      const stats = await ScreenTimeModule.getWeeklyStats();
      if (!stats) return;

      if (sync) syncToBackend(stats, token);

      const days = getCurrentWeekDays();
      const today = new Date();
      today.setHours(0, 0, 0, 0);

      const result = [];
      for (let i = 0; i < 7; i++) {
        const date = days[i];
        const isFuture = date > today;
        if (isFuture) {
          result.push({
            date,
            dayLabel: DAYS[i],
            dayFull: DAY_FULL[i],
            totalMs: 0,
            apps: [],
            isFuture: true,
          });
        } else {
          const daysAgo = Math.round((today - date) / (1000 * 60 * 60 * 24));
          const apps = await ScreenTimeModule.getAppsForDay(daysAgo);
          result.push({
            date,
            dayLabel: DAYS[i],
            dayFull: DAY_FULL[i],
            totalMs: stats.weeklyData[6 - daysAgo] || 0,
            apps: apps || [],
            isFuture: false,
          });
        }
      }

      const todayDay = today.getDay();
      const todayIdx = todayDay === 0 ? 6 : todayDay - 1;
      setTodayIndex(todayIdx);
      setWeekData(result);
    } catch (err) {
      console.error("fetchScreenTimeData error:", err);
    }
  }

  async function syncToBackend(stats, token) {
    try {
      const deviceId = Application.getAndroidId();
      await fetch(`${API_BASE}/v1/stats/sync-stats`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          screen_time: Math.floor(stats.totalTimeMs / 60000),
          recorded_at: new Date().toISOString(),
          device_id: deviceId,
        }),
      });
    } catch (err) {
      console.log("Sync error:", err);
    }
  }

  async function refetchScreenTime() {
    const token = await SecureStore.getItemAsync("jwt_token");
    await fetchScreenTimeData(token);
  }

  async function fetchGroups(token) {
    try {
      const res = await fetch(`${API_BASE}/v1/groups/get_user_groups`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (res.ok) {
        const data = await res.json();
        setGroupsData(Array.isArray(data.data) ? data.data : []);
      }
    } catch (err) {
      console.log("fetchGroups error:", err);
    }
  }

  async function refetchGroups() {
    const token = await SecureStore.getItemAsync("jwt_token");
    await fetchGroups(token);
  }

  return (
    <AppDataContext.Provider
      value={{
        isReady,
        isAuthenticated,
        weekData,
        todayIndex,
        hasPermission,
        groupsData,
        initializeApp,
        refetchScreenTime,
        refetchGroups,
      }}
    >
      {children}
    </AppDataContext.Provider>
  );
}

export const useAppData = () => useContext(AppDataContext);
