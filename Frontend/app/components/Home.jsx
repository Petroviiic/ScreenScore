import React, { useState, useEffect, useCallback } from "react";
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  AppState,
  Dimensions,
  ActivityIndicator,
  StatusBar,
  Linking,
} from "react-native";
import { useFocusEffect } from "@react-navigation/native";
import { styles, rankColor } from "@/assets/styles/home.styles";
import { useAppData } from "@/contexts/AppDataContext";

const { width } = Dimensions.get("window");

// ─── Helpers ───────────────────────────────────────────────────────────────────

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

// ─── Main Component ────────────────────────────────────────────────────────────

export default function Home() {
  const { weekData, todayIndex, hasPermission, refetchScreenTime } =
    useAppData();
  const [selectedDay, setSelectedDay] = useState(null);

  // Postavi selectedDay na today kad weekData stigne
  useEffect(() => {
    if (weekData.length > 0 && selectedDay === null) {
      setSelectedDay(todayIndex);
    }
  }, [weekData, todayIndex]);

  // Refetch kada se korisnik vrati na tab
  useFocusEffect(
    useCallback(() => {
      refetchScreenTime();
    }, [])
  );

  // Refetch kada se app vrati u foreground
  useEffect(() => {
    const sub = AppState.addEventListener("change", (state) => {
      if (state === "active") refetchScreenTime();
    });
    return () => sub.remove();
  }, []);

  const requestPermission = useCallback(async () => {
    await Linking.sendIntent("android.settings.USAGE_ACCESS_SETTINGS");
  }, []);

  // ── Derived values ─────────────────────────────────────────────────────────
  const maxMs = Math.max(...weekData.map((d) => d.totalMs), 1);
  const selected = weekData[selectedDay] ?? null;
  const pastDays = weekData.filter((d) => !d.isFuture && d.totalMs > 0);
  const weekAvgMs =
    pastDays.length > 0
      ? pastDays.reduce((s, d) => s + d.totalMs, 0) / pastDays.length
      : 0;

  // ── No permission ──────────────────────────────────────────────────────────
  if (!hasPermission && weekData.length === 0) {
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

  // ── Loading ────────────────────────────────────────────────────────────────
  if (weekData.length === 0) {
    return (
      <View style={styles.loadingContainer}>
        <StatusBar barStyle="light-content" />
        <ActivityIndicator size="large" color="#7C6EF5" />
        <Text style={styles.loadingText}>Analyzing your week…</Text>
      </View>
    );
  }

  // ── Main UI ────────────────────────────────────────────────────────────────
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
