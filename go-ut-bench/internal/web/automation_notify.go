package web

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go-ut-bench/internal/store"
)

type webhookChannelConfig struct {
	URL        string            `json:"url"`
	Headers    map[string]string `json:"headers,omitempty"`
	AuthHeader string            `json:"auth_header,omitempty"`
	AuthEnv    string            `json:"auth_env,omitempty"`
}

type smtpChannelConfig struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Username    string   `json:"username,omitempty"`
	Password    string   `json:"password,omitempty"`
	UsernameEnv string   `json:"username_env,omitempty"`
	PasswordEnv string   `json:"password_env,omitempty"`
	From        string   `json:"from"`
	To          []string `json:"to"`
}

type botWebhookChannelConfig struct {
	URL       string            `json:"url"`
	Secret    string            `json:"secret,omitempty"`
	SecretEnv string            `json:"secret_env,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

func (s *Server) dispatchAutomationNotifications(ctx context.Context, schedule store.AutomationSchedule, run store.AutomationRun) {
	var policy automationNotifyPolicy
	_ = json.Unmarshal([]byte(schedule.NotifyPolicyJSON), &policy)
	channelIDs := policyNotificationChannels(policy)
	if len(channelIDs) == 0 {
		return
	}
	success := run.Status == string(StatusCompleted)
	canceled := run.Status == string(StatusCanceled)
	if (success && !policy.OnSuccess) || (canceled && !policy.OnCanceled) || (!success && !canceled && !policy.OnFailure) {
		return
	}
	maxAttempts := policy.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 1
	}
	if maxAttempts > 5 {
		maxAttempts = 5
	}
	payload := map[string]any{
		"type":                "utbench.automation_run.finished",
		"schedule_id":         schedule.ScheduleID,
		"schedule_name":       schedule.Name,
		"automation_run_id":   run.AutomationRunID,
		"run_id":              run.RunID,
		"status":              run.Status,
		"error":               run.Error,
		"started_at_utc":      run.StartedAtUTC,
		"ended_at_utc":        run.EndedAtUTC,
		"report_html_path":    filepath.Join(s.outputRoot, "runs", run.RunID, "report", "report.html"),
		"report_summary_path": filepath.Join(s.outputRoot, "runs", run.RunID, "report", "report_summary.json"),
		"summary":             decodeSummary(run.SummaryJSON),
	}
	for _, channelID := range channelIDs {
		channel, err := s.db.GetNotificationChannel(ctx, channelID)
		if err != nil || !channel.Enabled {
			continue
		}
		delivery := store.NotificationDelivery{
			AutomationRunID: run.AutomationRunID,
			ChannelID:       channel.ChannelID,
			Status:          "pending",
			RequestSummary:  summarizeNotificationRequest(channel, payload),
			CreatedAtUTC:    time.Now().UTC().Format(time.RFC3339Nano),
		}
		responseSummary, attempts, err := s.sendNotificationWithRetry(ctx, channel, payload, maxAttempts)
		delivery.Attempt = attempts
		delivery.ResponseSummary = responseSummary
		if err != nil {
			delivery.Status = "failed"
			delivery.Error = err.Error()
		} else {
			delivery.Status = "sent"
			delivery.DeliveredAtUTC = time.Now().UTC().Format(time.RFC3339Nano)
		}
		_, _ = s.db.CreateNotificationDelivery(ctx, delivery)
	}
}

func policyNotificationChannels(policy automationNotifyPolicy) []string {
	seen := make(map[string]bool)
	var out []string
	for _, id := range append(policy.ChannelIDs, policy.Channels...) {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

func (s *Server) sendNotificationWithRetry(ctx context.Context, channel store.NotificationChannel, payload map[string]any, maxAttempts int) (string, int, error) {
	var lastErr error
	var responseSummary string
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		responseSummary, lastErr = s.sendNotification(ctx, channel, payload)
		if lastErr == nil {
			return responseSummary, attempt, nil
		}
		if attempt < maxAttempts {
			select {
			case <-ctx.Done():
				return responseSummary, attempt, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
	}
	return responseSummary, maxAttempts, lastErr
}

func (s *Server) sendNotification(ctx context.Context, channel store.NotificationChannel, payload map[string]any) (string, error) {
	switch strings.ToLower(channel.Type) {
	case "webhook":
		var cfg webhookChannelConfig
		if err := json.Unmarshal([]byte(channel.ConfigJSON), &cfg); err != nil {
			return "", err
		}
		if strings.TrimSpace(cfg.URL) == "" {
			return "", fmt.Errorf("webhook url is empty")
		}
		headers := map[string]string{"Content-Type": "application/json"}
		for k, v := range cfg.Headers {
			headers[k] = v
		}
		if cfg.AuthHeader != "" && cfg.AuthEnv != "" {
			if token := os.Getenv(cfg.AuthEnv); token != "" {
				headers[cfg.AuthHeader] = token
			}
		}
		return postJSONNotification(ctx, cfg.URL, headers, reportHTMLNotificationPayload(payload))
	case "feishu", "lark":
		var cfg botWebhookChannelConfig
		if err := json.Unmarshal([]byte(channel.ConfigJSON), &cfg); err != nil {
			return "", err
		}
		body := map[string]any{"msg_type": "text", "content": map[string]string{"text": notificationText(payload)}}
		if secret := resolveSecret(cfg.Secret, cfg.SecretEnv); secret != "" {
			ts := strconv.FormatInt(time.Now().Unix(), 10)
			body["timestamp"] = ts
			body["sign"] = feishuSign(ts, secret)
		}
		return postJSONNotification(ctx, cfg.URL, cfg.Headers, body)
	case "dingtalk":
		var cfg botWebhookChannelConfig
		if err := json.Unmarshal([]byte(channel.ConfigJSON), &cfg); err != nil {
			return "", err
		}
		targetURL := cfg.URL
		if secret := resolveSecret(cfg.Secret, cfg.SecretEnv); secret != "" {
			var err error
			targetURL, err = dingtalkSignedURL(targetURL, secret)
			if err != nil {
				return "", err
			}
		}
		body := map[string]any{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"title": "UTBench automation",
				"text":  notificationMarkdown(payload),
			},
		}
		return postJSONNotification(ctx, targetURL, cfg.Headers, body)
	case "wecom":
		var cfg botWebhookChannelConfig
		if err := json.Unmarshal([]byte(channel.ConfigJSON), &cfg); err != nil {
			return "", err
		}
		body := map[string]any{
			"msgtype":  "markdown",
			"markdown": map[string]string{"content": notificationMarkdown(payload)},
		}
		return postJSONNotification(ctx, cfg.URL, cfg.Headers, body)
	case "email_smtp":
		var cfg smtpChannelConfig
		if err := json.Unmarshal([]byte(channel.ConfigJSON), &cfg); err != nil {
			return "", err
		}
		if cfg.Port == 0 {
			cfg.Port = 587
		}
		if cfg.Host == "" || cfg.From == "" || len(cfg.To) == 0 {
			return "", fmt.Errorf("smtp host/from/to is required")
		}
		auth, err := smtpAuthFromConfig(cfg)
		if err != nil {
			return "", err
		}
		addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
		subject := fmt.Sprintf("UTBench automation %s: %s", payload["status"], payload["run_id"])
		reportPath := fmt.Sprint(payload["report_html_path"])
		msg := buildSMTPReportMessage(cfg.From, cfg.To, subject, smtpReportPreviewHTML(payload), reportPath)
		if err := sendSMTPMail(ctx, addr, cfg.Host, auth, cfg.From, cfg.To, msg); err != nil {
			return "", err
		}
		return "smtp sent to " + strings.Join(cfg.To, ","), nil
	default:
		return "", fmt.Errorf("unsupported notification channel type: %s", channel.Type)
	}
}

func smtpAuthFromConfig(cfg smtpChannelConfig) (smtp.Auth, error) {
	username := strings.TrimSpace(cfg.Username)
	password := cfg.Password
	usernameEnv := strings.TrimSpace(cfg.UsernameEnv)
	passwordEnv := strings.TrimSpace(cfg.PasswordEnv)

	if username == "" && strings.Contains(usernameEnv, "@") {
		username = usernameEnv
		password = passwordEnv
		usernameEnv = ""
		passwordEnv = ""
	}
	if username == "" && password == "" && usernameEnv == "" && passwordEnv == "" {
		return nil, nil
	}
	if username != "" || password != "" {
		if username == "" || password == "" {
			return nil, fmt.Errorf("smtp username and password must both be configured")
		}
		return smtp.PlainAuth("", username, password, cfg.Host), nil
	}
	if usernameEnv == "" || passwordEnv == "" {
		return nil, fmt.Errorf("smtp username_env and password_env must both be configured")
	}
	username = strings.TrimSpace(os.Getenv(usernameEnv))
	password = os.Getenv(passwordEnv)
	if username == "" {
		return nil, fmt.Errorf("smtp username environment variable %q is not set; set it before starting UTBench", usernameEnv)
	}
	if password == "" {
		return nil, fmt.Errorf("smtp password environment variable %q is not set; set it before starting UTBench", passwordEnv)
	}
	return smtp.PlainAuth("", username, password, cfg.Host), nil
}

func smtpReportPreviewHTML(payload map[string]any) string {
	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:20px;background:#f6f8fb;color:#111827;font-family:Segoe UI,Arial,sans-serif;">
  <div style="max-width:680px;margin:0 auto;background:#ffffff;border:1px solid #e5e7eb;border-radius:10px;padding:20px;">
    <div style="font-size:12px;font-weight:700;color:#0f766e;letter-spacing:.04em;text-transform:uppercase;">UTBench Report</div>
    <h2 style="margin:8px 0 14px;font-size:20px;line-height:1.25;color:#111827;">任务报告已生成</h2>
    <table style="width:100%%;border-collapse:collapse;font-size:14px;line-height:1.55;">
      <tr><td style="padding:6px 0;color:#64748b;width:96px;">计划</td><td style="padding:6px 0;font-weight:600;">%s</td></tr>
      <tr><td style="padding:6px 0;color:#64748b;">Run ID</td><td style="padding:6px 0;font-family:Consolas,monospace;">%s</td></tr>
      <tr><td style="padding:6px 0;color:#64748b;">状态</td><td style="padding:6px 0;font-weight:700;">%s</td></tr>
      <tr><td style="padding:6px 0;color:#64748b;">报告</td><td style="padding:6px 0;font-family:Consolas,monospace;word-break:break-all;">%s</td></tr>
      <tr><td style="padding:6px 0;color:#64748b;">错误</td><td style="padding:6px 0;word-break:break-word;">%s</td></tr>
    </table>
    <p style="margin:16px 0 0;color:#475569;font-size:13px;line-height:1.6;">完整交互式报告已作为 <strong>report.html</strong> 附件发送。邮件客户端通常会禁用脚本和复杂样式，请下载或用浏览器打开附件查看完整图表与移动端布局。</p>
  </div>
</body>
</html>`,
		html.EscapeString(fmt.Sprint(payload["schedule_name"])),
		html.EscapeString(fmt.Sprint(payload["run_id"])),
		html.EscapeString(fmt.Sprint(payload["status"])),
		html.EscapeString(fmt.Sprint(payload["report_html_path"])),
		html.EscapeString(fmt.Sprint(payload["error"])),
	)
}

