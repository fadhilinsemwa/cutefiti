/*! Micro front-end */
"use strict";
(self['webpackChunkassessments'] = self['webpackChunkassessments'] || []).push([["mupayapi/routes/page"], {
"./src/routes/mupayapi/routes/index.css?ef90": (function (module, __webpack_exports__, __webpack_require__) {
__webpack_require__.r(__webpack_exports__);
// extracted by css-extract-rspack-plugin

    if(true) {
      (function() {
        var localsJsonString = undefined;
        // 1730795010571
        var cssReload = (__webpack_require__("./node_modules/.pnpm/@rspack+core@1.0.5_@swc+helpers@0.5.13/node_modules/@rspack/core/dist/builtin-plugin/css-extract/hmr/hotModuleReplacement.js")/* .cssReload */.cssReload)(module.id, {});
        // only invalidate when locals change
        if (
          module.hot.data &&
          module.hot.data.value &&
          module.hot.data.value !== localsJsonString
        ) {
          module.hot.invalidate();
        } else {
          module.hot.accept();
        }
        module.hot.dispose(function(data) {
          data.value = localsJsonString;
          cssReload();
        });
      })();
    }
  

}),
"./src/routes/mupayapi/routes/page.tsx": (function (module, __webpack_exports__, __webpack_require__) {
__webpack_require__.r(__webpack_exports__);
/* harmony import */var react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__("./node_modules/.pnpm/react@18.2.0/node_modules/react/jsx-dev-runtime.js");
/* harmony import */var _index_css__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__("./src/routes/mupayapi/routes/index.css?ef90");
/* provided dependency */ var $ReactRefreshRuntime$ = __webpack_require__("./node_modules/.pnpm/@rspack+plugin-react-refresh@1.0.0_react-refresh@0.14.2/node_modules/@rspack/plugin-react-refresh/client/reactRefresh.js");
// src/routes/page.tsx
var _this = undefined;


var HomePage = function() {
    return /*#__PURE__*/ (0,react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxDEV)("div", {
        className: "home-page",
        children: /*#__PURE__*/ (0,react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxDEV)("h1", {
            children: "Welcome to the Home Page"
        }, void 0, false, {
            fileName: "/home/fadhili/Documents/apps/ex/exf/src/routes/mupayapi/routes/page.tsx",
            lineNumber: 7,
            columnNumber: 7
        }, _this)
    }, void 0, false, {
        fileName: "/home/fadhili/Documents/apps/ex/exf/src/routes/mupayapi/routes/page.tsx",
        lineNumber: 6,
        columnNumber: 5
    }, _this);
};
_c = HomePage;
/* harmony default export */ __webpack_exports__["default"] = (HomePage);
var _c;
$RefreshReg$(_c, "HomePage");

function $RefreshSig$() {
  return $ReactRefreshRuntime$.createSignatureFunctionForTransform();
}
function $RefreshReg$(type, id) {
  $ReactRefreshRuntime$.register(type, module.id + "_" + id);
}
Promise.resolve().then(function() {
  $ReactRefreshRuntime$.refresh(module.id, module.hot);
});


}),

}]);
//# sourceMappingURL=page.js.map