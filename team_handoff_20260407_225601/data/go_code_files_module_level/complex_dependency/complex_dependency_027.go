func (g *genFakeForType) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	pkg := filepath.Base(t.Name.Package)
	tags, err := util.ParseClientGenTags(append(t.SecondClosestCommentLines, t.CommentLines...))
	if err != nil {
		return err
	}
	canonicalGroup := g.group
	if canonicalGroup == "core" {
		canonicalGroup = ""
	}

	groupName := g.group
	if g.group == "core" {
		groupName = ""
	}

	// allow user to define a group name that's different from the one parsed from the directory.
	p := c.Universe.Package(path.Vendorless(g.inputPackage))
	if override := types.ExtractCommentTags("+", p.Comments)["groupName"]; override != nil {
		groupName = override[0]
	}

	const pkgClientGoTesting = "k8s.io/client-go/testing"
	m := map[string]interface{}{
		"type":                 t,
		"inputType":            t,
		"resultType":           t,
		"subresourcePath":      "",
		"package":              pkg,
		"Package":              namer.IC(pkg),
		"namespaced":           !tags.NonNamespaced,
		"Group":                namer.IC(g.group),
		"GroupGoName":          g.groupGoName,
		"Version":              namer.IC(g.version),
		"group":                canonicalGroup,
		"groupName":            groupName,
		"version":              g.version,
		"DeleteOptions":        c.Universe.Type(types.Name{Package: "k8s.io/apimachinery/pkg/apis/meta/v1", Name: "DeleteOptions"}),
		"ListOptions":          c.Universe.Type(types.Name{Package: "k8s.io/apimachinery/pkg/apis/meta/v1", Name: "ListOptions"}),
		"GetOptions":           c.Universe.Type(types.Name{Package: "k8s.io/apimachinery/pkg/apis/meta/v1", Name: "GetOptions"}),
		"Everything":           c.Universe.Function(types.Name{Package: "k8s.io/apimachinery/pkg/labels", Name: "Everything"}),
		"GroupVersionResource": c.Universe.Type(types.Name{Package: "k8s.io/apimachinery/pkg/runtime/schema", Name: "GroupVersionResource"}),
		"GroupVersionKind":     c.Universe.Type(types.Name{Package: "k8s.io/apimachinery/pkg/runtime/schema", Name: "GroupVersionKind"}),
		"PatchType":            c.Universe.Type(types.Name{Package: "k8s.io/apimachinery/pkg/types", Name: "PatchType"}),
		"watchInterface":       c.Universe.Type(types.Name{Package: "k8s.io/apimachinery/pkg/watch", Name: "Interface"}),

		"NewRootListAction":              c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootListAction"}),
		"NewListAction":                  c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewListAction"}),
		"NewRootGetAction":               c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootGetAction"}),
		"NewGetAction":                   c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewGetAction"}),
		"NewRootDeleteAction":            c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootDeleteAction"}),
		"NewDeleteAction":                c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewDeleteAction"}),
		"NewRootDeleteCollectionAction":  c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootDeleteCollectionAction"}),
		"NewDeleteCollectionAction":      c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewDeleteCollectionAction"}),
		"NewRootUpdateAction":            c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootUpdateAction"}),
		"NewUpdateAction":                c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewUpdateAction"}),
		"NewRootCreateAction":            c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootCreateAction"}),
		"NewCreateAction":                c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewCreateAction"}),
		"NewRootWatchAction":             c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootWatchAction"}),
		"NewWatchAction":                 c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewWatchAction"}),
		"NewCreateSubresourceAction":     c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewCreateSubresourceAction"}),
		"NewRootCreateSubresourceAction": c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootCreateSubresourceAction"}),
		"NewUpdateSubresourceAction":     c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewUpdateSubresourceAction"}),
		"NewGetSubresourceAction":        c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewGetSubresourceAction"}),
		"NewRootGetSubresourceAction":    c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootGetSubresourceAction"}),
		"NewRootUpdateSubresourceAction": c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootUpdateSubresourceAction"}),
		"NewRootPatchAction":             c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootPatchAction"}),
		"NewPatchAction":                 c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewPatchAction"}),
		"NewRootPatchSubresourceAction":  c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewRootPatchSubresourceAction"}),
		"NewPatchSubresourceAction":      c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "NewPatchSubresourceAction"}),
		"ExtractFromListOptions":         c.Universe.Function(types.Name{Package: pkgClientGoTesting, Name: "ExtractFromListOptions"}),
	}

	if tags.NonNamespaced {
		sw.Do(structNonNamespaced, m)
	} else {
		sw.Do(structNamespaced, m)
	}

	if tags.NoVerbs {
		return sw.Error()
	}
	sw.Do(resource, m)
	sw.Do(kind, m)

	if tags.HasVerb("get") {
		sw.Do(getTemplate, m)
	}
	if tags.HasVerb("list") {
		if hasObjectMeta(t) {
			sw.Do(listUsingOptionsTemplate, m)
		} else {
			sw.Do(listTemplate, m)
		}
	}
	if tags.HasVerb("watch") {
		sw.Do(watchTemplate, m)
	}

	if tags.HasVerb("create") {
		sw.Do(createTemplate, m)
	}
	if tags.HasVerb("update") {
		sw.Do(updateTemplate, m)
	}
	if tags.HasVerb("updateStatus") && genStatus(t) {
		sw.Do(updateStatusTemplate, m)
	}
	if tags.HasVerb("delete") {
		sw.Do(deleteTemplate, m)
	}
	if tags.HasVerb("deleteCollection") {
		sw.Do(deleteCollectionTemplate, m)
	}
	if tags.HasVerb("patch") {
		sw.Do(patchTemplate, m)
	}

	// generate extended client methods
	for _, e := range tags.Extensions {
		inputType := *t
		resultType := *t
		if len(e.InputTypeOverride) > 0 {
			if name, pkg := e.Input(); len(pkg) > 0 {
				newType := c.Universe.Type(types.Name{Package: pkg, Name: name})
				inputType = *newType
			} else {
				inputType.Name.Name = e.InputTypeOverride
			}
		}
		if len(e.ResultTypeOverride) > 0 {
			if name, pkg := e.Result(); len(pkg) > 0 {
				newType := c.Universe.Type(types.Name{Package: pkg, Name: name})
				resultType = *newType
			} else {
				resultType.Name.Name = e.ResultTypeOverride
			}
		}
		m["inputType"] = &inputType
		m["resultType"] = &resultType
		m["subresourcePath"] = e.SubResourcePath

		if e.HasVerb("get") {
			if e.IsSubresource() {
				sw.Do(adjustTemplate(e.VerbName, e.VerbType, getSubresourceTemplate), m)
			} else {
				sw.Do(adjustTemplate(e.VerbName, e.VerbType, getTemplate), m)
			}
		}

		if e.HasVerb("list") {

			sw.Do(adjustTemplate(e.VerbName, e.VerbType, listTemplate), m)
		}

		// TODO: Figure out schemantic for watching a sub-resource.
		if e.HasVerb("watch") {
			sw.Do(adjustTemplate(e.VerbName, e.VerbType, watchTemplate), m)
		}

		if e.HasVerb("create") {
			if e.IsSubresource() {
				sw.Do(adjustTemplate(e.VerbName, e.VerbType, createSubresourceTemplate), m)
			} else {
				sw.Do(adjustTemplate(e.VerbName, e.VerbType, createTemplate), m)
			}
		}

		if e.HasVerb("update") {
			if e.IsSubresource() {
				sw.Do(adjustTemplate(e.VerbName, e.VerbType, updateSubresourceTemplate), m)
			} else {
				sw.Do(adjustTemplate(e.VerbName, e.VerbType, updateTemplate), m)
			}
		}

		// TODO: Figure out schemantic for deleting a sub-resource (what arguments
		// are passed, does it need two names? etc.
		if e.HasVerb("delete") {
			sw.Do(adjustTemplate(e.VerbName, e.VerbType, deleteTemplate), m)
		}

		if e.HasVerb("patch") {
			sw.Do(adjustTemplate(e.VerbName, e.VerbType, patchTemplate), m)
		}
	}

	return sw.Error()
}
