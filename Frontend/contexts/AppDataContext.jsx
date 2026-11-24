import { createContext, useContext, useRef, useState } from "react";
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

const DAY_MS = 1000 * 60 * 60 * 24;

// weekOffset: 0 = ova sedmica, 1 = prosla, 2 = pretprosla...
export function getWeekDays(weekOffset = 0) {
  const today = new Date();
  const dayOfWeek = today.getDay();
  const daysFromMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
  const monday = new Date(today);
  monday.setDate(today.getDate() - daysFromMonday - weekOffset * 7);
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

      const result = await buildWeek(0, stats.weeklyData);

      const today = new Date();
      const todayDay = today.getDay();
      const todayIdx = todayDay === 0 ? 6 : todayDay - 1;
      setTodayIndex(todayIdx);
      setWeekData(result);
    } catch (err) {
      console.error("fetchScreenTimeData error:", err);
    }
  }

  // Pravi podatke za 7 dana jedne sedmice (pon-ned).
  // weeklyData (iz getWeeklyStats) postoji samo za zadnjih 7 dana; za starije
  // dane ukupno vrijeme racunamo kao zbir vremena aplikacija tog dana.
  async function buildWeek(weekOffset, weeklyData = null) {
    const days = getWeekDays(weekOffset);
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const result = [];
    for (let i = 0; i < 7; i++) {
      const date = days[i];
      if (date > today) {
        result.push({
          date,
          dayLabel: DAYS[i],
          dayFull: DAY_FULL[i],
          totalMs: 0,
          apps: [],
          isFuture: true,
        });
        continue;
      }

      const daysAgo = Math.round((today - date) / DAY_MS);
      const apps = (await ScreenTimeModule.getAppsForDay(daysAgo)) || [];
      const totalMs =
        weeklyData && daysAgo <= 6
          ? weeklyData[6 - daysAgo] || 0
          : apps.reduce((sum, a) => sum + a.totalTimeMs, 0);

      result.push({
        date,
        dayLabel: DAYS[i],
        dayFull: DAY_FULL[i],
        totalMs,
        apps,
        isFuture: false,
      });
    }
    return result;
  }

  // Prosle sedmice se vise ne mijenjaju, pa ih pamtimo. Kljuc je datum
  // ponedjeljka (ne offset), da kes ostane tacan i kad pocne nova sedmica.
  const pastWeeksCache = useRef({});

  async function fetchPastWeek(weekOffset) {
    const key = getWeekDays(weekOffset)[0].toDateString();
    if (pastWeeksCache.current[key]) return pastWeeksCache.current[key];
    try {
      const data = await buildWeek(weekOffset);
      pastWeeksCache.current[key] = data;
      return data;
    } catch (err) {
      console.log("fetchPastWeek error:", err);
      return null;
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
        fetchPastWeek,
        refetchGroups,
      }}
    >
      {children}
    </AppDataContext.Provider>
  );
}

export const useAppData = () => useContext(AppDataContext);
