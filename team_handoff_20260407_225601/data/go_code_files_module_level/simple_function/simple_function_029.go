func (parser *Parser) ParseGeneralAPIInfo(mainAPIFile string) error {
	fileSet := token.NewFileSet()
	fileTree, err := goparser.ParseFile(fileSet, mainAPIFile, nil, goparser.ParseComments)
	if err != nil {
		return errors.Wrap(err, "cannot parse soure files")
	}

	parser.swagger.Swagger = "2.0"
	securityMap := map[string]*spec.SecurityScheme{}

	// templated defaults
	parser.swagger.Info.Version = "{{.Version}}"
	parser.swagger.Info.Title = "{{.Title}}"
	parser.swagger.Info.Description = "{{.Description}}"
	parser.swagger.Host = "{{.Host}}"
	parser.swagger.BasePath = "{{.BasePath}}"

	if fileTree.Comments != nil {
		for _, comment := range fileTree.Comments {
			comments := strings.Split(comment.Text(), "\n")
			previousAttribute := ""
			for _, commentLine := range comments {
				attribute := strings.ToLower(strings.Split(commentLine, " ")[0])
				multilineBlock := false
				if previousAttribute == attribute {
					multilineBlock = true
				}
				switch attribute {
				case "@version":
					parser.swagger.Info.Version = strings.TrimSpace(commentLine[len(attribute):])
				case "@title":
					parser.swagger.Info.Title = strings.TrimSpace(commentLine[len(attribute):])
				case "@description":
					if parser.swagger.Info.Description == "{{.Description}}" {
						parser.swagger.Info.Description = strings.TrimSpace(commentLine[len(attribute):])
					} else if multilineBlock {
						parser.swagger.Info.Description += "\n" + strings.TrimSpace(commentLine[len(attribute):])
					}
				case "@termsofservice":
					parser.swagger.Info.TermsOfService = strings.TrimSpace(commentLine[len(attribute):])
				case "@contact.name":
					parser.swagger.Info.Contact.Name = strings.TrimSpace(commentLine[len(attribute):])
				case "@contact.email":
					parser.swagger.Info.Contact.Email = strings.TrimSpace(commentLine[len(attribute):])
				case "@contact.url":
					parser.swagger.Info.Contact.URL = strings.TrimSpace(commentLine[len(attribute):])
				case "@license.name":
					parser.swagger.Info.License.Name = strings.TrimSpace(commentLine[len(attribute):])
				case "@license.url":
					parser.swagger.Info.License.URL = strings.TrimSpace(commentLine[len(attribute):])
				case "@host":
					parser.swagger.Host = strings.TrimSpace(commentLine[len(attribute):])
				case "@basepath":
					parser.swagger.BasePath = strings.TrimSpace(commentLine[len(attribute):])
				case "@schemes":
					parser.swagger.Schemes = getSchemes(commentLine)
				case "@tag.name":
					commentInfo := strings.TrimSpace(commentLine[len(attribute):])
					parser.swagger.Tags = append(parser.swagger.Tags, spec.Tag{
						TagProps: spec.TagProps{
							Name: strings.TrimSpace(commentInfo),
						},
					})
				case "@tag.description":
					commentInfo := strings.TrimSpace(commentLine[len(attribute):])
					tag := parser.swagger.Tags[len(parser.swagger.Tags)-1]
					tag.TagProps.Description = commentInfo
					replaceLastTag(parser.swagger.Tags, tag)
				case "@tag.docs.url":
					commentInfo := strings.TrimSpace(commentLine[len(attribute):])
					tag := parser.swagger.Tags[len(parser.swagger.Tags)-1]
					tag.TagProps.ExternalDocs = &spec.ExternalDocumentation{
						URL: commentInfo,
					}
					replaceLastTag(parser.swagger.Tags, tag)

				case "@tag.docs.description":
					commentInfo := strings.TrimSpace(commentLine[len(attribute):])
					tag := parser.swagger.Tags[len(parser.swagger.Tags)-1]
					if tag.TagProps.ExternalDocs == nil {
						return errors.New("@tag.docs.description needs to come after a @tags.docs.url")
					}
					tag.TagProps.ExternalDocs.Description = commentInfo
					replaceLastTag(parser.swagger.Tags, tag)
				}
				previousAttribute = attribute
			}

			for i := 0; i < len(comments); i++ {
				attribute := strings.ToLower(strings.Split(comments[i], " ")[0])
				switch attribute {
				case "@securitydefinitions.basic":
					securityMap[strings.TrimSpace(comments[i][len(attribute):])] = spec.BasicAuth()
				case "@securitydefinitions.apikey":
					attrMap := map[string]string{}
					for _, v := range comments[i+1:] {
						securityAttr := strings.ToLower(strings.Split(v, " ")[0])
						if securityAttr == "@in" || securityAttr == "@name" {
							attrMap[securityAttr] = strings.TrimSpace(v[len(securityAttr):])
						}
						// next securityDefinitions
						if strings.Index(securityAttr, "@securitydefinitions.") == 0 {
							break
						}
					}
					if len(attrMap) != 2 {
						return errors.New("@securitydefinitions.apikey is @name and @in required")
					}
					securityMap[strings.TrimSpace(comments[i][len(attribute):])] = spec.APIKeyAuth(attrMap["@name"], attrMap["@in"])
				case "@securitydefinitions.oauth2.application":
					attrMap := map[string]string{}
					scopes := map[string]string{}
					for _, v := range comments[i+1:] {
						securityAttr := strings.ToLower(strings.Split(v, " ")[0])
						if securityAttr == "@tokenurl" {
							attrMap[securityAttr] = strings.TrimSpace(v[len(securityAttr):])
						} else {
							isExists, err := isExistsScope(securityAttr)
							if err != nil {
								return err
							}
							if isExists {
								scopScheme, err := getScopeScheme(securityAttr)
								if err != nil {
									return err
								}
								scopes[scopScheme] = v[len(securityAttr):]
							}
						}
						// next securityDefinitions
						if strings.Index(securityAttr, "@securitydefinitions.") == 0 {
							break
						}
					}
					if len(attrMap) != 1 {
						return errors.New("@securitydefinitions.oauth2.application is @tokenUrl required")
					}
					securityScheme := spec.OAuth2Application(attrMap["@tokenurl"])
					for scope, description := range scopes {
						securityScheme.AddScope(scope, description)
					}
					securityMap[strings.TrimSpace(comments[i][len(attribute):])] = securityScheme
				case "@securitydefinitions.oauth2.implicit":
					attrMap := map[string]string{}
					scopes := map[string]string{}
					for _, v := range comments[i+1:] {
						securityAttr := strings.ToLower(strings.Split(v, " ")[0])
						if securityAttr == "@authorizationurl" {
							attrMap[securityAttr] = strings.TrimSpace(v[len(securityAttr):])
						} else {
							isExists, err := isExistsScope(securityAttr)
							if err != nil {
								return err
							}
							if isExists {
								scopScheme, err := getScopeScheme(securityAttr)
								if err != nil {
									return err
								}
								scopes[scopScheme] = v[len(securityAttr):]
							}
						}
						// next securityDefinitions
						if strings.Index(securityAttr, "@securitydefinitions.") == 0 {
							break
						}
					}
					if len(attrMap) != 1 {
						return errors.New("@securitydefinitions.oauth2.implicit is @authorizationUrl required")
					}
					securityScheme := spec.OAuth2Implicit(attrMap["@authorizationurl"])
					for scope, description := range scopes {
						securityScheme.AddScope(scope, description)
					}
					securityMap[strings.TrimSpace(comments[i][len(attribute):])] = securityScheme
				case "@securitydefinitions.oauth2.password":
					attrMap := map[string]string{}
					scopes := map[string]string{}
					for _, v := range comments[i+1:] {
						securityAttr := strings.ToLower(strings.Split(v, " ")[0])
						if securityAttr == "@tokenurl" {
							attrMap[securityAttr] = strings.TrimSpace(v[len(securityAttr):])
						} else {
							isExists, err := isExistsScope(securityAttr)
							if err != nil {
								return err
							}
							if isExists {
								scopScheme, err := getScopeScheme(securityAttr)
								if err != nil {
									return err
								}
								scopes[scopScheme] = v[len(securityAttr):]
							}
						}
						// next securityDefinitions
						if strings.Index(securityAttr, "@securitydefinitions.") == 0 {
							break
						}
					}
					if len(attrMap) != 1 {
						return errors.New("@securitydefinitions.oauth2.password is @tokenUrl required")
					}
					securityScheme := spec.OAuth2Password(attrMap["@tokenurl"])
					for scope, description := range scopes {
						securityScheme.AddScope(scope, description)
					}
					securityMap[strings.TrimSpace(comments[i][len(attribute):])] = securityScheme
				case "@securitydefinitions.oauth2.accesscode":
					attrMap := map[string]string{}
					scopes := map[string]string{}
					for _, v := range comments[i+1:] {
						securityAttr := strings.ToLower(strings.Split(v, " ")[0])
						if securityAttr == "@tokenurl" || securityAttr == "@authorizationurl" {
							attrMap[securityAttr] = strings.TrimSpace(v[len(securityAttr):])
						} else {
							isExists, err := isExistsScope(securityAttr)
							if err != nil {
								return err
							}
							if isExists {
								scopScheme, err := getScopeScheme(securityAttr)
								if err != nil {
									return err
								}
								scopes[scopScheme] = v[len(securityAttr):]
							}
						}
						// next securityDefinitions
						if strings.Index(securityAttr, "@securitydefinitions.") == 0 {
							break
						}
					}
					if len(attrMap) != 2 {
						return errors.New("@securitydefinitions.oauth2.accessCode is @tokenUrl and @authorizationUrl required")
					}
					securityScheme := spec.OAuth2AccessToken(attrMap["@authorizationurl"], attrMap["@tokenurl"])
					for scope, description := range scopes {
						securityScheme.AddScope(scope, description)
					}
					securityMap[strings.TrimSpace(comments[i][len(attribute):])] = securityScheme
				}
			}
		}
	}
	if len(securityMap) > 0 {
		parser.swagger.SecurityDefinitions = securityMap
	}

	return nil
}
