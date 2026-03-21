import React, { useState, useEffect, useCallback, useRef } from "react";
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  NativeModules,
  Platform,
  AppState,
  Dimensions,
  ActivityIndicator,
  StatusBar,
  Linking,
} from "react-native";
import { useFocusEffect } from "@react-navigation/native";
import * as Application from "expo-application";
const { ScreenTimeModule } = NativeModules;
import { styles, rankColor } from "@/assets/styles/home.styles";
const { width } = Dimensions.get("window");
import * as SecureStore from "expo-secure-store";
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

// ─── Helpers ──────────────────────────────────────────────────────────────────

function formatDuration(ms) {
  const totalMinutes = Math.floor(ms / 60000);
  const hours = Math.floor(totalMinutes / 60);
  const minutes = totalMinutes % 60;
  if (hours === 0) return `${minutes}m`;
  if (minutes === 0) return `${hours}h`;
  return `${hours}h ${minutes}m`;
}

function formatDurationLong(ms) {
  const totalMinutes = Math.floor(ms / 60000);
  const hours = Math.floor(totalMinutes / 60);
  const minutes = totalMinutes % 60;
  if (hours === 0) return `${minutes} min`;
  if (minutes === 0) return `${hours} hr`;
  return `${hours} hr ${minutes} min`;
}

