func (r *replacer) getSubstitution(key string) string {
	// search custom replacements first
	if value, ok := r.customReplacements[key]; ok {
		return value
	}

	// search request headers then
	if key[1] == '>' {
		want := key[2 : len(key)-1]
		for key, values := range r.request.Header {
			// Header placeholders (case-insensitive)
			if strings.EqualFold(key, want) {
				return strings.Join(values, ",")
			}
		}
	}
	// search response headers then
	if r.responseRecorder != nil && key[1] == '<' {
		want := key[2 : len(key)-1]
		for key, values := range r.responseRecorder.Header() {
			// Header placeholders (case-insensitive)
			if strings.EqualFold(key, want) {
				return strings.Join(values, ",")
			}
		}
	}
	// next check for cookies
	if key[1] == '~' {
		name := key[2 : len(key)-1]
		if cookie, err := r.request.Cookie(name); err == nil {
			return cookie.Value
		}
	}
	// next check for query argument
	if key[1] == '?' {
		query := r.request.URL.Query()
		name := key[2 : len(key)-1]
		return query.Get(name)
	}

	// search default replacements in the end
	switch key {
	case "{method}":
		return r.request.Method
	case "{scheme}":
		if r.request.TLS != nil {
			return "https"
		}
		return "http"
	case "{hostname}":
		name, err := os.Hostname()
		if err != nil {
			return r.emptyValue
		}
		return name
	case "{host}":
		return r.request.Host
	case "{hostonly}":
		host, _, err := net.SplitHostPort(r.request.Host)
		if err != nil {
			return r.request.Host
		}
		return host
	case "{path}":
		u, _ := r.request.Context().Value(OriginalURLCtxKey).(url.URL)
		return u.Path
	case "{path_escaped}":
		u, _ := r.request.Context().Value(OriginalURLCtxKey).(url.URL)
		return url.QueryEscape(u.Path)
	case "{request_id}":
		reqid, _ := r.request.Context().Value(RequestIDCtxKey).(string)
		return reqid
	case "{rewrite_path}":
		return r.request.URL.Path
	case "{rewrite_path_escaped}":
		return url.QueryEscape(r.request.URL.Path)
	case "{query}":
		u, _ := r.request.Context().Value(OriginalURLCtxKey).(url.URL)
		return u.RawQuery
	case "{query_escaped}":
		u, _ := r.request.Context().Value(OriginalURLCtxKey).(url.URL)
		return url.QueryEscape(u.RawQuery)
	case "{fragment}":
		u, _ := r.request.Context().Value(OriginalURLCtxKey).(url.URL)
		return u.Fragment
	case "{proto}":
		return r.request.Proto
	case "{remote}":
		host, _, err := net.SplitHostPort(r.request.RemoteAddr)
		if err != nil {
			return r.request.RemoteAddr
		}
		return host
	case "{port}":
		_, port, err := net.SplitHostPort(r.request.RemoteAddr)
		if err != nil {
			return r.emptyValue
		}
		return port
	case "{uri}":
		u, _ := r.request.Context().Value(OriginalURLCtxKey).(url.URL)
		return u.RequestURI()
	case "{uri_escaped}":
		u, _ := r.request.Context().Value(OriginalURLCtxKey).(url.URL)
		return url.QueryEscape(u.RequestURI())
	case "{rewrite_uri}":
		return r.request.URL.RequestURI()
	case "{rewrite_uri_escaped}":
		return url.QueryEscape(r.request.URL.RequestURI())
	case "{when}":
		return now().Format(timeFormat)
	case "{when_iso_local}":
		return now().Format(timeFormatISO)
	case "{when_iso}":
		return now().UTC().Format(timeFormatISOUTC)
	case "{when_unix}":
		return strconv.FormatInt(now().Unix(), 10)
	case "{when_unix_ms}":
		return strconv.FormatInt(nanoToMilliseconds(now().UnixNano()), 10)
	case "{file}":
		_, file := path.Split(r.request.URL.Path)
		return file
	case "{dir}":
		dir, _ := path.Split(r.request.URL.Path)
		return dir
	case "{request}":
		dump, err := httputil.DumpRequest(r.request, false)
		if err != nil {
			return r.emptyValue
		}
		return requestReplacer.Replace(string(dump))
	case "{request_body}":
		if !canLogRequest(r.request) {
			return r.emptyValue
		}
		_, err := ioutil.ReadAll(r.request.Body)
		if err != nil {
			if err == ErrMaxBytesExceeded {
				return r.emptyValue
			}
		}
		return requestReplacer.Replace(r.requestBody.String())
	case "{mitm}":
		if val, ok := r.request.Context().Value(caddy.CtxKey("mitm")).(bool); ok {
			if val {
				return "likely"
			}
			return "unlikely"
		}
		return "unknown"
	case "{status}":
		if r.responseRecorder == nil {
			return r.emptyValue
		}
		return strconv.Itoa(r.responseRecorder.status)
	case "{size}":
		if r.responseRecorder == nil {
			return r.emptyValue
		}
		return strconv.Itoa(r.responseRecorder.size)
	case "{latency}":
		if r.responseRecorder == nil {
			return r.emptyValue
		}
		return roundDuration(time.Since(r.responseRecorder.start)).String()
	case "{latency_ms}":
		if r.responseRecorder == nil {
			return r.emptyValue
		}
		elapsedDuration := time.Since(r.responseRecorder.start)
		return strconv.FormatInt(convertToMilliseconds(elapsedDuration), 10)
	case "{tls_protocol}":
		if r.request.TLS != nil {
			if name, err := caddytls.GetSupportedProtocolName(r.request.TLS.Version); err == nil {
				return name
			} else {
				return "tls" // this should never happen, but guard in case
			}
		}
		return r.emptyValue // because not using a secure channel
	case "{tls_cipher}":
		if r.request.TLS != nil {
			if name, err := caddytls.GetSupportedCipherName(r.request.TLS.CipherSuite); err == nil {
				return name
			} else {
				return "UNKNOWN" // this should never happen, but guard in case
			}
		}
		return r.emptyValue
	case "{tls_client_escaped_cert}":
		cert := r.getPeerCert()
		if cert != nil {
			pemBlock := pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
			}
			return url.QueryEscape(string(pem.EncodeToMemory(&pemBlock)))
		}
		return r.emptyValue
	case "{tls_client_fingerprint}":
		cert := r.getPeerCert()
		if cert != nil {
			return fmt.Sprintf("%x", sha256.Sum256(cert.Raw))
		}
		return r.emptyValue
	case "{tls_client_i_dn}":
		cert := r.getPeerCert()
		if cert != nil {
			return cert.Issuer.String()
		}
		return r.emptyValue
	case "{tls_client_raw_cert}":
		cert := r.getPeerCert()
		if cert != nil {
			return string(cert.Raw)
		}
		return r.emptyValue
	case "{tls_client_s_dn}":
		cert := r.getPeerCert()
		if cert != nil {
			return cert.Subject.String()
		}
		return r.emptyValue
	case "{tls_client_serial}":
		cert := r.getPeerCert()
		if cert != nil {
			return fmt.Sprintf("%x", cert.SerialNumber)
		}
		return r.emptyValue
	case "{tls_client_v_end}":
		cert := r.getPeerCert()
		if cert != nil {
			return cert.NotAfter.In(time.UTC).Format("Jan 02 15:04:05 2006 MST")
		}
		return r.emptyValue
	case "{tls_client_v_remain}":
		cert := r.getPeerCert()
		if cert != nil {
			now := time.Now().In(time.UTC)
			days := int64(cert.NotAfter.Sub(now).Seconds() / 86400)
			return strconv.FormatInt(days, 10)
		}
		return r.emptyValue
	case "{tls_client_v_start}":
		cert := r.getPeerCert()
		if cert != nil {
			return cert.NotBefore.Format("Jan 02 15:04:05 2006 MST")
		}
		return r.emptyValue
	case "{server_port}":
		_, port, err := net.SplitHostPort(r.request.Host)
		if err != nil {
			if r.request.TLS != nil {
				return "443"
			} else {
				return "80"
			}
		}
		return port
	default:
		// {labelN}
		if strings.HasPrefix(key, "{label") {
			nStr := key[6 : len(key)-1] // get the integer N in "{labelN}"
			n, err := strconv.Atoi(nStr)
			if err != nil || n < 1 {
				return r.emptyValue
			}
			labels := strings.Split(r.request.Host, ".")
			if n > len(labels) {
				return r.emptyValue
			}
			return labels[n-1]
		}
	}

	return r.emptyValue
}
