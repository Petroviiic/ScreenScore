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
