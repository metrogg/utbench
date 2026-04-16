func (h *LangHandler) Handle(ctx context.Context, conn jsonrpc2.JSONRPC2, req *jsonrpc2.Request) (result interface{}, err error) {
	// Prevent any uncaught panics from taking the entire server down.
	defer func() {
		if perr := util.Panicf(recover(), "%v", req.Method); perr != nil {
			err = perr
		}
	}()

	var cancelManager *cancel
	h.mu.Lock()
	cancelManager = h.cancel
	if req.Method != "initialize" && h.init == nil {
		h.mu.Unlock()
		return nil, errors.New("server must be initialized")
	}
	h.mu.Unlock()
	if err := h.CheckReady(); err != nil {
		if req.Method == "exit" {
			err = nil
		}
		return nil, err
	}

	if conn, ok := conn.(*jsonrpc2.Conn); ok && conn != nil {
		h.InitTracer(conn)
	}
	span, ctx, err := h.SpanForRequest(ctx, "lang", req, opentracing.Tags{"mode": "go"})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogEvent(fmt.Sprintf("error: %v", err))
		}
		span.Finish()
	}()

	// Notifications don't have an ID, so they can't be cancelled
	if cancelManager != nil && !req.Notif {
		var cancel func()
		ctx, cancel = cancelManager.WithCancel(ctx, req.ID)
		defer cancel()
	}

	switch req.Method {
	case "initialize":
		if h.init != nil {
			return nil, errors.New("language server is already initialized")
		}
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params InitializeParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}

		// HACK: RootPath is not a URI, but historically we treated it
		// as such. Convert it to a file URI
		if params.RootPath != "" && !util.IsURI(lsp.DocumentURI(params.RootPath)) {
			params.RootPath = string(util.PathToURI(params.RootPath))
		}

		if err := h.reset(&params); err != nil {
			return nil, err
		}

		// PERF: Kick off a workspace/symbol in the background to warm up the server
		if yes, _ := strconv.ParseBool(envWarmupOnInitialize); yes {
			go func() {
				ctx, cancel := context.WithDeadline(ctx, time.Now().Add(30*time.Second))
				defer cancel()
				_, _ = h.handleWorkspaceSymbol(ctx, conn, req, lspext.WorkspaceSymbolParams{
					Query: "",
					Limit: 100,
				})
			}()
		}

		// set the configured linter
		if h.config.DiagnosticsEnabled && h.config.LintTool != lintToolNone {
			switch h.config.LintTool {
			case lintToolGolint:
				h.linter = golint{}
			}

			if h.linter == nil {
				log.Printf("warning: lint tool %s not supported", h.config.LintTool)
			} else if err := h.linter.IsInstalled(ctx, h.BuildContext(ctx)); err != nil {
				h.linter = nil
				log.Printf("warning: lint tool (%s) initialize err: %s", h.config.LintTool, err)
			} else {
				// kick off a lint of the entire workspace
				go func() {
					ctx, cancel := context.WithDeadline(ctx, time.Now().Add(30*time.Second))
					defer cancel()
					err := h.lintWorkspace(ctx, h.BuildContext(ctx), conn)
					if err != nil {
						log.Printf("warning: failed to lint workspace: %s", err)
					}
				}()
			}
		}

		kind := lsp.TDSKIncremental
		var completionOp *lsp.CompletionOptions
		if h.config.GocodeCompletionEnabled {
			completionOp = &lsp.CompletionOptions{TriggerCharacters: []string{"."}}
		}
		return lsp.InitializeResult{
			Capabilities: lsp.ServerCapabilities{
				TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
					Kind: &kind,
				},
				CompletionProvider:           completionOp,
				DefinitionProvider:           true,
				TypeDefinitionProvider:       true,
				DocumentFormattingProvider:   true,
				DocumentSymbolProvider:       true,
				HoverProvider:                true,
				ReferencesProvider:           true,
				RenameProvider:               true,
				WorkspaceSymbolProvider:      true,
				ImplementationProvider:       true,
				XWorkspaceReferencesProvider: true,
				XDefinitionProvider:          true,
				XWorkspaceSymbolByProperties: true,
				SignatureHelpProvider:        &lsp.SignatureHelpOptions{TriggerCharacters: []string{"(", ","}},
			},
		}, nil

	case "initialized":
		// A notification that the client is ready to receive requests. Ignore
		return nil, nil

	case "shutdown":
		h.ShutDown()
		return nil, nil

	case "exit":
		if c, ok := conn.(*jsonrpc2.Conn); ok {
			c.Close()
		}
		return nil, nil

	case "$/cancelRequest":
		// notification, don't send back results/errors
		if req.Params == nil {
			return nil, nil
		}
		var params lsp.CancelParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, nil
		}
		if cancelManager == nil {
			return nil, nil
		}
		cancelManager.Cancel(jsonrpc2.ID{
			Num:      params.ID.Num,
			Str:      params.ID.Str,
			IsString: params.ID.IsString,
		})
		return nil, nil

	case "textDocument/hover":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.TextDocumentPositionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleHover(ctx, conn, req, params)

	case "textDocument/definition":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.TextDocumentPositionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleDefinition(ctx, conn, req, params)

	case "textDocument/typeDefinition":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.TextDocumentPositionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleTypeDefinition(ctx, conn, req, params)

	case "textDocument/xdefinition":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.TextDocumentPositionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleXDefinition(ctx, conn, req, params)

	case "textDocument/completion":
		if !h.config.GocodeCompletionEnabled {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound,
				Message: fmt.Sprintf("completion is disabled. Enable with flag `-gocodecompletion`")}
		}
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.CompletionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleTextDocumentCompletion(ctx, conn, req, params)

	case "textDocument/references":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.ReferenceParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleTextDocumentReferences(ctx, conn, req, params)

	case "textDocument/implementation":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.TextDocumentPositionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleTextDocumentImplementation(ctx, conn, req, params)

	case "textDocument/documentSymbol":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.DocumentSymbolParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleTextDocumentSymbol(ctx, conn, req, params)

	case "textDocument/signatureHelp":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.TextDocumentPositionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleTextDocumentSignatureHelp(ctx, conn, req, params)

	case "textDocument/formatting":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.DocumentFormattingParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleTextDocumentFormatting(ctx, conn, req, params)

	case "workspace/symbol":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lspext.WorkspaceSymbolParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleWorkspaceSymbol(ctx, conn, req, params)

	case "workspace/xreferences":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lspext.WorkspaceReferencesParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleWorkspaceReferences(ctx, conn, req, params)

	case "textDocument/rename":
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		var params lsp.RenameParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleRename(ctx, conn, req, params)
	default:
		if isFileSystemRequest(req.Method) {
			uri, fileChanged, err := h.handleFileSystemRequest(ctx, req)
			if fileChanged {
				// a file changed, so we must re-typecheck and re-enumerate symbols
				h.resetCaches(true)
			}
			if uri != "" {
				// a user is viewing this path, hint to add it to the cache
				// (unless we're primarily using binary package cache .a
				// files).
				if !h.config.UseBinaryPkgCache || (h.config.DiagnosticsEnabled && req.Method == "textDocument/didSave") {
					go h.typecheck(ctx, conn, uri, lsp.Position{})
				}

				if h.config.DiagnosticsEnabled && h.linter != nil && req.Method == "textDocument/didSave" {
					go func() {
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()
						err := h.lintPackage(ctx, h.BuildContext(ctx), conn, uri)
						if err != nil {
							log.Printf("warning: failed to lint package: %s", err)
						}
					}()
				}
			}
			return nil, err
		}

		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	}
}
