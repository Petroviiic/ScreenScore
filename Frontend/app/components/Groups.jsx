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
import { useAppData } from "@/contexts/AppDataContext";
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

// ID ulogovanog korisnika citamo iz JWT tokena ("sub" claim), jer backend
// nema posebnu rutu koja vraca trenutnog korisnika.
async function getCurrentUserId() {
  try {
    const token = await SecureStore.getItemAsync("jwt_token");
    if (!token) return null;
    const payload = token.split(".")[1].replace(/-/g, "+").replace(/_/g, "/");
    const padded = payload + "=".repeat((4 - (payload.length % 4)) % 4);
    const sub = JSON.parse(atob(padded)).sub;
    return sub != null ? Number(sub) : null;
  } catch {
    return null;
  }
}

function todayDateString() {
  const d = new Date();
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const dd = String(d.getDate()).padStart(2, "0");
  return `${d.getFullYear()}-${mm}-${dd}`;
}

// Backend nema rutu za listu clanova grupe, pa clanove uzimamo iz
// get-group-stats. Ta ruta vraca samo clanove koji imaju statistiku za taj dan.
async function fetchGroupMembers(groupId) {
  const res = await authFetch(`/v1/stats/get-group-stats`, {
    method: "POST",
    body: JSON.stringify({
      group_id: groupId,
      desired_date: todayDateString(),
    }),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data?.error || `Status ${res.status}`);
  return Array.isArray(data.data) ? data.data : [];
}

// ─── Validacija (ista pravila kao na backendu: groups.go isValidString) ──────
const GROUP_NAME_MAX = 15;
const INVITE_CODE_MAX = 10;
const ALPHANUMERIC = /^[a-zA-Z0-9]+$/;

function validateGroupName(name) {
  if (name.length < 3 || name.length > GROUP_NAME_MAX)
    return `Group name must be 3–${GROUP_NAME_MAX} characters long.`;
  if (!ALPHANUMERIC.test(name))
    return "Group name can only contain letters (a–z) and numbers, without spaces.";
  return null;
}

function validateInviteCode(code) {
  if (
    code.length < 3 ||
    code.length > INVITE_CODE_MAX ||
    !ALPHANUMERIC.test(code)
  )
    return "Invalid invite code. Check the code and try again.";
  return null;
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
function MemberRow({ member, currentUserId, onKick }) {
  const isMe = member.id === currentUserId;
  // Backend dozvoljava svakom clanu da izbaci drugog (grupe nemaju vlasnika),
  // pa Kick prikazujemo kod svih osim kod samog sebe.
  const canKick = !isMe && currentUserId != null;

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
        {isMe && (
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
  const [members, setMembers] = useState(null);
  const [membersLoading, setMembersLoading] = useState(false);
  const [membersError, setMembersError] = useState(null);

  const loadMembers = useCallback(async () => {
    setMembersLoading(true);
    setMembersError(null);
    try {
      setMembers(await fetchGroupMembers(group.id));
    } catch (err) {
      setMembersError(err.message);
    } finally {
      setMembersLoading(false);
    }
  }, [group.id]);

  // Clanove ucitavamo tek kad se kartica otvori
  useEffect(() => {
    if (expanded) loadMembers();
  }, [expanded, loadMembers]);

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
            {members ? `${members.length} active today` : "Tap to see members"}
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
          {membersLoading && !members ? (
            <ActivityIndicator color="#7C6EF5" style={{ marginVertical: 12 }} />
          ) : membersError ? (
            <Text style={styles.groupMeta}>
              Could not load members: {membersError}
            </Text>
          ) : members && members.length === 0 ? (
            <Text style={styles.groupMeta}>
              No members have synced screen time today.
            </Text>
          ) : (
            (members || []).map((member) => (
              <MemberRow
                key={member.id}
                member={member}
                currentUserId={currentUserId}
                onKick={(m) => onKick(group, m, loadMembers)}
              />
            ))
          )}
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
    const success = await onConfirm(value.trim());
    setLoading(false);
    // brisi unos samo ako je uspjelo, da korisnik moze ispraviti gresku
    if (success) setValue("");
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
  const { groupsData, refetchGroups } = useAppData();
  const [groups, setGroups] = useState(groupsData);
  const [refreshing, setRefreshing] = useState(false);
  const [currentUserId, setCurrentUserId] = useState(null);

  useEffect(() => {
    getCurrentUserId().then(setCurrentUserId);
  }, []);

  // Modals
  const [showCreate, setShowCreate] = useState(false);
  const [showJoin, setShowJoin] = useState(false);

  // ── Fetch groups ─────────────────────────────────────────────────────────────

  useEffect(() => {
    setGroups(groupsData);
  }, [groupsData]);

  const handleRefresh = async () => {
    setRefreshing(true);
    await refetchGroups();
    setRefreshing(false);
  };
  // ── Create ───────────────────────────────────────────────────────────────────
  const handleCreate = async (name) => {
    const validationError = validateGroupName(name);
    if (validationError) {
      Alert.alert("Invalid name", validationError);
      return;
    }
    try {
      const res = await authFetch(
        `/v1/groups/create/${encodeURIComponent(name)}`,
        { method: "POST" },
      );
      if (res.ok) {
        setShowCreate(false);
        refetchGroups();
        return true;
      } else {
        const data = await res.json().catch(() => ({}));
        Alert.alert("Error", data?.error || `Status ${res.status}`);
      }
    } catch (err) {
      Alert.alert("Error", err.message);
    }
  };

  // ── Join ─────────────────────────────────────────────────────────────────────
  const handleJoin = async (code) => {
    const validationError = validateInviteCode(code);
    if (validationError) {
      Alert.alert("Invalid code", validationError);
      return;
    }
    try {
      const res = await authFetch(
        `/v1/groups/join/${encodeURIComponent(code)}`,
        { method: "POST" },
      );
      if (res.ok) {
        setShowJoin(false);
        refetchGroups();
        return true;
      } else {
        const data = await res.json().catch(() => ({}));
        Alert.alert("Error", data?.error || `Status ${res.status}`);
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
                { method: "POST" },
              );
              if (res.ok) refetchGroups();
              else {
                const data = await res.json().catch(() => ({}));
                Alert.alert("Error", data?.error || `Status ${res.status}`);
              }
            } catch (err) {
              Alert.alert("Error", err.message);
            }
          },
        },
      ],
    );
  };

  // ── Kick ─────────────────────────────────────────────────────────────────────
  // onDone: ponovo ucitaj clanove te kartice nakon uspjesnog izbacivanja
  const handleKick = (group, member, onDone) => {
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
              if (res.ok) onDone?.();
              else {
                const data = await res.json().catch(() => ({}));
                Alert.alert("Error", data?.error || `Status ${res.status}`);
              }
            } catch (err) {
              Alert.alert("Error", err.message);
            }
          },
        },
      ],
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
      {
        <ScrollView
          contentContainerStyle={[
            styles.scroll,
            groups.length === 0 && styles.scrollCenter,
          ]}
          showsVerticalScrollIndicator={false}
          refreshControl={
            <RefreshControl
              refreshing={refreshing}
              onRefresh={() => handleRefresh()}
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
      }

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
