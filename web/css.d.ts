// Next types only cover *.module.css. Plain stylesheets are imported for
// their side effect (see app/layout.tsx), which TypeScript rejects when
// noUncheckedSideEffectImports is enabled.
declare module "*.css";