// Returns days from Monday to Sunday of the current week
function getCurrentWeekDays() {
  const today = new Date();
  const dayOfWeek = today.getDay(); // 0=Sun, 1=Mon, ...
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

// ─── Main Component ───────────────────────────────────────────────────────────
export default function Home() {
  const [hasPermission, setHasPermission] = useState(false);
  const [loading, setLoading] = useState(true);
  const [weekData, setWeekData] = useState([]);
  const [selectedDay, setSelectedDay] = useState(null);
  const [todayIndex, setTodayIndex] = useState(0);

  // ── Check & request permission ──────────────────────────────────────────────
  const checkPermission = useCallback(async () => {
    try {
      const granted = await ScreenTimeModule.checkPermission();
      console.log("checkPermission result:", granted);
      setHasPermission(granted);
      return granted;
    } catch (err) {
      console.log("Permission error:", err);
      setHasPermission(false);
      return false;
    }
  }, []);

  const requestPermission = useCallback(async () => {
    await Linking.sendIntent("android.settings.USAGE_ACCESS_SETTINGS");
  }, []);

  // ── Send data to backend ──────────────────────
  const isSyncing = useRef(false);

  const syncToBackend = useCallback(async (stats) => {
    if (isSyncing.current) return;
    isSyncing.current = true;
    const token = await SecureStore.getItemAsync("jwt_token");
    const deviceId = Application.getAndroidId();
    try {
      await fetch("http://192.168.1.14:3000/v1/stats/sync-stats", {
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
    } finally {
      isSyncing.current = false;
    }
  }, []);

  // ── Fetch screen-time data for current week (Mon-Sun) ──────────────────────
  const fetchWeekData = useCallback(async () => {
    setLoading(true);
    try {
      const stats = await ScreenTimeModule.getWeeklyStats();

      if (!stats) {
        setLoading(false);
        return;
      }

      // Send data when user opens app
      await syncToBackend(stats);

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

      // Today index — 0=Mon, 6=Sun
      const todayDay = today.getDay();
      const todayIdx = todayDay === 0 ? 6 : todayDay - 1;
      setTodayIndex(todayIdx);
      setSelectedDay(todayIdx);
      setWeekData(result);
    } catch (err) {
      console.error("Failed to fetch usage stats:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  // ── Lifecycle ───────────────────────────────────────────────────────────────
  useFocusEffect(
    useCallback(() => {
      (async () => {
        const ok = await checkPermission();
        if (ok) fetchWeekData();
        else setLoading(false);
      })();
    }, [])
  );

  useEffect(() => {
    const sub = AppState.addEventListener("change", async (state) => {
      if (state === "active") {
        const ok = await checkPermission();
        if (ok) fetchWeekData();
      }
    });
    return () => sub.remove();
  }, []);

  // ── Derived values ──────────────────────────────────────────────────────────
  const maxMs = Math.max(...weekData.map((d) => d.totalMs), 1);
  const selected = weekData[selectedDay] ?? null;

  // Average only over past days (not future)
  const pastDays = weekData.filter((d) => !d.isFuture && d.totalMs > 0);
  const weekAvgMs =
    pastDays.length > 0
      ? pastDays.reduce((s, d) => s + d.totalMs, 0) / pastDays.length
      : 0;

  // ── No permission screen ────────────────────────────────────────────────────
  if (!hasPermission && !loading) {
    return (
      <View style={styles.permissionContainer}>
        <StatusBar barStyle="light-content" />
        <Text style={styles.permissionIcon}>🔒</Text>
        <Text style={styles.permissionTitle}>Permission Required</Text>
        <Text style={styles.permissionBody}>
          ScreenScore needs access to your usage data to show screen time
          statistics. Your data stays on your device.
        </Text>
        <TouchableOpacity
          style={styles.permissionBtn}
          onPress={requestPermission}
        >
          <Text style={styles.permissionBtnText}>Grant Access</Text>
        </TouchableOpacity>
      </View>
    );
  }

  // ── Loading screen ──────────────────────────────────────────────────────────
  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <StatusBar barStyle="light-content" />
        <ActivityIndicator size="large" color="#7C6EF5" />
        <Text style={styles.loadingText}>Analyzing your week…</Text>
      </View>
    );
  }

  // ── Main UI ─────────────────────────────────────────────────────────────────
  return (
    <View style={styles.root}>
      <StatusBar barStyle="light-content" />
      <ScrollView
        contentContainerStyle={styles.scroll}
        showsVerticalScrollIndicator={false}
      >
        {/* ── Header ── */}
        <View style={styles.header}>
          <Text style={styles.headerSub}>This week</Text>
          <Text style={styles.headerTitle}>Screen Time</Text>
          <View style={styles.avgPill}>
            <Text style={styles.avgText}>
              Daily avg · {formatDurationLong(weekAvgMs)}
            </Text>
          </View>
        </View>

        {/* ── Bar Chart ── */}
        <View style={styles.card}>
          <View style={styles.chartRow}>
            {weekData.map((day, i) => {
              const isSelected = i === selectedDay;
              const isToday = i === todayIndex;
              const barHeight = Math.max((day.totalMs / maxMs) * 130, 4);

              return (
                <TouchableOpacity
                  key={i}
                  style={styles.barCol}
                  onPress={() => !day.isFuture && setSelectedDay(i)}
                  activeOpacity={day.isFuture ? 1 : 0.75}
                >
                  <Text
                    style={[
                      styles.barValue,
                      isSelected && styles.barValueSelected,
                    ]}
                  >
                    {day.totalMs > 0 ? formatDuration(day.totalMs) : ""}
                  </Text>
                  <View style={styles.barTrack}>
                    <View
                      style={[
                        styles.bar,
                        { height: barHeight },
                        isSelected && styles.barSelected,
                        isToday && !isSelected && styles.barToday,
                        day.isFuture && styles.barFuture,
                      ]}
                    />
                  </View>
                  <Text
                    style={[
                      styles.barLabel,
                      isSelected && styles.barLabelSelected,
                    ]}
                  >
                    {day.dayLabel}
                  </Text>
                  {isToday && <View style={styles.todayDot} />}
                </TouchableOpacity>
              );
            })}
          </View>
        </View>

        {/* ── Selected Day Detail ── */}
        {selected && !selected.isFuture && (
          <View style={styles.detailSection}>
            <View style={styles.detailHeader}>
              <View>
                <Text style={styles.detailDay}>
                  {selectedDay === todayIndex ? "Today" : selected.dayFull}
                </Text>
                <Text style={styles.detailTotal}>
                  {formatDurationLong(selected.totalMs)}
                </Text>
              </View>
              {selected.totalMs > weekAvgMs ? (
                <View style={styles.badge}>
                  <Text style={styles.badgeText}>↑ Above avg</Text>
                </View>
              ) : (
                <View style={[styles.badge, styles.badgeGood]}>
                  <Text style={[styles.badgeText, styles.badgeTextGood]}>
                    ↓ Below avg
                  </Text>
                </View>
              )}
            </View>

            {/* App list */}
            {selected.apps.length === 0 ? (
              <Text style={styles.noDataText}>No app usage recorded.</Text>
            ) : (
              selected.apps.map((app, idx) => {
                const pct =
                  selected.totalMs > 0
                    ? (app.totalTimeMs / selected.totalMs) * 100
                    : 0;
                const appName =
                  app.appName ||
                  (app.packageName
                    ? app.packageName.split(".").pop()
                    : "Unknown");

                return (
                  <View key={app.packageName} style={styles.appRow}>
                    <View style={styles.appLeft}>
                      <View
                        style={[
                          styles.appRank,
                          { backgroundColor: rankColor(idx) },
                        ]}
                      >
                        <Text style={styles.appRankText}>{idx + 1}</Text>
                      </View>
                      <View style={styles.appInfo}>
                        <Text style={styles.appName} numberOfLines={1}>
                          {appName.charAt(0).toUpperCase() + appName.slice(1)}
                        </Text>
                        <Text style={styles.appPkg} numberOfLines={1}>
                          {app.packageName}
                        </Text>
                        <View style={styles.progressTrack}>
                          <View
                            style={[
                              styles.progressFill,
                              {
                                width: `${pct}%`,
                                backgroundColor: rankColor(idx),
                              },
                            ]}
                          />
                        </View>
                      </View>
                    </View>
                    <Text style={styles.appTime}>
                      {formatDuration(app.totalTimeMs)}
                    </Text>
                  </View>
                );
              })
            )}
          </View>
        )}

        {/* ── Week Summary Cards ── */}
        <View style={styles.summaryRow}>
          <View style={[styles.summaryCard, { flex: 1, marginRight: 8 }]}>
            <Text style={styles.summaryLabel}>Best day</Text>
            <Text style={styles.summaryValue}>
              {pastDays.length > 0
                ? pastDays.reduce((a, b) => (a.totalMs < b.totalMs ? a : b))
                    .dayLabel
                : "—"}
            </Text>
            <Text style={styles.summarySubValue}>
              {pastDays.length > 0
                ? formatDuration(Math.min(...pastDays.map((d) => d.totalMs)))
                : ""}
            </Text>
          </View>
          <View style={[styles.summaryCard, { flex: 1, marginLeft: 8 }]}>
            <Text style={styles.summaryLabel}>Most used</Text>
            <Text style={styles.summaryValue}>
              {(() => {
                const allApps = {};
                weekData.forEach((d) =>
                  d.apps.forEach((a) => {
                    allApps[a.packageName] =
                      (allApps[a.packageName] || 0) + a.totalTimeMs;
                  })
                );
                const top = Object.entries(allApps).sort(
                  (a, b) => b[1] - a[1]
                )[0];
                if (!top) return "—";
                const name = top[0].split(".").pop();
                return name.charAt(0).toUpperCase() + name.slice(1);
              })()}
            </Text>
            <Text style={styles.summarySubValue}>this week</Text>
          </View>
        </View>

        <View style={styles.footer}>
          <Text style={styles.footerText}>
            Data sourced from Android UsageStats
          </Text>
        </View>
      </ScrollView>
    </View>
  );
}
