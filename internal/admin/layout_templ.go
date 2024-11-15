// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package admin

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func AdminLayout(title string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/admin/layout.templ`, Line: 9, Col: 17}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><script src=\"https://unpkg.com/htmx.org@1.9.6\"></script><script src=\"https://cdn.tailwindcss.com\"></script><script src=\"/static/js/admin.js\"></script><style>\n\t\t@keyframes slide-in {\n\t\t\tfrom {\n\t\t\t\topacity: 0;\n\t\t\t\ttransform: translateY(-10px);\n\t\t\t}\n\n\t\t\tto {\n\t\t\t\topacity: 1;\n\t\t\t\ttransform: translateY(0);\n\t\t\t}\n\t\t}\n\n\t\t.animate-slide-in {\n\t\t\tanimation: slide-in 0.3s ease-out forwards;\n\t\t}\n\n\t\t.player-enter {\n\t\t\topacity: 0;\n\t\t\ttransform: translateY(-10px);\n\t\t}\n\n\t\t.player-enter-active {\n\t\t\topacity: 1;\n\t\t\ttransform: translateY(0);\n\t\t\ttransition: opacity 300ms, transform 300ms;\n\t\t}\n\n\t\t.player-exit {\n\t\t\topacity: 1;\n\t\t}\n\n\t\t.player-exit-active {\n\t\t\topacity: 0;\n\t\t\ttransform: translateY(-10px);\n\t\t\ttransition: opacity 300ms, transform 300ms;\n\t\t}\n\t</style></head><body><div class=\"min-h-screen bg-gray-100\"><nav class=\"bg-white shadow-lg\"><div class=\"max-w-7xl mx-auto px-4\"><div class=\"flex justify-between h-16\"><div class=\"flex\"><div class=\"flex-shrink-0 flex items-center\"><span class=\"text-xl font-bold\">Admin Dashboard</span></div></div></div></div></nav>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
