import { StyleSheet } from "react-native";

export const styles = StyleSheet.create({
  root: {
    flex: 1,
    backgroundColor: "#0D0D14",
  },
  scroll: {
    paddingBottom: 40,
  },

  // ── Permission ──
  permissionContainer: {
    flex: 1,
    backgroundColor: "#0D0D14",
    alignItems: "center",
    justifyContent: "center",
    paddingHorizontal: 36,
  },
  permissionIcon: { fontSize: 52, marginBottom: 20 },
  permissionTitle: {
    color: "#FFFFFF",
    fontSize: 24,
    fontWeight: "700",
    marginBottom: 12,
    letterSpacing: 0.3,
  },
  permissionBody: {
    color: "#9090A8",
    fontSize: 15,
    textAlign: "center",
    lineHeight: 22,
    marginBottom: 32,
  },
  permissionBtn: {
    backgroundColor: "#7C6EF5",
    paddingHorizontal: 36,
    paddingVertical: 14,
    borderRadius: 100,
  },
  permissionBtnText: {
    color: "#FFFFFF",
    fontSize: 16,
    fontWeight: "700",
    letterSpacing: 0.4,
  },

  // ── Loading ──
  loadingContainer: {
    flex: 1,
    backgroundColor: "#0D0D14",
    alignItems: "center",
    justifyContent: "center",
    gap: 16,
  },
  loadingText: {
    color: "#9090A8",
    fontSize: 14,
    marginTop: 8,
  },

  // ── Header ──
  header: {
    paddingTop: 64,
    paddingHorizontal: 24,
    paddingBottom: 24,
  },
  headerSub: {
    color: "#7C6EF5",
    fontSize: 13,
    fontWeight: "600",
    letterSpacing: 1.8,
    textTransform: "uppercase",
    marginBottom: 6,
  },
  headerTitle: {
    color: "#FFFFFF",
    fontSize: 34,
    fontWeight: "800",
    letterSpacing: -0.5,
    marginBottom: 14,
  },
  avgPill: {
    alignSelf: "flex-start",
    backgroundColor: "#1A1A28",
    borderRadius: 100,
    paddingHorizontal: 14,
    paddingVertical: 6,
    borderWidth: 1,
    borderColor: "#2A2A3D",
  },
  avgText: {
    color: "#9090A8",
    fontSize: 13,
    fontWeight: "500",
  },

  // ── Chart Card ──
  card: {
    marginHorizontal: 20,
    backgroundColor: "#13131F",
    borderRadius: 20,
    padding: 20,
    borderWidth: 1,
    borderColor: "#1F1F2E",
    marginBottom: 20,
  },
  chartRow: {
    flexDirection: "row",
    alignItems: "flex-end",
    justifyContent: "space-between",
    paddingTop: 24,
  },
  barCol: {
    alignItems: "center",
    flex: 1,
  },
  barValue: {
    color: "transparent",
    fontSize: 9,
    fontWeight: "700",
    marginBottom: 4,
    height: 14,
  },
  barValueSelected: {
    color: "#7C6EF5",
  },
  barTrack: {
    width: 26,
    height: 130,
    backgroundColor: "#1F1F30",
    borderRadius: 8,
    justifyContent: "flex-end",
    overflow: "hidden",
  },
  bar: {
    width: "100%",
    backgroundColor: "#2D2D48",
    borderRadius: 8,
  },
  barSelected: {
    backgroundColor: "#7C6EF5",
  },
  barToday: {
    backgroundColor: "#3D3D60",
  },
  barLabel: {
    color: "#5A5A78",
    fontSize: 11,
    fontWeight: "600",
    marginTop: 8,
  },
  barLabelSelected: {
    color: "#7C6EF5",
  },
  todayDot: {
    width: 5,
    height: 5,
    borderRadius: 3,
    backgroundColor: "#7C6EF5",
    marginTop: 4,
  },

  // ── Detail Section ──
  detailSection: {
    marginHorizontal: 20,
    marginBottom: 20,
    backgroundColor: "#13131F",
    borderRadius: 20,
    padding: 20,
    borderWidth: 1,
    borderColor: "#1F1F2E",
  },
  detailHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: 20,
  },
  detailDay: {
    color: "#9090A8",
    fontSize: 13,
    fontWeight: "600",
    letterSpacing: 0.5,
    textTransform: "uppercase",
    marginBottom: 4,
  },
  detailTotal: {
    color: "#FFFFFF",
    fontSize: 28,
    fontWeight: "800",
    letterSpacing: -0.5,
  },
  badge: {
    backgroundColor: "#3A1A1A",
    borderRadius: 100,
    paddingHorizontal: 12,
    paddingVertical: 6,
  },
  badgeText: {
    color: "#F5866E",
    fontSize: 12,
    fontWeight: "700",
  },
  badgeGood: {
    backgroundColor: "#1A3A2A",
  },
  badgeTextGood: {
    color: "#6EF5A8",
  },
  noDataText: {
    color: "#5A5A78",
    fontSize: 14,
    textAlign: "center",
    paddingVertical: 20,
  },
  appRow: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    marginBottom: 16,
  },
  appLeft: {
    flexDirection: "row",
    alignItems: "center",
    flex: 1,
    marginRight: 12,
  },
  appRank: {
    width: 30,
    height: 30,
    borderRadius: 10,
    alignItems: "center",
    justifyContent: "center",
    marginRight: 12,
  },
  appRankText: {
    color: "#000000",
    fontSize: 13,
    fontWeight: "800",
  },
  appInfo: {
    flex: 1,
  },
  appName: {
    color: "#E0E0F0",
    fontSize: 14,
    fontWeight: "700",
    marginBottom: 2,
  },
  appPkg: {
    color: "#5A5A78",
    fontSize: 10,
    marginBottom: 6,
  },
  progressTrack: {
    height: 3,
    backgroundColor: "#1F1F30",
    borderRadius: 2,
    overflow: "hidden",
  },
  progressFill: {
    height: "100%",
    borderRadius: 2,
  },
  appTime: {
    color: "#FFFFFF",
    fontSize: 14,
    fontWeight: "700",
    minWidth: 44,
    textAlign: "right",
  },

  // ── Summary Cards ──
  summaryRow: {
    flexDirection: "row",
    marginHorizontal: 20,
    marginBottom: 20,
  },
  summaryCard: {
    backgroundColor: "#13131F",
    borderRadius: 20,
    padding: 18,
    borderWidth: 1,
    borderColor: "#1F1F2E",
  },
  summaryLabel: {
    color: "#5A5A78",
    fontSize: 11,
    fontWeight: "600",
    textTransform: "uppercase",
    letterSpacing: 1,
    marginBottom: 8,
  },
  summaryValue: {
    color: "#FFFFFF",
    fontSize: 20,
    fontWeight: "800",
    letterSpacing: -0.3,
    marginBottom: 2,
  },
  summarySubValue: {
    color: "#7C6EF5",
    fontSize: 12,
    fontWeight: "600",
  },

  // ── Footer ──
  footer: {
    alignItems: "center",
    paddingTop: 8,
  },
  footerText: {
    color: "#2A2A3D",
    fontSize: 11,
    fontWeight: "500",
  },

  barFuture: {
    backgroundColor: "#1A1A28",
    opacity: 0.4,
  },
});

