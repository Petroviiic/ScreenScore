import React, { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  StatusBar,
  ScrollView,
  Alert,
} from "react-native";
import * as SecureStore from "expo-secure-store";
import { router } from "expo-router";
import { registerStyles as styles } from "@/assets/styles/home.styles";
import * as Application from "expo-application";
const API_URL = "http://192.168.1.14:3000";

export default function Register() {
  const [form, setForm] = useState({
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
  });
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const update = (field, value) =>
    setForm((prev) => ({ ...prev, [field]: value }));

  const validate = () => {
    if (!form.username.trim()) return "Please enter a username.";
    if (form.username.trim().length < 3)
      return "Username must be at least 3 characters.";
    if (!form.email.trim()) return "Please enter an email address.";
    if (!form.email.includes("@")) return "Please enter a valid email address.";
    if (!form.password) return "Please enter a password.";
    if (form.password.length < 8)
      return "Password must be at least 8 characters.";
    if (form.password !== form.confirmPassword)
      return "Passwords do not match.";
    return null;
  };

  const handleRegister = async () => {
    const error = validate();
    if (error) {
      Alert.alert("Error", error);
      return;
    }

    setLoading(true);
    try {
      const deviceId = Application.getAndroidId();

      const response = await fetch(`${API_URL}/v1/users/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: form.username.trim(),
          email: form.email.trim().toLowerCase(),
          password: form.password,
          device_id: deviceId,
          push_notification_token: "",
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        Alert.alert("Error", data.error || "Registration failed.");
        return;
      }

      router.replace("/");
    } catch (err) {
      Alert.alert("Error", "Unable to connect to the server.");
    } finally {
      setLoading(false);
    }
  };

  const fields = [
    {
      key: "username",
      label: "USERNAME",
      placeholder: "johndoe",
      keyboardType: "default",
      autoCapitalize: "none",
    },
    {
      key: "email",
      label: "EMAIL",
      placeholder: "your@email.com",
      keyboardType: "email-address",
      autoCapitalize: "none",
    },
  ];

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === "ios" ? "padding" : undefined}
    >
      <StatusBar barStyle="light-content" backgroundColor="#0A0A0F" />
      <ScrollView
        contentContainerStyle={styles.scroll}
        keyboardShouldPersistTaps="handled"
        showsVerticalScrollIndicator={false}
      >
        {/* Header */}
        <View style={styles.header}>
          <TouchableOpacity
            style={styles.backButton}
            onPress={() => router.back()}
          >
            <Text style={styles.backIcon}>←</Text>
          </TouchableOpacity>
          <View style={styles.logoContainer}>
            <Text style={styles.logoIcon}>⏱</Text>
          </View>
          <Text style={styles.appName}>ScreenScore</Text>
          <Text style={styles.tagline}>
            Create an account and start tracking
          </Text>
        </View>

        {/* Form */}
        <View style={styles.form}>
          <Text style={styles.formTitle}>Sign Up</Text>

          {fields.map(
            ({ key, label, placeholder, keyboardType, autoCapitalize }) => (
              <View key={key} style={styles.inputWrapper}>
                <Text style={styles.inputLabel}>{label}</Text>
                <TextInput
                  style={styles.input}
                  placeholder={placeholder}
                  placeholderTextColor="#3A3A4A"
                  value={form[key]}
                  onChangeText={(v) => update(key, v)}
                  keyboardType={keyboardType}
                  autoCapitalize={autoCapitalize}
                  autoCorrect={false}
                />
              </View>
            )
          )}

          <View style={styles.inputWrapper}>
            <Text style={styles.inputLabel}>PASSWORD</Text>
            <View style={styles.passwordContainer}>
              <TextInput
                style={[styles.input, styles.passwordInput]}
                placeholder="at least 8 characters"
                placeholderTextColor="#3A3A4A"
                value={form.password}
                onChangeText={(v) => update("password", v)}
                secureTextEntry={!showPassword}
              />
              <TouchableOpacity
                style={styles.eyeButton}
                onPress={() => setShowPassword(!showPassword)}
              >
                <Text style={styles.eyeIcon}>{showPassword ? "🙈" : "👁"}</Text>
              </TouchableOpacity>
            </View>
          </View>

          <View style={styles.inputWrapper}>
            <Text style={styles.inputLabel}>CONFIRM PASSWORD</Text>
            <TextInput
              style={[
                styles.input,
                form.confirmPassword &&
                  form.password !== form.confirmPassword &&
                  styles.inputError,
              ]}
              placeholder="repeat password"
              placeholderTextColor="#3A3A4A"
              value={form.confirmPassword}
              onChangeText={(v) => update("confirmPassword", v)}
              secureTextEntry={!showPassword}
            />
            {form.confirmPassword && form.password !== form.confirmPassword && (
              <Text style={styles.errorHint}>Passwords do not match</Text>
            )}
          </View>

          <TouchableOpacity
            style={[styles.button, loading && styles.buttonDisabled]}
            onPress={handleRegister}
            disabled={loading}
            activeOpacity={0.85}
          >
            {loading ? (
              <ActivityIndicator color="#0A0A0F" />
            ) : (
              <Text style={styles.buttonText}>CREATE ACCOUNT</Text>
            )}
          </TouchableOpacity>

          <TouchableOpacity
            style={styles.loginLink}
            onPress={() => router.back()}
            activeOpacity={0.7}
          >
            <Text style={styles.loginLinkText}>
              Already have an account?{" "}
              <Text style={styles.loginLinkAccent}>Sign In</Text>
            </Text>
          </TouchableOpacity>
        </View>

        <View style={styles.bottomSpacer} />
      </ScrollView>
    </KeyboardAvoidingView>
  );
}
