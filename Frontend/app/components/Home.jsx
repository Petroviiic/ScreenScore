import React, { useState, useEffect, useCallback, useRef } from "react";
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  AppState,
  ActivityIndicator,
  StatusBar,
  Linking,
  PanResponder,
  Animated,
} from "react-native";
import { useFocusEffect } from "@react-navigation/native";
import { styles, rankColor } from "@/assets/styles/home.styles";
import { useAppData, getWeekDays } from "@/contexts/AppDataContext";

// Koliko sedmica unazad korisnik moze ici. Android cuva detaljnu istoriju
// koristenja ograniceno vrijeme, pa dalje od ovoga nema smisla.
const MAX_WEEKS_BACK = 4;
const SWIPE_THRESHOLD = 50;

const MONTHS = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];

function weekLabel(weekOffset) {
  if (weekOffset === 0) return "This week";
  if (weekOffset === 1) return "Last week";
  const days = getWeekDays(weekOffset);
  const fmt = (d) => `${MONTHS[d.getMonth()]} ${d.getDate()}`;
  return `${fmt(days[0])} – ${fmt(days[6])}`;
}

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
  const {
    weekData,
    todayIndex,
    hasPermission,
    refetchScreenTime,
    fetchPastWeek,
  } = useAppData();
  const [selectedDay, setSelectedDay] = useState(null);

  // ── Sedmice: 0 = ova sedmica, 1 = prosla... ──────────────────────────────
  const [weekOffset, setWeekOffset] = useState(0);
  const [pastWeek, setPastWeek] = useState(null); // podaci kad je weekOffset > 0
  const weekOffsetRef = useRef(0); // trenutni offset za PanResponder
  const translateX = useRef(new Animated.Value(0)).current;

  // Postavi selectedDay na today kad weekData stigne
  useEffect(() => {
    if (weekData.length > 0 && selectedDay === null) {
      setSelectedDay(todayIndex);
    }
  }, [weekData, todayIndex]);

  // Ucitaj izabranu proslu sedmicu (iz kesa je instant) i unaprijed ucitaj
  // jos jednu stariju, da sljedeci swipe bude brz.
  // fetchPastWeek namjerno nije u dependency listi: iz contexta stize kao nova
  // funkcija pri svakom renderu, pa bi se efekat vrtio u krug.
  useEffect(() => {
    if (weekOffset === 0) {
      setPastWeek(null);
      return;
    }
    let cancelled = false;
    setPastWeek(null);
    fetchPastWeek(weekOffset).then((data) => {
      if (cancelled) return;
      setPastWeek(data ?? []);
      if (weekOffset < MAX_WEEKS_BACK) fetchPastWeek(weekOffset + 1);
    });
    return () => {
      cancelled = true;
    };
  }, [weekOffset]);

  // delta: +1 = starija sedmica, -1 = novija
  const changeWeek = useCallback(
    (delta) => {
      const next = weekOffsetRef.current + delta;
      if (next < 0 || next > MAX_WEEKS_BACK) {
        // na granici: samo vrati graf na mjesto
        Animated.spring(translateX, {
          toValue: 0,
          useNativeDriver: true,
        }).start();
        return;
      }
      weekOffsetRef.current = next;
      setWeekOffset(next);
      // za proslu sedmicu oznaci nedjelju, za ovu danasnji dan
      setSelectedDay(next === 0 ? todayIndex : 6);
      // nova sedmica "uleti" u smjeru pokreta prsta
      // (starija dolazi s lijeva, novija s desna)
      translateX.setValue(delta > 0 ? -60 : 60);
      Animated.spring(translateX, {
        toValue: 0,
        useNativeDriver: true,
        speed: 20,
      }).start();
    },
    [todayIndex, translateX],
  );
  // PanResponder se pravi jednom, pa preko ref-a uvijek zove najnoviji changeWeek
  const changeWeekRef = useRef(changeWeek);
  useEffect(() => {
    changeWeekRef.current = changeWeek;
  }, [changeWeek]);

  // Prst desno = prethodna sedmica, prst lijevo = sljedeca (novija), kao u kalendaru.
  // Hvata samo horizontalne pokrete, pa vertikalni scroll radi normalno.
  const panResponder = useRef(
    PanResponder.create({
      onMoveShouldSetPanResponder: (_, g) =>
        Math.abs(g.dx) > 15 && Math.abs(g.dx) > Math.abs(g.dy) * 1.5,
      onPanResponderTerminationRequest: () => false,
      onPanResponderMove: (_, g) => translateX.setValue(g.dx * 0.4),
      onPanResponderRelease: (_, g) => {
        if (g.dx >= SWIPE_THRESHOLD) changeWeekRef.current(1);
        else if (g.dx <= -SWIPE_THRESHOLD) changeWeekRef.current(-1);
        else
          Animated.spring(translateX, {
            toValue: 0,
            useNativeDriver: true,
          }).start();
      },
      onPanResponderTerminate: () =>
        Animated.spring(translateX, {
          toValue: 0,
          useNativeDriver: true,
        }).start(),
    }),
  ).current;

  // Refetch kada se korisnik vrati na tab
  useFocusEffect(
    useCallback(() => {
      refetchScreenTime();
    }, []),
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
  // sedmica koja se trenutno prikazuje
  const shownWeek = weekOffset === 0 ? weekData : (pastWeek ?? []);
  const isLoadingWeek = weekOffset > 0 && pastWeek === null;
  const isCurrentWeek = weekOffset === 0;
  const weekHasData = shownWeek.some((d) => d.totalMs > 0);

  const maxMs = Math.max(...shownWeek.map((d) => d.totalMs), 1);
  const selected = shownWeek[selectedDay] ?? null;
  const pastDays = shownWeek.filter((d) => !d.isFuture && d.totalMs > 0);
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
          <Text style={styles.headerSub}>{weekLabel(weekOffset)}</Text>
          <Text style={styles.headerTitle}>Screen Time</Text>
          <View style={styles.avgPill}>
            <Text style={styles.avgText}>
              Daily avg · {formatDurationLong(weekAvgMs)}
            </Text>
          </View>
        </View>

        {/* ── Bar Chart (swipe lijevo/desno mijenja sedmicu) ── */}
        <View style={styles.card} {...panResponder.panHandlers}>
          {/* Navigacija sedmica: strelice za one koji ne swipe-uju */}
          <View style={styles.weekNav}>
            <TouchableOpacity
              style={styles.weekNavBtn}
              onPress={() => changeWeek(1)}
              disabled={weekOffset >= MAX_WEEKS_BACK}
              hitSlop={10}
            >
              <Text
                style={[
                  styles.weekNavArrow,
                  weekOffset >= MAX_WEEKS_BACK && styles.weekNavArrowDisabled,
                ]}
              >
                ‹
              </Text>
            </TouchableOpacity>
            <Text style={styles.weekNavLabel}>{weekLabel(weekOffset)}</Text>
            <TouchableOpacity
              style={styles.weekNavBtn}
              onPress={() => changeWeek(-1)}
              disabled={isCurrentWeek}
              hitSlop={10}
            >
              <Text
                style={[
                  styles.weekNavArrow,
                  isCurrentWeek && styles.weekNavArrowDisabled,
                ]}
              >
                ›
              </Text>
            </TouchableOpacity>
          </View>

          {isLoadingWeek ? (
            <View style={styles.chartLoading}>
              <ActivityIndicator color="#7C6EF5" />
            </View>
          ) : (
            <Animated.View
              style={[styles.chartRow, { transform: [{ translateX }] }]}
            >
              {shownWeek.map((day, i) => {
                const isSelected = i === selectedDay;
                const isToday = isCurrentWeek && i === todayIndex;
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
            </Animated.View>
          )}

          {!isLoadingWeek && !isCurrentWeek && !weekHasData && (
            <Text style={styles.noWeekDataText}>
              No usage data for this week. Android keeps detailed usage history
              only for a limited time.
            </Text>
          )}
        </View>

        {/* ── Selected Day Detail ── */}
        {!isLoadingWeek && selected && !selected.isFuture && (
          <View style={styles.detailSection}>
            <View style={styles.detailHeader}>
              <View>
                <Text style={styles.detailDay}>
                  {isCurrentWeek && selectedDay === todayIndex
                    ? "Today"
                    : selected.dayFull}
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
                shownWeek.forEach((d) =>
                  d.apps.forEach((a) => {
                    allApps[a.packageName] =
                      (allApps[a.packageName] || 0) + a.totalTimeMs;
                  }),
                );
                const top = Object.entries(allApps).sort(
                  (a, b) => b[1] - a[1],
                )[0];
                if (!top) return "—";
                const name = top[0].split(".").pop();
                return name.charAt(0).toUpperCase() + name.slice(1);
              })()}
            </Text>
            <Text style={styles.summarySubValue}>
              {weekLabel(weekOffset).toLowerCase()}
            </Text>
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