export const rankColor = (idx) => {
  const palette = [
    "#7C6EF5",
    "#F5866E",
    "#6EC9F5",
    "#F5D56E",
    "#6EF5A8",
    "#F56EC9",
  ];
  return palette[idx % palette.length];
};

export const loginStyles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#0A0A0F",
    justifyContent: "center",
    paddingHorizontal: 28,
  },
  header: {
    alignItems: "center",
    marginBottom: 40,
  },
  logoContainer: {
    width: 68,
    height: 68,
    borderRadius: 20,
    backgroundColor: "#12121A",
    borderWidth: 1,
    borderColor: "#00E5A0",
    alignItems: "center",
    justifyContent: "center",
    marginBottom: 14,
  },
  logoIcon: {
    fontSize: 30,
  },
  appName: {
    fontSize: 28,
    fontWeight: "800",
    color: "#FFFFFF",
    letterSpacing: 1.5,
    marginBottom: 6,
  },
  tagline: {
    fontSize: 13,
    color: "#555568",
    letterSpacing: 0.4,
  },
  form: {
    backgroundColor: "#12121A",
    borderRadius: 20,
    padding: 24,
    borderWidth: 1,
    borderColor: "#1E1E2E",
  },
  formTitle: {
    fontSize: 18,
    fontWeight: "700",
    color: "#FFFFFF",
    marginBottom: 24,
  },
  inputWrapper: {
    marginBottom: 18,
  },
  inputLabel: {
    fontSize: 10,
    fontWeight: "700",
    color: "#00E5A0",
    letterSpacing: 1.5,
    marginBottom: 8,
  },
  input: {
    backgroundColor: "#0A0A0F",
    borderWidth: 1,
    borderColor: "#1E1E2E",
    borderRadius: 12,
    paddingHorizontal: 16,
    paddingVertical: 14,
    color: "#FFFFFF",
    fontSize: 15,
  },
  passwordContainer: {
    position: "relative",
  },
  passwordInput: {
    paddingRight: 50,
  },
  eyeButton: {
    position: "absolute",
    right: 14,
    top: 14,
  },
  eyeIcon: {
    fontSize: 18,
  },
  button: {
    backgroundColor: "#00E5A0",
    borderRadius: 12,
    paddingVertical: 16,
    alignItems: "center",
    marginTop: 8,
  },
  buttonDisabled: {
    opacity: 0.6,
  },
  buttonText: {
    color: "#0A0A0F",
    fontSize: 14,
    fontWeight: "800",
    letterSpacing: 1.5,
  },
  divider: {
    flexDirection: "row",
    alignItems: "center",
    marginVertical: 20,
    gap: 10,
  },
  dividerLine: {
    flex: 1,
    height: 1,
    backgroundColor: "#1E1E2E",
  },
  dividerText: {
    color: "#3A3A4A",
    fontSize: 13,
  },
  secondaryButton: {
    alignItems: "center",
  },
  secondaryButtonText: {
    color: "#555568",
    fontSize: 14,
  },
  secondaryButtonAccent: {
    color: "#00E5A0",
    fontWeight: "700",
  },
});

