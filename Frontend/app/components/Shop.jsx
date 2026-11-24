import React, {
  useState,
  useCallback,
  useMemo,
  useRef,
} from "react";
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  RefreshControl,
  Alert,
  StatusBar,
  Animated,
} from "react-native";
import * as SecureStore from "expo-secure-store";
import { useFocusEffect } from "@react-navigation/native";
import { useShop } from "@/contexts/ShopContext";
const API_URL = "https://shenika-ovarian-unpiratically.ngrok-free.dev";
import { shopStyles as styles } from "@/assets/styles/home.styles";

const RARITY = {
  common: {
    label: "COMMON",
    color: "#9CA3AF",
    glow: "#6B7280",
    border: "#374151",
    badge: "#1F2937",
    icon: "◆",
  },
  rare: {
    label: "RARE",
    color: "#60A5FA",
    glow: "#3B82F6",
    border: "#1D4ED8",
    badge: "#1E3A8A",
    icon: "◈",
  },
  legendary: {
    label: "LEGENDARY",
    color: "#F59E0B",
    glow: "#D97706",
    border: "#92400E",
    badge: "#78350F",
    icon: "✦",
  },
};

const FILTERS = ["all", "common", "rare", "legendary"];

// van komponente da bi bila ista funkcija pri svakom renderu
const keyExtractor = (item) => String(item.id);

const MessageCard = React.memo(({ item, onBuy }) => {
  const rarity = RARITY[item.rarity] || RARITY.common;
  const owned = !item.is_active;
  const scaleAnim = useRef(new Animated.Value(1)).current;

  const handlePressIn = () => {
    Animated.spring(scaleAnim, {
      toValue: 0.97,
      useNativeDriver: true,
      speed: 50,
    }).start();
  };

  const handlePressOut = () => {
    Animated.spring(scaleAnim, {
      toValue: 1,
      useNativeDriver: true,
      speed: 50,
    }).start();
  };

  return (
    <Animated.View style={{ transform: [{ scale: scaleAnim }] }}>
      <TouchableOpacity
        activeOpacity={0.85}
        onPressIn={handlePressIn}
        onPressOut={handlePressOut}
        onPress={() => !owned && onBuy(item)}
        disabled={owned}
        style={[
          styles.card,
          { borderColor: owned ? "#2D2D2D" : rarity.border },
          owned && styles.cardOwned,
        ]}
      >
        <View
          style={[
            styles.rarityBar,
            { backgroundColor: owned ? "#2D2D2D" : rarity.color },
          ]}
        />
        <View style={styles.cardHeader}>
          <View
            style={[
              styles.rarityBadge,
              { backgroundColor: owned ? "#1A1A1A" : rarity.badge },
            ]}
          >
            <Text
              style={[
                styles.rarityIcon,
                { color: owned ? "#444" : rarity.color },
              ]}
            >
              {rarity.icon}
            </Text>
            <Text
              style={[
                styles.rarityLabel,
                { color: owned ? "#444" : rarity.color },
              ]}
            >
              {rarity.label}
            </Text>
          </View>
          {owned && (
            <View style={styles.ownedBadge}>
              <Text style={styles.ownedText}>✓ OWNED</Text>
            </View>
          )}
        </View>
        <Text style={[styles.messageText, owned && styles.messageTextOwned]}>
          "{item.message}"
        </Text>
        <View style={styles.cardFooter}>
          {!owned ? (
            <TouchableOpacity
              style={[styles.buyButton, { borderColor: rarity.color }]}
              onPress={() => onBuy(item)}
            >
              <Text style={styles.coinIcon}>🪙</Text>
              <Text style={[styles.buyButtonText, { color: rarity.color }]}>
                {item.price.toLocaleString()}
              </Text>
            </TouchableOpacity>
          ) : (
            <View style={styles.buyButton}>
              <Text style={styles.ownedFooterText}>Already unlocked</Text>
            </View>
          )}
        </View>
      </TouchableOpacity>
    </Animated.View>
  );
});