func buildSMTPReportMessage(from string, to []string, subject string, htmlBody string, reportPath string) []byte {
	boundary := "utbench-report-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	var msg bytes.Buffer
	msg.WriteString("From: " + from + "\r\n")
	msg.WriteString("To: " + strings.Join(to, ",") + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: multipart/mixed; boundary=\"" + boundary + "\"\r\n\r\n")

	msg.WriteString("--" + boundary + "\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	writeBase64Lines(&msg, []byte(htmlBody))
	msg.WriteString("\r\n")

	if data, filename, ok := readReportHTMLAttachment(reportPath); ok {
		msg.WriteString("--" + boundary + "\r\n")
		msg.WriteString("Content-Type: text/html; charset=UTF-8; name=\"" + filename + "\"\r\n")
		msg.WriteString("Content-Disposition: attachment; filename=\"" + filename + "\"\r\n")
		msg.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
		writeBase64Lines(&msg, data)
		msg.WriteString("\r\n")
	}

	msg.WriteString("--" + boundary + "--\r\n")
	return msg.Bytes()
}

func readReportHTMLAttachment(path string) ([]byte, string, bool) {
	if strings.TrimSpace(path) == "" {
		return nil, "", false
	}
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return nil, "", false
	}
	filename := filepath.Base(path)
	if filename == "." || filename == string(filepath.Separator) || strings.TrimSpace(filename) == "" {
		filename = "report.html"
	}
	return data, filename, true
}

