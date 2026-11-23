package docs

// Este archivo NO es generado por `swag init`.
//
// Tema oscuro para Swagger UI (servido en /index.html y /es/index.html) con la
// misma paleta que el panel web (tema "dark" por defecto de Vuetify):
//
//	fondo      #121212   superficie #212121   superficie clara #424242
//	texto      #e6e6e6   texto sec. #b0b0b0   enlaces          #64b5f6
//	primario   #2196f3   (botón "Execute" #1976d2 para contraste AA con blanco)
//	GET #64b5f6 · POST #4caf50 · PUT #fb8c00 · PATCH #54b6b2 · DELETE #cf6679
//
// Todas las reglas usan el prefijo "#swagger-ui .swagger-ui" (incluye un ID),
// por lo que siempre ganan sobre swagger-ui.css (que solo usa clases) sin
// depender del orden de carga ni de "prefers-color-scheme": el tema es oscuro
// tanto si el sistema operativo está en modo claro como en modo oscuro.
// "color-scheme: dark" además evita que el navegador aplique su propio
// oscurecimiento automático (p. ej. "Auto Dark Mode" de Chrome), que es lo que
// producía la mezcla ilegible de fondos negros con texto gris oscuro.

// SwaggerThemeCSS es la hoja de estilos que se pasa a swagger.Config.CustomStyle
// (como template.CSS(docs.SwaggerThemeCSS)).
const SwaggerThemeCSS = `
:root {
  color-scheme: dark;
  --hs-bg: #121212;
  --hs-surface: #212121;
  --hs-surface-2: #1a1a1a;
  --hs-surface-light: #424242;
  --hs-code-bg: #0d0d0d;
  --hs-border: rgba(255, 255, 255, 0.12);
  --hs-border-strong: #6b6b6b;
  --hs-text: #e6e6e6;
  --hs-text-muted: #b0b0b0;
  --hs-text-subtle: #9e9e9e;
  --hs-link: #64b5f6;
  --hs-link-hover: #90caf9;
  --hs-primary: #2196f3;
  --hs-primary-strong: #1976d2;
  --hs-success: #66bb6a;
  --hs-error: #cf6679;
  --hs-error-text: #ef9a9a;
  --hs-get: #64b5f6;
  --hs-post: #4caf50;
  --hs-put: #fb8c00;
  --hs-patch: #54b6b2;
  --hs-delete: #cf6679;
  --hs-head: #b39ddb;
  --hs-options: #90a4ae;
  --hs-deprecated: #9e9e9e;
  --hs-on-method: #121212;
}

html,
body {
  background: var(--hs-bg) !important;
  color: var(--hs-text);
  scrollbar-color: var(--hs-surface-light) var(--hs-bg);
}

::selection {
  background: rgba(33, 150, 243, 0.4);
  color: #ffffff;
}

::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-track {
  background: var(--hs-bg);
}

::-webkit-scrollbar-thumb {
  background: var(--hs-surface-light);
  border: 2px solid var(--hs-bg);
  border-radius: 6px;
}

::-webkit-scrollbar-thumb:hover {
  background: #5a5a5a;
}

/* ---------------------------------------------------------------- base ---- */
#swagger-ui .swagger-ui {
  color: var(--hs-text);
  background: var(--hs-bg);
}

/* Barra superior de Swagger (logo "Swagger / Supported by SMARTBEAR" y el
   campo para explorar otra URL): no aporta nada en una instalación propia y
   chocaba con el selector de idioma fijo, así que se oculta. */
#swagger-ui .swagger-ui .topbar {
  display: none;
}

#swagger-ui .swagger-ui a.nostyle,
#swagger-ui .swagger-ui a.nostyle:visited {
  color: inherit;
}

#swagger-ui .swagger-ui .markdown a,
#swagger-ui .swagger-ui .renderedMarkdown a,
#swagger-ui .swagger-ui .info a,
#swagger-ui .swagger-ui .scopes h2 a,
#swagger-ui .swagger-ui .opblock-external-docs-wrapper a,
#swagger-ui .swagger-ui .response-col_links a {
  color: var(--hs-link);
}

#swagger-ui .swagger-ui .markdown a:hover,
#swagger-ui .swagger-ui .renderedMarkdown a:hover,
#swagger-ui .swagger-ui .info a:hover {
  color: var(--hs-link-hover);
}

/* Texto que swagger-ui pinta en gris oscuro (#3b4151) por defecto. */
#swagger-ui .swagger-ui section h3,
#swagger-ui .swagger-ui label,
#swagger-ui .swagger-ui .opblock-tag,
#swagger-ui .swagger-ui .opblock .opblock-section-header h4,
#swagger-ui .swagger-ui .opblock .opblock-section-header > label,
#swagger-ui .swagger-ui .opblock .opblock-summary-operation-id,
#swagger-ui .swagger-ui .opblock .opblock-summary-path,
#swagger-ui .swagger-ui .opblock .opblock-summary-path__deprecated,
#swagger-ui .swagger-ui .opblock-description-wrapper,
#swagger-ui .swagger-ui .opblock-description-wrapper h4,
#swagger-ui .swagger-ui .opblock-description-wrapper p,
#swagger-ui .swagger-ui .opblock-external-docs-wrapper,
#swagger-ui .swagger-ui .opblock-external-docs-wrapper h4,
#swagger-ui .swagger-ui .opblock-external-docs-wrapper p,
#swagger-ui .swagger-ui .opblock-title_normal,
#swagger-ui .swagger-ui .opblock-title_normal h4,
#swagger-ui .swagger-ui .opblock-title_normal p,
#swagger-ui .swagger-ui .responses-inner h4,
#swagger-ui .swagger-ui .responses-inner h5,
#swagger-ui .swagger-ui .response-col_status,
#swagger-ui .swagger-ui .response-col_links,
#swagger-ui .swagger-ui .parameter__name,
#swagger-ui .swagger-ui .scheme-container .schemes > label,
#swagger-ui .swagger-ui .scheme-container .schemes > .schemes-server-container > label,
#swagger-ui .swagger-ui .servers > label,
#swagger-ui .swagger-ui .servers-title,
#swagger-ui .swagger-ui .schemes-title,
#swagger-ui .swagger-ui table.headers td,
#swagger-ui .swagger-ui .checkbox p,
#swagger-ui .swagger-ui .markdown p,
#swagger-ui .swagger-ui .markdown pre,
#swagger-ui .swagger-ui .markdown li,
#swagger-ui .swagger-ui .renderedMarkdown p,
#swagger-ui .swagger-ui .renderedMarkdown pre,
#swagger-ui .swagger-ui .renderedMarkdown li,
#swagger-ui .swagger-ui .model,
#swagger-ui .swagger-ui .model-title,
#swagger-ui .swagger-ui .errors-wrapper .errors h4,
#swagger-ui .swagger-ui .errors-wrapper hgroup h4,
#swagger-ui .swagger-ui .loading-container .loading:after,
#swagger-ui .swagger-ui .json-schema-2020-12__title,
#swagger-ui .swagger-ui .json-schema-2020-12-property .json-schema-2020-12__title,
#swagger-ui .swagger-ui .json-schema-2020-12-keyword__name--primary,
#swagger-ui .swagger-ui .json-schema-2020-12-keyword__value--primary,
#swagger-ui .swagger-ui .json-schema-2020-12__attribute {
  color: var(--hs-text);
}

/* Texto secundario. */
#swagger-ui .swagger-ui .opblock-tag small,
#swagger-ui .swagger-ui .opblock .opblock-summary-description,
#swagger-ui .swagger-ui .parameter__type,
#swagger-ui .swagger-ui .parameter__extension,
#swagger-ui .swagger-ui .parameter__enum,
#swagger-ui .swagger-ui .prop-format,
#swagger-ui .swagger-ui .prop-enum,
#swagger-ui .swagger-ui .model .property,
#swagger-ui .swagger-ui .model .property.primitive,
#swagger-ui .swagger-ui .model .external-docs,
#swagger-ui .swagger-ui table.model tr.description,
#swagger-ui .swagger-ui table.model tr.extension,
#swagger-ui .swagger-ui table.headers .header-example,
#swagger-ui .swagger-ui span.token-string,
#swagger-ui .swagger-ui span.token-not-formatted,
#swagger-ui .swagger-ui .errors-wrapper .errors small,
#swagger-ui .swagger-ui section.models h5,
#swagger-ui .swagger-ui .json-schema-2020-12-keyword__name--secondary,
#swagger-ui .swagger-ui .json-schema-2020-12-keyword__value,
#swagger-ui .swagger-ui .json-schema-2020-12-keyword__value--secondary {
  color: var(--hs-text-muted);
}

#swagger-ui .swagger-ui .parameter__in,
#swagger-ui .swagger-ui .response-col_status .response-undocumented,
#swagger-ui .swagger-ui .response-col_links .response-undocumented,
#swagger-ui .swagger-ui .model .deprecated span,
#swagger-ui .swagger-ui .model .deprecated td {
  color: var(--hs-text-subtle) !important;
}

#swagger-ui .swagger-ui .prop-type,
#swagger-ui .swagger-ui .json-schema-2020-12__attribute--primary {
  color: #90caf9;
}

#swagger-ui .swagger-ui .parameter__name.required span,
#swagger-ui .swagger-ui .parameter__name.required:after,
#swagger-ui .swagger-ui .parameter__deprecated,
#swagger-ui .swagger-ui table.model tr.property-row .star {
  color: var(--hs-error-text);
}

#swagger-ui .swagger-ui .markdown code,
#swagger-ui .swagger-ui .renderedMarkdown code,
#swagger-ui .swagger-ui .version-pragma__message code {
  background: rgba(255, 255, 255, 0.08);
  color: #ce93d8;
}

/* --------------------------------------------------------------- info ---- */
#swagger-ui .swagger-ui .info .title,
#swagger-ui .swagger-ui .info h1,
#swagger-ui .swagger-ui .info h2,
#swagger-ui .swagger-ui .info h3,
#swagger-ui .swagger-ui .info h4,
#swagger-ui .swagger-ui .info h5 {
  color: #ffffff;
}

#swagger-ui .swagger-ui .info p,
#swagger-ui .swagger-ui .info li,
#swagger-ui .swagger-ui .info table,
#swagger-ui .swagger-ui .info .description,
#swagger-ui .swagger-ui .info .renderedMarkdown p {
  color: var(--hs-text);
}

#swagger-ui .swagger-ui .info .base-url,
#swagger-ui .swagger-ui .info .url {
  color: var(--hs-text-muted);
}

#swagger-ui .swagger-ui .info .title small {
  background: var(--hs-surface-light);
}

#swagger-ui .swagger-ui .info .title small.version-stamp {
  background-color: #2e7d32;
}

#swagger-ui .swagger-ui .info .title small pre {
  color: #ffffff;
}

/* ------------------------------------------- esquemas + botón Authorize ---- */
#swagger-ui .swagger-ui .scheme-container {
  background: var(--hs-surface);
  border-top: 1px solid var(--hs-border);
  border-bottom: 1px solid var(--hs-border);
  box-shadow: none;
}

#swagger-ui .swagger-ui .btn {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.35);
  box-shadow: none;
  color: var(--hs-text);
}

#swagger-ui .swagger-ui .btn:hover {
  background: rgba(255, 255, 255, 0.06);
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.12);
}

#swagger-ui .swagger-ui .btn:focus-visible,
#swagger-ui .swagger-ui select:focus-visible,
#swagger-ui .swagger-ui .opblock-summary-control:focus-visible {
  outline: 2px solid var(--hs-primary);
  outline-offset: 2px;
}

#swagger-ui .swagger-ui .btn.authorize {
  border-color: var(--hs-success);
  color: var(--hs-success);
}

#swagger-ui .swagger-ui .btn.authorize svg {
  fill: var(--hs-success);
}

#swagger-ui .swagger-ui .btn.execute {
  background-color: var(--hs-primary-strong);
  border-color: var(--hs-primary-strong);
  color: #ffffff;
}

#swagger-ui .swagger-ui .btn.execute:hover {
  background-color: #1565c0;
  border-color: #1565c0;
}

#swagger-ui .swagger-ui .btn.cancel {
  border-color: var(--hs-error);
  color: var(--hs-error-text);
}

#swagger-ui .swagger-ui .btn.cancel:hover {
  background: rgba(207, 102, 121, 0.12);
}

#swagger-ui .swagger-ui button.invalid {
  background: rgba(207, 102, 121, 0.15);
  border-color: var(--hs-error);
}

#swagger-ui .swagger-ui .download-contents,
#swagger-ui .swagger-ui .copy-to-clipboard {
  background: var(--hs-surface-light);
  color: #ffffff;
}

#swagger-ui .swagger-ui .copy-to-clipboard:active {
  background: #5a5a5a;
}

/* Iconos SVG que swagger-ui deja en negro (relleno por defecto). */
#swagger-ui .swagger-ui .arrow,
#swagger-ui .swagger-ui .opblock-control-arrow svg,
#swagger-ui .swagger-ui .expand-operation svg,
#swagger-ui .swagger-ui .expand-methods svg,
#swagger-ui .swagger-ui .models-control svg,
#swagger-ui .swagger-ui .model-box-control svg,
#swagger-ui .swagger-ui .opblock-tag svg,
#swagger-ui .swagger-ui section.models h4 svg,
#swagger-ui .swagger-ui .authorization__btn svg,
#swagger-ui .swagger-ui .json-schema-2020-12-accordion__icon svg,
#swagger-ui .swagger-ui .dialog-ux .modal-ux-header .close-modal svg {
  fill: var(--hs-text-muted);
}

#swagger-ui .swagger-ui .expand-methods:hover svg,
#swagger-ui .swagger-ui .authorization__btn:hover svg,
#swagger-ui .swagger-ui .dialog-ux .modal-ux-header .close-modal:hover svg {
  fill: var(--hs-text);
}

#swagger-ui .swagger-ui .model-toggle:after {
  filter: invert(0.85);
}

/* ------------------------------------------------- modal de autorización ---- */
#swagger-ui .swagger-ui .dialog-ux .backdrop-ux {
  background: rgba(0, 0, 0, 0.75);
}

#swagger-ui .swagger-ui .dialog-ux .modal-ux {
  background: var(--hs-surface);
  border: 1px solid var(--hs-border);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.6);
  color: var(--hs-text);
}

#swagger-ui .swagger-ui .dialog-ux .modal-ux-header {
  border-bottom: 1px solid var(--hs-border);
}

#swagger-ui .swagger-ui .dialog-ux .modal-ux-header h3,
#swagger-ui .swagger-ui .dialog-ux .modal-ux-content h4,
#swagger-ui .swagger-ui .dialog-ux .modal-ux-content p,
#swagger-ui .swagger-ui .auth-container h4,
#swagger-ui .swagger-ui .auth-container h5,
#swagger-ui .swagger-ui .auth-container h6,
#swagger-ui .swagger-ui .auth-container label,
#swagger-ui .swagger-ui .auth-container p,
#swagger-ui .swagger-ui .scopes h2 {
  color: var(--hs-text);
}

#swagger-ui .swagger-ui .auth-container {
  border-bottom: 1px solid var(--hs-border);
}

#swagger-ui .swagger-ui .auth-container code {
  background: rgba(255, 255, 255, 0.08);
  color: #ce93d8;
}

#swagger-ui .swagger-ui .auth-container .errors {
  background: rgba(207, 102, 121, 0.15);
  color: var(--hs-error-text);
}

/* --------------------------------------------- grupos y operaciones ---- */
#swagger-ui .swagger-ui .opblock-tag {
  border-bottom: 1px solid var(--hs-border);
}

#swagger-ui .swagger-ui .opblock-tag:hover,
#swagger-ui .swagger-ui section.models h4:hover {
  background: rgba(255, 255, 255, 0.04);
}

#swagger-ui .swagger-ui .opblock {
  background: var(--hs-surface);
  border: 1px solid var(--hs-border);
  box-shadow: none;
}

#swagger-ui .swagger-ui .opblock .opblock-summary-method {
  color: var(--hs-on-method);
  text-shadow: none;
}

#swagger-ui .swagger-ui .opblock.is-open .opblock-summary {
  border-bottom-color: var(--hs-border);
}

#swagger-ui .swagger-ui .opblock .opblock-section-header {
  background: rgba(255, 255, 255, 0.04);
  box-shadow: none;
  border-bottom: 1px solid var(--hs-border);
}

#swagger-ui .swagger-ui .opblock .tab-header .tab-item.active h4 span:after {
  background: var(--hs-primary);
}

#swagger-ui .swagger-ui .opblock.opblock-get {
  background: rgba(100, 181, 246, 0.07);
  border-color: rgba(100, 181, 246, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-get .opblock-summary-method,
#swagger-ui .swagger-ui .opblock.opblock-get .tab-header .tab-item.active h4 span:after {
  background: var(--hs-get);
}

#swagger-ui .swagger-ui .opblock.opblock-get .opblock-summary {
  border-color: rgba(100, 181, 246, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-post {
  background: rgba(76, 175, 80, 0.07);
  border-color: rgba(76, 175, 80, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-post .opblock-summary-method,
#swagger-ui .swagger-ui .opblock.opblock-post .tab-header .tab-item.active h4 span:after {
  background: var(--hs-post);
}

#swagger-ui .swagger-ui .opblock.opblock-post .opblock-summary {
  border-color: rgba(76, 175, 80, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-put {
  background: rgba(251, 140, 0, 0.07);
  border-color: rgba(251, 140, 0, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-put .opblock-summary-method,
#swagger-ui .swagger-ui .opblock.opblock-put .tab-header .tab-item.active h4 span:after {
  background: var(--hs-put);
}

#swagger-ui .swagger-ui .opblock.opblock-put .opblock-summary {
  border-color: rgba(251, 140, 0, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-patch {
  background: rgba(84, 182, 178, 0.07);
  border-color: rgba(84, 182, 178, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-patch .opblock-summary-method,
#swagger-ui .swagger-ui .opblock.opblock-patch .tab-header .tab-item.active h4 span:after {
  background: var(--hs-patch);
}

#swagger-ui .swagger-ui .opblock.opblock-patch .opblock-summary {
  border-color: rgba(84, 182, 178, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-delete {
  background: rgba(207, 102, 121, 0.07);
  border-color: rgba(207, 102, 121, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-delete .opblock-summary-method,
#swagger-ui .swagger-ui .opblock.opblock-delete .tab-header .tab-item.active h4 span:after {
  background: var(--hs-delete);
}

#swagger-ui .swagger-ui .opblock.opblock-delete .opblock-summary {
  border-color: rgba(207, 102, 121, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-head {
  background: rgba(179, 157, 219, 0.07);
  border-color: rgba(179, 157, 219, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-head .opblock-summary-method,
#swagger-ui .swagger-ui .opblock.opblock-head .tab-header .tab-item.active h4 span:after {
  background: var(--hs-head);
}

#swagger-ui .swagger-ui .opblock.opblock-head .opblock-summary {
  border-color: rgba(179, 157, 219, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-options {
  background: rgba(144, 164, 174, 0.07);
  border-color: rgba(144, 164, 174, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-options .opblock-summary-method,
#swagger-ui .swagger-ui .opblock.opblock-options .tab-header .tab-item.active h4 span:after {
  background: var(--hs-options);
}

#swagger-ui .swagger-ui .opblock.opblock-options .opblock-summary {
  border-color: rgba(144, 164, 174, 0.55);
}

#swagger-ui .swagger-ui .opblock.opblock-deprecated {
  background: rgba(158, 158, 158, 0.06);
  border-color: rgba(158, 158, 158, 0.45);
  opacity: 0.75;
}

#swagger-ui .swagger-ui .opblock.opblock-deprecated .opblock-summary-method,
#swagger-ui .swagger-ui .opblock.opblock-deprecated .tab-header .tab-item.active h4 span:after {
  background: var(--hs-deprecated);
}

#swagger-ui .swagger-ui .opblock.opblock-deprecated .opblock-summary {
  border-color: rgba(158, 158, 158, 0.45);
}

/* Pestañas "Example Value | Model" y similares. */
#swagger-ui .swagger-ui .tab li {
  color: var(--hs-text-muted);
}

#swagger-ui .swagger-ui .tab li.active {
  color: var(--hs-text);
}

#swagger-ui .swagger-ui .tab li:first-of-type:after {
  background: var(--hs-border);
}

/* Tablas de parámetros y respuestas. */
#swagger-ui .swagger-ui table thead tr td,
#swagger-ui .swagger-ui table thead tr th {
  border-bottom: 1px solid var(--hs-border);
  color: var(--hs-text-muted);
}

#swagger-ui .swagger-ui .response-control-media-type__accept-message {
  color: #81c784;
}

#swagger-ui .swagger-ui .response-control-media-type--accept-controller select {
  border-color: var(--hs-success);
}

/* --------------------------------------------- campos de formulario ---- */
#swagger-ui .swagger-ui input[type="text"],
#swagger-ui .swagger-ui input[type="password"],
#swagger-ui .swagger-ui input[type="search"],
#swagger-ui .swagger-ui input[type="email"],
#swagger-ui .swagger-ui input[type="file"],
#swagger-ui .swagger-ui .filter .operation-filter-input,
#swagger-ui .swagger-ui textarea {
  background: var(--hs-surface-2);
  border: 1px solid var(--hs-border-strong);
  color: var(--hs-text);
}

#swagger-ui .swagger-ui input::placeholder,
#swagger-ui .swagger-ui textarea::placeholder {
  color: var(--hs-text-subtle);
  opacity: 1;
}

#swagger-ui .swagger-ui input:focus,
#swagger-ui .swagger-ui textarea:focus,
#swagger-ui .swagger-ui select:focus {
  border-color: var(--hs-primary);
  outline: none;
  box-shadow: 0 0 0 1px var(--hs-primary);
}

#swagger-ui .swagger-ui select {
  background: var(--hs-surface-2)
    url("data:image/svg+xml;charset=utf-8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20'><path fill='%23e6e6e6' d='M13.418 7.859a.695.695 0 0 1 .978 0 .68.68 0 0 1 0 .969l-3.908 3.83a.697.697 0 0 1-.979 0l-3.908-3.83a.68.68 0 0 1 0-.969.695.695 0 0 1 .978 0L10 11z'/></svg>")
    right 10px center no-repeat;
  background-size: 20px;
  border: 1px solid var(--hs-border-strong);
  box-shadow: none;
  color: var(--hs-text);
}

#swagger-ui .swagger-ui select[multiple] {
  background: var(--hs-surface-2);
}

#swagger-ui .swagger-ui select option {
  background: var(--hs-surface);
  color: var(--hs-text);
}

#swagger-ui .swagger-ui input[disabled],
#swagger-ui .swagger-ui select[disabled],
#swagger-ui .swagger-ui textarea[disabled] {
  background-color: var(--hs-surface-2);
  border-color: var(--hs-surface-light);
  color: var(--hs-text-subtle);
}

#swagger-ui .swagger-ui textarea.curl {
  background: var(--hs-code-bg);
  border: 1px solid var(--hs-border);
  color: var(--hs-text);
}

#swagger-ui .swagger-ui input.invalid,
#swagger-ui .swagger-ui select.invalid,
#swagger-ui .swagger-ui textarea.invalid {
  background-color: rgba(207, 102, 121, 0.15);
  border-color: var(--hs-error);
}

/* ------------------------------------------ bloques de código / JSON ---- */
#swagger-ui .swagger-ui .microlight,
#swagger-ui .swagger-ui .highlight-code > .microlight,
#swagger-ui .swagger-ui .opblock-body pre.microlight {
  background: var(--hs-code-bg) !important;
  border: 1px solid var(--hs-border);
  color: var(--hs-text) !important;
}

#swagger-ui .swagger-ui .microlight code {
  background: transparent !important;
}

/* ------------------------------------------------------------ modelos ---- */
#swagger-ui .swagger-ui section.models {
  border: 1px solid var(--hs-border);
}

#swagger-ui .swagger-ui section.models.is-open h4 {
  border-bottom: 1px solid var(--hs-border);
}

#swagger-ui .swagger-ui section.models h4 {
  color: var(--hs-text);
}

#swagger-ui .swagger-ui section.models .model-container {
  background: rgba(255, 255, 255, 0.03);
}

#swagger-ui .swagger-ui section.models .model-container:hover {
  background: rgba(255, 255, 255, 0.06);
}

#swagger-ui .swagger-ui .model-box {
  background: rgba(255, 255, 255, 0.04);
}

#swagger-ui .swagger-ui section.models .model-box {
  background: none;
}

#swagger-ui .swagger-ui .model-hint {
  background: var(--hs-surface-light);
  color: var(--hs-text);
}

#swagger-ui .swagger-ui .json-schema-2020-12-keyword__children,
#swagger-ui .swagger-ui .json-schema-2020-12-keyword--\$vocabulary ul {
  border-left-color: var(--hs-border);
}

/* -------------------------------------------- errores y carga ---- */
#swagger-ui .swagger-ui .errors-wrapper {
  background: rgba(207, 102, 121, 0.1);
  border-color: var(--hs-error);
}

#swagger-ui .swagger-ui .loading-container .loading:before {
  border-color: rgba(255, 255, 255, 0.1);
  border-top-color: rgba(255, 255, 255, 0.6);
}

/* ------------------------------------------- selector de idioma ---- */
/* El script de container.go crea un <div> fijo arriba a la derecha con
   estilos en línea (fondo blanco). Se cubre tanto por clase (si se le asigna
   "hs-lang-switcher") como por posición en el DOM (es el único <div> hijo de
   <body> además de #swagger-ui), con !important para vencer al estilo en línea. */
body > div:not(#swagger-ui),
.hs-lang-switcher {
  background: var(--hs-surface) !important;
  border: 1px solid var(--hs-border) !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5) !important;
  color: var(--hs-text-muted) !important;
}

body > div:not(#swagger-ui) a,
.hs-lang-switcher a {
  color: var(--hs-link) !important;
  text-decoration: none;
}

body > div:not(#swagger-ui) a:hover,
.hs-lang-switcher a:hover {
  color: var(--hs-link-hover) !important;
  text-decoration: underline;
}

body > div:not(#swagger-ui) b,
.hs-lang-switcher b {
  color: #ffffff !important;
}
`

// SwaggerThemeScript es un JavaScript opcional para anteponer a CustomScript.
// Declara el esquema de color oscuro con <meta> (además del "color-scheme" del
// CSS) para que Chrome no aplique su oscurecimiento automático, y añade la
// marca "darkreader-lock" para que la extensión Dark Reader no recoloree la
// página (ya es oscura).
const SwaggerThemeScript = `
(function () {
	var metas = [["color-scheme", "dark"], ["darkreader-lock", ""]];
	metas.forEach(function (pair) {
		var meta = document.createElement("meta");
		meta.name = pair[0];
		meta.content = pair[1];
		document.head.appendChild(meta);
	});
})();
`