export default function ShopScreen() {
  const { shopData, fetchShop: refetchShop } = useShop();
  const [activeFilter, setActiveFilter] = useState("all");
  const [refreshing, setRefreshing] = useState(false);
  const [userPoints, setUserPoints] = useState(null);

  const counts = useMemo(
    () => ({
      all: shopData.all.length,
      common: shopData.common.length,
      rare: shopData.rare.length,
      legendary: shopData.legendary.length,
    }),
    [shopData]
  );

  // Tiho osvjezavanje pri svakom ulasku u tab: stari podaci se prikazu odmah,
  // novi stignu u pozadini (bez spinnera)
  useFocusEffect(
    useCallback(() => {
      refetchShop();
    }, [refetchShop])
  );

  const handleRefresh = useCallback(async () => {
    setRefreshing(true);
    await refetchShop();
    setRefreshing(false);
  }, [refetchShop]);

  const confirmBuy = useCallback(
    async (item) => {
      try {
        const token = await SecureStore.getItemAsync("jwt_token");
        const res = await fetch(
          `${API_URL}/v1/notifications/purchase/${item.id}`,
          {
            method: "POST",
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "application/json",
            },
          }
        );

        if (!res.ok) {
          const err = await res.json();
          Alert.alert("Purchase Failed", err.error ?? "Something went wrong.");
          return;
        }

        Alert.alert("🎉 Purchased!", `"${item.message}" is now yours.`);
        await refetchShop();
      } catch (e) {
        Alert.alert("Error", e.message);
      }
    },
    [refetchShop]
  );

  const handleBuy = useCallback(
    (item) => {
      Alert.alert(
        "Purchase Message",
        `Buy "${item.message}" for ${item.price.toLocaleString()} points?`,
        [
          { text: "Cancel", style: "cancel" },
          { text: "Buy", style: "default", onPress: () => confirmBuy(item) },
        ]
      );
    },
    [confirmBuy]
  );

  const renderItemHandlers = useMemo(
    () =>
      Object.fromEntries(
        FILTERS.map((f) => [
          f,
          ({ item }) => <MessageCard item={item} onBuy={handleBuy} />,
        ])
      ),
    [handleBuy]
  );

  return (
    <View style={styles.root}>
      <StatusBar barStyle="light-content" backgroundColor="#0A0A0A" />

      <View style={styles.header}>
        <View>
          <Text style={styles.headerTitle}>Message Shop</Text>
          <Text style={styles.headerSub}>Unlock preset messages</Text>
        </View>
        {userPoints !== null && (
          <View style={styles.pointsChip}>
            <Text style={styles.coinEmoji}>🪙</Text>
            <Text style={styles.pointsText}>{userPoints.toLocaleString()}</Text>
          </View>
        )}
      </View>

      <View style={styles.filterRow}>
        {FILTERS.map((f) => {
          const r = RARITY[f];
          const isActive = activeFilter === f;
          return (
            <TouchableOpacity
              key={f}
              style={[
                styles.filterTab,
                isActive && {
                  backgroundColor: f === "all" ? "#1F2937" : r.badge,
                  borderColor: f === "all" ? "#6B7280" : r.color,
                },
              ]}
              onPress={() => setActiveFilter(f)}
            >
              {f !== "all" && (
                <Text
                  style={[
                    styles.filterIcon,
                    { color: isActive ? r.color : "#555" },
                  ]}
                >
                  {r.icon}
                </Text>
              )}
              <Text
                style={[
                  styles.filterLabel,
                  isActive && { color: f === "all" ? "#E5E7EB" : r.color },
                ]}
              >
                {f === "all" ? "All" : r.label}
              </Text>
              <View
                style={[
                  styles.filterCount,
                  isActive && {
                    backgroundColor: f === "all" ? "#374151" : r.badge,
                  },
                ]}
              >
                <Text
                  style={[
                    styles.filterCountText,
                    isActive && { color: f === "all" ? "#9CA3AF" : r.color },
                  ]}
                >
                  {counts[f]}
                </Text>
              </View>
            </TouchableOpacity>
          );
        })}
      </View>

      <View style={listStyles.container}>
        {FILTERS.map((f) => {
          const data = shopData[f];
          const isActive = activeFilter === f;
          return (
            <View
              key={f}
              style={[listStyles.pane, { opacity: isActive ? 1 : 0 }]}
              pointerEvents={isActive ? "auto" : "none"}
            >
              {data.length === 0 ? (
                <View style={styles.emptyState}>
                  <Text style={styles.emptyIcon}>🛒</Text>
                  <Text style={styles.emptyTitle}>Nothing here</Text>
                  <Text style={styles.emptySub}>
                    No {f !== "all" ? RARITY[f]?.label.toLowerCase() : ""}{" "}
                    messages available.
                  </Text>
                </View>
              ) : (
                <FlatList
                  data={data}
                  keyExtractor={keyExtractor}
                  renderItem={renderItemHandlers[f]}
                  contentContainerStyle={styles.list}
                  showsVerticalScrollIndicator={false}
                  initialNumToRender={10}
                  maxToRenderPerBatch={10}
                  removeClippedSubviews={false}
                  windowSize={5}
                  // RefreshControl mora uvijek postojati: na Androidu njegovo
                  // dodavanje/uklanjanje ponovo montira cijelu listu
                  refreshControl={
                    <RefreshControl
                      refreshing={refreshing}
                      onRefresh={handleRefresh}
                      tintColor="#F59E0B"
                    />
                  }
                />
              )}
            </View>
          );
        })}
      </View>
    </View>
  );
}

const listStyles = StyleSheet.create({
  container: {
    flex: 1,
    position: "relative",
  },
  pane: {
    ...StyleSheet.absoluteFillObject,
  },
});
