import React, { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  StatusBar,
  StyleSheet,
  Alert,
} from "react-native";
import * as SecureStore from "expo-secure-store";
import { groupsStyles as styles } from "@/assets/styles/home.styles";

const API_URL = "https://shenika-ovarian-unpiratically.ngrok-free.dev";

async function authFetch(path, options = {}) {
  const token = await SecureStore.getItemAsync("jwt_token");
  return fetch(`${API_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...(options.headers || {}),
    },
  });
}

// ─── Section wrapper ──────────────────────────────────────────────────────────
function Section({ title, children }) {
  return (
    <View style={styles.section}>
      <Text style={styles.sectionTitle}>{title}</Text>
      {children}
    </View>
  );
}

// ─── Main Component ───────────────────────────────────────────────────────────
export default function Groups() {
  // Create group
  const [groupName, setGroupName] = useState("");
  const [creating, setCreating] = useState(false);

  // Join group
  const [inviteCode, setInviteCode] = useState("");
  const [joining, setJoining] = useState(false);

  // Kick user
  const [kickGroupId, setKickGroupId] = useState("");
  const [kickUserId, setKickUserId] = useState("");
  const [kicking, setKicking] = useState(false);

  // Leave group
  const [leaveGroupId, setLeaveGroupId] = useState("");
  const [leaving, setLeaving] = useState(false);

  // ── Create ──────────────────────────────────────────────────────────────────
  const handleCreate = async () => {
    if (!groupName.trim()) return Alert.alert("Error", "Enter a group name.");
    setCreating(true);
    try {
      const res = await authFetch(
        `/v1/groups/create/${encodeURIComponent(groupName.trim())}`,
        {
          method: "POST",
        }
      );
      const data = await res.json().catch(() => ({}));
      if (res.ok) {
        Alert.alert("Success", `Group "${groupName}" created!`);
        setGroupName("");
      } else {
        Alert.alert("Error", data?.message || `Status ${res.status}`);
      }
    } catch (err) {
      Alert.alert("Error", err.message);
    } finally {
      setCreating(false);
    }
  };

  // ── Join ────────────────────────────────────────────────────────────────────
  const handleJoin = async () => {
    if (!inviteCode.trim())
      return Alert.alert("Error", "Enter an invite code.");
    setJoining(true);
    try {
      const res = await authFetch(
        `/v1/groups/join/${encodeURIComponent(inviteCode.trim())}`,
        {
          method: "POST",
        }
      );
      const data = await res.json().catch(() => ({}));
      if (res.ok) {
        Alert.alert("Success", "Joined group!");
        setInviteCode("");
      } else {
        Alert.alert("Error", data?.message || `Status ${res.status}`);
      }
    } catch (err) {
      Alert.alert("Error", err.message);
    } finally {
      setJoining(false);
    }
  };

  // ── Kick ────────────────────────────────────────────────────────────────────
  const handleKick = async () => {
    if (!kickGroupId.trim() || !kickUserId.trim())
      return Alert.alert("Error", "Fill in both Group ID and User ID.");
    const userId = parseInt(kickUserId, 10);
    if (isNaN(userId)) return Alert.alert("Error", "User ID must be a number.");
    setKicking(true);
    try {
      const res = await authFetch(`/v1/groups/kick`, {
        method: "POST",
        body: JSON.stringify({
          group_id: kickGroupId.trim(),
          user_to_kick_id: userId,
        }),
      });
      const data = await res.json().catch(() => ({}));
      if (res.ok) {
        Alert.alert("Success", "User kicked.");
        setKickGroupId("");
        setKickUserId("");
      } else {
        Alert.alert("Error", data?.message || `Status ${res.status}`);
      }
    } catch (err) {
      Alert.alert("Error", err.message);
    } finally {
      setKicking(false);
    }
  };

  // ── Leave ───────────────────────────────────────────────────────────────────
  const handleLeave = async () => {
    if (!leaveGroupId.trim()) return Alert.alert("Error", "Enter a Group ID.");
    setLeaving(true);
    try {
      const res = await authFetch(
        `/v1/groups/leave/${encodeURIComponent(leaveGroupId.trim())}`,
        {
          method: "POST",
        }
      );
      const data = await res.json().catch(() => ({}));
      if (res.ok) {
        Alert.alert("Success", "Left group.");
        setLeaveGroupId("");
      } else {
        Alert.alert("Error", data?.message || `Status ${res.status}`);
      }
    } catch (err) {
      Alert.alert("Error", err.message);
    } finally {
      setLeaving(false);
    }
  };

  // ── UI ──────────────────────────────────────────────────────────────────────
  return (
    <View style={styles.root}>
      <StatusBar barStyle="light-content" />
      <ScrollView
        contentContainerStyle={styles.scroll}
        showsVerticalScrollIndicator={false}
      >
        <Text style={styles.pageTitle}>Groups</Text>

        {/* ── Create ── */}
        <Section title="Create Group">
          <TextInput
            style={styles.input}
            placeholder="Group name"
            placeholderTextColor="#666"
            value={groupName}
            onChangeText={setGroupName}
          />
          <TouchableOpacity
            style={styles.btn}
            onPress={handleCreate}
            disabled={creating}
          >
            {creating ? (
              <ActivityIndicator color="#fff" />
            ) : (
              <Text style={styles.btnText}>Create</Text>
            )}
          </TouchableOpacity>
        </Section>

        {/* ── Join ── */}
        <Section title="Join via Invite Code">
          <TextInput
            style={styles.input}
            placeholder="Invite code"
            placeholderTextColor="#666"
            value={inviteCode}
            onChangeText={setInviteCode}
            autoCapitalize="none"
          />
          <TouchableOpacity
            style={styles.btn}
            onPress={handleJoin}
            disabled={joining}
          >
            {joining ? (
              <ActivityIndicator color="#fff" />
            ) : (
              <Text style={styles.btnText}>Join</Text>
            )}
          </TouchableOpacity>
        </Section>

        {/* ── Kick ── */}
        <Section title="Kick Member">
          <TextInput
            style={styles.input}
            placeholder="Group ID"
            placeholderTextColor="#666"
            value={kickGroupId}
            onChangeText={setKickGroupId}
            autoCapitalize="none"
          />
          <TextInput
            style={[styles.input, { marginTop: 8 }]}
            placeholder="User ID (number)"
            placeholderTextColor="#666"
            value={kickUserId}
            onChangeText={setKickUserId}
            keyboardType="numeric"
          />
          <TouchableOpacity
            style={[styles.btn, styles.btnDanger]}
            onPress={handleKick}
            disabled={kicking}
          >
            {kicking ? (
              <ActivityIndicator color="#fff" />
            ) : (
              <Text style={styles.btnText}>Kick</Text>
            )}
          </TouchableOpacity>
        </Section>

        {/* ── Leave ── */}
        <Section title="Leave Group">
          <TextInput
            style={styles.input}
            placeholder="Group ID"
            placeholderTextColor="#666"
            value={leaveGroupId}
            onChangeText={setLeaveGroupId}
            autoCapitalize="none"
          />
          <TouchableOpacity
            style={[styles.btn, styles.btnDanger]}
            onPress={handleLeave}
            disabled={leaving}
          >
            {leaving ? (
              <ActivityIndicator color="#fff" />
            ) : (
              <Text style={styles.btnText}>Leave</Text>
            )}
          </TouchableOpacity>
        </Section>
      </ScrollView>
    </View>
  );
}
