import React, { useState, useEffect, useCallback } from "react";
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
  Modal,
  RefreshControl,
  Animated,
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

// ─── Icons (simple SVG-style unicode replacements) ────────────────────────────
const Icon = ({ name, size = 16, color = "#fff" }) => {
  const icons = {
    plus: "+",
    users: "👥",
    crown: "👑",
    copy: "⧉",
    leave: "→",
    kick: "✕",
    close: "✕",
    refresh: "↻",
  };
  return (
    <Text style={{ fontSize: size, color, lineHeight: size + 4 }}>
      {icons[name] || "?"}
    </Text>
  );
};

// ─── Empty state ──────────────────────────────────────────────────────────────
function EmptyGroups({ onJoin, onCreate }) {
  return (
    <View style={styles.emptyContainer}>
      <Text style={styles.emptyIcon}>🏠</Text>
      <Text style={styles.emptyTitle}>No groups yet</Text>
      <Text style={styles.emptySubtitle}>
        Create a group or join one with an invite code
      </Text>
      <View style={styles.emptyActions}>
        <TouchableOpacity style={styles.emptyBtn} onPress={onCreate}>
          <Text style={styles.emptyBtnText}>Create group</Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[styles.emptyBtn, styles.emptyBtnOutline]}
          onPress={onJoin}
        >
          <Text style={[styles.emptyBtnText, { color: "#7C6EF5" }]}>
            Join with code
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}

// ─── Member row ───────────────────────────────────────────────────────────────
function MemberRow({ member, isOwner, currentUserId, onKick }) {
  const isMe = member.id === currentUserId;
  const canKick = isOwner && !isMe;

  return (
    <View style={styles.memberRow}>
      <View style={styles.memberAvatar}>
        <Text style={styles.memberAvatarText}>
          {(member.username || member.email || "?")[0].toUpperCase()}
        </Text>
      </View>
      <View style={styles.memberInfo}>
        <Text style={styles.memberName}>
          {member.username || member.email || `User #${member.id}`}
        </Text>
        {member.is_owner && (
          <View style={styles.ownerBadge}>
            <Text style={styles.ownerBadgeText}>owner</Text>
          </View>
        )}
        {isMe && !member.is_owner && (
          <View style={[styles.ownerBadge, styles.meBadge]}>
            <Text style={styles.ownerBadgeText}>you</Text>
          </View>
        )}
      </View>
      {canKick && (
        <TouchableOpacity style={styles.kickBtn} onPress={() => onKick(member)}>
          <Text style={styles.kickBtnText}>Kick</Text>
        </TouchableOpacity>
      )}
    </View>
  );
}

// ─── Group card ───────────────────────────────────────────────────────────────
function GroupCard({ group, currentUserId, onLeave, onKick }) {
  const [expanded, setExpanded] = useState(false);
  const isOwner = group.owner_id === currentUserId;

  const copyInviteCode = () => {
    Alert.alert("Invite Code", group.invite_code, [{ text: "OK" }]);
  };

  return (
    <View style={styles.groupCard}>
      {/* Card header */}
      <TouchableOpacity
        style={styles.groupCardHeader}
        onPress={() => setExpanded((v) => !v)}
        activeOpacity={0.75}
      >
        <View style={styles.groupIconWrap}>
          <Text style={styles.groupIconText}>
            {(group.name || "G")[0].toUpperCase()}
          </Text>
        </View>
        <View style={styles.groupHeaderInfo}>
          <Text style={styles.groupName}>{group.name}</Text>
          <Text style={styles.groupMeta}>
            {group.members?.length ?? 0} member
            {(group.members?.length ?? 0) !== 1 ? "s" : ""}
            {isOwner ? " · owner" : ""}
          </Text>
        </View>
        <Text style={[styles.chevron, expanded && styles.chevronOpen]}>›</Text>
      </TouchableOpacity>

      {/* Invite code row */}
      <TouchableOpacity style={styles.inviteRow} onPress={copyInviteCode}>
        <Text style={styles.inviteLabel}>Invite code</Text>
        <Text style={styles.inviteCode}>{group.invite_code}</Text>
        <Text style={styles.inviteCopy}>tap to copy</Text>
      </TouchableOpacity>

      {/* Members list (expanded) */}
      {expanded && (
        <View style={styles.membersList}>
          <View style={styles.membersDivider} />
          {(group.members || []).map((member) => (
            <MemberRow
              key={member.id}
              member={member}
              isOwner={isOwner}
              currentUserId={currentUserId}
              onKick={(m) => onKick(group, m)}
            />
          ))}
        </View>
      )}

      {/* Leave button */}
      <TouchableOpacity style={styles.leaveBtn} onPress={() => onLeave(group)}>
        <Text style={styles.leaveBtnText}>Leave group</Text>
      </TouchableOpacity>
    </View>
  );
}

