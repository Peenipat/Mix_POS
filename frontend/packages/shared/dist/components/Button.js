var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __rest = (this && this.__rest) || function (s, e) {
    var t = {};
    for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0)
        t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function")
        for (var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++) {
            if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i]))
                t[p[i]] = s[p[i]];
        }
    return t;
};
import { jsx as _jsx } from "react/jsx-runtime";
export function Button(_a) {
    var _b = _a.variant, variant = _b === void 0 ? "primary" : _b, _c = _a.className, className = _c === void 0 ? "" : _c, props = __rest(_a, ["variant", "className"]);
    var base = "px-4 py-2 rounded text-white";
    var colors = variant === "primary"
        ? "bg-blue-600 hover:bg-blue-700"
        : "bg-gray-600 hover:bg-gray-700";
    // แก้สัญลักษณ์ template literal ไม่ต้อง escape
    return (_jsx("button", __assign({ className: "".concat(base, " ").concat(colors, " ").concat(className) }, props)));
}