export const registerStyles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#0A0A0F",
  },
  scroll: {
    paddingHorizontal: 28,
    paddingTop: 20,
  },
  header: {
    alignItems: "center",
    marginBottom: 32,
  },
  backButton: {
    alignSelf: "flex-start",
    padding: 8,
    marginBottom: 16,
  },
  backIcon: {
    color: "#00E5A0",
    fontSize: 22,
    fontWeight: "600",
  },
  logoContainer: {
    width: 60,
    height: 60,
    borderRadius: 18,
    backgroundColor: "#12121A",
    borderWidth: 1,
    borderColor: "#00E5A0",
    alignItems: "center",
    justifyContent: "center",
    marginBottom: 12,
  },
  logoIcon: {
    fontSize: 26,
  },
  appName: {
    fontSize: 24,
    fontWeight: "800",
    color: "#FFFFFF",
    letterSpacing: 1.5,
    marginBottom: 4,
  },
  tagline: {
    fontSize: 13,
    color: "#555568",
    letterSpacing: 0.3,
  },
  form: {
    backgroundColor: "#12121A",
    borderRadius: 20,
    padding: 24,
    borderWidth: 1,
    borderColor: "#1E1E2E",
  },
  formTitle: {
    fontSize: 18,
    fontWeight: "700",
    color: "#FFFFFF",
    marginBottom: 24,
  },
  inputWrapper: {
    marginBottom: 18,
  },
  inputLabel: {
    fontSize: 10,
    fontWeight: "700",
    color: "#00E5A0",
    letterSpacing: 1.5,
    marginBottom: 8,
  },
  input: {
    backgroundColor: "#0A0A0F",
    borderWidth: 1,
    borderColor: "#1E1E2E",
    borderRadius: 12,
    paddingHorizontal: 16,
    paddingVertical: 14,
    color: "#FFFFFF",
    fontSize: 15,
  },
  inputError: {
    borderColor: "#FF4D6D",
  },
  errorHint: {
    color: "#FF4D6D",
    fontSize: 11,
    marginTop: 5,
    marginLeft: 4,
  },
  passwordContainer: {
    position: "relative",
  },
  passwordInput: {
    paddingRight: 50,
  },
  eyeButton: {
    position: "absolute",
    right: 14,
    top: 14,
  },
  eyeIcon: {
    fontSize: 18,
  },
  button: {
    backgroundColor: "#00E5A0",
    borderRadius: 12,
    paddingVertical: 16,
    alignItems: "center",
    marginTop: 8,
  },
  buttonDisabled: {
    opacity: 0.6,
  },
  buttonText: {
    color: "#0A0A0F",
    fontSize: 14,
    fontWeight: "800",
    letterSpacing: 1.5,
  },
  loginLink: {
    alignItems: "center",
    marginTop: 20,
  },
  loginLinkText: {
    color: "#555568",
    fontSize: 14,
  },
  loginLinkAccent: {
    color: "#00E5A0",
    fontWeight: "700",
  },
  bottomSpacer: {
    height: 40,
  },
});