// ─── Modal ────────────────────────────────────────────────────────────────────
function ActionModal({
  visible,
  title,
  placeholder,
  onConfirm,
  onClose,
  confirmLabel = "Confirm",
  danger = false,
}) {
  const [value, setValue] = useState("");
  const [loading, setLoading] = useState(false);

  const handleConfirm = async () => {
    if (!value.trim()) return;
    setLoading(true);
    await onConfirm(value.trim());
    setLoading(false);
    setValue("");
  };

  const handleClose = () => {
    setValue("");
    onClose();
  };

  return (
    <Modal visible={visible} transparent animationType="fade">
      <View style={styles.modalOverlay}>
        <View style={styles.modalBox}>
          <View style={styles.modalHeader}>
            <Text style={styles.modalTitle}>{title}</Text>
            <TouchableOpacity onPress={handleClose} style={styles.modalClose}>
              <Text style={styles.modalCloseText}>✕</Text>
            </TouchableOpacity>
          </View>
          <TextInput
            style={styles.modalInput}
            placeholder={placeholder}
            placeholderTextColor="#555"
            value={value}
            onChangeText={setValue}
            autoFocus
          />
          <TouchableOpacity
            style={[styles.modalBtn, danger && styles.modalBtnDanger]}
            onPress={handleConfirm}
            disabled={loading}
          >
            {loading ? (
              <ActivityIndicator color="#fff" size="small" />
            ) : (
              <Text style={styles.modalBtnText}>{confirmLabel}</Text>
            )}
          </TouchableOpacity>
        </View>
      </View>
    </Modal>
  );
}

