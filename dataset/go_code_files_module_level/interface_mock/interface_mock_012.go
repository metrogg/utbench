func setup(c *cli.Context) (err error) {
	interact.Run(&interact.Interact{
		Before: func(context interact.Context) error {
			context.SetErr(realize.Red.Bold("INVALID INPUT"))
			context.SetPrfx(realize.Output, realize.Yellow.Regular("[")+time.Now().Format("15:04:05")+realize.Yellow.Regular("]")+realize.Yellow.Bold("[")+strings.ToUpper(realize.RPrefix)+realize.Yellow.Bold("]"))
			return nil
		},
		Questions: []*interact.Question{
			{
				Before: func(d interact.Context) error {
					if _, err := os.Stat(realize.RFile); err != nil {
						d.Skip()
					}
					d.SetDef(false, realize.Green.Regular("(n)"))
					return nil
				},
				Quest: interact.Quest{
					Options: realize.Yellow.Regular("[y/n]"),
					Msg:     "Would you want to overwrite existing " + realize.Magenta.Regular(realize.RPrefix) + " config?",
				},
				Action: func(d interact.Context) interface{} {
					val, err := d.Ans().Bool()
					if err != nil {
						return d.Err()
					} else if val {
						r := realize.Realize{}
						r.Server = realize.Server{Parent: &r, Status: false, Open: false, Port: realize.Port, Host: realize.Host}
					}
					return nil
				},
			},
			{
				Before: func(d interact.Context) error {
					d.SetDef(false, realize.Green.Regular("(n)"))
					return nil
				},
				Quest: interact.Quest{
					Options: realize.Yellow.Regular("[y/n]"),
					Msg:     "Would you want to customize settings?",
					Resolve: func(d interact.Context) bool {
						val, _ := d.Ans().Bool()
						return val
					},
				},
				Subs: []*interact.Question{
					{
						Before: func(d interact.Context) error {
							d.SetDef(0, realize.Green.Regular("(os default)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[int]"),
							Msg:     "Set max number of open files (root required)",
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Int()
							if err != nil {
								return d.Err()
							}
							r.Settings.FileLimit = int32(val)
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Force polling watcher?",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef(100, realize.Green.Regular("(100ms)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[int]"),
									Msg:     "Set polling interval",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().Int()
									if err != nil {
										return d.Err()
									}
									r.Settings.Legacy.Interval = time.Duration(int(val)) * time.Millisecond
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Settings.Legacy.Force = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable logging files",
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Settings.Files.Errors = realize.Resource{Name: realize.FileErr, Status: val}
							r.Settings.Files.Outputs = realize.Resource{Name: realize.FileOut, Status: val}
							r.Settings.Files.Logs = realize.Resource{Name: realize.FileLog, Status: val}
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable web server",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef(realize.Port, realize.Green.Regular("("+strconv.Itoa(realize.Port)+")"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[int]"),
									Msg:     "Server port",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().Int()
									if err != nil {
										return d.Err()
									}
									r.Server.Port = int(val)
									return nil
								},
							},
							{
								Before: func(d interact.Context) error {
									d.SetDef(realize.Host, realize.Green.Regular("("+realize.Host+")"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Server host",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									r.Server.Host = val
									return nil
								},
							},
							{
								Before: func(d interact.Context) error {
									d.SetDef(false, realize.Green.Regular("(n)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[y/n]"),
									Msg:     "Open in current browser",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().Bool()
									if err != nil {
										return d.Err()
									}
									r.Server.Open = val
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Server.Status = val
							return nil
						},
					},
				},
				Action: func(d interact.Context) interface{} {
					_, err := d.Ans().Bool()
					if err != nil {
						return d.Err()
					}
					return nil
				},
			},
			{
				Before: func(d interact.Context) error {
					d.SetDef(true, realize.Green.Regular("(y)"))
					d.SetEnd("!")
					return nil
				},
				Quest: interact.Quest{
					Options: realize.Yellow.Regular("[y/n]"),
					Msg:     "Would you want to " + realize.Magenta.Regular("add a new project") + "? (insert '!' to stop)",
					Resolve: func(d interact.Context) bool {
						val, _ := d.Ans().Bool()
						if val {
							r.Schema.Add(r.Schema.New(c))
						}
						return val
					},
				},
				Subs: []*interact.Question{
					{
						Before: func(d interact.Context) error {
							d.SetDef(realize.Wdir(), realize.Green.Regular("("+realize.Wdir()+")"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[string]"),
							Msg:     "Project name",
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().String()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Name = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							dir := realize.Wdir()
							d.SetDef(dir, realize.Green.Regular("("+dir+")"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[string]"),
							Msg:     "Project path",
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().String()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Path = filepath.Clean(val)
							return nil
						},
					},

					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable go vet",
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(none)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Vet additional arguments",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									if val != "" {
										r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Vet.Args = append(r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Vet.Args, val)
									}
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Vet.Status = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable go fmt",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(none)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Fmt additional arguments",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									if val != "" {
										r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Fmt.Args = append(r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Fmt.Args, val)
									}
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Fmt.Status = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable go test",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(none)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Test additional arguments",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									if val != "" {
										r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Test.Args = append(r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Test.Args, val)
									}
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Test.Status = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable go clean",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(none)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Clean additional arguments",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									if val != "" {
										r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Clean.Args = append(r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Clean.Args, val)
									}
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Clean.Status = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable go generate",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(none)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Generate additional arguments",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									if val != "" {
										r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Generate.Args = append(r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Generate.Args, val)
									}
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Generate.Status = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(true, realize.Green.Regular("(y)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable go install",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(none)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Install additional arguments",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									if val != "" {
										r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Install.Args = append(r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Install.Args, val)
									}
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Install.Status = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable go build",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(none)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Build additional arguments",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									if val != "" {
										r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Build.Args = append(r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Build.Args, val)
									}
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Build.Status = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(true, realize.Green.Regular("(y)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Enable go run",
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].Tools.Run.Status = val
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Customize watching paths",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								if val {
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Paths = r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Paths[:len(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Paths)-1]
								}
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetEnd("!")
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Insert a path to watch (insert '!' to stop)",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Paths = append(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Paths, val)
									d.Reload()
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							_, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Customize ignore paths",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								if val {
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Ignore = r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Ignore[:len(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Ignore)-1]
								}
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetEnd("!")
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Insert a path to ignore (insert '!' to stop)",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Ignore = append(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Ignore, val)
									d.Reload()
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							_, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(n)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Add an additional argument",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									d.SetEnd("!")
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Add another argument (insert '!' to stop)",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Args = append(r.Schema.Projects[len(r.Schema.Projects)-1].Args, val)
									d.Reload()
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							_, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(none)"))
							d.SetEnd("!")
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Add a 'before' custom command (insert '!' to stop)",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Insert a command",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts = append(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts, realize.Command{Type: "before", Cmd: val})
									return nil
								},
							},
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(n)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Launch from a specific path",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts[len(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts)-1].Path = val
									return nil
								},
							},
							{
								Before: func(d interact.Context) error {
									d.SetDef(false, realize.Green.Regular("(n)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[y/n]"),
									Msg:     "Tag as global command",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().Bool()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts[len(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts)-1].Global = val
									return nil
								},
							},
							{
								Before: func(d interact.Context) error {
									d.SetDef(false, realize.Green.Regular("(n)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[y/n]"),
									Msg:     "Display command output",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().Bool()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts[len(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts)-1].Output = val
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							if val {
								d.Reload()
							}
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef(false, realize.Green.Regular("(none)"))
							d.SetEnd("!")
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[y/n]"),
							Msg:     "Add an 'after' custom commands  (insert '!' to stop)",
							Resolve: func(d interact.Context) bool {
								val, _ := d.Ans().Bool()
								return val
							},
						},
						Subs: []*interact.Question{
							{
								Before: func(d interact.Context) error {
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Insert a command",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts = append(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts, realize.Command{Type: "after", Cmd: val})
									return nil
								},
							},
							{
								Before: func(d interact.Context) error {
									d.SetDef("", realize.Green.Regular("(n)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[string]"),
									Msg:     "Launch from a specific path",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().String()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts[len(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts)-1].Path = val
									return nil
								},
							},
							{
								Before: func(d interact.Context) error {
									d.SetDef(false, realize.Green.Regular("(n)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[y/n]"),
									Msg:     "Tag as global command",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().Bool()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts[len(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts)-1].Global = val
									return nil
								},
							},
							{
								Before: func(d interact.Context) error {
									d.SetDef(false, realize.Green.Regular("(n)"))
									return nil
								},
								Quest: interact.Quest{
									Options: realize.Yellow.Regular("[y/n]"),
									Msg:     "Display command output",
								},
								Action: func(d interact.Context) interface{} {
									val, err := d.Ans().Bool()
									if err != nil {
										return d.Err()
									}
									r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts[len(r.Schema.Projects[len(r.Schema.Projects)-1].Watcher.Scripts)-1].Output = val
									return nil
								},
							},
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().Bool()
							if err != nil {
								return d.Err()
							}
							if val {
								d.Reload()
							}
							return nil
						},
					},
					{
						Before: func(d interact.Context) error {
							d.SetDef("", realize.Green.Regular("(none)"))
							return nil
						},
						Quest: interact.Quest{
							Options: realize.Yellow.Regular("[string]"),
							Msg:     "Set an error output pattern",
						},
						Action: func(d interact.Context) interface{} {
							val, err := d.Ans().String()
							if err != nil {
								return d.Err()
							}
							r.Schema.Projects[len(r.Schema.Projects)-1].ErrPattern = val
							return nil
						},
					},
				},
				Action: func(d interact.Context) interface{} {
					if val, err := d.Ans().Bool(); err != nil {
						return d.Err()
					} else if val {
						d.Reload()
					}
					return nil
				},
			},
		},
		After: func(d interact.Context) error {
			if val, _ := d.Qns().Get(0).Ans().Bool(); val {
				err := r.Settings.Remove(realize.RFile)
				if err != nil {
					return err
				}
			}
			return nil
		},
	})
	// create config
	err = r.Settings.Write(r)
	if err != nil {
		return err
	}
	log.Println(r.Prefix(realize.Green.Bold("Config successfully created")))
	return nil
}