func writeBase64Lines(dst *bytes.Buffer, data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	for len(encoded) > 76 {
		dst.WriteString(encoded[:76])
		dst.WriteString("\r\n")
		encoded = encoded[76:]
	}
	if encoded != "" {
		dst.WriteString(encoded)
		dst.WriteString("\r\n")
	}
}

func postJSONNotification(ctx context.Context, targetURL string, headers map[string]string, payload any) (string, error) {
	if strings.TrimSpace(targetURL) == "" {
		return "", fmt.Errorf("notification url is empty")
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		if strings.TrimSpace(k) != "" {
			req.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	excerpt := readResponseExcerpt(resp.Body, 600)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Sprintf("HTTP %d %s", resp.StatusCode, excerpt), fmt.Errorf("webhook HTTP %d", resp.StatusCode)
	}
	return fmt.Sprintf("HTTP %d %s", resp.StatusCode, excerpt), nil
}

func sendSMTPMail(ctx context.Context, addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Minute)
		defer cancel()
	}
	deadline, _ := ctx.Deadline()
	dialer := &net.Dialer{Timeout: time.Until(deadline)}
	var (
		conn net.Conn
		err  error
	)
	if strings.HasSuffix(addr, ":465") {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12})
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("smtp connect %s failed: %w; check outbound network/proxy/firewall and try port 465 or 587", addr, err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(deadline); err != nil {
		return fmt.Errorf("smtp set deadline failed: %w", err)
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp handshake with %s failed: %w; check TLS mode, port, and whether the network/proxy allows SMTP", addr, err)
	}
	defer c.Close()
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("smtp STARTTLS failed on %s: %w", addr, err)
		}
	}
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err := c.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth failed: %w; for Gmail use an app password, not the account login password", err)
			}
		}
	}
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}
	for _, addr := range to {
		if err := c.Rcpt(addr); err != nil {
			return fmt.Errorf("smtp RCPT TO %s failed: %w", addr, err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA failed: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return fmt.Errorf("smtp write message failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp finish message failed: %w", err)
	}
	if err := c.Quit(); err != nil {
		return fmt.Errorf("smtp QUIT failed: %w", err)
	}
	return nil
}

func resolveSecret(value, envName string) string {
	if strings.TrimSpace(envName) != "" {
		return os.Getenv(strings.TrimSpace(envName))
	}
	return strings.TrimSpace(value)
}

func feishuSign(timestamp, secret string) string {
	stringToSign := timestamp + "\n" + secret
	mac := hmac.New(sha256.New, []byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func dingtalkSignedURL(rawURL, secret string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	stringToSign := timestamp + "\n" + secret
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(stringToSign))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	q := u.Query()
	q.Set("timestamp", timestamp)
	q.Set("sign", sign)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func notificationText(payload map[string]any) string {
	return fmt.Sprintf("UTBench report.html\nSchedule: %s\nRun: %s\nStatus: %s\nReport: %s\nError: %s",
		payload["schedule_name"], payload["run_id"], payload["status"], payload["report_html_path"], payload["error"])
}

func notificationMarkdown(payload map[string]any) string {
	return fmt.Sprintf("### UTBench report.html\n- Schedule: %s\n- Run: `%s`\n- Status: **%s**\n- Report: %s\n- Error: %s",
		payload["schedule_name"], payload["run_id"], payload["status"], payload["report_html_path"], payload["error"])
}

func reportHTMLNotificationPayload(payload map[string]any) map[string]any {
	return map[string]any{
		"type":              "utbench.report_html",
		"schedule_id":       payload["schedule_id"],
		"schedule_name":     payload["schedule_name"],
		"automation_run_id": payload["automation_run_id"],
		"run_id":            payload["run_id"],
		"status":            payload["status"],
		"error":             payload["error"],
		"report_html_path":  payload["report_html_path"],
		"sent_at_utc":       time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func readReportHTMLBody(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func summarizeNotificationRequest(channel store.NotificationChannel, payload map[string]any) string {
	target := channel.Type
	if strings.EqualFold(channel.Type, "webhook") {
		var cfg webhookChannelConfig
		if json.Unmarshal([]byte(channel.ConfigJSON), &cfg) == nil {
			target = "webhook " + redactURL(cfg.URL)
		}
	} else if strings.EqualFold(channel.Type, "feishu") || strings.EqualFold(channel.Type, "lark") || strings.EqualFold(channel.Type, "dingtalk") || strings.EqualFold(channel.Type, "wecom") {
		var cfg botWebhookChannelConfig
		if json.Unmarshal([]byte(channel.ConfigJSON), &cfg) == nil {
			target = channel.Type + " " + redactURL(cfg.URL)
		}
	}
	raw, _ := json.Marshal(map[string]any{
		"channel": target,
		"type":    payload["type"],
		"run_id":  payload["run_id"],
		"status":  payload["status"],
	})
	return string(raw)
}

func redactURL(raw string) string {
	lower := strings.ToLower(raw)
	for _, marker := range []string{"token=", "access_token=", "key=", "secret="} {
		idx := strings.Index(lower, marker)
		if idx >= 0 {
			return raw[:idx+len(marker)] + "***"
		}
	}
	return raw
}

func readResponseExcerpt(r io.Reader, limit int64) string {
	if r == nil {
		return ""
	}
	data, _ := io.ReadAll(io.LimitReader(r, limit))
	return strings.TrimSpace(string(data))
}

func (s *Server) buildAutomationSummary(runID, status, errText string) string {
	summaryPath := filepath.Join(s.outputRoot, "runs", runID, "report", "report_summary.json")
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		raw, _ := json.Marshal(map[string]any{"status": status, "error": errText})
		return string(raw)
	}
	var summary any
	if err := json.Unmarshal(data, &summary); err != nil {
		raw, _ := json.Marshal(map[string]any{"status": status, "error": errText, "report_summary_path": summaryPath})
		return string(raw)
	}
	raw, _ := json.Marshal(summary)
	return string(raw)
}

func decodeSummary(raw string) any {
	var out any
	if raw != "" && json.Unmarshal([]byte(raw), &out) == nil {
		return out
	}
	return map[string]any{}
}