export const groupsStyles = StyleSheet.create({
  root: {
    flex: 1,
    backgroundColor: "#0F0F14",
  },

  // Header
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: 20,
    paddingTop: 56,
    paddingBottom: 16,
    backgroundColor: "#0F0F14",
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: "700",
    color: "#fff",
    letterSpacing: -0.5,
  },
  headerActions: {
    flexDirection: "row",
    gap: 8,
  },
  headerBtn: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 20,
    backgroundColor: "#1A1A24",
    borderWidth: 1,
    borderColor: "#2A2A38",
  },
  headerBtnPrimary: {
    backgroundColor: "#7C6EF5",
    borderColor: "#7C6EF5",
  },
  headerBtnText: {
    color: "#aaa",
    fontSize: 14,
    fontWeight: "600",
  },

  // Scroll
  scroll: {
    padding: 16,
    paddingBottom: 40,
  },
  scrollCenter: {
    flexGrow: 1,
    justifyContent: "center",
  },

  // Loading
  loadingContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },

  // Empty state
  emptyContainer: {
    alignItems: "center",
    paddingVertical: 48,
    paddingHorizontal: 32,
  },
  emptyIcon: {
    fontSize: 52,
    marginBottom: 16,
  },
  emptyTitle: {
    fontSize: 20,
    fontWeight: "700",
    color: "#fff",
    marginBottom: 8,
  },
  emptySubtitle: {
    fontSize: 14,
    color: "#666",
    textAlign: "center",
    lineHeight: 20,
    marginBottom: 28,
  },
  emptyActions: {
    flexDirection: "row",
    gap: 10,
  },
  emptyBtn: {
    backgroundColor: "#7C6EF5",
    paddingHorizontal: 20,
    paddingVertical: 12,
    borderRadius: 12,
  },
  emptyBtnOutline: {
    backgroundColor: "transparent",
    borderWidth: 1,
    borderColor: "#7C6EF5",
  },
  emptyBtnText: {
    color: "#fff",
    fontWeight: "600",
    fontSize: 14,
  },

  // Group card
  groupCard: {
    backgroundColor: "#1A1A24",
    borderRadius: 16,
    marginBottom: 12,
    overflow: "hidden",
    borderWidth: 1,
    borderColor: "#252532",
  },
  groupCardHeader: {
    flexDirection: "row",
    alignItems: "center",
    padding: 16,
  },
  groupIconWrap: {
    width: 44,
    height: 44,
    borderRadius: 12,
    backgroundColor: "#7C6EF520",
    alignItems: "center",
    justifyContent: "center",
    marginRight: 12,
  },
  groupIconText: {
    fontSize: 20,
    fontWeight: "700",
    color: "#7C6EF5",
  },
  groupHeaderInfo: {
    flex: 1,
  },
  groupName: {
    fontSize: 16,
    fontWeight: "700",
    color: "#fff",
    marginBottom: 2,
  },
  groupMeta: {
    fontSize: 13,
    color: "#666",
  },
  chevron: {
    fontSize: 22,
    color: "#555",
    transform: [{ rotate: "0deg" }],
  },
  chevronOpen: {
    transform: [{ rotate: "90deg" }],
  },

  // Invite code
  inviteRow: {
    flexDirection: "row",
    alignItems: "center",
    paddingHorizontal: 16,
    paddingVertical: 10,
    backgroundColor: "#13131C",
    marginHorizontal: 12,
    marginBottom: 4,
    borderRadius: 10,
  },
  inviteLabel: {
    fontSize: 12,
    color: "#555",
    marginRight: 8,
    fontWeight: "500",
  },
  inviteCode: {
    flex: 1,
    fontSize: 13,
    color: "#7C6EF5",
    fontWeight: "600",
    letterSpacing: 1,
    fontFamily: "monospace",
  },
  inviteCopy: {
    fontSize: 11,
    color: "#444",
  },

  // Members
  membersDivider: {
    height: 1,
    backgroundColor: "#252532",
    marginHorizontal: 16,
    marginBottom: 8,
  },
  membersList: {
    paddingBottom: 4,
  },
  memberRow: {
    flexDirection: "row",
    alignItems: "center",
    paddingHorizontal: 16,
    paddingVertical: 10,
  },
  memberAvatar: {
    width: 34,
    height: 34,
    borderRadius: 10,
    backgroundColor: "#252532",
    alignItems: "center",
    justifyContent: "center",
    marginRight: 10,
  },
  memberAvatarText: {
    fontSize: 14,
    fontWeight: "700",
    color: "#aaa",
  },
  memberInfo: {
    flex: 1,
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
  },
  memberName: {
    fontSize: 14,
    color: "#ddd",
    fontWeight: "500",
  },
  ownerBadge: {
    backgroundColor: "#7C6EF520",
    paddingHorizontal: 7,
    paddingVertical: 2,
    borderRadius: 6,
  },
  meBadge: {
    backgroundColor: "#00E5A020",
  },
  ownerBadgeText: {
    fontSize: 10,
    color: "#7C6EF5",
    fontWeight: "600",
    textTransform: "uppercase",
    letterSpacing: 0.5,
  },
  kickBtn: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 8,
    backgroundColor: "#E0555520",
    borderWidth: 1,
    borderColor: "#E0555540",
  },
  kickBtnText: {
    fontSize: 12,
    color: "#E05555",
    fontWeight: "600",
  },

  // Leave button
  leaveBtn: {
    margin: 12,
    marginTop: 8,
    paddingVertical: 12,
    borderRadius: 10,
    backgroundColor: "#E0555510",
    borderWidth: 1,
    borderColor: "#E0555530",
    alignItems: "center",
  },
  leaveBtnText: {
    fontSize: 14,
    color: "#E05555",
    fontWeight: "600",
  },

  // Modal
  modalOverlay: {
    flex: 1,
    backgroundColor: "#000000cc",
    justifyContent: "flex-end",
  },
  modalBox: {
    backgroundColor: "#1A1A24",
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    padding: 24,
    paddingBottom: 40,
    borderTopWidth: 1,
    borderColor: "#2A2A38",
  },
  modalHeader: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 20,
  },
  modalTitle: {
    flex: 1,
    fontSize: 18,
    fontWeight: "700",
    color: "#fff",
  },
  modalClose: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: "#252532",
    alignItems: "center",
    justifyContent: "center",
  },
  modalCloseText: {
    fontSize: 14,
    color: "#888",
  },
  modalInput: {
    backgroundColor: "#0F0F14",
    borderRadius: 12,
    paddingHorizontal: 16,
    paddingVertical: 14,
    color: "#fff",
    fontSize: 16,
    borderWidth: 1,
    borderColor: "#2A2A38",
    marginBottom: 14,
  },
  modalBtn: {
    backgroundColor: "#7C6EF5",
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: "center",
  },
  modalBtnDanger: {
    backgroundColor: "#E05555",
  },
  modalBtnText: {
    color: "#fff",
    fontWeight: "700",
    fontSize: 16,
  },
});