// ─── Main Component ───────────────────────────────────────────────────────────
export default function Groups() {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [currentUserId, setCurrentUserId] = useState(null);

  // Modals
  const [showCreate, setShowCreate] = useState(false);
  const [showJoin, setShowJoin] = useState(false);

  // ── Fetch groups ─────────────────────────────────────────────────────────────

  const fetchGroups = useCallback(async (isRefresh = false) => {
    if (isRefresh) setRefreshing(true);
    else setLoading(true);
    try {
      const res = await authFetch("/v1/groups/get_user_groups");
      if (res.ok) {
        const data = await res.json();
        const arr = Array.isArray(data.data) ? data.data : [];
        setGroups(arr);
      }
    } catch (err) {
      console.log("fetch groups error:", err);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  useEffect(() => {
    fetchGroups();
  }, []);

  // ── Create ───────────────────────────────────────────────────────────────────
  const handleCreate = async (name) => {
    try {
      const res = await authFetch(
        `/v1/groups/create/${encodeURIComponent(name)}`,
        { method: "POST" }
      );
      if (res.ok) {
        setShowCreate(false);
        fetchGroups();
      } else {
        const data = await res.json().catch(() => ({}));
        Alert.alert("Error", data?.message || `Status ${res.status}`);
      }
    } catch (err) {
      Alert.alert("Error", err.message);
    }
  };

  // ── Join ─────────────────────────────────────────────────────────────────────
  const handleJoin = async (code) => {
    try {
      const res = await authFetch(
        `/v1/groups/join/${encodeURIComponent(code)}`,
        { method: "POST" }
      );
      if (res.ok) {
        setShowJoin(false);
        fetchGroups();
      } else {
        const data = await res.json().catch(() => ({}));
        Alert.alert("Error", data?.message || `Status ${res.status}`);
      }
    } catch (err) {
      Alert.alert("Error", err.message);
    }
  };

  // ── Leave ────────────────────────────────────────────────────────────────────

  const handleLeave = (group) => {
    Alert.alert(
      "Leave group",
      `Are you sure you want to leave "${group.name}"?`,
      [
        { text: "Cancel", style: "cancel" },
        {
          text: "Leave",
          style: "destructive",
          onPress: async () => {
            try {
              const res = await authFetch(
                `/v1/groups/leave/${encodeURIComponent(group.id)}`,
                { method: "POST" }
              );
              if (res.ok) fetchGroups();
              else {
                const data = await res.json().catch(() => ({}));
                Alert.alert("Error", data?.message || `Status ${res.status}`);
              }
            } catch (err) {
              Alert.alert("Error", err.message);
            }
          },
        },
      ]
    );
  };

  // ── Kick ─────────────────────────────────────────────────────────────────────
  const handleKick = (group, member) => {
    Alert.alert(
      "Kick member",
      `Remove ${
        member.username || member.email || `User #${member.id}`
      } from "${group.name}"?`,
      [
        { text: "Cancel", style: "cancel" },
        {
          text: "Kick",
          style: "destructive",
          onPress: async () => {
            try {
              const res = await authFetch(`/v1/groups/kick`, {
                method: "POST",
                body: JSON.stringify({
                  group_id: group.id,
                  user_to_kick_id: member.id,
                }),
              });
              if (res.ok) fetchGroups();
              else {
                const data = await res.json().catch(() => ({}));
                Alert.alert("Error", data?.message || `Status ${res.status}`);
              }
            } catch (err) {
              Alert.alert("Error", err.message);
            }
          },
        },
      ]
    );
  };

  // ── UI ───────────────────────────────────────────────────────────────────────
  return (
    <View style={styles.root}>
      <StatusBar barStyle="light-content" />

      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Groups</Text>
        <View style={styles.headerActions}>
          <TouchableOpacity
            style={styles.headerBtn}
            onPress={() => setShowJoin(true)}
          >
            <Text style={styles.headerBtnText}>Join</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.headerBtn, styles.headerBtnPrimary]}
            onPress={() => setShowCreate(true)}
          >
            <Text style={[styles.headerBtnText, { color: "#fff" }]}>+ New</Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Content */}
      {loading ? (
        <View style={styles.loadingContainer}>
          <ActivityIndicator color="#7C6EF5" size="large" />
        </View>
      ) : (
        <ScrollView
          contentContainerStyle={[
            styles.scroll,
            groups.length === 0 && styles.scrollCenter,
          ]}
          showsVerticalScrollIndicator={false}
          refreshControl={
            <RefreshControl
              refreshing={refreshing}
              onRefresh={() => fetchGroups(true)}
              tintColor="#7C6EF5"
            />
          }
        >
          {groups.length === 0 ? (
            <EmptyGroups
              onCreate={() => setShowCreate(true)}
              onJoin={() => setShowJoin(true)}
            />
          ) : (
            groups.map((group) => (
              <GroupCard
                key={group.id}
                group={group}
                currentUserId={currentUserId}
                onLeave={handleLeave}
                onKick={handleKick}
              />
            ))
          )}
        </ScrollView>
      )}

      {/* Create modal */}
      <ActionModal
        visible={showCreate}
        title="Create group"
        placeholder="Group name"
        confirmLabel="Create"
        onConfirm={handleCreate}
        onClose={() => setShowCreate(false)}
      />

      {/* Join modal */}
      <ActionModal
        visible={showJoin}
        title="Join a group"
        placeholder="Invite code"
        confirmLabel="Join"
        onConfirm={handleJoin}
        onClose={() => setShowJoin(false)}
      />
    </View>
  );
}
