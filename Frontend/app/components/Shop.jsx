import React, { useState, useCallback, useMemo, useRef } from "react";
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

// badge = boja rijetkosti sa providnoscu (isti stil kao ikonice u Groups tabu)
const RARITY = {
  common: {
    label: "COMMON",
    name: "Common",
    color: "#9CA3AF",
    badge: "#9CA3AF22",
    icon: "◆",
  },
  rare: {
    label: "RARE",
    name: "Rare",
    color: "#60A5FA",
    badge: "#60A5FA22",
    icon: "◈",
  },
  legendary: {
    label: "LEGENDARY",
    name: "Legendary",
    color: "#F59E0B",
    badge: "#F59E0B22",
    icon: "✦",
  },
};

const ALL_FILTER = { name: "All", color: "#7C6EF5", badge: "#7C6EF526" };
const FILTERS = ["all", "common", "rare", "legendary"];
const OWNED_GREY = "#33334A";

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
        style={[styles.card, owned && styles.cardOwned]}
      >
        <View
          style={[
            styles.rarityBar,
            { backgroundColor: owned ? OWNED_GREY : rarity.color },
          ]}
        />
        <View style={styles.cardHeader}>
          <View
            style={[
              styles.rarityBadge,
              { backgroundColor: owned ? "#22222E" : rarity.badge },
            ]}
          >
            <Text
              style={[
                styles.rarityIcon,
                { color: owned ? OWNED_GREY : rarity.color },
              ]}
            >
              {rarity.icon}
            </Text>
            <Text
              style={[
                styles.rarityLabel,
                { color: owned ? OWNED_GREY : rarity.color },
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
              style={styles.buyButton}
              onPress={() => onBuy(item)}
            >
              <Text style={styles.coinIcon}>🪙</Text>
              <Text style={styles.buyButtonText}>
                {item.price.toLocaleString()}
              </Text>
            </TouchableOpacity>
          ) : (
            <View style={[styles.buyButton, styles.buyButtonOwned]}>
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
    [shopData],
  );

  // Tiho osvjezavanje pri svakom ulasku u tab: stari podaci se prikazu odmah,
  // novi stignu u pozadini (bez spinnera)
  useFocusEffect(
    useCallback(() => {
      refetchShop();
    }, [refetchShop]),
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
          },
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
    [refetchShop],
  );

  const handleBuy = useCallback(
    (item) => {
      Alert.alert(
        "Purchase Message",
        `Buy "${item.message}" for ${item.price.toLocaleString()} points?`,
        [
          { text: "Cancel", style: "cancel" },
          { text: "Buy", style: "default", onPress: () => confirmBuy(item) },
        ],
      );
    },
    [confirmBuy],
  );

  const renderItemHandlers = useMemo(
    () =>
      Object.fromEntries(
        FILTERS.map((f) => [
          f,
          ({ item }) => <MessageCard item={item} onBuy={handleBuy} />,
        ]),
      ),
    [handleBuy],
  );

  return (
    <View style={styles.root}>
      <StatusBar barStyle="light-content" backgroundColor="#0F0F14" />

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
          const r = f === "all" ? ALL_FILTER : RARITY[f];
          const isActive = activeFilter === f;
          return (
            <TouchableOpacity
              key={f}
              style={[
                styles.filterTab,
                isActive && { backgroundColor: r.badge },
              ]}
              onPress={() => setActiveFilter(f)}
            >
              {/* Naziv u jednom redu, broj ispod njega, da sve stane i na uskom ekranu */}
              <Text
                numberOfLines={1}
                adjustsFontSizeToFit
                style={[styles.filterLabel, isActive && { color: r.color }]}
              >
                {r.name}
              </Text>
              <Text
                numberOfLines={1}
                style={[styles.filterCountText, isActive && { color: r.color }]}
              >
                {counts[f]}
              </Text>
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
                      tintColor="#7C6EF5"
                      colors={["#7C6EF5"]}
                      progressBackgroundColor="#1A1A24"
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
