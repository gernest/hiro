(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([[31],{qpEo:function(e,a,t){"use strict";var l=t("g09b");Object.defineProperty(a,"__esModule",{value:!0}),a.default=void 0,t("jCWc");var s=l(t("kPKH"));t("14J3");var d=l(t("BMrR"));t("IzEo");var n=l(t("bx4M"));t("Znn+");var r=l(t("ZTPi"));t("iQDF");for(var i=l(t("+eQT")),u=t("Y2fQ"),f=l(t("q1tI")),m=l(t("ZhIB")),c=t("ucLW"),o=l(t("WD6q")),g=i.default.RangePicker,E=r.default.TabPane,y=[],k=0;k<7;k+=1)y.push({title:(0,u.formatMessage)({id:"dashboard-analysis.analysis.test"},{no:k}),total:323234});var h=function(e){var a=e.rangePickerValue,t=e.salesData,l=e.isActive,i=e.handleRangePickerChange,k=e.loading,h=e.selectDate;return f.default.createElement(n.default,{loading:k,bordered:!1,bodyStyle:{padding:0}},f.default.createElement("div",{className:o.default.salesCard},f.default.createElement(r.default,{tabBarExtraContent:f.default.createElement("div",{className:o.default.salesExtraWrap},f.default.createElement("div",{className:o.default.salesExtra},f.default.createElement("a",{className:l("today"),onClick:function(){return h("today")}},f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.all-day",defaultMessage:"All Day"})),f.default.createElement("a",{className:l("week"),onClick:function(){return h("week")}},f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.all-week",defaultMessage:"All Week"})),f.default.createElement("a",{className:l("month"),onClick:function(){return h("month")}},f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.all-month",defaultMessage:"All Month"})),f.default.createElement("a",{className:l("year"),onClick:function(){return h("year")}},f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.all-year",defaultMessage:"All Year"}))),f.default.createElement(g,{value:a,onChange:i,style:{width:256}})),size:"large",tabBarStyle:{marginBottom:24}},f.default.createElement(E,{tab:f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.sales",defaultMessage:"Sales"}),key:"sales"},f.default.createElement(d.default,{type:"flex"},f.default.createElement(s.default,{xl:16,lg:12,md:12,sm:24,xs:24},f.default.createElement("div",{className:o.default.salesBar},f.default.createElement(c.Bar,{height:295,title:f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.sales-trend",defaultMessage:"Sales Trend"}),data:t}))),f.default.createElement(s.default,{xl:8,lg:12,md:12,sm:24,xs:24},f.default.createElement("div",{className:o.default.salesRank},f.default.createElement("h4",{className:o.default.rankingTitle},f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.sales-ranking",defaultMessage:"Sales Ranking"})),f.default.createElement("ul",{className:o.default.rankingList},y.map(function(e,a){return f.default.createElement("li",{key:e.title},f.default.createElement("span",{className:"".concat(o.default.rankingItemNumber," ").concat(a<3?o.default.active:"")},a+1),f.default.createElement("span",{className:o.default.rankingItemTitle,title:e.title},e.title),f.default.createElement("span",{className:o.default.rankingItemValue},(0,m.default)(e.total).format("0,0")))})))))),f.default.createElement(E,{tab:f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.visits",defaultMessage:"Visits"}),key:"views"},f.default.createElement(d.default,null,f.default.createElement(s.default,{xl:16,lg:12,md:12,sm:24,xs:24},f.default.createElement("div",{className:o.default.salesBar},f.default.createElement(c.Bar,{height:292,title:f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.visits-trend",defaultMessage:"Visits Trend"}),data:t}))),f.default.createElement(s.default,{xl:8,lg:12,md:12,sm:24,xs:24},f.default.createElement("div",{className:o.default.salesRank},f.default.createElement("h4",{className:o.default.rankingTitle},f.default.createElement(u.FormattedMessage,{id:"dashboard-analysis.analysis.visits-ranking",defaultMessage:"Visits Ranking"})),f.default.createElement("ul",{className:o.default.rankingList},y.map(function(e,a){return f.default.createElement("li",{key:e.title},f.default.createElement("span",{className:"".concat(o.default.rankingItemNumber," ").concat(a<3?o.default.active:"")},a+1),f.default.createElement("span",{className:o.default.rankingItemTitle,title:e.title},e.title),f.default.createElement("span",null,(0,m.default)(e.total).format("0,0")))})))))))))},b=h;a.default=b}}]);