(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([[22],{Vl52:function(e,t,n){"use strict";var r=n("g09b");Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var a=r(n("p0pE")),u=r(n("d6i3")),i=n("7DNP"),c=n("WHbk"),o=n("sH+L"),s={namespace:"userLogin",state:{status:void 0},effects:{login:u.default.mark(function e(t,n){var r,a,s,p,l,f,d,h;return u.default.wrap(function(e){while(1)switch(e.prev=e.next){case 0:return r=t.payload,a=n.call,s=n.put,e.next=4,a(c.fakeAccountLogin,r);case 4:return p=e.sent,e.next=7,s({type:"changeLoginStatus",payload:p});case 7:if("ok"!==p.status){e.next=22;break}if(l=new URL(window.location.href),f=(0,o.getPageQuery)(),d=f.redirect,!d){e.next=20;break}if(h=new URL(d),h.origin!==l.origin){e.next=18;break}d=d.substr(l.origin.length),d.match(/^\/.*#/)&&(d=d.substr(d.indexOf("#")+1)),e.next=20;break;case 18:return window.location.href=d,e.abrupt("return");case 20:return e.next=22,s(i.routerRedux.replace(d||"/"));case 22:case"end":return e.stop()}},e)}),getCaptcha:u.default.mark(function e(t,n){var r,a;return u.default.wrap(function(e){while(1)switch(e.prev=e.next){case 0:return r=t.payload,a=n.call,e.next=4,a(c.getFakeCaptcha,r);case 4:case"end":return e.stop()}},e)})},reducers:{changeLoginStatus:function(e,t){var n=t.payload;return(0,o.setAuthority)(n.currentAuthority),(0,a.default)({},e,{status:n.status,type:n.type})}}},p=s;t.default=p},WHbk:function(e,t,n){"use strict";var r=n("g09b");Object.defineProperty(t,"__esModule",{value:!0}),t.fakeAccountLogin=c,t.getFakeCaptcha=s;var a=r(n("d6i3")),u=r(n("1l/V")),i=r(n("sy1d"));function c(e){return o.apply(this,arguments)}function o(){return o=(0,u.default)(a.default.mark(function e(t){return a.default.wrap(function(e){while(1)switch(e.prev=e.next){case 0:return e.abrupt("return",(0,i.default)("/api/login/account",{method:"POST",data:t}));case 1:case"end":return e.stop()}},e)})),o.apply(this,arguments)}function s(e){return p.apply(this,arguments)}function p(){return p=(0,u.default)(a.default.mark(function e(t){return a.default.wrap(function(e){while(1)switch(e.prev=e.next){case 0:return e.abrupt("return",(0,i.default)("/api/login/captcha?mobile=".concat(t)));case 1:case"end":return e.stop()}},e)})),p.apply(this,arguments)}},"sH+L":function(e,t,n){"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.getPageQuery=a,t.setAuthority=u;var r=n("Qyje");function a(){return(0,r.parse)(window.location.href.split("?")[1])}function u(e){var t="string"===typeof e?[e]:e;return localStorage.setItem("antd-pro-authority",JSON.stringify(t))}}}]);