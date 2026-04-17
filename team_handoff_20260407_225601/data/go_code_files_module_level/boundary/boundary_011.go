func (d *SQLiteDriver) Open(dsn string) (driver.Conn, error) {
	if C.sqlite3_threadsafe() == 0 {
		return nil, errors.New("sqlite library was not compiled for thread-safe operation")
	}

	var pkey string

	// Options
	var loc *time.Location
	authCreate := false
	authUser := ""
	authPass := ""
	authCrypt := ""
	authSalt := ""
	mutex := C.int(C.SQLITE_OPEN_FULLMUTEX)
	txlock := "BEGIN"

	// PRAGMA's
	autoVacuum := -1
	busyTimeout := 5000
	caseSensitiveLike := -1
	deferForeignKeys := -1
	foreignKeys := -1
	ignoreCheckConstraints := -1
	journalMode := "DELETE"
	lockingMode := "NORMAL"
	queryOnly := -1
	recursiveTriggers := -1
	secureDelete := "DEFAULT"
	synchronousMode := "NORMAL"
	writableSchema := -1

	pos := strings.IndexRune(dsn, '?')
	if pos >= 1 {
		params, err := url.ParseQuery(dsn[pos+1:])
		if err != nil {
			return nil, err
		}

		// Authentication
		if _, ok := params["_auth"]; ok {
			authCreate = true
		}
		if val := params.Get("_auth_user"); val != "" {
			authUser = val
		}
		if val := params.Get("_auth_pass"); val != "" {
			authPass = val
		}
		if val := params.Get("_auth_crypt"); val != "" {
			authCrypt = val
		}
		if val := params.Get("_auth_salt"); val != "" {
			authSalt = val
		}

		// _loc
		if val := params.Get("_loc"); val != "" {
			switch strings.ToLower(val) {
			case "auto":
				loc = time.Local
			default:
				loc, err = time.LoadLocation(val)
				if err != nil {
					return nil, fmt.Errorf("Invalid _loc: %v: %v", val, err)
				}
			}
		}

		// _mutex
		if val := params.Get("_mutex"); val != "" {
			switch strings.ToLower(val) {
			case "no":
				mutex = C.SQLITE_OPEN_NOMUTEX
			case "full":
				mutex = C.SQLITE_OPEN_FULLMUTEX
			default:
				return nil, fmt.Errorf("Invalid _mutex: %v", val)
			}
		}

		// _txlock
		if val := params.Get("_txlock"); val != "" {
			switch strings.ToLower(val) {
			case "immediate":
				txlock = "BEGIN IMMEDIATE"
			case "exclusive":
				txlock = "BEGIN EXCLUSIVE"
			case "deferred":
				txlock = "BEGIN"
			default:
				return nil, fmt.Errorf("Invalid _txlock: %v", val)
			}
		}

		// Auto Vacuum (_vacuum)
		//
		// https://www.sqlite.org/pragma.html#pragma_auto_vacuum
		//
		pkey = "" // Reset pkey
		if _, ok := params["_auto_vacuum"]; ok {
			pkey = "_auto_vacuum"
		}
		if _, ok := params["_vacuum"]; ok {
			pkey = "_vacuum"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "none":
				autoVacuum = 0
			case "1", "full":
				autoVacuum = 1
			case "2", "incremental":
				autoVacuum = 2
			default:
				return nil, fmt.Errorf("Invalid _auto_vacuum: %v, expecting value of '0 NONE 1 FULL 2 INCREMENTAL'", val)
			}
		}

		// Busy Timeout (_busy_timeout)
		//
		// https://www.sqlite.org/pragma.html#pragma_busy_timeout
		//
		pkey = "" // Reset pkey
		if _, ok := params["_busy_timeout"]; ok {
			pkey = "_busy_timeout"
		}
		if _, ok := params["_timeout"]; ok {
			pkey = "_timeout"
		}
		if val := params.Get(pkey); val != "" {
			iv, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid _busy_timeout: %v: %v", val, err)
			}
			busyTimeout = int(iv)
		}

		// Case Sensitive Like (_cslike)
		//
		// https://www.sqlite.org/pragma.html#pragma_case_sensitive_like
		//
		pkey = "" // Reset pkey
		if _, ok := params["_case_sensitive_like"]; ok {
			pkey = "_case_sensitive_like"
		}
		if _, ok := params["_cslike"]; ok {
			pkey = "_cslike"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				caseSensitiveLike = 0
			case "1", "yes", "true", "on":
				caseSensitiveLike = 1
			default:
				return nil, fmt.Errorf("Invalid _case_sensitive_like: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Defer Foreign Keys (_defer_foreign_keys | _defer_fk)
		//
		// https://www.sqlite.org/pragma.html#pragma_defer_foreign_keys
		//
		pkey = "" // Reset pkey
		if _, ok := params["_defer_foreign_keys"]; ok {
			pkey = "_defer_foreign_keys"
		}
		if _, ok := params["_defer_fk"]; ok {
			pkey = "_defer_fk"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				deferForeignKeys = 0
			case "1", "yes", "true", "on":
				deferForeignKeys = 1
			default:
				return nil, fmt.Errorf("Invalid _defer_foreign_keys: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Foreign Keys (_foreign_keys | _fk)
		//
		// https://www.sqlite.org/pragma.html#pragma_foreign_keys
		//
		pkey = "" // Reset pkey
		if _, ok := params["_foreign_keys"]; ok {
			pkey = "_foreign_keys"
		}
		if _, ok := params["_fk"]; ok {
			pkey = "_fk"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				foreignKeys = 0
			case "1", "yes", "true", "on":
				foreignKeys = 1
			default:
				return nil, fmt.Errorf("Invalid _foreign_keys: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Ignore CHECK Constrains (_ignore_check_constraints)
		//
		// https://www.sqlite.org/pragma.html#pragma_ignore_check_constraints
		//
		if val := params.Get("_ignore_check_constraints"); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				ignoreCheckConstraints = 0
			case "1", "yes", "true", "on":
				ignoreCheckConstraints = 1
			default:
				return nil, fmt.Errorf("Invalid _ignore_check_constraints: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Journal Mode (_journal_mode | _journal)
		//
		// https://www.sqlite.org/pragma.html#pragma_journal_mode
		//
		pkey = "" // Reset pkey
		if _, ok := params["_journal_mode"]; ok {
			pkey = "_journal_mode"
		}
		if _, ok := params["_journal"]; ok {
			pkey = "_journal"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToUpper(val) {
			case "DELETE", "TRUNCATE", "PERSIST", "MEMORY", "OFF":
				journalMode = strings.ToUpper(val)
			case "WAL":
				journalMode = strings.ToUpper(val)

				// For WAL Mode set Synchronous Mode to 'NORMAL'
				// See https://www.sqlite.org/pragma.html#pragma_synchronous
				synchronousMode = "NORMAL"
			default:
				return nil, fmt.Errorf("Invalid _journal: %v, expecting value of 'DELETE TRUNCATE PERSIST MEMORY WAL OFF'", val)
			}
		}

		// Locking Mode (_locking)
		//
		// https://www.sqlite.org/pragma.html#pragma_locking_mode
		//
		pkey = "" // Reset pkey
		if _, ok := params["_locking_mode"]; ok {
			pkey = "_locking_mode"
		}
		if _, ok := params["_locking"]; ok {
			pkey = "_locking"
		}
		if val := params.Get("_locking"); val != "" {
			switch strings.ToUpper(val) {
			case "NORMAL", "EXCLUSIVE":
				lockingMode = strings.ToUpper(val)
			default:
				return nil, fmt.Errorf("Invalid _locking_mode: %v, expecting value of 'NORMAL EXCLUSIVE", val)
			}
		}

		// Query Only (_query_only)
		//
		// https://www.sqlite.org/pragma.html#pragma_query_only
		//
		if val := params.Get("_query_only"); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				queryOnly = 0
			case "1", "yes", "true", "on":
				queryOnly = 1
			default:
				return nil, fmt.Errorf("Invalid _query_only: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Recursive Triggers (_recursive_triggers)
		//
		// https://www.sqlite.org/pragma.html#pragma_recursive_triggers
		//
		pkey = "" // Reset pkey
		if _, ok := params["_recursive_triggers"]; ok {
			pkey = "_recursive_triggers"
		}
		if _, ok := params["_rt"]; ok {
			pkey = "_rt"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				recursiveTriggers = 0
			case "1", "yes", "true", "on":
				recursiveTriggers = 1
			default:
				return nil, fmt.Errorf("Invalid _recursive_triggers: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Secure Delete (_secure_delete)
		//
		// https://www.sqlite.org/pragma.html#pragma_secure_delete
		//
		if val := params.Get("_secure_delete"); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				secureDelete = "OFF"
			case "1", "yes", "true", "on":
				secureDelete = "ON"
			case "fast":
				secureDelete = "FAST"
			default:
				return nil, fmt.Errorf("Invalid _secure_delete: %v, expecting boolean value of '0 1 false true no yes off on fast'", val)
			}
		}

		// Synchronous Mode (_synchronous | _sync)
		//
		// https://www.sqlite.org/pragma.html#pragma_synchronous
		//
		pkey = "" // Reset pkey
		if _, ok := params["_synchronous"]; ok {
			pkey = "_synchronous"
		}
		if _, ok := params["_sync"]; ok {
			pkey = "_sync"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToUpper(val) {
			case "0", "OFF", "1", "NORMAL", "2", "FULL", "3", "EXTRA":
				synchronousMode = strings.ToUpper(val)
			default:
				return nil, fmt.Errorf("Invalid _synchronous: %v, expecting value of '0 OFF 1 NORMAL 2 FULL 3 EXTRA'", val)
			}
		}

		// Writable Schema (_writeable_schema)
		//
		// https://www.sqlite.org/pragma.html#pragma_writeable_schema
		//
		if val := params.Get("_writable_schema"); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				writableSchema = 0
			case "1", "yes", "true", "on":
				writableSchema = 1
			default:
				return nil, fmt.Errorf("Invalid _writable_schema: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		if !strings.HasPrefix(dsn, "file:") {
			dsn = dsn[:pos]
		}
	}

	var db *C.sqlite3
	name := C.CString(dsn)
	defer C.free(unsafe.Pointer(name))
	rv := C._sqlite3_open_v2(name, &db,
		mutex|C.SQLITE_OPEN_READWRITE|C.SQLITE_OPEN_CREATE,
		nil)
	if rv != 0 {
		if db != nil {
			C.sqlite3_close_v2(db)
		}
		return nil, Error{Code: ErrNo(rv)}
	}
	if db == nil {
		return nil, errors.New("sqlite succeeded without returning a database")
	}

	rv = C.sqlite3_busy_timeout(db, C.int(busyTimeout))
	if rv != C.SQLITE_OK {
		C.sqlite3_close_v2(db)
		return nil, Error{Code: ErrNo(rv)}
	}

	exec := func(s string) error {
		cs := C.CString(s)
		rv := C.sqlite3_exec(db, cs, nil, nil, nil)
		C.free(unsafe.Pointer(cs))
		if rv != C.SQLITE_OK {
			return lastError(db)
		}
		return nil
	}

	// USER AUTHENTICATION
	//
	// User Authentication is always performed even when
	// sqlite_userauth is not compiled in, because without user authentication
	// the authentication is a no-op.
	//
	// Workflow
	//	- Authenticate
	//		ON::SUCCESS		=> Continue
	//		ON::SQLITE_AUTH => Return error and exit Open(...)
	//
	//  - Activate User Authentication
	//		Check if the user wants to activate User Authentication.
	//		If so then first create a temporary AuthConn to the database
	//		This is possible because we are already successfully authenticated.
	//
	//	- Check if `sqlite_user`` table exists
	//		YES				=> Add the provided user from DSN as Admin User and
	//						   activate user authentication.
	//		NO				=> Continue
	//

	// Create connection to SQLite
	conn := &SQLiteConn{db: db, loc: loc, txlock: txlock}

	// Password Cipher has to be registered before authentication
	if len(authCrypt) > 0 {
		switch strings.ToUpper(authCrypt) {
		case "SHA1":
			if err := conn.RegisterFunc("sqlite_crypt", CryptEncoderSHA1, true); err != nil {
				return nil, fmt.Errorf("CryptEncoderSHA1: %s", err)
			}
		case "SSHA1":
			if len(authSalt) == 0 {
				return nil, fmt.Errorf("_auth_crypt=ssha1, requires _auth_salt")
			}
			if err := conn.RegisterFunc("sqlite_crypt", CryptEncoderSSHA1(authSalt), true); err != nil {
				return nil, fmt.Errorf("CryptEncoderSSHA1: %s", err)
			}
		case "SHA256":
			if err := conn.RegisterFunc("sqlite_crypt", CryptEncoderSHA256, true); err != nil {
				return nil, fmt.Errorf("CryptEncoderSHA256: %s", err)
			}
		case "SSHA256":
			if len(authSalt) == 0 {
				return nil, fmt.Errorf("_auth_crypt=ssha256, requires _auth_salt")
			}
			if err := conn.RegisterFunc("sqlite_crypt", CryptEncoderSSHA256(authSalt), true); err != nil {
				return nil, fmt.Errorf("CryptEncoderSSHA256: %s", err)
			}
		case "SHA384":
			if err := conn.RegisterFunc("sqlite_crypt", CryptEncoderSHA384, true); err != nil {
				return nil, fmt.Errorf("CryptEncoderSHA384: %s", err)
			}
		case "SSHA384":
			if len(authSalt) == 0 {
				return nil, fmt.Errorf("_auth_crypt=ssha384, requires _auth_salt")
			}
			if err := conn.RegisterFunc("sqlite_crypt", CryptEncoderSSHA384(authSalt), true); err != nil {
				return nil, fmt.Errorf("CryptEncoderSSHA384: %s", err)
			}
		case "SHA512":
			if err := conn.RegisterFunc("sqlite_crypt", CryptEncoderSHA512, true); err != nil {
				return nil, fmt.Errorf("CryptEncoderSHA512: %s", err)
			}
		case "SSHA512":
			if len(authSalt) == 0 {
				return nil, fmt.Errorf("_auth_crypt=ssha512, requires _auth_salt")
			}
			if err := conn.RegisterFunc("sqlite_crypt", CryptEncoderSSHA512(authSalt), true); err != nil {
				return nil, fmt.Errorf("CryptEncoderSSHA512: %s", err)
			}
		}
	}

	// Preform Authentication
	if err := conn.Authenticate(authUser, authPass); err != nil {
		return nil, err
	}

	// Register: authenticate
	// Authenticate will perform an authentication of the provided username
	// and password against the database.
	//
	// If a database contains the SQLITE_USER table, then the
	// call to Authenticate must be invoked with an
	// appropriate username and password prior to enable read and write
	//access to the database.
	//
	// Return SQLITE_OK on success or SQLITE_ERROR if the username/password
	// combination is incorrect or unknown.
	//
	// If the SQLITE_USER table is not present in the database file, then
	// this interface is a harmless no-op returnning SQLITE_OK.
	if err := conn.RegisterFunc("authenticate", conn.authenticate, true); err != nil {
		return nil, err
	}
	//
	// Register: auth_user_add
	// auth_user_add can be used (by an admin user only)
	// to create a new user. When called on a no-authentication-required
	// database, this routine converts the database into an authentication-
	// required database, automatically makes the added user an
	// administrator, and logs in the current connection as that user.
	// The AuthUserAdd only works for the "main" database, not
	// for any ATTACH-ed databases. Any call to AuthUserAdd by a
	// non-admin user results in an error.
	if err := conn.RegisterFunc("auth_user_add", conn.authUserAdd, true); err != nil {
		return nil, err
	}
	//
	// Register: auth_user_change
	// auth_user_change can be used to change a users
	// login credentials or admin privilege.  Any user can change their own
	// login credentials. Only an admin user can change another users login
	// credentials or admin privilege setting. No user may change their own
	// admin privilege setting.
	if err := conn.RegisterFunc("auth_user_change", conn.authUserChange, true); err != nil {
		return nil, err
	}
	//
	// Register: auth_user_delete
	// auth_user_delete can be used (by an admin user only)
	// to delete a user. The currently logged-in user cannot be deleted,
	// which guarantees that there is always an admin user and hence that
	// the database cannot be converted into a no-authentication-required
	// database.
	if err := conn.RegisterFunc("auth_user_delete", conn.authUserDelete, true); err != nil {
		return nil, err
	}

	// Register: auth_enabled
	// auth_enabled can be used to check if user authentication is enabled
	if err := conn.RegisterFunc("auth_enabled", conn.authEnabled, true); err != nil {
		return nil, err
	}

	// Auto Vacuum
	// Moved auto_vacuum command, the user preference for auto_vacuum needs to be implemented directly after
	// the authentication and before the sqlite_user table gets created if the user
	// decides to activate User Authentication because
	// auto_vacuum needs to be set before any tables are created
	// and activating user authentication creates the internal table `sqlite_user`.
	if autoVacuum > -1 {
		if err := exec(fmt.Sprintf("PRAGMA auto_vacuum = %d;", autoVacuum)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	// Check if user wants to activate User Authentication
	if authCreate {
		// Before going any further, we need to check that the user
		// has provided an username and password within the DSN.
		// We are not allowed to continue.
		if len(authUser) < 0 {
			return nil, fmt.Errorf("Missing '_auth_user' while user authentication was requested with '_auth'")
		}
		if len(authPass) < 0 {
			return nil, fmt.Errorf("Missing '_auth_pass' while user authentication was requested with '_auth'")
		}

		// Check if User Authentication is Enabled
		authExists := conn.AuthEnabled()
		if !authExists {
			if err := conn.AuthUserAdd(authUser, authPass, true); err != nil {
				return nil, err
			}
		}
	}

	// Case Sensitive LIKE
	if caseSensitiveLike > -1 {
		if err := exec(fmt.Sprintf("PRAGMA case_sensitive_like = %d;", caseSensitiveLike)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	// Defer Foreign Keys
	if deferForeignKeys > -1 {
		if err := exec(fmt.Sprintf("PRAGMA defer_foreign_keys = %d;", deferForeignKeys)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	// Forgein Keys
	if foreignKeys > -1 {
		if err := exec(fmt.Sprintf("PRAGMA foreign_keys = %d;", foreignKeys)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	// Ignore CHECK Constraints
	if ignoreCheckConstraints > -1 {
		if err := exec(fmt.Sprintf("PRAGMA ignore_check_constraints = %d;", ignoreCheckConstraints)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	// Journal Mode
	// Because default Journal Mode is DELETE this PRAGMA can always be executed.
	if err := exec(fmt.Sprintf("PRAGMA journal_mode = %s;", journalMode)); err != nil {
		C.sqlite3_close_v2(db)
		return nil, err
	}

	// Locking Mode
	// Because the default is NORMAL and this is not changed in this package
	// by using the compile time SQLITE_DEFAULT_LOCKING_MODE this PRAGMA can always be executed
	if err := exec(fmt.Sprintf("PRAGMA locking_mode = %s;", lockingMode)); err != nil {
		C.sqlite3_close_v2(db)
		return nil, err
	}

	// Query Only
	if queryOnly > -1 {
		if err := exec(fmt.Sprintf("PRAGMA query_only = %d;", queryOnly)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	// Recursive Triggers
	if recursiveTriggers > -1 {
		if err := exec(fmt.Sprintf("PRAGMA recursive_triggers = %d;", recursiveTriggers)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	// Secure Delete
	//
	// Because this package can set the compile time flag SQLITE_SECURE_DELETE with a build tag
	// the default value for secureDelete var is 'DEFAULT' this way
	// you can compile with secure_delete 'ON' and disable it for a specific database connection.
	if secureDelete != "DEFAULT" {
		if err := exec(fmt.Sprintf("PRAGMA secure_delete = %s;", secureDelete)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	// Synchronous Mode
	//
	// Because default is NORMAL this statement is always executed
	if err := exec(fmt.Sprintf("PRAGMA synchronous = %s;", synchronousMode)); err != nil {
		C.sqlite3_close_v2(db)
		return nil, err
	}

	// Writable Schema
	if writableSchema > -1 {
		if err := exec(fmt.Sprintf("PRAGMA writable_schema = %d;", writableSchema)); err != nil {
			C.sqlite3_close_v2(db)
			return nil, err
		}
	}

	if len(d.Extensions) > 0 {
		if err := conn.loadExtensions(d.Extensions); err != nil {
			conn.Close()
			return nil, err
		}
	}

	if d.ConnectHook != nil {
		if err := d.ConnectHook(conn); err != nil {
			conn.Close()
			return nil, err
		}
	}
	runtime.SetFinalizer(conn, (*SQLiteConn).Close)
	return conn, nil
}
